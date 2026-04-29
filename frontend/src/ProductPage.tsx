import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import "./ProductPage.css";
import LoadingSpinner from "./LoadingSpinner";

type CategoryType = {
  id: string;
  title: string;
  description?: string;
};

type ProductType = {
  id: string;
  name: string;
  brand: string;
  price: number;
  quantity: number;
  description?: string;
  category_id?: string;
  category?: {
    id: string;
    title: string;
  };
};

function ProductPage() {
  const { id } = useParams();

  const [product, setProduct] = useState<ProductType | null>(null);
  const [categories, setCategories] = useState<CategoryType[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [message, setMessage] = useState("");
  const [isEditing, setIsEditing] = useState(false);

  const [showNewCategory, setShowNewCategory] = useState(false);
  const [newCategoryTitle, setNewCategoryTitle] = useState("");

  useEffect(() => {
    if (!id) return;

    Promise.all([
      fetch(`http://localhost:8080/products/${id}`).then((res) => {
        if (!res.ok) {
          throw new Error("Failed to fetch product");
        }
        return res.json();
      }),
      fetch("http://localhost:8080/categories").then((res) => {
        if (!res.ok) {
          throw new Error("Failed to fetch categories");
        }
        return res.json();
      }),
    ])
      .then(([productData, categoryData]) => {
        setProduct({
          ...productData,
          category_id:
            productData.category_id || productData.category?.id || "",
        });

        setCategories(categoryData);
        setLoading(false);
      })
      .catch((err) => {
        setError(err.message);
        setLoading(false);
      });
  }, [id]);

  const handleChange = (
    e: React.ChangeEvent<
      HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement
    >,
  ) => {
    if (!product) return;

    const { name, value } = e.target;

    setProduct({
      ...product,
      [name]: name === "price" || name === "quantity" ? Number(value) : value,
    });
  };

  const handleCreateCategory = () => {
    if (!newCategoryTitle.trim()) {
      setMessage("❌ Category title is required");
      return;
    }

    fetch("http://localhost:8080/categories", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        title: newCategoryTitle,
        description: "",
      }),
    })
      .then(async (res) => {
        if (!res.ok) {
          const text = await res.text();
          throw new Error(text);
        }
        return res.json();
      })
      .then((newCategory) => {
        setCategories([...categories, newCategory]);

        if (product) {
          setProduct({
            ...product,
            category_id: newCategory.id,
          });
        }

        setNewCategoryTitle("");
        setShowNewCategory(false);
        setMessage("✅ Category created");
      })
      .catch((err) => {
        setMessage("❌ " + err.message);
      });
  };

  const handleSave = () => {
    if (!product) return;

    const payload = {
      ...product,
      category: undefined,
    };

    fetch(`http://localhost:8080/products/${id}`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(payload),
    })
      .then(async (res) => {
        if (!res.ok) {
          const text = await res.text();
          throw new Error(text);
        }
        return res.json();
      })
      .then((updated) => {
        setProduct({
          ...updated,
          category_id: updated.category_id || updated.category?.id || "",
        });

        setMessage("✅ Product updated successfully");
        setIsEditing(false);
      })
      .catch((err) => {
        setMessage("❌ " + err.message);
      });
  };

  if (loading) return <LoadingSpinner />;
  if (error) return <p>Error: {error}</p>;
  if (!product) return <p>Product not found</p>;

  return (
    <div className="product-page">
      {!isEditing ? (
        <>
          <h2>{product.name}</h2>

          <p>
            <strong>Brand:</strong> {product.brand}
          </p>

          <p>
            <strong>Price:</strong> ₹{product.price}
          </p>

          <p>
            <strong>Stock:</strong> {product.quantity}
          </p>

          {product.category && (
            <p>
              <strong>Category:</strong>{" "}
              <a href={`/categories/${product.category.id}`}>
                {product.category.title}
              </a>
            </p>
          )}

          {product.description && product.description.trim() !== "" && (
            <p>
              <strong>Description:</strong> {product.description}
            </p>
          )}

          <button
            className="page-button edit-button"
            onClick={() => setIsEditing(true)}
          >
            Edit Product
          </button>
        </>
      ) : (
        <>
          <h2>Edit Product</h2>

          <label>Name:</label>
          <br />
          <input name="name" value={product.name} onChange={handleChange} />
          <br />
          <br />

          <label>Brand:</label>
          <br />
          <input name="brand" value={product.brand} onChange={handleChange} />
          <br />
          <br />

          <label>Price:</label>
          <br />
          <input
            type="number"
            name="price"
            value={product.price}
            onChange={handleChange}
          />
          <br />
          <br />

          <label>Quantity:</label>
          <br />
          <input
            type="number"
            name="quantity"
            value={product.quantity}
            onChange={handleChange}
          />
          <br />
          <br />

          <label>Category:</label>
          <br />
          <select
            name="category_id"
            value={product.category_id || ""}
            onChange={handleChange}
          >
            <option value="">Select category</option>

            {categories.map((cat) => (
              <option key={cat.id} value={cat.id}>
                {cat.title}
              </option>
            ))}
          </select>

          <br />
          <br />

          {!showNewCategory ? (
            <button
              className="page-button edit-button"
              onClick={() => setShowNewCategory(true)}
            >
              + New Category
            </button>
          ) : (
            <>
              <input
                placeholder="Category title"
                value={newCategoryTitle}
                onChange={(e) => setNewCategoryTitle(e.target.value)}
              />

              <button
                className="page-button save-button"
                onClick={handleCreateCategory}
                style={{ marginLeft: "10px" }}
              >
                Create
              </button>

              <button
                className="page-button cancel-button"
                onClick={() => {
                  setShowNewCategory(false);
                  setNewCategoryTitle("");
                }}
              >
                Cancel
              </button>
            </>
          )}

          <br />
          <br />

          <label>Description:</label>
          <br />
          <textarea
            name="description"
            value={product.description || ""}
            onChange={handleChange}
          />
          <br />
          <br />

          <button className="page-button save-button" onClick={handleSave}>
            Save Product
          </button>

          <button
            className="page-button cancel-button"
            onClick={() => setIsEditing(false)}
          >
            Cancel Edit
          </button>
        </>
      )}

      <p>{message}</p>
    </div>
  );
}

export default ProductPage;
