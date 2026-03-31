'use client';

import useGlobalStore from '@/config/use-global';
import { buttonVariants } from '@workspace/ui/components/button';
import { useTranslations } from 'next-intl';
import Image from 'next/legacy/image';
import Link from 'next/link';
import LanguageSwitch from './language-switch';
import ThemeSwitch from './theme-switch';
import { UserNav } from './user-nav';

export default function Header() {
  const t = useTranslations('common');

  const { common, user } = useGlobalStore();
  const { site } = common;
  const Logo = (
    <Link href='/' className='flex items-center gap-2 text-lg font-bold'>
      {site.site_logo && (
        <Image src={site.site_logo} width={48} height={48} alt='logo' unoptimized />
      )}
      <span className=''>{site.site_name}</span>
    </Link>
  );
  return (
    <header className='sticky top-0 z-50 border-b backdrop-blur-md'>
      <div className='container flex h-16 items-center justify-between'>
        <nav className='flex-col gap-6 text-lg font-medium md:flex md:flex-row md:items-center md:gap-5 md:text-sm lg:gap-6'>
          {Logo}
        </nav>
        <div className='flex flex-1 items-center justify-end gap-2'>
          <LanguageSwitch />
          <ThemeSwitch />
          <UserNav />
          {!user && (
            <Link
              href='/auth'
              className={buttonVariants({
                size: 'sm',
              })}
            >
              {t('login')}
            </Link>
          )}
        </div>
      </div>
    </header>
  );
}
