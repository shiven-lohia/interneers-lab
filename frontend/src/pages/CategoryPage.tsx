import { useEffect, useLayoutEffect, useRef, useState } from "react";
import { useParams, useNavigate, Link } from "react-router-dom";
import PageShell from "../components/layout/PageShell";
import ProductCard from "../components/product/ProductCard";
import ViewToggle from "../components/ui/ViewToggle";
import LoadingSpinner from "../components/ui/LoadingSpinner";
import ErrorMessage from "../components/ui/ErrorMessage";
import { getProductsByCategory } from "../api/products";
import { getCategory } from "../api/categories";
import type { Product, Category } from "../types";
import "./CategoryPage.css";

type SortCol = "name" | "brand" | "price" | "stock";
interface SortState {
  col: SortCol;
  dir: "asc" | "desc";
}

function CategoryPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

  const [category, setCategory] = useState<Category | null>(null);
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

  const sortedProducts = [...products].sort((a, b) => {
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

  const cols: { key: SortCol; label: string; numeric?: boolean }[] = [
    { key: "name", label: "Name" },
    { key: "brand", label: "Brand" },
    { key: "price", label: "Price", numeric: true },
    { key: "stock", label: "Stock", numeric: true },
  ];

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
          <div className="category-products__heading-right">
            <span className="category-products__count">{products.length}</span>
            {products.length > 0 && (
              <ViewToggle view={view} onChange={handleViewChange} />
            )}
          </div>
        </h2>

        {products.length === 0 ? (
          <p className="category-products__empty">
            No products in this category.
          </p>
        ) : view === "grid" ? (
          <div className="category-products__grid">
            {products.map((product) => (
              <ProductCard
                key={product.id}
                product={product}
                onClick={() => navigate(`/products/${product.id}`)}
              />
            ))}
          </div>
        ) : (
          <div className="category-products__table-wrap">
            <table className="category-products__table">
              <thead>
                <tr>
                  {cols.map(({ key, label, numeric }) => (
                    <th
                      key={key}
                      className={`category-products__th${numeric ? " category-products__th--num" : ""}${sort.col === key ? " category-products__th--active" : ""}`}
                      onClick={() => handleSort(key)}
                      onMouseDown={(e) => e.preventDefault()}
                    >
                      {label}
                      {sort.col === key && (
                        <span className="category-products__sort-icon">
                          {sort.dir === "asc" ? " ↑" : " ↓"}
                        </span>
                      )}
                    </th>
                  ))}
                </tr>
              </thead>
              <tbody>
                {sortedProducts.map((product) => (
                  <tr
                    key={product.id}
                    className="category-products__tr"
                    onClick={() => navigate(`/products/${product.id}`)}
                  >
                    <td className="category-products__td category-products__td--name">
                      {product.name}
                    </td>
                    <td className="category-products__td">{product.brand}</td>
                    <td className="category-products__td category-products__td--num">
                      ₹{product.price.toLocaleString("en-IN")}
                    </td>
                    <td className="category-products__td category-products__td--num">
                      {product.quantity}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </section>
    </PageShell>
  );
}

export default CategoryPage;
