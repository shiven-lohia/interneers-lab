import { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from "recharts";
import { getPriceDistribution } from "../../api/reports";
import type { PriceDistributionReport } from "../../types";
import Button from "../ui/Button";
import LoadingSpinner from "../ui/LoadingSpinner";
import ErrorMessage from "../ui/ErrorMessage";
import { buildCSV, downloadCSV } from "../../utils/csv";
import "./PriceDistributionTab.css";

const DEFAULT_BUCKETS = "50,100,500,1000";
const STORAGE_KEY = "reports_buckets";

const BUCKET_COLORS = [
  "oklch(72% 0.09 155)",
  "oklch(65% 0.11 155)",
  "oklch(57% 0.12 155)",
  "oklch(50% 0.13 155)",
  "oklch(43% 0.14 155)",
  "oklch(36% 0.12 155)",
];

function PriceDistributionTab() {
  const [report, setReport] = useState<PriceDistributionReport | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [bucketsInput, setBucketsInput] = useState(
    () => sessionStorage.getItem(STORAGE_KEY) ?? DEFAULT_BUCKETS,
  );
  const [filterError, setFilterError] = useState("");

  const fetchReport = (buckets?: string) => {
    setLoading(true);
    setError("");
    getPriceDistribution(buckets || undefined)
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
    fetchReport(sessionStorage.getItem(STORAGE_KEY) ?? DEFAULT_BUCKETS);
  }, []);

  const handleApply = () => {
    setFilterError("");
    const raw = bucketsInput.trim();
    if (raw === "") {
      fetchReport();
      return;
    }
    const parts = raw.split(",").map((s) => s.trim());
    const nums = parts.map(Number);
    if (nums.some(isNaN)) {
      setFilterError(
        "Bucket edges must be comma-separated numbers (e.g. 50,100,500,1000).",
      );
      return;
    }
    if (nums.some((n) => n <= 0)) {
      setFilterError("All bucket edges must be positive.");
      return;
    }
    for (let i = 1; i < nums.length; i++) {
      if (nums[i] <= nums[i - 1]) {
        setFilterError("Bucket edges must be strictly increasing.");
        return;
      }
    }
    fetchReport(raw);
  };

  const chartData = report
    ? report.categories.map((cat) => {
        const row: Record<string, string | number> = {
          category_title: cat.category_title,
        };
        report.buckets.forEach((bucket, i) => {
          row[bucket.label] = cat.counts[i];
        });
        return row;
      })
    : [];

  const handleDownload = () => {
    if (!report) return;
    const header = ["Category", ...report.buckets.map((b) => b.label)];
    const rows: (string | number)[][] = [
      header,
      ...report.categories.map((cat) => [cat.category_title, ...cat.counts]),
    ];
    downloadCSV("price-distribution.csv", buildCSV(rows));
  };

  return (
    <div className="report-tab">
      <div className="report-tab__filters">
        <div className="report-tab__filter-group">
          <label className="report-tab__filter-label" htmlFor="bucket-edges">
            Bucket edges
          </label>
          <input
            id="bucket-edges"
            className="report-tab__filter-input report-tab__filter-input--wide"
            type="text"
            placeholder={DEFAULT_BUCKETS}
            value={bucketsInput}
            onChange={(e) => {
              setBucketsInput(e.target.value);
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
      <p className="report-tab__filter-hint">
        Comma-separated price thresholds (e.g.{" "}
        <code className="report-tab__code">50,100,500,1000</code>). Products are
        grouped into: 0–first, first–second, …, last+.
      </p>

      {filterError && <p className="report-tab__filter-error">{filterError}</p>}

      {loading && <LoadingSpinner />}
      {error && <ErrorMessage message={error} />}

      {!loading && !error && report && (
        <>
          {report.categories.length === 0 ? (
            <p className="report-tab__empty">No categories found.</p>
          ) : (
            <>
              <div className="report-tab__chart" style={{ height: 320 }}>
                <ResponsiveContainer width="100%" height="100%">
                  <BarChart
                    data={chartData}
                    margin={{ top: 4, right: 20, left: 0, bottom: 40 }}
                  >
                    <XAxis
                      dataKey="category_title"
                      tick={{ fontSize: 12 }}
                      angle={-35}
                      textAnchor="end"
                      interval={0}
                    />
                    <YAxis allowDecimals={false} tick={{ fontSize: 12 }} />
                    <Tooltip contentStyle={{ fontSize: 13 }} />
                    <Legend wrapperStyle={{ fontSize: 12, paddingTop: 8 }} />
                    {report.buckets.map((bucket, i) => (
                      <Bar
                        key={bucket.label}
                        dataKey={bucket.label}
                        fill={BUCKET_COLORS[i % BUCKET_COLORS.length]}
                        maxBarSize={24}
                        radius={[3, 3, 0, 0]}
                      />
                    ))}
                  </BarChart>
                </ResponsiveContainer>
              </div>

              <table className="report-tab__table">
                <thead>
                  <tr>
                    <th>Category</th>
                    {report.buckets.map((b) => (
                      <th key={b.label} className="report-tab__table-num">
                        {b.label}
                      </th>
                    ))}
                  </tr>
                </thead>
                <tbody>
                  {report.categories.map((cat) => (
                    <tr key={cat.category_id}>
                      <td>
                        <Link
                          to={`/categories/${cat.category_id}`}
                          className="report-tab__table-link"
                        >
                          {cat.category_title}
                        </Link>
                      </td>
                      {cat.counts.map((count, i) => (
                        <td key={i} className="report-tab__table-num">
                          {count}
                        </td>
                      ))}
                    </tr>
                  ))}
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

export default PriceDistributionTab;
