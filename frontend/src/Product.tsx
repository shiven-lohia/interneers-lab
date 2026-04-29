import "./Product.css";

type ProductProps = {
  name: string;
  brand: string;
  price: number;
  quantity: number;
};

function Product({ name, brand, price, quantity }: ProductProps) {
  return (
    <div className="product-card">
      <h2>{name}</h2>
      <p>
        <strong>Brand:</strong> {brand}
      </p>
      <p>
        <strong>Price:</strong> ₹{price}
      </p>
      <p>
        <strong>Quantity:</strong> {quantity}
      </p>
    </div>
  );
}

export default Product;
