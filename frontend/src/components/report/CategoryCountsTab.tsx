import { useState, useEffect } from "react";
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  Tooltip,
  ResponsiveContainer,
} from "recharts";
import { getCategoryCounts } from "../../api/reports";
import type { CategoryCountsReport } from "../../types";
import Button from "../ui/Button";
import LoadingSpinner from "../ui/LoadingSpinner";
import ErrorMessage from "../ui/ErrorMessage";
import { buildCSV, downloadCSV } from "../../utils/csv";
import "./CategoryCountsTab.css";

const SAGE = "oklch(72% 0.09 155)";

function CategoryCountsTab() {
  const [report, setReport] = useState<CategoryCountsReport | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [minCount, setMinCount] = useState("");
  const [maxCount, setMaxCount] = useState("");
  const [filterError, setFilterError] = useState("");

  const fetch = (min?: number, max?: number) => {
    setLoading(true);
    setError("");
    getCategoryCounts(min, max)
      .then((data) => {
        setReport(data);
        setLoading(false);
      })
      .catch((err: Error) => {
        setError(err.message);
        setLoading(false);
      });
  };

  useEffect(() => {
    fetch();
  }, []);

  const handleApply = () => {
    setFilterError("");
    const min = minCount !== "" ? Number(minCount) : undefined;
    const max = maxCount !== "" ? Number(maxCount) : undefined;
    if (
      (min !== undefined && !Number.isInteger(min)) ||
      (min !== undefined && min < 0)
    ) {
      setFilterError("Min count must be a non-negative integer.");
      return;
    }
    if (
      (max !== undefined && !Number.isInteger(max)) ||
      (max !== undefined && max < 0)
    ) {
      setFilterError("Max count must be a non-negative integer.");
      return;
    }
    if (min !== undefined && max !== undefined && min > max) {
      setFilterError("Min count must not exceed max count.");
      return;
    }
    fetch(min, max);
  };

  const handleDownload = () => {
    if (!report) return;
    const rows: (string | number)[][] = [
      ["Category", "Product Count"],
      ...report.categories.map((c) => [c.category_title, c.product_count]),
    ];
    if (report.uncategorized_count > 0) {
      rows.push(["Uncategorized", report.uncategorized_count]);
    }
    downloadCSV("category-counts.csv", buildCSV(rows));
  };

  return (
    <div className="report-tab">
      <div className="report-tab__filters">
        <div className="report-tab__filter-group">
          <label className="report-tab__filter-label" htmlFor="min-count">
            Min count
          </label>
          <input
            id="min-count"
            className="report-tab__filter-input"
            type="number"
            min={0}
            placeholder="No minimum"
            value={minCount}
            onChange={(e) => setMinCount(e.target.value)}
          />
        </div>
        <div className="report-tab__filter-group">
          <label className="report-tab__filter-label" htmlFor="max-count">
            Max count
          </label>
          <input
            id="max-count"
            className="report-tab__filter-input"
            type="number"
            min={0}
            placeholder="No maximum"
            value={maxCount}
            onChange={(e) => setMaxCount(e.target.value)}
          />
        </div>
        <Button
          variant="primary"
          className="report-tab__apply-btn"
          onClick={handleApply}
          disabled={loading}
        >
          Apply
        </Button>
      </div>

      {filterError && <p className="report-tab__filter-error">{filterError}</p>}

      {loading && <LoadingSpinner />}
      {error && <ErrorMessage message={error} />}

      {!loading && !error && report && (
        <>
          {report.categories.length === 0 ? (
            <p className="report-tab__empty">
              No categories match the current filters.
            </p>
          ) : (
            <>
              <div
                className="report-tab__chart"
                style={{
                  height: Math.max(220, report.categories.length * 44),
                }}
              >
                <ResponsiveContainer width="100%" height="100%">
                  <BarChart
                    layout="vertical"
                    data={report.categories}
                    margin={{ top: 4, right: 40, left: 0, bottom: 4 }}
                  >
                    <XAxis
                      type="number"
                      allowDecimals={false}
                      tick={{ fontSize: 12 }}
                    />
                    <YAxis
                      type="category"
                      dataKey="category_title"
                      width={130}
                      tick={{ fontSize: 13 }}
                    />
                    <Tooltip
                      formatter={(value) => [value, "Products"]}
                      contentStyle={{ fontSize: 13 }}
                    />
                    <Bar
                      dataKey="product_count"
                      fill={SAGE}
                      radius={[0, 4, 4, 0]}
                      maxBarSize={28}
                    />
                  </BarChart>
                </ResponsiveContainer>
              </div>

              <table className="report-tab__table">
                <thead>
                  <tr>
                    <th>Category</th>
                    <th className="report-tab__table-num">Product Count</th>
                  </tr>
                </thead>
                <tbody>
                  {report.categories.map((c) => (
                    <tr key={c.category_id}>
                      <td>{c.category_title}</td>
                      <td className="report-tab__table-num">
                        {c.product_count}
                      </td>
                    </tr>
                  ))}
                  {report.uncategorized_count > 0 && (
                    <tr className="report-tab__table-row--muted">
                      <td>Uncategorized</td>
                      <td className="report-tab__table-num">
                        {report.uncategorized_count}
                      </td>
                    </tr>
                  )}
                </tbody>
              </table>
            </>
          )}

          <div className="report-tab__actions">
            <Button
              variant="secondary"
              onClick={handleDownload}
              disabled={report.categories.length === 0}
            >
              Download CSV
            </Button>
          </div>
        </>
      )}
    </div>
  );
}

export default CategoryCountsTab;
