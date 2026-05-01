import { useState, useRef } from "react";
import { useNavigate } from "react-router-dom";
import PageShell from "../components/layout/PageShell";
import { bulkCreateProducts } from "../api/products";
import type { BulkImportResult } from "../types";
import "./BulkImportPage.css";

function BulkImportPage() {
  const navigate = useNavigate();
  const [file, setFile] = useState<File | null>(null);
  const [dragging, setDragging] = useState(false);
  const [uploading, setUploading] = useState(false);
  const [result, setResult] = useState<BulkImportResult | null>(null);
  const [zoneError, setZoneError] = useState("");
  const inputRef = useRef<HTMLInputElement>(null);

  const acceptFile = (f: File) => {
    if (f.name.endsWith(".csv")) {
      setFile(f);
      setZoneError("");
    } else {
      setZoneError("Only .csv files are accepted.");
    }
  };

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    setDragging(true);
  };
  const handleDragLeave = () => setDragging(false);
  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    setDragging(false);
    const dropped = e.dataTransfer.files[0];
    if (dropped) acceptFile(dropped);
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const selected = e.target.files?.[0];
    if (selected) acceptFile(selected);
  };

  const handleUpload = async () => {
    if (!file) return;
    setUploading(true);
    setResult(null);
    setZoneError("");
    try {
      const data = await bulkCreateProducts(file);
      setResult(data);
      setFile(null);
      if (inputRef.current) inputRef.current.value = "";
    } catch (err) {
      setZoneError(err instanceof Error ? err.message : "Upload failed.");
    } finally {
      setUploading(false);
    }
  };

  const zoneClass = [
    "bulk-import__zone",
    dragging ? "bulk-import__zone--drag" : "",
    uploading ? "bulk-import__zone--uploading" : "",
  ]
    .filter(Boolean)
    .join(" ");

  return (
    <PageShell>
      <div className="bulk-import__header">
        <div>
          <h1 className="bulk-import__title">Bulk Import Products</h1>
          <p className="bulk-import__subtitle">
            Upload a CSV to add multiple products at once.
          </p>
        </div>
        <button
          className="bulk-import__back-btn"
          onClick={() => navigate("/products")}
        >
          Back to products
        </button>
      </div>

      <div className="bulk-import__format">
        <span className="bulk-import__format-label">Required format</span>
        <code className="bulk-import__format-code">
          name, description, category_id, price, quantity, brand
        </code>
      </div>

      <div
        className={zoneClass}
        onDragOver={handleDragOver}
        onDragLeave={handleDragLeave}
        onDrop={handleDrop}
        onClick={() => !uploading && inputRef.current?.click()}
        role="button"
        aria-label="Upload CSV file"
        tabIndex={0}
        onKeyDown={(e) => {
          if (e.key === "Enter" || e.key === " ") inputRef.current?.click();
        }}
      >
        <input
          ref={inputRef}
          type="file"
          accept=".csv"
          className="bulk-import__input"
          onChange={handleFileChange}
        />
        {file ? (
          <span className="bulk-import__zone-filename">{file.name}</span>
        ) : (
          <span className="bulk-import__zone-prompt">
            {uploading ? "Uploading…" : "Drop a CSV here, or choose a file"}
          </span>
        )}
      </div>

      {zoneError && <p className="bulk-import__zone-error">{zoneError}</p>}

      <button
        className="bulk-import__upload-btn"
        disabled={!file || uploading}
        onClick={handleUpload}
      >
        {uploading ? "Uploading…" : "Upload"}
      </button>

      {result && (
        <div className="bulk-import__results">
          {result.created.length > 0 && (
            <section className="bulk-import__result-section bulk-import__result-section--created">
              <h2 className="bulk-import__result-heading">
                Created{" "}
                <span className="bulk-import__result-count">
                  {result.created.length}
                </span>
              </h2>
              <ul className="bulk-import__result-list">
                {result.created.map((p) => (
                  <li key={p.id} className="bulk-import__result-item">
                    <button
                      className="bulk-import__result-link"
                      onClick={() => navigate(`/products/${p.id}`)}
                    >
                      {p.name}
                    </button>
                    <span className="bulk-import__result-meta">
                      {p.brand} &middot; {p.price}
                    </span>
                  </li>
                ))}
              </ul>
            </section>
          )}

          {result.errors.length > 0 && (
            <section className="bulk-import__result-section bulk-import__result-section--errors">
              <h2 className="bulk-import__result-heading">
                Not imported{" "}
                <span className="bulk-import__result-count">
                  {result.errors.length}
                </span>
              </h2>
              <ul className="bulk-import__result-list">
                {result.errors.map((e, i) => (
                  <li
                    key={i}
                    className="bulk-import__result-item bulk-import__result-item--error"
                  >
                    Row {e.index + 1}: {e.reason}
                  </li>
                ))}
              </ul>
            </section>
          )}

          {result.created.length === 0 && result.errors.length === 0 && (
            <p className="bulk-import__result-empty">
              No products were imported.
            </p>
          )}
        </div>
      )}
    </PageShell>
  );
}

export default BulkImportPage;
