// @ts-ignore
/* eslint-disable */
import request from '@/utils/request';

/** Query User Affiliate Count GET /v1/public/user/affiliate/count */
export async function queryUserAffiliate(options?: { [key: string]: any }) {
  return request<API.Response & { data?: API.QueryUserAffiliateCountResponse }>(
    '/v1/public/user/affiliate/count',
    {
      method: 'GET',
      ...(options || {}),
    },
  );
}

/** Query User Affiliate List GET /v1/public/user/affiliate/list */
export async function queryUserAffiliateList(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.QueryUserAffiliateListParams,
  options?: { [key: string]: any },
) {
  return request<API.Response & { data?: API.QueryUserAffiliateListResponse }>(
    '/v1/public/user/affiliate/list',
    {
      method: 'GET',
      params: {
        ...params,
      },
      ...(options || {}),
    },
  );
}

/** Query User Balance Log GET /v1/public/user/balance_log */
export async function queryUserBalanceLog(options?: { [key: string]: any }) {
  return request<API.Response & { data?: API.QueryUserBalanceLogListResponse }>(
    '/v1/public/user/balance_log',
    {
      method: 'GET',
      ...(options || {}),
    },
  );
}

/** Update Bind Email PUT /v1/public/user/bind_email */
export async function updateBindEmail(
  body: API.UpdateBindEmailRequest,
  options?: { [key: string]: any },
) {
  return request<API.Response & { data?: any }>('/v1/public/user/bind_email', {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** Update Bind Mobile PUT /v1/public/user/bind_mobile */
export async function updateBindMobile(
  body: API.UpdateBindMobileRequest,
  options?: { [key: string]: any },
) {
  return request<API.Response & { data?: any }>('/v1/public/user/bind_mobile', {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** Bind OAuth POST /v1/public/user/bind_oauth */
export async function bindOAuth(body: API.BindOAuthRequest, options?: { [key: string]: any }) {
  return request<API.Response & { data?: API.BindOAuthResponse }>('/v1/public/user/bind_oauth', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** Bind OAuth Callback POST /v1/public/user/bind_oauth/callback */
export async function bindOAuthCallback(
  body: API.BindOAuthCallbackRequest,
  options?: { [key: string]: any },
) {
  return request<API.Response & { data?: any }>('/v1/public/user/bind_oauth/callback', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** Bind Telegram GET /v1/public/user/bind_telegram */
export async function bindTelegram(options?: { [key: string]: any }) {
  return request<API.Response & { data?: API.BindTelegramResponse }>(
    '/v1/public/user/bind_telegram',
    {
      method: 'GET',
      ...(options || {}),
    },
  );
}

/** Query User Commission Log GET /v1/public/user/commission_log */
export async function queryUserCommissionLog(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.QueryUserCommissionLogParams,
  options?: { [key: string]: any },
) {
  return request<API.Response & { data?: API.QueryUserCommissionLogListResponse }>(
    '/v1/public/user/commission_log',
    {
      method: 'GET',
      params: {
        ...params,
      },
      ...(options || {}),
    },
  );
}

/** Get Device List GET /v1/public/user/devices */
export async function getDeviceList(options?: { [key: string]: any }) {
  return request<API.Response & { data?: API.GetDeviceListResponse }>('/v1/public/user/devices', {
    method: 'GET',
    ...(options || {}),
  });
}

/** Query User Info GET /v1/public/user/info */
export async function queryUserInfo(options?: { [key: string]: any }) {
  return request<API.Response & { data?: API.User }>('/v1/public/user/info', {
    method: 'GET',
    ...(options || {}),
  });
}

/** Get Login Log GET /v1/public/user/login_log */
export async function getLoginLog(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.GetLoginLogParams,
  options?: { [key: string]: any },
) {
  return request<API.Response & { data?: API.GetLoginLogResponse }>('/v1/public/user/login_log', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** Update User Notify PUT /v1/public/user/notify */
export async function updateUserNotify(
  body: API.UpdateUserNotifyRequest,
  options?: { [key: string]: any },
) {
  return request<API.Response & { data?: any }>('/v1/public/user/notify', {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** Get OAuth Methods GET /v1/public/user/oauth_methods */
export async function getOAuthMethods(options?: { [key: string]: any }) {
  return request<API.Response & { data?: API.GetOAuthMethodsResponse }>(
    '/v1/public/user/oauth_methods',
    {
      method: 'GET',
      ...(options || {}),
    },
  );
}

/** Update User Password PUT /v1/public/user/password */
export async function updateUserPassword(
  body: API.UpdateUserPasswordRequest,
  options?: { [key: string]: any },
) {
  return request<API.Response & { data?: any }>('/v1/public/user/password', {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** Query User Subscribe GET /v1/public/user/subscribe */
export async function queryUserSubscribe(options?: { [key: string]: any }) {
  return request<API.Response & { data?: API.QueryUserSubscribeListResponse }>(
    '/v1/public/user/subscribe',
    {
      method: 'GET',
      ...(options || {}),
    },
  );
}

/** Get Subscribe Log GET /v1/public/user/subscribe_log */
export async function getSubscribeLog(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.GetSubscribeLogParams,
  options?: { [key: string]: any },
) {
  return request<API.Response & { data?: API.GetSubscribeLogResponse }>(
    '/v1/public/user/subscribe_log',
    {
      method: 'GET',
      params: {
        ...params,
      },
      ...(options || {}),
    },
  );
}

/** Reset User Subscribe Token PUT /v1/public/user/subscribe_token */
export async function resetUserSubscribeToken(
  body: API.ResetUserSubscribeTokenRequest,
  options?: { [key: string]: any },
) {
  return request<API.Response & { data?: any }>('/v1/public/user/subscribe_token', {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** Unbind Device PUT /v1/public/user/unbind_device */
export async function unbindDevice(
  body: API.UnbindDeviceRequest,
  options?: { [key: string]: any },
) {
  return request<API.Response & { data?: any }>('/v1/public/user/unbind_device', {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** Unbind OAuth POST /v1/public/user/unbind_oauth */
export async function unbindOAuth(body: API.UnbindOAuthRequest, options?: { [key: string]: any }) {
  return request<API.Response & { data?: any }>('/v1/public/user/unbind_oauth', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** Unbind Telegram POST /v1/public/user/unbind_telegram */
export async function unbindTelegram(options?: { [key: string]: any }) {
  return request<API.Response & { data?: any }>('/v1/public/user/unbind_telegram', {
    method: 'POST',
    ...(options || {}),
  });
}

/** 此处后端没有提供注释 POST /v1/public/user/unsubscribe */
export async function unsubscribe(body: API.UnsubscribeRequest, options?: { [key: string]: any }) {
  return request<API.Response & { data?: any }>('/v1/public/user/unsubscribe', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** Pre Unsubscribe POST /v1/public/user/unsubscribe/pre */
export async function preUnsubscribe(
  body: API.PreUnsubscribeRequest,
  options?: { [key: string]: any },
) {
  return request<API.Response & { data?: API.PreUnsubscribeResponse }>(
    '/v1/public/user/unsubscribe/pre',
    {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      data: body,
      ...(options || {}),
    },
  );
}

/** Verify Email POST /v1/public/user/verify_email */
export async function verifyEmail(body: API.VerifyEmailRequest, options?: { [key: string]: any }) {
  return request<API.Response & { data?: any }>('/v1/public/user/verify_email', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}
