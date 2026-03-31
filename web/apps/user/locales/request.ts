import { locales, NEXT_PUBLIC_DEFAULT_LANGUAGE } from '@/config/constants';
import { getRequestConfig } from 'next-intl/server';
import { cookies, headers } from 'next/headers';

export default getRequestConfig(async () => {
  const browserLocale = (await headers()).get('Accept-Language')?.split(',')?.[0] || '';
  const defaultLocale = locales.includes(browserLocale) ? browserLocale : '';
  const cookieLocale = (await cookies()).get('locale')?.value || '';

  const locale = cookieLocale || defaultLocale || NEXT_PUBLIC_DEFAULT_LANGUAGE;

  const messages = {
    menu: (await import(`./${locale}/menu.json`)).default,
    auth: (await import(`./${locale}/auth.json`)).default,
    common: (await import(`./${locale}/common.json`)).default,
    layout: (await import(`./${locale}/layout.json`)).default,
    index: (await import(`./${locale}/index.json`)).default,
    dashboard: (await import(`./${locale}/dashboard.json`)).default,
    profile: (await import(`./${locale}/profile.json`)).default,
    subscribe: (await import(`./${locale}/subscribe.json`)).default,
    order: (await import(`./${locale}/order.json`)).default,
    payment: (await import(`./${locale}/payment.json`)).default,
    wallet: (await import(`./${locale}/wallet.json`)).default,
    ticket: (await import(`./${locale}/ticket.json`)).default,
    document: (await import(`./${locale}/document.json`)).default,
    affiliate: (await import(`./${locale}/affiliate.json`)).default,
  };

  return {
    locale,
    messages,
  };
});
