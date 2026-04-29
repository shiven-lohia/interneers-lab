import "./App.css";
import Navbar from "./Navbar";
import ProductList from "./ProductList";
import ProductPage from "./ProductPage";
import CategoryListPage from "./CategoryListPage";
import CategoryPage from "./CategoryPage";

import { Routes, Route, Navigate } from "react-router-dom";

function App() {
  return (
    <div>
      <Navbar />

      <Routes>
        <Route path="/" element={<Navigate to="/products" />} />

        <Route path="/products" element={<ProductList />} />
        <Route path="/products/:id" element={<ProductPage />} />

        <Route path="/categories" element={<CategoryListPage />} />
        <Route path="/categories/:id" element={<CategoryPage />} />
      </Routes>
    </div>
  );
}

export default App;
