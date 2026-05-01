import { apiGet, apiPut } from "./client";
import type { Product, ProductFormData } from "../types";

export const getProducts = () => apiGet<Product[]>("/products");
export const getProduct = (id: string) => apiGet<Product>(`/products/${id}`);
export const getProductsByCategory = (categoryId: string) =>
  apiGet<Product[]>(`/products?category_id=${categoryId}`);
export const updateProduct = (id: string, data: ProductFormData) =>
  apiPut<Product>(`/products/${id}`, data);
