import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import PageShell from "../components/layout/PageShell";
import CategoryCard from "../components/category/CategoryCard";
import LoadingSpinner from "../components/ui/LoadingSpinner";
import ErrorMessage from "../components/ui/ErrorMessage";
import { getCategories } from "../api/categories";
import type { Category } from "../types";
import "./CategoryListPage.css";

function CategoryListPage() {
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const navigate = useNavigate();

  useEffect(() => {
    getCategories()
      .then((data) => {
        setCategories(data);
        setLoading(false);
      })
      .catch((err) => {
        setError(err.message);
        setLoading(false);
      });
  }, []);

  if (loading) return <LoadingSpinner />;

  const isEmpty = categories.length === 0;

  return (
    <PageShell>
      <div className="category-list__header">
        <h1 className="category-list__title">Categories</h1>
        <button
          className="category-list__add-btn"
          onClick={() => navigate("/categories/new")}
        >
          Add category
        </button>
      </div>

      {error && <ErrorMessage message={error} />}

      {isEmpty && !error && (
        <div className="category-list__empty">
          <p className="category-list__empty-message">No categories yet.</p>
          <button
            className="category-list__add-btn"
            onClick={() => navigate("/categories/new")}
          >
            Add your first category
          </button>
        </div>
      )}

      {!isEmpty && (
        <div className="category-list__items">
          {categories.map((cat) => (
            <CategoryCard
              key={cat.id}
              category={cat}
              onClick={() => navigate(`/categories/${cat.id}`)}
            />
          ))}
        </div>
      )}
    </PageShell>
  );
}

export default CategoryListPage;
