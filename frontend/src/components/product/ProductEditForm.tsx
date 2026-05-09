import { useState, type ChangeEvent } from "react";
import Button from "../ui/Button";
import ErrorMessage from "../ui/ErrorMessage";
import "./ProductEditForm.css";
import { createCategory } from "../../api/categories";
import type { Category, ProductFormData } from "../../types";

interface ProductEditFormProps {
  product: ProductFormData;
  categories: Category[];
  onChange: (
    e: ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>,
  ) => void;
  onSave: () => void;
  onCancel: () => void;
  onCategoryCreated: (category: Category) => void;
  onDelete: () => Promise<void>;
  message: string;
}

function ProductEditForm({
  product,
  categories,
  onChange,
  onSave,
  onCancel,
  onCategoryCreated,
  onDelete,
  message,
}: ProductEditFormProps) {
  const [showNewCategory, setShowNewCategory] = useState(false);
  const [newCategoryTitle, setNewCategoryTitle] = useState("");
  const [catError, setCatError] = useState("");
  const [confirmDelete, setConfirmDelete] = useState(false);
  const [deleting, setDeleting] = useState(false);
  const [deleteError, setDeleteError] = useState("");

  const handleDeleteConfirm = async () => {
    setDeleting(true);
    try {
      await onDelete();
    } catch (err: unknown) {
      setDeleteError(err instanceof Error ? err.message : "Failed to delete");
      setDeleting(false);
      setConfirmDelete(false);
    }
  };

  const handleCreateCategory = () => {
    if (!newCategoryTitle.trim()) {
      setCatError("Category title is required");
      return;
    }

    createCategory({ title: newCategoryTitle, description: "" })
      .then((newCat) => {
        onCategoryCreated(newCat);
        setNewCategoryTitle("");
        setShowNewCategory(false);
        setCatError("");
      })
      .catch((err) => setCatError(err.message));
  };

  return (
    <div className="product-edit-form">
      <h2>Edit Product</h2>

      <div className="form-grid">
        <div className="form-field">
          <label htmlFor="name">Name</label>
          <input
            id="name"
            name="name"
            value={product.name}
            onChange={onChange}
          />
        </div>

        <div className="form-field">
          <label htmlFor="brand">Brand</label>
          <input
            id="brand"
            name="brand"
            value={product.brand}
            onChange={onChange}
          />
        </div>

        <div className="form-field">
          <label htmlFor="price">Price</label>
          <input
            id="price"
            type="number"
            name="price"
            value={product.price}
            onChange={onChange}
          />
        </div>

        <div className="form-field">
          <label htmlFor="quantity">Quantity</label>
          <input
            id="quantity"
            type="number"
            name="quantity"
            value={product.quantity}
            onChange={onChange}
          />
        </div>

        <div className="form-field">
          <label htmlFor="category_id">Category</label>
          <select
            id="category_id"
            name="category_id"
            value={product.category_id || ""}
            onChange={onChange}
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
            value={product.description || ""}
            onChange={onChange}
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
        <Button variant="primary" onClick={onSave}>
          Save Product
        </Button>
        <Button variant="secondary" onClick={onCancel}>
          Cancel
        </Button>
      </div>

      {message && <p className="form-message">{message}</p>}

      <div className="form-delete">
        {!confirmDelete ? (
          <button
            className="delete-link"
            onClick={() => setConfirmDelete(true)}
          >
            Delete product
          </button>
        ) : (
          <span className="delete-confirm">
            Delete this product permanently?{" "}
            <button
              className="delete-confirm__yes"
              onClick={handleDeleteConfirm}
              disabled={deleting}
            >
              {deleting ? "Deleting…" : "Yes, delete"}
            </button>
            {" · "}
            <button
              className="delete-confirm__cancel"
              onClick={() => {
                setConfirmDelete(false);
                setDeleteError("");
              }}
            >
              Cancel
            </button>
          </span>
        )}
        {deleteError && <p className="delete-error">{deleteError}</p>}
      </div>
    </div>
  );
}

export default ProductEditForm;
