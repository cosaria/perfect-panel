'use client';
import { findNavByUrl } from '@/config/navs';
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from '@workspace/ui/components/breadcrumb';
import { Separator } from '@workspace/ui/components/separator';
import { SidebarTrigger } from '@workspace/ui/components/sidebar';
import { useTranslations } from 'next-intl';
import { usePathname } from 'next/navigation';
import { Fragment, useMemo } from 'react';
import LanguageSwitch from './language-switch';
import ThemeSwitch from './theme-switch';
import TimezoneSwitch from './timezone-switch';
import { UserNav } from './user-nav';

export function Header() {
  const t = useTranslations('menu');
  const pathname = usePathname();
  const items = useMemo(() => findNavByUrl(pathname), [pathname]);
  return (
    <header className='bg-background sticky top-0 z-50 flex h-14 shrink-0 items-center gap-2'>
      <div className='flex flex-1 items-center gap-2 px-3'>
        <SidebarTrigger />
        <Separator orientation='vertical' className='mr-2 h-4' />
        <Breadcrumb>
          <BreadcrumbList>
            {items.map((item, index) => {
              return (
                <Fragment key={item?.title}>
                  {index !== items.length - 1 && (
                    <BreadcrumbItem>
                      <BreadcrumbLink href={item?.url || '/dashboard'}>
                        {t(item?.title)}
                      </BreadcrumbLink>
                    </BreadcrumbItem>
                  )}
                  {index < items.length - 1 && <BreadcrumbSeparator />}
                  {index === items.length - 1 && <BreadcrumbPage>{t(item?.title)}</BreadcrumbPage>}
                </Fragment>
              );
            })}
          </BreadcrumbList>
        </Breadcrumb>
      </div>
      <div className='flex items-center gap-2 px-3'>
        <LanguageSwitch />
        <TimezoneSwitch />
        <ThemeSwitch />
        <UserNav />
      </div>
    </header>
  );
}
