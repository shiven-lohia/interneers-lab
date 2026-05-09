import { apiGet } from "./client";
import type {
  CategoryCountsReport,
  PriceDistributionReport,
  LowStockReport,
} from "../types";

export const getCategoryCounts = (
  minCount?: number,
  maxCount?: number,
): Promise<CategoryCountsReport> => {
  const qs = new URLSearchParams();
  if (minCount !== undefined) qs.set("min_count", String(minCount));
  if (maxCount !== undefined) qs.set("max_count", String(maxCount));
  const suffix = qs.toString() ? `?${qs}` : "";
  return apiGet<CategoryCountsReport>(`/reports/category-counts${suffix}`);
};

export const getPriceDistribution = (
  buckets?: string,
): Promise<PriceDistributionReport> => {
  const suffix = buckets ? `?buckets=${encodeURIComponent(buckets)}` : "";
  return apiGet<PriceDistributionReport>(
    `/reports/price-distribution${suffix}`,
  );
};

export const getLowStock = (threshold?: number): Promise<LowStockReport> => {
  const suffix = threshold !== undefined ? `?threshold=${threshold}` : "";
  return apiGet<LowStockReport>(`/reports/low-stock${suffix}`);
};
