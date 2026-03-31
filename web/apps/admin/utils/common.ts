import { locales, NEXT_PUBLIC_DEFAULT_LANGUAGE } from '@/config/constants';
import { isBrowser } from '@workspace/ui/utils';
import { intlFormat } from 'date-fns';
import Cookies from 'universal-cookie';

const cookies = new Cookies(null, {
  path: '/',
  maxAge: 365 * 24 * 60 * 60,
});

export function getLocale() {
  const browserLocale = navigator.language?.split(',')?.[0] || '';
  const defaultLocale = locales.includes(browserLocale) ? browserLocale : '';
  const cookies = new Cookies(null, { path: '/' });
  const cookieLocale = cookies.get('locale') || '';
  const locale = cookieLocale || defaultLocale || NEXT_PUBLIC_DEFAULT_LANGUAGE;
  return locale;
}

export function setLocale(value: string) {
  return cookies.set('locale', value);
}

export function setAuthorization(value: string) {
  return cookies.set('Authorization', value);
}

export function getAuthorization(value?: string) {
  const Authorization = isBrowser() ? cookies.get('Authorization') : value;
  if (!Authorization) return;
  return Authorization;
}

export function setRedirectUrl(value?: string) {
  if (value) {
    sessionStorage.setItem('redirect-url', value);
  }
}

export function getRedirectUrl() {
  return sessionStorage.getItem('redirect-url') ?? '/dashboard';
}

export function Logout() {
  if (!isBrowser()) return;
  cookies.remove('Authorization');
  const pathname = location.pathname;
  if (!['', '/'].includes(pathname)) {
    setRedirectUrl(location.pathname);
    location.href = `/`;
  }
}

export function formatDate(date?: Date | number, showTime: boolean = true) {
  if (!date) return;
  const timeZone = localStorage.getItem('timezone') || 'UTC';
  return intlFormat(date, {
    year: 'numeric',
    month: 'numeric',
    day: 'numeric',
    ...(showTime && {
      hour: 'numeric',
      minute: 'numeric',
      second: 'numeric',
    }),
    hour12: false,
    timeZone,
  });
}
