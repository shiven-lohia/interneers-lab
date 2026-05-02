import "./ViewToggle.css";

interface ViewToggleProps {
  view: "grid" | "list";
  onChange: (v: "grid" | "list") => void;
}

function GridIcon() {
  return (
    <svg
      width="14"
      height="14"
      viewBox="0 0 14 14"
      fill="none"
      aria-hidden="true"
    >
      <rect x="0" y="0" width="6" height="6" rx="1" fill="currentColor" />
      <rect x="8" y="0" width="6" height="6" rx="1" fill="currentColor" />
      <rect x="0" y="8" width="6" height="6" rx="1" fill="currentColor" />
      <rect x="8" y="8" width="6" height="6" rx="1" fill="currentColor" />
    </svg>
  );
}

function ListIcon() {
  return (
    <svg
      width="14"
      height="14"
      viewBox="0 0 14 14"
      fill="none"
      aria-hidden="true"
    >
      <rect x="0" y="1" width="14" height="2" rx="1" fill="currentColor" />
      <rect x="0" y="6" width="14" height="2" rx="1" fill="currentColor" />
      <rect x="0" y="11" width="14" height="2" rx="1" fill="currentColor" />
    </svg>
  );
}

function ViewToggle({ view, onChange }: ViewToggleProps) {
  return (
    <div className="view-toggle" role="group" aria-label="View mode">
      <button
        className={`view-toggle__btn${view === "grid" ? " view-toggle__btn--active" : ""}`}
        onClick={() => onChange("grid")}
        aria-label="Grid view"
        aria-pressed={view === "grid"}
      >
        <GridIcon />
      </button>
      <button
        className={`view-toggle__btn${view === "list" ? " view-toggle__btn--active" : ""}`}
        onClick={() => onChange("list")}
        aria-label="List view"
        aria-pressed={view === "list"}
      >
        <ListIcon />
      </button>
    </div>
  );
}

export default ViewToggle;
