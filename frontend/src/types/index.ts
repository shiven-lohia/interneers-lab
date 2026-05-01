export interface Category {
  id: string;
  title: string;
  description: string;
}

export interface Product {
  id: string;
  name: string;
  description: string;
  category_id: string;
  category?: Category;
  price: number;
  brand: string;
  quantity: number;
}

export type ProductFormData = Omit<Product, "category">;

export interface BulkImportError {
  index: number;
  reason: string;
}

export interface BulkImportResult {
  created: Product[];
  errors: BulkImportError[];
}
