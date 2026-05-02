import { useState, type ChangeEvent } from "react";
import { useNavigate, Link } from "react-router-dom";
import PageShell from "../components/layout/PageShell";
import Button from "../components/ui/Button";
import ErrorMessage from "../components/ui/ErrorMessage";
import { createCategory } from "../api/categories";
import "../components/product/ProductEditForm.css";
import "./AddCategoryPage.css";

function AddCategoryPage() {
  const navigate = useNavigate();
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [error, setError] = useState("");
  const [saving, setSaving] = useState(false);

  const handleChange = (
    e: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>,
  ) => {
    const { name, value } = e.target;
    if (name === "title") setTitle(value);
    if (name === "description") setDescription(value);
  };

  const handleSave = async () => {
    if (!title.trim()) {
      setError("Title is required");
      return;
    }
    setSaving(true);
    setError("");
    try {
      await createCategory({ title, description });
      navigate("/categories");
    } catch (err: unknown) {
      setError(
        err instanceof Error ? err.message : "Failed to create category",
      );
      setSaving(false);
    }
  };

  return (
    <PageShell>
      <Link to="/categories" className="back-link">
        ← Categories
      </Link>

      <div className="add-category-form">
        <h1>New Category</h1>

        {error && <ErrorMessage message={error} />}

        <div className="form-grid">
          <div className="form-field form-field--full">
            <label htmlFor="title">Title</label>
            <input
              id="title"
              name="title"
              value={title}
              onChange={handleChange}
              autoFocus
            />
          </div>

          <div className="form-field form-field--full">
            <label htmlFor="description">Description</label>
            <textarea
              id="description"
              name="description"
              value={description}
              onChange={handleChange}
              rows={3}
            />
          </div>
        </div>

        <div className="form-actions">
          <Button variant="primary" onClick={handleSave} disabled={saving}>
            {saving ? "Creating…" : "Create Category"}
          </Button>
          <Button variant="secondary" onClick={() => navigate("/categories")}>
            Cancel
          </Button>
        </div>
      </div>
    </PageShell>
  );
}

export default AddCategoryPage;
