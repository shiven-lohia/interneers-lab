import "./Navbar.css";
import { Link } from "react-router-dom";

function Navbar() {
  return (
    <nav className="navbar">
      <h2 className="logo">My Store</h2>

      <div className="nav-links">
        <Link to="/products">Products</Link>
        <Link to="/categories">Categories</Link>
      </div>
    </nav>
  );
}

export default Navbar;
