// @ts-ignore
/* eslint-disable */
import request from '@/utils/request';

/** Get System Log GET /v1/admin/tool/log */
export async function getSystemLog(options?: { [key: string]: any }) {
  return request<API.Response & { data?: API.LogResponse }>('/v1/admin/tool/log', {
    method: 'GET',
    ...(options || {}),
  });
}

/** Restart System GET /v1/admin/tool/restart */
export async function restartSystem(options?: { [key: string]: any }) {
  return request<API.Response & { data?: any }>('/v1/admin/tool/restart', {
    method: 'GET',
    ...(options || {}),
  });
}

/** Get Version GET /v1/admin/tool/version */
export async function getVersion(options?: { [key: string]: any }) {
  return request<API.Response & { data?: API.VersionResponse }>('/v1/admin/tool/version', {
    method: 'GET',
    ...(options || {}),
  });
}
