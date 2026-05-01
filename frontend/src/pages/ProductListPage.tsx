import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import PageShell from "../components/layout/PageShell";
import ProductCard from "../components/product/ProductCard";
import LoadingSpinner from "../components/ui/LoadingSpinner";
import ErrorMessage from "../components/ui/ErrorMessage";
import { getProducts } from "../api/products";
import type { Product } from "../types";
import "./ProductListPage.css";

function ProductListPage() {
  const navigate = useNavigate();
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    getProducts()
      .then((data) => {
        setProducts(data);
        setLoading(false);
      })
      .catch((err) => {
        setError(err.message);
        setLoading(false);
      });
  }, []);

  if (loading) return <LoadingSpinner />;

  const grouped: Record<string, Product[]> = {};
  products.forEach((p) => {
    const name = p.category?.title || "Uncategorized";
    if (!grouped[name]) grouped[name] = [];
    grouped[name].push(p);
  });

  const isEmpty = products.length === 0;

  return (
    <PageShell>
      <div className="product-list__header">
        <h1 className="product-list__title">Products</h1>
        <div className="product-list__header-actions">
          <button
            className="product-list__import-btn"
            onClick={() => navigate("/products/bulk")}
          >
            Import CSV
          </button>
          <button
            className="product-list__add-btn"
            onClick={() => navigate("/products/new")}
          >
            Add product
          </button>
        </div>
      </div>

      {error && <ErrorMessage message={error} />}

      {isEmpty && !error && (
        <div className="product-list__empty">
          <p className="product-list__empty-message">No products yet.</p>
          <button
            className="product-list__add-btn"
            onClick={() => navigate("/products/new")}
          >
            Add your first product
          </button>
        </div>
      )}

      {Object.entries(grouped).map(([categoryName, items]) => (
        <section key={categoryName} className="category-section">
          <h2 className="category-section__heading">{categoryName}</h2>
          <div className="category-section__grid">
            {items.map((product) => (
              <ProductCard
                key={product.id}
                product={product}
                onClick={() => navigate(`/products/${product.id}`)}
              />
            ))}
          </div>
        </section>
      ))}
    </PageShell>
  );
}

export default ProductListPage;
