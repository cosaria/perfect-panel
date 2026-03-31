// @ts-expect-error
/* eslint-disable */
import request from "@/utils/request";

/** Apple Login Callback POST /v1/auth/oauth/callback/apple */
export async function appleLoginCallback(
  body: {
    code: string;
    id_token: string;
    state: string;
  },
  options?: Record<string, unknown>,
) {
  const formData = new FormData();

  for (const [key, item] of Object.entries(body)) {
    if (item === undefined || item === null) {
      continue;
    }

    if (typeof item === "object" && !(item instanceof File) && !(item instanceof Blob)) {
      if (Array.isArray(item)) {
        for (const value of item) {
          formData.append(key, value == null ? "" : String(value));
        }
      } else {
        formData.append(key, JSON.stringify(item));
      }
      continue;
    }

    formData.append(key, item);
  }

  return request<API.Response & { data?: unknown }>("/v1/auth/oauth/callback/apple", {
    method: "POST",
    data: formData,
    ...(options || {}),
  });
}

/** OAuth login POST /v1/auth/oauth/login */
export async function oAuthLogin(body: API.OAthLoginRequest, options?: Record<string, unknown>) {
  return request<API.Response & { data?: API.OAuthLoginResponse }>("/v1/auth/oauth/login", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    data: body,
    ...(options || {}),
  });
}

/** OAuth login get token POST /v1/auth/oauth/login/token */
export async function oAuthLoginGetToken(
  body: API.OAuthLoginGetTokenRequest,
  options?: Record<string, unknown>,
) {
  return request<API.Response & { data?: API.LoginResponse }>("/v1/auth/oauth/login/token", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    data: body,
    ...(options || {}),
  });
}
