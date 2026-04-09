import { VITE_ADMIN_PATH } from "@/config/constants";
import { normalizeAdminPath } from "@/utils/admin-path";

export function getAdminRouterBasename() {
  return normalizeAdminPath(VITE_ADMIN_PATH);
}
