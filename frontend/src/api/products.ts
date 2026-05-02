import { apiGet, apiPost, apiPut, apiPostMultipart, apiDelete } from "./client";
import type { Product, ProductFormData, BulkImportResult } from "../types";

export const getProducts = () => apiGet<Product[]>("/products");
export const createProduct = (data: ProductFormData) =>
  apiPost<Product>("/products", data);
export const getProduct = (id: string) => apiGet<Product>(`/products/${id}`);
export const getProductsByCategory = (categoryId: string) =>
  apiGet<Product[]>(`/products?category_id=${categoryId}`);
export const updateProduct = (id: string, data: ProductFormData) =>
  apiPut<Product>(`/products/${id}`, data);

export const deleteProduct = (id: string) => apiDelete(`/products/${id}`);

export const bulkCreateProducts = (file: File): Promise<BulkImportResult> => {
  const formData = new FormData();
  formData.append("file", file);
  return apiPostMultipart<BulkImportResult>("/products/bulk", formData);
};
