import { NEXT_PUBLIC_ADMIN_PATH } from "@/config/constants";
import { normalizeAdminPath } from "@/utils/admin-path";

export function getAdminRouterBasename() {
  return normalizeAdminPath(NEXT_PUBLIC_ADMIN_PATH);
}
