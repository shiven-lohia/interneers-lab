import { useEffect, useState } from "react";
import { useParams, useNavigate, Link } from "react-router-dom";
import PageShell from "../components/layout/PageShell";
import ProductCard from "../components/product/ProductCard";
import LoadingSpinner from "../components/ui/LoadingSpinner";
import ErrorMessage from "../components/ui/ErrorMessage";
import { getProductsByCategory } from "../api/products";
import { getCategory } from "../api/categories";
import type { Product, Category } from "../types";
import "./CategoryPage.css";

function CategoryPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

  const [category, setCategory] = useState<Category | null>(null);
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    if (!id) return;

    Promise.all([getCategory(id), getProductsByCategory(id)])
      .then(([categoryData, productsData]) => {
        setCategory(categoryData);
        setProducts(productsData);
        setLoading(false);
      })
      .catch((err) => {
        setError(err.message);
        setLoading(false);
      });
  }, [id]);

  if (loading) return <LoadingSpinner />;
  if (error)
    return (
      <PageShell>
        <ErrorMessage message={error} />
      </PageShell>
    );
  if (!category)
    return (
      <PageShell>
        <p>Category not found</p>
      </PageShell>
    );

  const hasDescription =
    category.description && category.description.trim() !== "";

  return (
    <PageShell>
      <Link to="/categories" className="back-link">
        ← All categories
      </Link>

      <div className="category-masthead">
        <h1 className="category-masthead__title">{category.title}</h1>
        {hasDescription && (
          <p className="category-masthead__description">
            {category.description}
          </p>
        )}
      </div>

      <section className="category-products">
        <h2 className="category-products__heading">
          Products
          <span className="category-products__count">{products.length}</span>
        </h2>

        {products.length === 0 ? (
          <p className="category-products__empty">
            No products in this category.
          </p>
        ) : (
          <div className="category-products__grid">
            {products.map((product) => (
              <ProductCard
                key={product.id}
                product={product}
                onClick={() => navigate(`/products/${product.id}`)}
              />
            ))}
          </div>
        )}
      </section>
    </PageShell>
  );
}

export default CategoryPage;
