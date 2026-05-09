import { useState, useEffect, type ChangeEvent } from "react";
import { useNavigate, Link } from "react-router-dom";
import PageShell from "../components/layout/PageShell";
import Button from "../components/ui/Button";
import ErrorMessage from "../components/ui/ErrorMessage";
import { createProduct } from "../api/products";
import { getCategories, createCategory } from "../api/categories";
import type { Category, ProductFormData } from "../types";
import "../components/product/ProductEditForm.css";
import "./AddProductPage.css";

const EMPTY_FORM: ProductFormData = {
  id: "",
  name: "",
  brand: "",
  price: 0,
  quantity: 0,
  category_id: "",
  description: "",
};

function AddProductPage() {
  const navigate = useNavigate();
  const [form, setForm] = useState<ProductFormData>(EMPTY_FORM);
  const [categories, setCategories] = useState<Category[]>([]);
  const [error, setError] = useState("");
  const [saving, setSaving] = useState(false);
  const [showNewCategory, setShowNewCategory] = useState(false);
  const [newCategoryTitle, setNewCategoryTitle] = useState("");
  const [catError, setCatError] = useState("");

  useEffect(() => {
    getCategories()
      .then(setCategories)
      .catch(() => {});
  }, []);

  const handleChange = (
    e: ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>,
  ) => {
    const { name, value } = e.target;
    setForm((prev) => ({
      ...prev,
      [name]: name === "price" || name === "quantity" ? Number(value) : value,
    }));
  };

  const handleSave = async () => {
    if (!form.name.trim()) {
      setError("Name is required");
      return;
    }
    if (!form.brand.trim()) {
      setError("Brand is required");
      return;
    }
    setSaving(true);
    setError("");
    try {
      await createProduct(form);
      navigate("/products");
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : "Failed to create product");
      setSaving(false);
    }
  };

  const handleCreateCategory = () => {
    if (!newCategoryTitle.trim()) {
      setCatError("Category title is required");
      return;
    }
    createCategory({ title: newCategoryTitle, description: "" })
      .then((newCat) => {
        setCategories((prev) => [...prev, newCat]);
        setForm((prev) => ({ ...prev, category_id: newCat.id }));
        setNewCategoryTitle("");
        setShowNewCategory(false);
        setCatError("");
      })
      .catch((err) => setCatError(err.message));
  };

  return (
    <PageShell>
      <Link to="/products" className="back-link">
        ← Products
      </Link>

      <div className="add-product-form">
        <h1>New Product</h1>

        {error && <ErrorMessage message={error} />}

        <div className="form-grid">
          <div className="form-field">
            <label htmlFor="name">Name</label>
            <input
              id="name"
              name="name"
              value={form.name}
              onChange={handleChange}
              autoFocus
            />
          </div>

          <div className="form-field">
            <label htmlFor="brand">Brand</label>
            <input
              id="brand"
              name="brand"
              value={form.brand}
              onChange={handleChange}
            />
          </div>

          <div className="form-field">
            <label htmlFor="price">Price</label>
            <input
              id="price"
              type="number"
              name="price"
              min="0"
              value={form.price}
              onChange={handleChange}
            />
          </div>

          <div className="form-field">
            <label htmlFor="quantity">Quantity</label>
            <input
              id="quantity"
              type="number"
              name="quantity"
              min="0"
              value={form.quantity}
              onChange={handleChange}
            />
          </div>

          <div className="form-field">
            <label htmlFor="category_id">Category</label>
            <select
              id="category_id"
              name="category_id"
              value={form.category_id || ""}
              onChange={handleChange}
            >
              <option value="">Select category</option>
              {categories.map((cat) => (
                <option key={cat.id} value={cat.id}>
                  {cat.title}
                </option>
              ))}
            </select>
          </div>

          <div className="form-field form-field--full">
            <label htmlFor="description">Description</label>
            <textarea
              id="description"
              name="description"
              value={form.description || ""}
              onChange={handleChange}
              rows={3}
            />
          </div>
        </div>

        {!showNewCategory ? (
          <Button variant="secondary" onClick={() => setShowNewCategory(true)}>
            + New Category
          </Button>
        ) : (
          <div className="new-category">
            <input
              placeholder="Category title"
              value={newCategoryTitle}
              onChange={(e) => setNewCategoryTitle(e.target.value)}
            />
            <Button variant="primary" onClick={handleCreateCategory}>
              Create
            </Button>
            <Button
              variant="secondary"
              onClick={() => {
                setShowNewCategory(false);
                setNewCategoryTitle("");
                setCatError("");
              }}
            >
              Cancel
            </Button>
            {catError && <ErrorMessage message={catError} />}
          </div>
        )}

        <div className="form-actions">
          <Button variant="primary" onClick={handleSave} disabled={saving}>
            {saving ? "Creating…" : "Create Product"}
          </Button>
          <Button variant="secondary" onClick={() => navigate("/products")}>
            Cancel
          </Button>
        </div>
      </div>
    </PageShell>
  );
}

export default AddProductPage;
