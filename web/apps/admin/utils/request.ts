import { NEXT_PUBLIC_API_URL, NEXT_PUBLIC_SITE_URL } from '@/config/constants';
import { getTranslations } from '@/locales/utils';
import { isBrowser } from '@workspace/ui/utils';
import axios, { InternalAxiosRequestConfig } from 'axios';
import { toast } from 'sonner';
import { getAuthorization, Logout } from './common';

async function handleError(response: any) {
  const code = response.data?.code;
  if ([40002, 40003, 40004, 40005].includes(code)) return Logout();
  if (response?.config?.skipErrorHandler) return;
  if (!isBrowser()) return;

  const t = await getTranslations('common');
  const message =
    t(`request.${code}`) !== `request.${code}`
      ? t(`request.${code}`)
      : response.data?.message || response.message;

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
  async (error: Error) => {
    await handleError(error);
    return Promise.reject(error);
  },
);

export default requset;
