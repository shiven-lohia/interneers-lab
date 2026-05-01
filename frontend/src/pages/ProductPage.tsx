import { useEffect, useState, type ChangeEvent } from "react";
import { useParams, Link } from "react-router-dom";
import PageShell from "../components/layout/PageShell";
import LoadingSpinner from "../components/ui/LoadingSpinner";
import ErrorMessage from "../components/ui/ErrorMessage";
import Button from "../components/ui/Button";
import ProductEditForm from "../components/product/ProductEditForm";
import { getProduct, updateProduct } from "../api/products";
import { getCategories } from "../api/categories";
import type { Product, Category, ProductFormData } from "../types";
import "./ProductPage.css";

function ProductPage() {
  const { id } = useParams<{ id: string }>();

  const [product, setProduct] = useState<Product | null>(null);
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [message, setMessage] = useState("");
  const [isEditing, setIsEditing] = useState(false);

  useEffect(() => {
    if (!id) return;

    Promise.all([getProduct(id), getCategories()])
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
    e: ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>,
  ) => {
    if (!product) return;
    const { name, value } = e.target;
    setProduct({
      ...product,
      [name]: name === "price" || name === "quantity" ? Number(value) : value,
    });
  };

  const handleSave = () => {
    if (!product || !id) return;
    const payload: ProductFormData = {
      id: product.id,
      name: product.name,
      description: product.description,
      category_id: product.category_id,
      price: product.price,
      brand: product.brand,
      quantity: product.quantity,
    };
    updateProduct(id, payload)
      .then((updated) => {
        setProduct({
          ...updated,
          category_id: updated.category_id || updated.category?.id || "",
        });
        setMessage("Saved");
        setIsEditing(false);
      })
      .catch((err) => setMessage(err.message));
  };

  const handleCategoryCreated = (newCat: Category) => {
    setCategories((prev) => [...prev, newCat]);
    if (product) setProduct({ ...product, category_id: newCat.id });
  };

  if (loading) return <LoadingSpinner />;
  if (error)
    return (
      <PageShell>
        <ErrorMessage message={error} />
      </PageShell>
    );
  if (!product)
    return (
      <PageShell>
        <p>Product not found</p>
      </PageShell>
    );

  return (
    <PageShell>
      <Link to="/products" className="back-link">
        ← All products
      </Link>

      {isEditing ? (
        <ProductEditForm
          product={product as ProductFormData}
          categories={categories}
          onChange={handleChange}
          onSave={handleSave}
          onCancel={() => {
            setIsEditing(false);
            setMessage("");
          }}
          onCategoryCreated={handleCategoryCreated}
          message={message}
        />
      ) : (
        <div className="product-detail">
          <h1 className="product-detail__name">{product.name}</h1>
          <dl className="product-detail__fields">
            <dt>Brand</dt>
            <dd>{product.brand}</dd>
            <dt>Price</dt>
            <dd>₹{product.price.toLocaleString("en-IN")}</dd>
            <dt>Stock</dt>
            <dd>{product.quantity}</dd>
            {product.category && (
              <>
                <dt>Category</dt>
                <dd>
                  <Link
                    to={`/categories/${product.category.id}`}
                    className="detail-link"
                  >
                    {product.category.title}
                  </Link>
                </dd>
              </>
            )}
            {product.description && product.description.trim() !== "" && (
              <>
                <dt>Description</dt>
                <dd className="product-detail__description">
                  {product.description}
                </dd>
              </>
            )}
          </dl>
          {message && <p className="product-detail__message">{message}</p>}
          <Button onClick={() => setIsEditing(true)}>Edit product</Button>
        </div>
      )}
    </PageShell>
  );
}

export default ProductPage;
