'use client';

import { Display } from '@/components/display';
import { Empty } from '@/components/empty';
import { ProList } from '@/components/pro-list';
import useGlobalStore from '@/config/use-global';
import { queryUserAffiliate, queryUserAffiliateList } from '@/services/user/user';
import { useQuery } from '@tanstack/react-query';
import { Button } from '@workspace/ui/components/button';
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@workspace/ui/components/card';
import { formatDate, isBrowser } from '@workspace/ui/utils';
import { Copy } from 'lucide-react';
import { useTranslations } from 'next-intl';
import { useState } from 'react';
import { CopyToClipboard } from 'react-copy-to-clipboard';
import { toast } from 'sonner';

export default function Affiliate() {
  const t = useTranslations('affiliate');
  const { user, common } = useGlobalStore();
  const [sum, setSum] = useState<number>();
  const { data } = useQuery({
    queryKey: ['queryUserAffiliate'],
    queryFn: async () => {
      const response = await queryUserAffiliate();
      return response.data.data;
    },
  });

  return (
    <div className='flex flex-col gap-4'>
      <Card>
        <CardHeader>
          <CardTitle>{t('totalCommission')}</CardTitle>
          <CardDescription>{t('commissionInfo')}</CardDescription>
        </CardHeader>
        <CardContent>
          <div className='flex items-baseline gap-2'>
            <span className='text-3xl font-bold'>
              <Display type='currency' value={data?.total_commission} />
            </span>
            <span className='text-muted-foreground text-sm'>
              ({t('commissionRate')}:{' '}
              {user?.referral_percentage || common?.invite?.referral_percentage}%)
            </span>
          </div>
        </CardContent>
      </Card>
      <Card>
        <CardHeader className='flex flex-row items-center justify-between space-y-0'>
          <CardTitle className='text-lg font-medium'>{t('inviteCode')}</CardTitle>
        </CardHeader>
        <CardContent>
          <div className='flex items-center justify-between'>
            <code className='bg-muted rounded px-2 py-1 text-2xl font-bold'>
              {user?.refer_code}
            </code>
            {isBrowser() && (
              <CopyToClipboard
                text={`${location?.origin}/auth?invite=${user?.refer_code}`}
                onCopy={(text, result) => {
                  if (result) {
                    toast.success(t('copySuccess'));
                  }
                }}
              >
                <Button variant='secondary' size='sm' className='gap-2'>
                  <Copy className='h-4 w-4' />
                  {t('copyInviteLink')}
                </Button>
              </CopyToClipboard>
            )}
          </div>
        </CardContent>
      </Card>
      <ProList<API.UserAffiliate, Record<string, unknown>>
        request={async (pagination, filter) => {
          const response = await queryUserAffiliateList({ ...pagination, ...filter });
          setSum(response.data.data?.sum);
          return {
            list: response.data.data?.list || [],
            total: response.data.data?.total || 0,
          };
        }}
        header={{
          title: t('inviteRecords'),
        }}
        renderItem={(item) => {
          return (
            <Card className='overflow-hidden'>
              <CardContent className='p-3 text-sm'>
                <ul className='grid grid-cols-2 gap-3 *:flex *:flex-col'>
                  <li className='font-semibold'>
                    <span className='text-muted-foreground'>{t('userIdentifier')}</span>
                    <span>{item.identifier}</span>
                  </li>
                  <li className='font-semibold'>
                    <span className='text-muted-foreground'>{t('registrationTime')}</span>
                    <time>{formatDate(item.registered_at)}</time>
                  </li>
                </ul>
              </CardContent>
            </Card>
          );
        }}
        empty={<Empty />}
      />
    </div>
  );
}
