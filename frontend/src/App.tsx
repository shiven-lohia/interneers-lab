import Navbar from "./components/layout/Navbar";
import ProductListPage from "./pages/ProductListPage";
import ProductPage from "./pages/ProductPage";
import BulkImportPage from "./pages/BulkImportPage";
import CategoryListPage from "./pages/CategoryListPage";
import CategoryPage from "./pages/CategoryPage";
import ReportsPage from "./pages/ReportsPage";
import { Routes, Route, Navigate } from "react-router-dom";

function App() {
  return (
    <div>
      <Navbar />
      <Routes>
        <Route path="/" element={<Navigate to="/products" />} />
        <Route path="/products" element={<ProductListPage />} />
        <Route path="/products/bulk" element={<BulkImportPage />} />
        <Route path="/products/:id" element={<ProductPage />} />
        <Route path="/categories" element={<CategoryListPage />} />
        <Route path="/categories/:id" element={<CategoryPage />} />
        <Route path="/reports" element={<ReportsPage />} />
      </Routes>
    </div>
  );
}

export default App;
