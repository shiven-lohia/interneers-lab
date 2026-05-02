import { useEffect, useLayoutEffect, useRef, useState } from "react";
import { useNavigate } from "react-router-dom";
import PageShell from "../components/layout/PageShell";
import ProductCard from "../components/product/ProductCard";
import ViewToggle from "../components/ui/ViewToggle";
import LoadingSpinner from "../components/ui/LoadingSpinner";
import ErrorMessage from "../components/ui/ErrorMessage";
import { getProducts } from "../api/products";
import type { Product } from "../types";
import "./ProductListPage.css";

type SortCol = "name" | "brand" | "price" | "stock";
interface SortState {
  col: SortCol;
  dir: "asc" | "desc";
}

function sortItems(items: Product[], sort: SortState): Product[] {
  return [...items].sort((a, b) => {
    let av: string | number;
    let bv: string | number;
    switch (sort.col) {
      case "name":
        av = a.name.toLowerCase();
        bv = b.name.toLowerCase();
        break;
      case "brand":
        av = a.brand.toLowerCase();
        bv = b.brand.toLowerCase();
        break;
      case "price":
        av = a.price;
        bv = b.price;
        break;
      case "stock":
        av = a.quantity;
        bv = b.quantity;
        break;
      default:
        av = "";
        bv = "";
    }
    if (av < bv) return sort.dir === "asc" ? -1 : 1;
    if (av > bv) return sort.dir === "asc" ? 1 : -1;
    return 0;
  });
}

const LIST_COLS: { key: SortCol; label: string; numeric?: boolean }[] = [
  { key: "name", label: "Name" },
  { key: "brand", label: "Brand" },
  { key: "price", label: "Price", numeric: true },
  { key: "stock", label: "Stock", numeric: true },
];

function ProductListPage() {
  const navigate = useNavigate();
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [view, setView] = useState<"grid" | "list">(() => {
    return (
      (localStorage.getItem("productViewMode") as "grid" | "list") || "grid"
    );
  });
  const [sort, setSort] = useState<SortState>({ col: "name", dir: "asc" });
  const sortScrollRef = useRef(0);

  useLayoutEffect(() => {
    window.scrollTo(0, sortScrollRef.current);
  }, [sort]);

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

  const handleViewChange = (v: "grid" | "list") => {
    setView(v);
    localStorage.setItem("productViewMode", v);
  };

  const handleSort = (col: SortCol) => {
    sortScrollRef.current = window.scrollY;
    setSort((prev) =>
      prev.col === col
        ? { col, dir: prev.dir === "asc" ? "desc" : "asc" }
        : { col, dir: "asc" },
    );
  };

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
          <ViewToggle view={view} onChange={handleViewChange} />
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

      {!isEmpty &&
        Object.entries(grouped).map(([categoryName, items]) => (
          <section key={categoryName} className="category-section">
            <h2 className="category-section__heading">{categoryName}</h2>

            {view === "grid" ? (
              <div className="category-section__grid">
                {items.map((product) => (
                  <ProductCard
                    key={product.id}
                    product={product}
                    onClick={() => navigate(`/products/${product.id}`)}
                  />
                ))}
              </div>
            ) : (
              <div className="product-list__table-wrap">
                <table className="product-list__table">
                  <thead>
                    <tr>
                      {LIST_COLS.map(({ key, label, numeric }) => (
                        <th
                          key={key}
                          className={`product-list__th${numeric ? " product-list__th--num" : ""}${sort.col === key ? " product-list__th--active" : ""}`}
                          onClick={() => handleSort(key)}
                          onMouseDown={(e) => e.preventDefault()}
                        >
                          {label}
                          {sort.col === key && (
                            <span className="product-list__sort-icon">
                              {sort.dir === "asc" ? " ↑" : " ↓"}
                            </span>
                          )}
                        </th>
                      ))}
                    </tr>
                  </thead>
                  <tbody>
                    {sortItems(items, sort).map((product) => (
                      <tr
                        key={product.id}
                        className="product-list__tr"
                        onClick={() => navigate(`/products/${product.id}`)}
                      >
                        <td className="product-list__td product-list__td--name">
                          {product.name}
                        </td>
                        <td className="product-list__td">{product.brand}</td>
                        <td className="product-list__td product-list__td--num">
                          ₹{product.price.toLocaleString("en-IN")}
                        </td>
                        <td className="product-list__td product-list__td--num">
                          {product.quantity}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </section>
        ))}
    </PageShell>
  );
}

export default ProductListPage;
