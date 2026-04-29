import { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import Product from "./Product";
import LoadingSpinner from "./LoadingSpinner";

type ProductType = {
  id: string;
  name: string;
  brand: string;
  price: number;
  quantity: number;
};

type CategoryType = {
  id: string;
  title: string;
  description?: string;
};

function CategoryPage() {
  const { id } = useParams();
  const navigate = useNavigate();

  const [category, setCategory] = useState<CategoryType | null>(null);
  const [products, setProducts] = useState<ProductType[]>([]);

  useEffect(() => {
    if (!id) return;

    fetch(`http://localhost:8080/categories/${id}`)
      .then((res) => res.json())
      .then((data) => setCategory(data));

    fetch(`http://localhost:8080/products?category_id=${id}`)
      .then((res) => res.json())
      .then((data) => setProducts(data));
  }, [id]);

  if (!category) return <LoadingSpinner />;

  const hasDescription =
    category.description && category.description.trim() !== "";

  return (
    <div style={{ padding: "20px" }}>
      <h2>{category.title}</h2>

      <p
        style={{
          color: "#6b7280",
          fontStyle: "italic",
          marginTop: "4px",
        }}
      >
        {hasDescription ? category.description : "No description"}
      </p>

      <h3>Products</h3>

      {products.map((product) => (
        <div
          key={product.id}
          onClick={() => navigate(`/products/${product.id}`)}
          style={{
            cursor: "pointer",
            marginBottom: "15px",
          }}
        >
          <Product
            name={product.name}
            brand={product.brand}
            price={product.price}
            quantity={product.quantity}
          />
        </div>
      ))}
    </div>
  );
}

export default CategoryPage;
