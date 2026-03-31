'use client';

import { Display } from '@/components/display';
import Recharge from '@/components/subscribe/recharge';
import useGlobalStore from '@/config/use-global';
import { Button } from '@workspace/ui/components/button';
import { Card, CardContent, CardHeader, CardTitle } from '@workspace/ui/components/card';
import { Sidebar, SidebarContent } from '@workspace/ui/components/sidebar';
import { Icon } from '@workspace/ui/custom-components/icon';
import { isBrowser } from '@workspace/ui/utils';
import { useTranslations } from 'next-intl';
import CopyToClipboard from 'react-copy-to-clipboard';
import { toast } from 'sonner';

export function SidebarRight({ ...props }: React.ComponentProps<typeof Sidebar>) {
  const { user } = useGlobalStore();
  const t = useTranslations('layout');

  return (
    <Sidebar collapsible='none' side='right' {...props}>
      <SidebarContent>
        <Card>
          <CardHeader className='flex flex-row items-center justify-between space-y-0 p-3 pb-2'>
            <CardTitle className='text-sm font-medium'>{t('accountBalance')}</CardTitle>
            <Recharge variant='link' className='p-0' />
          </CardHeader>
          <CardContent className='p-3 text-2xl font-bold'>
            <Display type='currency' value={user?.balance} />
          </CardContent>
        </Card>
        <Card>
          <CardHeader className='space-y-0 p-3 pb-2'>
            <CardTitle className='text-sm font-medium'>{t('giftAmount')}</CardTitle>
          </CardHeader>
          <CardContent className='p-3 text-2xl font-bold'>
            <Display type='currency' value={user?.gift_amount} />
          </CardContent>
        </Card>
        <Card>
          <CardHeader className='space-y-0 p-3 pb-2'>
            <CardTitle className='text-sm font-medium'>{t('commission')}</CardTitle>
          </CardHeader>
          <CardContent className='p-3 text-2xl font-bold'>
            <Display type='currency' value={user?.commission} />
          </CardContent>
        </Card>
        {user?.refer_code && (
          <Card>
            <CardHeader className='flex flex-row items-center justify-between space-y-0 p-3 pb-2'>
              <CardTitle className='text-sm font-medium'>{t('inviteCode')}</CardTitle>
              <CopyToClipboard
                text={`${isBrowser() && location?.origin}/auth?invite=${user?.refer_code}`}
                onCopy={(text, result) => {
                  if (result) {
                    toast.success(t('copySuccess'));
                  }
                }}
              >
                <Button variant='ghost' className='size-5 p-0'>
                  <Icon icon='mdi:content-copy' className='text-primary text-2xl' />
                </Button>
              </CopyToClipboard>
            </CardHeader>
            <CardContent className='truncate p-3 font-bold'>{user?.refer_code}</CardContent>
          </Card>
        )}
      </SidebarContent>
    </Sidebar>
  );
}
