import React from "react";
import "./PageShell.css";

interface PageShellProps {
  children: React.ReactNode;
}

function PageShell({ children }: PageShellProps) {
  return <main className="page-shell">{children}</main>;
}

export default PageShell;
