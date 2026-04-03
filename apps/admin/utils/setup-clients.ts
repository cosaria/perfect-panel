import { isBrowser } from "@workspace/ui/utils";
import { toast } from "sonner";
import { client as adminClient } from "@/services/admin-api/client.gen";
import { client as commonClient } from "@/services/common-api/client.gen";
import { client as userClient } from "@/services/user-api/client.gen";
import { NEXT_PUBLIC_API_URL, NEXT_PUBLIC_SITE_URL } from "@/config/constants";
import { getAuthorization, Logout } from "./common";

const baseUrl = NEXT_PUBLIC_API_URL || NEXT_PUBLIC_SITE_URL || "";

function setupClient(client: typeof adminClient, serverPrefix: string) {
  client.setConfig({ baseUrl: `${baseUrl}${serverPrefix}` });

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

setupClient(adminClient, "/v1/admin");
setupClient(commonClient, "/v1/common");
// user-api 用绝对路径（/v1/auth/*），baseUrl 为 API 根
setupClient(userClient, "");
