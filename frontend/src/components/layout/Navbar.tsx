import "./Navbar.css";
import { Link, NavLink } from "react-router-dom";

function Navbar() {
  return (
    <nav className="navbar">
      <Link to="/" className="navbar__logo">
        My Store
      </Link>
      <div className="navbar__links">
        <NavLink to="/products">Products</NavLink>
        <NavLink to="/categories">Categories</NavLink>
      </div>
    </nav>
  );
}

export default Navbar;
