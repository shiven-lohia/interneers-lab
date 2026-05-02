import Navbar from "./components/layout/Navbar";
import ProductListPage from "./pages/ProductListPage";
import ProductPage from "./pages/ProductPage";
import BulkImportPage from "./pages/BulkImportPage";
import AddProductPage from "./pages/AddProductPage";
import CategoryListPage from "./pages/CategoryListPage";
import CategoryPage from "./pages/CategoryPage";
import AddCategoryPage from "./pages/AddCategoryPage";
import ReportsPage from "./pages/ReportsPage";
import { Routes, Route, Navigate } from "react-router-dom";
import { ThemeProvider } from "./context/ThemeContext";

function App() {
  return (
    <ThemeProvider>
      <div>
        <Navbar />
        <Routes>
          <Route path="/" element={<Navigate to="/products" />} />
          <Route path="/products" element={<ProductListPage />} />
          <Route path="/products/bulk" element={<BulkImportPage />} />
          <Route path="/products/new" element={<AddProductPage />} />
          <Route path="/products/:id" element={<ProductPage />} />
          <Route path="/categories" element={<CategoryListPage />} />
          <Route path="/categories/new" element={<AddCategoryPage />} />
          <Route path="/categories/:id" element={<CategoryPage />} />
          <Route path="/reports" element={<ReportsPage />} />
        </Routes>
      </div>
    </ThemeProvider>
  );
}

export default App;
