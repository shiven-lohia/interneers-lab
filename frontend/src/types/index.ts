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

// Reports

export interface CategoryCount {
  category_id: string;
  category_title: string;
  product_count: number;
}

export interface CountFilters {
  min_count: number | null;
  max_count: number | null;
}

export interface CategoryCountsReport {
  categories: CategoryCount[];
  uncategorized_count: number;
  filters: CountFilters;
}

export interface PriceBucket {
  label: string;
  min: number;
  max: number | null;
}

export interface CategoryPriceCounts {
  category_id: string;
  category_title: string;
  counts: number[];
}

export interface PriceDistributionReport {
  buckets: PriceBucket[];
  categories: CategoryPriceCounts[];
}

export interface LowStockProduct {
  id: string;
  name: string;
  quantity: number;
  category_id: string;
  category_title: string;
}

export interface LowStockCategory {
  category_id: string;
  category_title: string;
  low_stock_count: number;
  total_count: number;
  percentage: number;
}

export interface LowStockReport {
  threshold: number;
  products: LowStockProduct[];
  categories: LowStockCategory[];
}
