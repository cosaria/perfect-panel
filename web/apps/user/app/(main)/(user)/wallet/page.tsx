'use client';

import { Display } from '@/components/display';
import { ProList, ProListActions } from '@/components/pro-list';
import useGlobalStore from '@/config/use-global';
import { queryUserBalanceLog } from '@/services/user/user';
import { Card, CardContent } from '@workspace/ui/components/card';
import { useTranslations } from 'next-intl';
import { useRef } from 'react';

import { Empty } from '@/components/empty';
import Recharge from '@/components/subscribe/recharge';
import { formatDate } from '@workspace/ui/utils';

export default function Page() {
  const t = useTranslations('wallet');
  const { user } = useGlobalStore();
  const ref = useRef<ProListActions>(null);
  const totalAssets = (user?.balance || 0) + (user?.commission || 0) + (user?.gift_amount || 0);
  return (
    <>
      <Card className='mb-4'>
        <CardContent className='p-6'>
          <h2 className='text-foreground mb-4 text-2xl font-bold'>{t('assetOverview')}</h2>
          <div className='mb-4'>
            <div className='flex items-center justify-between'>
              <div>
                <p className='text-sm font-medium'>{t('totalAssets')}</p>
                <p className='text-3xl font-bold'>
                  <Display type='currency' value={totalAssets} />
                </p>
              </div>
              <Recharge />
            </div>
          </div>
          <div className='grid grid-cols-1 gap-6 md:grid-cols-3'>
            <div className='bg-secondary rounded-lg p-4 shadow-sm transition-all duration-300 hover:shadow-md'>
              <p className='text-secondary-foreground text-sm font-medium opacity-80'>
                {t('balance')}
              </p>
              <p className='text-secondary-foreground text-2xl font-bold'>
                <Display type='currency' value={user?.balance} />
              </p>
            </div>
            <div className='bg-secondary rounded-lg p-4 shadow-sm transition-all duration-300 hover:shadow-md'>
              <p className='text-secondary-foreground text-sm font-medium opacity-80'>
                {t('giftAmount')}
              </p>
              <p className='text-secondary-foreground text-2xl font-bold'>
                <Display type='currency' value={user?.gift_amount} />
              </p>
            </div>
            <div className='bg-secondary rounded-lg p-4 shadow-sm transition-all duration-300 hover:shadow-md'>
              <p className='text-secondary-foreground text-sm font-medium opacity-80'>
                {t('commission')}
              </p>
              <p className='text-secondary-foreground text-2xl font-bold'>
                <Display type='currency' value={user?.commission} />
              </p>
            </div>
          </div>
        </CardContent>
      </Card>
      <ProList<API.BalanceLog, Record<string, unknown>>
        action={ref}
        request={async (pagination, filter) => {
          const response = await queryUserBalanceLog({ ...pagination, ...filter });
          return {
            list: response.data.data?.list || [],
            total: response.data.data?.total || 0,
          };
        }}
        renderItem={(item) => {
          return (
            <Card className='overflow-hidden'>
              <CardContent className='p-3 text-sm'>
                <ul className='grid grid-cols-2 gap-3 *:flex *:flex-col lg:grid-cols-4'>
                  <li className='font-semibold'>
                    <span className='text-muted-foreground'>{t('createdAt')}</span>
                    <time>{formatDate(item.timestamp)}</time>
                  </li>
                  <li className='font-semibold'>
                    <span className='text-muted-foreground'>{t('type.0')}</span>
                    <span>{t(`type.${item.type}`)}</span>
                  </li>
                  <li className='font-semibold'>
                    <span className='text-muted-foreground'>{t('amount')}</span>
                    <span>
                      <Display type='currency' value={item.amount} />
                    </span>
                  </li>

                  <li>
                    <span className='text-muted-foreground'>{t('balance')}</span>
                    <span>
                      <Display type='currency' value={item.balance} />
                    </span>
                  </li>
                </ul>
              </CardContent>
            </Card>
          );
        }}
        empty={<Empty />}
      />
    </>
  );
}
