import Card from "../ui/Card";
import "./ProductCard.css";
import type { Product } from "../../types";

interface ProductCardProps {
  product: Pick<Product, "name" | "brand" | "price" | "quantity">;
  onClick?: () => void;
}

function ProductCard({ product, onClick }: ProductCardProps) {
  const { name, brand, price, quantity } = product;
  return (
    <Card onClick={onClick} className="product-card">
      <div className="product-card__body">
        <h3 className="product-card__name">{name}</h3>
        <p className="product-card__brand">{brand}</p>
      </div>
      <div className="product-card__footer">
        <span className="product-card__price">
          ₹{price.toLocaleString("en-IN")}
        </span>
        <span className="product-card__stock">{quantity} in stock</span>
      </div>
    </Card>
  );
}

export default ProductCard;
