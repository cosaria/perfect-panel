// @ts-expect-error
/* eslint-disable */
import request from "@/utils/request";

/** Get System Log GET /v1/admin/tool/log */
export async function getSystemLog(options?: Record<string, unknown>) {
  return request<API.Response & { data?: API.LogResponse }>("/v1/admin/tool/log", {
    method: "GET",
    ...(options || {}),
  });
}

/** Restart System GET /v1/admin/tool/restart */
export async function restartSystem(options?: Record<string, unknown>) {
  return request<API.Response & { data?: unknown }>("/v1/admin/tool/restart", {
    method: "GET",
    ...(options || {}),
  });
}

/** Get Version GET /v1/admin/tool/version */
export async function getVersion(options?: Record<string, unknown>) {
  return request<API.Response & { data?: API.VersionResponse }>("/v1/admin/tool/version", {
    method: "GET",
    ...(options || {}),
  });
}
