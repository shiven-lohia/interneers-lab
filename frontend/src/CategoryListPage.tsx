import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import LoadingSpinner from "./LoadingSpinner";

type CategoryType = {
  id: string;
  title: string;
  description?: string;
};

function CategoryListPage() {
  const [categories, setCategories] = useState<CategoryType[]>([]);
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();

  useEffect(() => {
    fetch("http://localhost:8080/categories")
      .then((res) => res.json())
      .then((data) => {
        setCategories(data);
        setLoading(false);
      });
  }, []);

  if (loading) return <LoadingSpinner />;

  return (
    <div style={{ padding: "20px" }}>
      <h2>Categories</h2>

      {categories.map((cat) => {
        const hasDescription = cat.description && cat.description.trim() !== "";

        return (
          <div
            key={cat.id}
            onClick={() => navigate(`/categories/${cat.id}`)}
            style={{
              padding: "12px",
              marginBottom: "10px",
              border: "1px solid #ddd",
              borderRadius: "8px",
              cursor: "pointer",
            }}
          >
            <strong>{cat.title}</strong>

            <p
              style={{
                color: "#6b7280",
                fontStyle: "italic",
                marginTop: "6px",
              }}
            >
              {hasDescription ? cat.description : "No description"}
            </p>
          </div>
        );
      })}
    </div>
  );
}

export default CategoryListPage;
