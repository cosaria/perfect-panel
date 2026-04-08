import { BrowserRouter, useRoutes } from "react-router-dom";
import { getAdminRouterBasename } from "./env";
import { routes } from "./routes";

function AdminRouteRenderer() {
  return useRoutes(routes);
}

export function AppRouter() {
  return (
    <BrowserRouter basename={getAdminRouterBasename()} unstable_useTransitions={false}>
      <AdminRouteRenderer />
    </BrowserRouter>
  );
}
