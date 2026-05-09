import { useState } from "react";
import PageShell from "../components/layout/PageShell";
import CategoryCountsTab from "../components/report/CategoryCountsTab";
import PriceDistributionTab from "../components/report/PriceDistributionTab";
import LowStockTab from "../components/report/LowStockTab";
import "./ReportsPage.css";

type Tab = "counts" | "prices" | "lowstock";

const TABS: { id: Tab; label: string }[] = [
  { id: "counts", label: "Category Counts" },
  { id: "prices", label: "Price Distribution" },
  { id: "lowstock", label: "Low Stock" },
];

function ReportsPage() {
  const [activeTab, setActiveTab] = useState<Tab>("counts");

  return (
    <PageShell>
      <div className="reports__header">
        <h1 className="reports__title">Reports</h1>
        <p className="reports__subtitle">
          On-demand analytics computed from current inventory data.
        </p>
      </div>

      <div className="reports__tabs" role="tablist">
        {TABS.map((tab) => (
          <button
            key={tab.id}
            role="tab"
            aria-selected={activeTab === tab.id}
            className={`reports__tab${activeTab === tab.id ? " reports__tab--active" : ""}`}
            onClick={() => setActiveTab(tab.id)}
          >
            {tab.label}
          </button>
        ))}
      </div>

      <div className="reports__content" role="tabpanel">
        {activeTab === "counts" && <CategoryCountsTab />}
        {activeTab === "prices" && <PriceDistributionTab />}
        {activeTab === "lowstock" && <LowStockTab />}
      </div>
    </PageShell>
  );
}

export default ReportsPage;
