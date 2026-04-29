import { useEffect, useState } from "react";
import Product from "./Product";
import { useNavigate } from "react-router-dom";
import "./ProductList.css";
import LoadingSpinner from "./LoadingSpinner";

type ProductType = {
  id: string;
  name: string;
  brand: string;
  price: number;
  quantity: number;
  description?: string;
  category?: {
    id: string;
    title: string;
  };
};

function ProductList() {
  const navigate = useNavigate();
  const [products, setProducts] = useState<ProductType[]>([]);
  const [selectedId, setSelectedId] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch("http://localhost:8080/products")
      .then((res) => res.json())
      .then((data) => {
        setProducts(data);
        setLoading(false);
      })
      .catch((err) => console.error("Error fetching products:", err));
  }, []);

  const handleClick = (id: string) => {
    if (selectedId === id) {
      setSelectedId(null);
    } else {
      setSelectedId(id);
    }
  };

  const groupedProducts: Record<string, ProductType[]> = {};

  products.forEach((product) => {
    const categoryName = product.category?.title || "Uncategorized";

    if (!groupedProducts[categoryName]) {
      groupedProducts[categoryName] = [];
    }

    groupedProducts[categoryName].push(product);
  });

  if (loading) return <LoadingSpinner />;

  return (
    <div className="product-list">
      {Object.entries(groupedProducts).map(
        ([categoryName, categoryProducts]) => (
          <div key={categoryName} className="category-section">
            <h2 className="category-heading">{categoryName}</h2>

            <div className="category-products">
              {categoryProducts.map((product) => (
                <div key={product.id} className="product-wrapper">
                  <div onClick={() => navigate(`/products/${product.id}`)}>
                    <Product
                      name={product.name}
                      brand={product.brand}
                      price={product.price}
                      quantity={product.quantity}
                    />
                  </div>

                  {selectedId === product.id && (
                    <div className="product-details">
                      <p>Stock: {product.quantity} available</p>

                      {product.category && (
                        <p>Category: {product.category.title}</p>
                      )}

                      {product.description &&
                        product.description.trim() !== "" && (
                          <p>Description: {product.description}</p>
                        )}
                    </div>
                  )}
                </div>
              ))}
            </div>
          </div>
        ),
      )}
    </div>
  );
}

export default ProductList;
