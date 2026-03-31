// @ts-expect-error
/* eslint-disable */
import request from "@/utils/request";

/** Get available payment methods GET /v1/public/payment/methods */
export async function getAvailablePaymentMethods(options?: Record<string, unknown>) {
  return request<API.Response & { data?: API.GetAvailablePaymentMethodsResponse }>(
    "/v1/public/payment/methods",
    {
      method: "GET",
      ...(options || {}),
    },
  );
}
