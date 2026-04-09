import { buildApiBaseUrl, isBrowser } from "@workspace/ui/utils";
import { toast } from "sonner";
import { VITE_API_URL, VITE_SITE_URL } from "@/config/constants";
import { client as adminClient } from "@/services/admin-api/client.gen";
import { client as commonClient } from "@/services/common-api/client.gen";
import { client as userClient } from "@/services/user-api/client.gen";
import { getAuthorization, Logout } from "./common";

function setupClient(client: typeof adminClient, serverPrefix: string) {
  client.setConfig({
    baseUrl: buildApiBaseUrl(VITE_API_URL, VITE_SITE_URL, serverPrefix),
  });

  client.interceptors.request.use((request) => {
    const token = getAuthorization();
    if (token) {
      request.headers.set("Authorization", token);
    }
    return request;
  });

  client.interceptors.response.use(async (response) => {
    if (response.ok) return response;

    if (response.status === 401) {
      Logout();
      return response;
    }

    if (isBrowser()) {
      try {
        const body = await response.clone().json();
        toast.error(body.detail || body.title || "Unknown error");
      } catch {
        toast.error(`Error: ${response.status}`);
      }
    }

    return response;
  });
}

setupClient(adminClient, "/api/v1/admin");
setupClient(commonClient, "/api/v1/common");
// user-api 已携带绝对路径（/api/v1/auth/*），baseUrl 保持站点根
setupClient(userClient, "");
