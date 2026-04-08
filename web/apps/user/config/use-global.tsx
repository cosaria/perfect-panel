import { extractDomain } from "@workspace/ui/utils";
import { create } from "zustand";
import { NEXT_PUBLIC_API_URL, NEXT_PUBLIC_SITE_URL } from "@/config/constants";
import type { GetGlobalConfigResponse } from "@/services/common-api/types.gen";
import { queryUserInfo } from "@/services/user-api/sdk.gen";
import type { User } from "@/services/user-api/types.gen";

export interface GlobalStore {
  common: GetGlobalConfigResponse;
  user?: User;
  setCommon: (common: Partial<GetGlobalConfigResponse>) => void;
  setUser: (user?: User) => void;
  getUserInfo: () => Promise<void>;
  getUserSubscribe: (uuid: string, type?: string) => string[];
  getAppSubLink: (url: string, schema?: string) => string;
}

function parseSchemaObjectExpression(expr: string, url: string, name: string) {
  const processedExpr = expr
    .replace(/\burl\b/g, JSON.stringify(url))
    .replace(/\bname\b/g, JSON.stringify(name));

  if (processedExpr.includes("server_remote")) {
    const serverRemoteValue = `${url}, tag=${name}`;
    return { server_remote: [serverRemoteValue] };
  }

  const normalizedExpr = processedExpr
    .replace(/([{,]\s*)([A-Za-z_][A-Za-z0-9_]*)\s*:/g, '$1"$2":')
    .replace(/'([^'\\]*(?:\\.[^'\\]*)*)'/g, (_match, value: string) => JSON.stringify(value));

  return JSON.parse(normalizedExpr) as Record<string, unknown>;
}

export const useGlobalStore = create<GlobalStore>((set, get) => ({
  common: {
    site: {
      host: "",
      site_name: "",
      site_desc: "",
      site_logo: "",
      keywords: "",
      custom_html: "",
      custom_data: "",
    },
    verify: {
      turnstile_site_key: "",
      enable_login_verify: false,
      enable_register_verify: false,
      enable_reset_password_verify: false,
    },
    auth: {
      mobile: {
        enable: false,
        enable_whitelist: false,
        whitelist: [],
      },
      email: {
        enable: false,
        enable_verify: false,
        enable_domain_suffix: false,
        domain_suffix_list: "",
      },
      register: {
        stop_register: false,
        enable_ip_register_limit: false,
        ip_register_limit: 0,
        ip_register_limit_duration: 0,
      },
      device: {
        enable: false,
        show_ads: false,
        enable_security: false,
        only_real_device: false,
      },
    },
    invite: {
      forced_invite: false,
      referral_percentage: 0,
      only_first_purchase: false,
    },
    currency: {
      currency_unit: "USD",
      currency_symbol: "$",
    },
    subscribe: {
      single_model: false,
      subscribe_path: "",
      subscribe_domain: "",
      pan_domain: false,
      user_agent_limit: false,
      user_agent_list: "",
    },
    verify_code: {
      verify_code_interval: 60,
    },
    oauth_methods: [],
    web_ad: false,
  },
  user: undefined,
  setCommon: (common) =>
    set((state) => ({
      common: {
        ...state.common,
        ...common,
      },
    })),
  setUser: (user) => set({ user }),
  getUserInfo: async () => {
    try {
      const { data } = await queryUserInfo();
      set({ user: data });
    } catch (error) {
      console.error("Failed to refresh user:", error);
    }
  },
  getUserSubscribe: (uuid: string, type?: string) => {
    const { pan_domain, subscribe_domain, subscribe_path } = get().common.subscribe || {};
    const domains = subscribe_domain
      ? subscribe_domain.split("\n")
      : [extractDomain(NEXT_PUBLIC_API_URL || NEXT_PUBLIC_SITE_URL || "", pan_domain)];

    return domains.map((domain) => {
      if (pan_domain) {
        if (type) return `https://${uuid}.${type}.${domain}`;
        return `https://${uuid}.${domain}`;
      } else {
        if (type) return `https://${domain}${subscribe_path}?token=${uuid}&type=${type}`;
        return `https://${domain}${subscribe_path}?token=${uuid}`;
      }
    });
  },
  getAppSubLink: (url: string, schema?: string) => {
    const name = get().common?.site?.site_name || "";

    if (!schema) return "url";
    try {
      let result = schema.replace(/\${url}/g, url).replace(/\${name}/g, name);

      const maxLoop = 10;
      let prev = "";
      let loop = 0;
      do {
        prev = result;
        result = result.replace(
          /\${encodeURIComponent\(JSON\.stringify\(([^)]+)\)\)}/g,
          (match, expr) => {
            try {
              const obj = parseSchemaObjectExpression(expr, url, name);
              return encodeURIComponent(JSON.stringify(obj));
            } catch {
              return match;
            }
          },
        );

        result = result.replace(/\${encodeURIComponent\(([^)]+)\)}/g, (match, expr) => {
          if (expr === "url") return encodeURIComponent(url);
          if (expr === "name") return encodeURIComponent(name);
          try {
            return encodeURIComponent(expr);
          } catch {
            return match;
          }
        });

        result = result.replace(/\${window\.btoa\(([^)]+)\)}/g, (match, expr) => {
          const btoa = typeof window !== "undefined" ? window.btoa : (str: string) => str;
          if (expr === "url") return btoa(url);
          if (expr === "name") return btoa(name);
          try {
            return btoa(expr);
          } catch {
            return match;
          }
        });

        result = result.replace(/\${JSON\.stringify\(([^}]+)\)}/g, (match, expr) => {
          try {
            return JSON.stringify(parseSchemaObjectExpression(expr, url, name));
          } catch {
            return match;
          }
        });
        loop++;
      } while (result !== prev && loop < maxLoop);
      return result;
    } catch (_error) {
      return "";
    }
  },
}));

export default useGlobalStore;
