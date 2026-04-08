import { createBrowserRouter } from "react-router-dom";
import { getUserRouterBasename } from "./env";
import { routes } from "./routes";

export const router = createBrowserRouter(routes, {
  basename: getUserRouterBasename(),
});
