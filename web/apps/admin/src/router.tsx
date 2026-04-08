import { createBrowserRouter } from "react-router-dom";
import { getAdminRouterBasename } from "./env";
import { routes } from "./routes";

export const router = createBrowserRouter(routes, {
  basename: getAdminRouterBasename(),
});
