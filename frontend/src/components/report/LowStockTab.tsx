import { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  Tooltip,
  ResponsiveContainer,
  Cell,
} from "recharts";
import { getLowStock } from "../../api/reports";
import type { LowStockReport } from "../../types";
import Button from "../ui/Button";
import LoadingSpinner from "../ui/LoadingSpinner";
import ErrorMessage from "../ui/ErrorMessage";
import { buildCSV, downloadCSV } from "../../utils/csv";
import "./LowStockTab.css";

const SAGE = "oklch(72% 0.09 155)";
const CLAY = "oklch(60% 0.16 25)";
const DEFAULT_THRESHOLD = 10;
const STORAGE_KEY = "reports_threshold";

function LowStockTab() {
  const [report, setReport] = useState<LowStockReport | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [thresholdInput, setThresholdInput] = useState(
    () => sessionStorage.getItem(STORAGE_KEY) ?? String(DEFAULT_THRESHOLD),
  );
  const [filterError, setFilterError] = useState("");

  const fetchReport = (threshold?: number) => {
    setLoading(true);
    setError("");
    getLowStock(threshold)
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
    const stored = sessionStorage.getItem(STORAGE_KEY);
    fetchReport(stored ? Number(stored) : DEFAULT_THRESHOLD);
  }, []);

  const handleApply = () => {
    setFilterError("");
    const v = Number(thresholdInput);
    if (!Number.isInteger(v) || v < 0) {
      setFilterError("Threshold must be a non-negative integer.");
      return;
    }
    fetchReport(v);
  };

  const handleDownloadProducts = () => {
    if (!report) return;
    const rows: (string | number)[][] = [
      ["Name", "Quantity", "Category"],
      ...report.products.map((p) => [
        p.name,
        p.quantity,
        p.category_title || "Uncategorized",
      ]),
    ];
    downloadCSV("low-stock-products.csv", buildCSV(rows));
  };

  const handleDownloadCategories = () => {
    if (!report) return;
    const rows: (string | number)[][] = [
      ["Category", "Low-Stock Count", "Total Count", "Percentage (%)"],
      ...report.categories.map((c) => [
        c.category_title,
        c.low_stock_count,
        c.total_count,
        c.percentage,
      ]),
    ];
    downloadCSV("low-stock-categories.csv", buildCSV(rows));
  };

  return (
    <div className="report-tab">
      <div className="report-tab__filters">
        <div className="report-tab__filter-group">
          <label className="report-tab__filter-label" htmlFor="threshold">
            Low-stock threshold (qty &lt;)
          </label>
          <input
            id="threshold"
            className="report-tab__filter-input"
            type="number"
            min={0}
            value={thresholdInput}
            onChange={(e) => {
              setThresholdInput(e.target.value);
              sessionStorage.setItem(STORAGE_KEY, e.target.value);
            }}
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
          {/* Low-stock products */}
          <section className="low-stock__section">
            <div className="low-stock__section-header">
              <h2 className="low-stock__section-title">
                Low-stock products
                <span className="low-stock__count">
                  {report.products.length}
                </span>
              </h2>
              <Button
                variant="secondary"
                onClick={handleDownloadProducts}
                disabled={report.products.length === 0}
              >
                Download CSV
              </Button>
            </div>

            {report.products.length === 0 ? (
              <p className="report-tab__empty">
                No products are below the threshold of {report.threshold}.
              </p>
            ) : (
              <table className="report-tab__table">
                <thead>
                  <tr>
                    <th>Product</th>
                    <th>Category</th>
                    <th className="report-tab__table-num">Quantity</th>
                  </tr>
                </thead>
                <tbody>
                  {report.products.map((p) => (
                    <tr key={p.id}>
                      <td>
                        <Link
                          to={`/products/${p.id}`}
                          className="report-tab__table-link"
                        >
                          {p.name}
                        </Link>
                      </td>
                      <td
                        className={p.category_title ? "" : "low-stock__muted"}
                      >
                        {p.category_id ? (
                          <Link
                            to={`/categories/${p.category_id}`}
                            className="report-tab__table-link"
                          >
                            {p.category_title}
                          </Link>
                        ) : (
                          "Uncategorized"
                        )}
                      </td>
                      <td className="report-tab__table-num low-stock__qty">
                        {p.quantity}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            )}
          </section>

          {/* Low-stock categories */}
          <section className="low-stock__section">
            <div className="low-stock__section-header">
              <h2 className="low-stock__section-title">
                Categories with &gt;10% low-stock
                <span className="low-stock__count">
                  {report.categories.length}
                </span>
              </h2>
              <Button
                variant="secondary"
                onClick={handleDownloadCategories}
                disabled={report.categories.length === 0}
              >
                Download CSV
              </Button>
            </div>

            {report.categories.length === 0 ? (
              <p className="report-tab__empty">
                No categories exceed 10% low-stock at this threshold.
              </p>
            ) : (
              <>
                <div
                  className="report-tab__chart"
                  style={{
                    height: Math.max(200, report.categories.length * 44),
                  }}
                >
                  <ResponsiveContainer width="100%" height="100%">
                    <BarChart
                      layout="vertical"
                      data={report.categories}
                      margin={{ top: 4, right: 60, left: 0, bottom: 4 }}
                    >
                      <XAxis
                        type="number"
                        domain={[0, 100]}
                        tickFormatter={(v) => `${v}%`}
                        tick={{ fontSize: 12 }}
                      />
                      <YAxis
                        type="category"
                        dataKey="category_title"
                        width={130}
                        tick={{ fontSize: 13 }}
                      />
                      <Tooltip
                        formatter={(value) => [`${value}%`, "Low-stock %"]}
                        contentStyle={{ fontSize: 13 }}
                      />
                      <Bar
                        dataKey="percentage"
                        maxBarSize={28}
                        radius={[0, 4, 4, 0]}
                      >
                        {report.categories.map((cat) => (
                          <Cell
                            key={cat.category_id}
                            fill={cat.percentage >= 30 ? CLAY : SAGE}
                          />
                        ))}
                      </Bar>
                    </BarChart>
                  </ResponsiveContainer>
                </div>

                <table className="report-tab__table">
                  <thead>
                    <tr>
                      <th>Category</th>
                      <th className="report-tab__table-num">Low-stock</th>
                      <th className="report-tab__table-num">Total</th>
                      <th className="report-tab__table-num">%</th>
                    </tr>
                  </thead>
                  <tbody>
                    {report.categories.map((cat) => (
                      <tr
                        key={cat.category_id}
                        className={
                          cat.percentage >= 30 ? "low-stock__row--danger" : ""
                        }
                      >
                        <td>
                          <Link
                            to={`/categories/${cat.category_id}`}
                            className="report-tab__table-link"
                          >
                            {cat.category_title}
                          </Link>
                        </td>
                        <td className="report-tab__table-num">
                          {cat.low_stock_count}
                        </td>
                        <td className="report-tab__table-num">
                          {cat.total_count}
                        </td>
                        <td className="report-tab__table-num">
                          {cat.percentage}%
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </>
            )}
          </section>
        </>
      )}
    </div>
  );
}

export default LowStockTab;
