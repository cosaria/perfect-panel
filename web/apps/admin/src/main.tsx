import "@workspace/ui/globals.css";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { AppRouter } from "./router";

const rootElement = document.getElementById("root");

if (!rootElement) {
  throw new Error("Missing #root element for admin app");
}

createRoot(rootElement).render(
  <StrictMode>
    <AppRouter />
  </StrictMode>,
);
