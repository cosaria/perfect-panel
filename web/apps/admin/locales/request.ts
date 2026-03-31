import { locales, NEXT_PUBLIC_DEFAULT_LANGUAGE } from '@/config/constants';
import { getRequestConfig } from 'next-intl/server';
import { cookies, headers } from 'next/headers';

export default getRequestConfig(async () => {
  const browserLocale = (await headers()).get('Accept-Language')?.split(',')?.[0] || '';
  const defaultLocale = locales.includes(browserLocale) ? browserLocale : '';
  const cookieLocale = (await cookies()).get('locale')?.value || '';

  const locale = cookieLocale || defaultLocale || NEXT_PUBLIC_DEFAULT_LANGUAGE;

  const messages = {
    'menu': (await import(`./${locale}/menu.json`)).default,
    'auth': (await import(`./${locale}/auth.json`)).default,
    'common': (await import(`./${locale}/common.json`)).default,
    'system': (await import(`./${locale}/system.json`)).default,
    'auth-control': (await import(`./${locale}/auth-control.json`)).default,
    'ads': (await import(`./${locale}/ads.json`)).default,
    'payment': (await import(`./${locale}/payment.json`)).default,
    'servers': (await import(`./${locale}/servers.json`)).default,
    'nodes': (await import(`./${locale}/nodes.json`)).default,
    'product': (await import(`./${locale}/product.json`)).default,
    'order': (await import(`./${locale}/order.json`)).default,
    'coupon': (await import(`./${locale}/coupon.json`)).default,
    'user': (await import(`./${locale}/user.json`)).default,
    'announcement': (await import(`./${locale}/announcement.json`)).default,
    'ticket': (await import(`./${locale}/ticket.json`)).default,
    'document': (await import(`./${locale}/document.json`)).default,
    'tool': (await import(`./${locale}/tool.json`)).default,
    'index': (await import(`./${locale}/index.json`)).default,
    'subscribe': (await import(`./${locale}/subscribe.json`)).default,
    'marketing': (await import(`./${locale}/marketing.json`)).default,
    'log': (await import(`./${locale}/log.json`)).default,
  };

  return {
    locale,
    messages,
  };
});
