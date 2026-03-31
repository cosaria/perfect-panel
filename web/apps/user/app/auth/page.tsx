'use client';

import { OAuthMethods } from '@/components/auth/oauth-methods';
import LanguageSwitch from '@/components/language-switch';
import ThemeSwitch from '@/components/theme-switch';
import useGlobalStore from '@/config/use-global';
import { DotLottieReact } from '@lottiefiles/dotlottie-react';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@workspace/ui/components/tabs';
import LoginLottie from '@workspace/ui/lotties/login.json';
import { useTranslations } from 'next-intl';
import Image from 'next/legacy/image';
import Link from 'next/link';
import EmailAuthForm from './email/auth-form';
import PhoneAuthForm from './phone/auth-form';

export default function Page() {
  const t = useTranslations('auth');
  const { common } = useGlobalStore();
  const { site, auth } = common;

  const AUTH_METHODS = [
    {
      key: 'email',
      enabled: auth.email.enable,
      children: <EmailAuthForm />,
    },
    {
      key: 'mobile',
      enabled: auth.mobile.enable,
      children: <PhoneAuthForm />,
    },
  ].filter((method) => method.enabled);

  return (
    <main className='bg-muted/50 flex h-full min-h-screen items-center'>
      <div className='flex size-full flex-auto flex-col lg:flex-row'>
        <div className='flex bg-cover bg-center lg:w-1/2 lg:flex-auto'>
          <div className='lg:py-15 md:px-15 flex w-full flex-col items-center justify-center px-5 py-7'>
            <Link className='mb-0 flex flex-col items-center lg:mb-12' href='/'>
              {site.site_logo && (
                <Image src={site.site_logo} height={48} width={48} alt='logo' unoptimized />
              )}
              <span className='text-2xl font-semibold'>{site.site_name}</span>
            </Link>
            <DotLottieReact
              data={LoginLottie}
              autoplay
              loop
              className='mx-auto hidden w-[275px] lg:block xl:w-[500px]'
            />
            <p className='hidden w-[275px] text-center md:w-1/2 lg:block xl:w-[500px]'>
              {site.site_desc}
            </p>
          </div>
        </div>
        <div className='flex flex-initial justify-center p-12 lg:flex-auto lg:justify-end'>
          <div className='lg:bg-background flex w-full flex-col items-center rounded-2xl md:w-[600px] md:p-10 lg:flex-auto lg:shadow'>
            <div className='flex w-full flex-col items-stretch justify-center md:w-[400px] lg:h-full'>
              <div className='flex flex-col justify-center lg:flex-auto'>
                <h1 className='mb-3 text-center text-2xl font-bold'>{t('verifyAccount')}</h1>
                <div className='text-muted-foreground mb-6 text-center font-medium'>
                  {t('verifyAccountDesc')}
                </div>
                {AUTH_METHODS.length === 1
                  ? AUTH_METHODS[0]?.children
                  : AUTH_METHODS[0] && (
                      <Tabs defaultValue={AUTH_METHODS[0].key}>
                        <TabsList className='mb-6 flex w-full *:flex-1'>
                          {AUTH_METHODS.map((item) => (
                            <TabsTrigger key={item.key} value={item.key}>
                              {t(`methods.${item.key}`)}
                            </TabsTrigger>
                          ))}
                        </TabsList>
                        {AUTH_METHODS.map((item) => (
                          <TabsContent key={item.key} value={item.key}>
                            {item.children}
                          </TabsContent>
                        ))}
                      </Tabs>
                    )}
              </div>
              <div className='py-8'>
                <OAuthMethods />
              </div>
              <div className='flex items-center justify-between'>
                <div className='flex items-center gap-5'>
                  <LanguageSwitch />
                  <ThemeSwitch />
                </div>
                <div className='text-primary flex gap-2 text-sm font-semibold'>
                  <Link href='/tos'>{t('tos')}</Link>
                  <span className='text-foreground/30'>|</span>
                  <Link href='/privacy-policy'>{t('privacyPolicy')}</Link>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>
  );
}
