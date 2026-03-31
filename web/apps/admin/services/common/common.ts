// @ts-ignore
/* eslint-disable */
import request from '@/utils/request';

/** Get Ads GET /v1/common/ads */
export async function getAds(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.GetAdsParams,
  options?: { [key: string]: any },
) {
  return request<API.Response & { data?: API.GetAdsResponse }>('/v1/common/ads', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** Check verification code POST /v1/common/check_verification_code */
export async function checkVerificationCode(
  body: API.CheckVerificationCodeRequest,
  options?: { [key: string]: any },
) {
  return request<API.Response & { data?: API.CheckVerificationCodeRespone }>(
    '/v1/common/check_verification_code',
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

/** Get Client GET /v1/common/client */
export async function getClient(options?: { [key: string]: any }) {
  return request<API.Response & { data?: API.GetSubscribeClientResponse }>('/v1/common/client', {
    method: 'GET',
    ...(options || {}),
  });
}

/** Get verification code POST /v1/common/send_code */
export async function sendEmailCode(body: API.SendCodeRequest, options?: { [key: string]: any }) {
  return request<API.Response & { data?: API.SendCodeResponse }>('/v1/common/send_code', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** Get sms verification code POST /v1/common/send_sms_code */
export async function sendSmsCode(body: API.SendSmsCodeRequest, options?: { [key: string]: any }) {
  return request<API.Response & { data?: API.SendCodeResponse }>('/v1/common/send_sms_code', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** Get global config GET /v1/common/site/config */
export async function getGlobalConfig(options?: { [key: string]: any }) {
  return request<API.Response & { data?: API.GetGlobalConfigResponse }>('/v1/common/site/config', {
    method: 'GET',
    ...(options || {}),
  });
}

/** Get Privacy Policy GET /v1/common/site/privacy */
export async function getPrivacyPolicy(options?: { [key: string]: any }) {
  return request<API.Response & { data?: API.PrivacyPolicyConfig }>('/v1/common/site/privacy', {
    method: 'GET',
    ...(options || {}),
  });
}

/** Get stat GET /v1/common/site/stat */
export async function getStat(options?: { [key: string]: any }) {
  return request<API.Response & { data?: API.GetStatResponse }>('/v1/common/site/stat', {
    method: 'GET',
    ...(options || {}),
  });
}

/** Get Tos Content GET /v1/common/site/tos */
export async function getTos(options?: { [key: string]: any }) {
  return request<API.Response & { data?: API.GetTosResponse }>('/v1/common/site/tos', {
    method: 'GET',
    ...(options || {}),
  });
}
