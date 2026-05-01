import { apiGet, apiPost } from "./client";
import type { Category } from "../types";

export const getCategories = () => apiGet<Category[]>("/categories");
export const getCategory = (id: string) =>
  apiGet<Category>(`/categories/${id}`);
export const createCategory = (data: Pick<Category, "title" | "description">) =>
  apiPost<Category>("/categories", data);
