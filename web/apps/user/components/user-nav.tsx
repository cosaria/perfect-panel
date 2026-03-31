'use client';

import { navs } from '@/config/navs';
import useGlobalStore from '@/config/use-global';
import { Logout } from '@/utils/common';
import { Avatar, AvatarFallback, AvatarImage } from '@workspace/ui/components/avatar';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@workspace/ui/components/dropdown-menu';
import { Icon } from '@workspace/ui/custom-components/icon';
import { useTranslations } from 'next-intl';
import { useRouter } from 'next/navigation';

export function UserNav() {
  const t = useTranslations('menu');
  const { user, setUser } = useGlobalStore();
  const router = useRouter();

  if (user) {
    return (
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <div className='bg-background hover:bg-accent flex cursor-pointer items-center gap-2 rounded-full border px-2 py-1.5 transition-colors duration-200'>
            <Avatar className='h-6 w-6'>
              <AvatarImage
                alt={user?.avatar ?? ''}
                src={user?.auth_methods?.[0]?.auth_identifier ?? ''}
                className='object-cover'
              />
              <AvatarFallback className='from-primary/90 to-primary text-background bg-gradient-to-br font-medium'>
                {user?.auth_methods?.[0]?.auth_identifier.toUpperCase().charAt(0)}
              </AvatarFallback>
            </Avatar>
            <span className='max-w-[40px] truncate text-sm sm:max-w-[100px]'>
              {user?.auth_methods?.[0]?.auth_identifier.split('@')[0]}
            </span>
            <Icon icon='mdi:chevron-down' className='text-muted-foreground size-4' />
          </div>
        </DropdownMenuTrigger>
        <DropdownMenuContent forceMount align='end' className='w-64'>
          <div className='flex items-center justify-start gap-2 p-2'>
            <Avatar className='h-10 w-10'>
              <AvatarImage
                alt={user?.avatar ?? ''}
                src={user?.avatar ?? ''}
                className='object-cover'
              />
              <AvatarFallback className='from-primary/90 to-primary text-background bg-gradient-to-br'>
                {user?.auth_methods?.[0]?.auth_identifier.toUpperCase().charAt(0)}
              </AvatarFallback>
            </Avatar>
            <div className='flex flex-col space-y-0.5'>
              <p className='text-sm font-medium leading-none'>
                {user?.auth_methods?.[0]?.auth_identifier.split('@')[0]}
              </p>
              <p className='text-muted-foreground text-xs'>
                {user?.auth_methods?.[0]?.auth_identifier}
              </p>
            </div>
          </div>
          <DropdownMenuSeparator />
          {navs.map((nav) => (
            <DropdownMenuGroup key={nav.title}>
              {(nav.items || [nav]).map((item) => (
                <DropdownMenuItem
                  key={item.title}
                  onClick={() => {
                    router.push(`${item.url}`);
                  }}
                  className='flex cursor-pointer items-center gap-2 py-2'
                >
                  <Icon className='text-muted-foreground size-4 flex-none' icon={item.icon!} />
                  <span className='flex-grow truncate'>{t(item.title)}</span>
                  <Icon
                    icon='lucide:chevron-right'
                    className='text-muted-foreground size-4 opacity-50'
                  />
                </DropdownMenuItem>
              ))}
            </DropdownMenuGroup>
          ))}
          <DropdownMenuSeparator />
          <DropdownMenuItem
            onClick={() => {
              Logout();
              setUser();
            }}
            className='text-destructive focus:text-destructive flex cursor-pointer items-center gap-2 py-2'
          >
            <Icon className='size-4 flex-none' icon='uil:exit' />
            <span className='flex-grow'>{t('logout')}</span>
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    );
  }
}
