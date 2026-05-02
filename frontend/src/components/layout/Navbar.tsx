import { useState, useEffect } from "react";
import "./Navbar.css";
import { Link, NavLink } from "react-router-dom";

function Navbar() {
  const [scrolled, setScrolled] = useState(false);

  useEffect(() => {
    const onScroll = () => setScrolled(window.scrollY > 0);
    window.addEventListener("scroll", onScroll, { passive: true });
    return () => window.removeEventListener("scroll", onScroll);
  }, []);

  return (
    <nav className={`navbar${scrolled ? " navbar--scrolled" : ""}`}>
      <Link to="/" className="navbar__logo">
        My Store
      </Link>
      <div className="navbar__links">
        <NavLink to="/products">Products</NavLink>
        <NavLink to="/categories">Categories</NavLink>
        <NavLink to="/reports">Reports</NavLink>
      </div>
    </nav>
  );
}

export default Navbar;
