import React from "react";
import { render, screen } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import Navbar from "../components/layout/Navbar";

test("renders navigation links", () => {
  render(
    <MemoryRouter>
      <Navbar />
    </MemoryRouter>,
  );
  expect(screen.getByText("My Store")).toBeInTheDocument();
  expect(screen.getByText("Products")).toBeInTheDocument();
  expect(screen.getByText("Categories")).toBeInTheDocument();
});
