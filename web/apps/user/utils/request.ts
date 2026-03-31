import { isBrowser } from "@workspace/ui/utils";
import axios, { type AxiosError, type AxiosResponse, type InternalAxiosRequestConfig } from "axios";
import { toast } from "sonner";
import { NEXT_PUBLIC_API_URL, NEXT_PUBLIC_SITE_URL } from "@/config/constants";
import { getTranslations } from "@/locales/utils";
import { getAuthorization, Logout } from "./common";

type ErrorResponseData = {
  code?: number;
  message?: string;
};

type RequestConfigWithErrorHandler = InternalAxiosRequestConfig & {
  skipErrorHandler?: boolean;
};

async function handleError(
  response: AxiosError<ErrorResponseData> | AxiosResponse<ErrorResponseData>,
) {
  const data = "data" in response ? response.data : response.response?.data;
  const config = ("data" in response ? response.config : response.config) as
    | RequestConfigWithErrorHandler
    | undefined;
  const code = data?.code;

  if (typeof code === "number" && [40002, 40003, 40004, 40005].includes(code)) {
    return Logout();
  }
  if (config?.skipErrorHandler) return;
  if (!isBrowser()) return;

  const t = await getTranslations("common");
  const message =
    t(`request.${code}`) !== `request.${code}`
      ? t(`request.${code}`)
      : data?.message || ("message" in response ? response.message : undefined);

  toast.error(message);
}

const requset = axios.create({
  baseURL: NEXT_PUBLIC_API_URL || NEXT_PUBLIC_SITE_URL,
  // timeout: 10000,
  // withCredentials: true,
});

requset.interceptors.request.use(
  async (
    config: InternalAxiosRequestConfig & {
      Authorization?: string;
      skipErrorHandler?: boolean;
    },
  ) => {
    const Authorization = getAuthorization(config.Authorization);
    if (Authorization) config.headers.Authorization = Authorization;
    return config;
  },
  (error: Error) => Promise.reject(error),
);

requset.interceptors.response.use(
  async (response) => {
    const { code } = response.data;
    if (code !== 200) {
      await handleError(response);
      throw response;
    }
    return response;
  },
  async (error: AxiosError<ErrorResponseData>) => {
    await handleError(error);
    return Promise.reject(error);
  },
);

export default requset;
