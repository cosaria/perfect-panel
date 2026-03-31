'use client';

import { Display } from '@/components/display';
import StripePayment from '@/components/payment/stripe';
import { SubscribeBilling } from '@/components/subscribe/billing';
import { SubscribeDetail } from '@/components/subscribe/detail';
import useGlobalStore from '@/config/use-global';
import { purchaseCheckout, queryPurchaseOrder } from '@/services/user/portal';
import { setAuthorization } from '@/utils/common';
import { useQuery } from '@tanstack/react-query';
import { Badge } from '@workspace/ui/components/badge';
import { Button } from '@workspace/ui/components/button';
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@workspace/ui/components/card';
import { Separator } from '@workspace/ui/components/separator';
import { Icon } from '@workspace/ui/custom-components/icon';
import { formatDate } from '@workspace/ui/utils';
import { useCountDown } from 'ahooks';
import { addMinutes, format } from 'date-fns';
import { useTranslations } from 'next-intl';
import Link from 'next/link';
import { QRCodeCanvas } from 'qrcode.react';
import { useEffect, useState } from 'react';

export default function Page() {
  const t = useTranslations('order');
  const { getUserInfo } = useGlobalStore();
  const [orderNo, setOrderNo] = useState<string>();
  const [enabled, setEnabled] = useState<boolean>(false);

  const { data } = useQuery({
    enabled: enabled,
    queryKey: ['queryPurchaseOrder', orderNo],
    queryFn: async () => {
      if (!orderNo) return;
      const params = localStorage.getItem(orderNo);
      const authParams = params ? JSON.parse(params) : {};
      const { data } = await queryPurchaseOrder({ order_no: orderNo!, ...authParams });
      if (data?.data?.status !== 1) {
        setEnabled(false);
        if (data?.data?.token) {
          setAuthorization(data?.data?.token);
          await new Promise((resolve) => setTimeout(resolve, 100));
          await getUserInfo();
        }
      }
      return data?.data;
    },
    refetchInterval: 3000,
  });

  const { data: payment } = useQuery({
    enabled: !!orderNo && data?.status === 1,
    queryKey: ['purchaseCheckout', orderNo],
    queryFn: async () => {
      const { data } = await purchaseCheckout({
        orderNo: orderNo!,
        returnUrl: window.location.href,
      });
      if (data.data?.type === 'url' && data.data?.checkout_url) {
        window.open(data.data.checkout_url, '_blank');
      }
      return data?.data;
    },
  });

  useEffect(() => {
    const searchParams = new URLSearchParams(location.search);
    if (searchParams.get('order_no')) {
      setOrderNo(searchParams.get('order_no')!);
      setEnabled(true);
    }
  }, []);

  const [countDown, formattedRes] = useCountDown({
    targetDate: data && format(addMinutes(data?.created_at, 15), "yyyy-MM-dd'T'HH:mm:ss.SSSxxx"),
  });

  const { hours, minutes, seconds } = formattedRes;

  const countdownDisplay =
    countDown > 0 ? (
      <>
        {hours.toString().length === 1 ? `0${hours}` : hours} :{' '}
        {minutes.toString().length === 1 ? `0${minutes}` : minutes} :{' '}
        {seconds.toString().length === 1 ? `0${seconds}` : seconds}
      </>
    ) : (
      <>{t('timeExpired')}</>
    );

  return (
    <main className='container lg:mt-16'>
      <div className='grid gap-4 xl:grid-cols-2'>
        <Card className='order-2 xl:order-1'>
          <CardHeader className='bg-muted/50 flex flex-row items-start'>
            <div className='grid gap-0.5'>
              <CardTitle className='flex flex-col text-lg'>
                {t('orderNumber')}
                <span>{data?.order_no}</span>
              </CardTitle>
              <CardDescription>
                {t('createdAt')}: {formatDate(data?.created_at)}
              </CardDescription>
            </div>
          </CardHeader>
          <CardContent className='grid gap-3 p-6 text-sm'>
            <div className='font-semibold'>{t('paymentMethod')}</div>
            <dl className='grid gap-3'>
              <div className='flex items-center justify-between'>
                <dt className='text-muted-foreground'>
                  <Badge>{data?.payment.name || data?.payment.platform}</Badge>
                </dt>
              </div>
            </dl>
            <Separator />

            {data?.status && [1, 2].includes(data.status) && (
              <SubscribeDetail
                subscribe={{
                  ...data?.subscribe,
                  quantity: data?.quantity,
                }}
              />
            )}
            {data?.status === 3 && (
              <>
                <div className='font-semibold'>{t('resetTraffic')}</div>
                <ul className='grid grid-cols-2 gap-3 *:flex *:items-center *:justify-between lg:grid-cols-1'>
                  <li className='flex items-center justify-between'>
                    <span className='text-muted-foreground line-clamp-2 flex-1'>
                      {t('resetPrice')}
                    </span>
                    <span>
                      <Display type='currency' value={data.amount} />
                    </span>
                  </li>
                </ul>
              </>
            )}

            {data?.status === 4 && (
              <>
                <div className='font-semibold'>{t('balanceRecharge')}</div>
                <ul className='grid grid-cols-2 gap-3 *:flex *:items-center *:justify-between lg:grid-cols-1'>
                  <li className='flex items-center justify-between'>
                    <span className='text-muted-foreground line-clamp-2 flex-1'>
                      {t('rechargeAmount')}
                    </span>
                    <span>
                      <Display type='currency' value={data.amount} />
                    </span>
                  </li>
                </ul>
              </>
            )}
            <Separator />
            <SubscribeBilling
              order={{
                ...data,
                unit_price: data?.subscribe?.unit_price,
              }}
            />
          </CardContent>
        </Card>
        <Card className='order-1 flex flex-auto items-center justify-center xl:order-2'>
          <CardContent className='py-16'>
            {data?.status && [2, 5].includes(data?.status) && (
              <div className='flex flex-col items-center gap-8 text-center'>
                <h3 className='text-2xl font-bold tracking-tight'>{t('paymentSuccess')}</h3>
                <Icon icon='mdi:success-circle-outline' className='text-7xl text-green-500' />
                <div className='flex gap-4'>
                  <Button asChild>
                    <Link href='/dashboard'>{t('subscribeNow')}</Link>
                  </Button>
                  <Button variant='outline'>
                    <Link href='/document'>{t('viewDocument')}</Link>
                  </Button>
                </div>
              </div>
            )}
            {data?.status === 1 && payment?.type === 'url' && (
              <div className='flex flex-col items-center gap-8 text-center'>
                <h3 className='text-2xl font-bold tracking-tight'>{t('waitingForPayment')}</h3>
                <p className='flex items-center text-3xl font-bold'>{countdownDisplay}</p>
                <Icon icon='mdi:access-time' className='text-muted-foreground text-7xl' />
                <div className='flex gap-4'>
                  <Button
                    onClick={() => {
                      if (payment?.checkout_url) {
                        window.location.href = payment?.checkout_url;
                      }
                    }}
                  >
                    {t('goToPayment')}
                  </Button>
                  <Button variant='outline'>
                    <Link href='/'>{t('productList')}</Link>
                  </Button>
                </div>
              </div>
            )}

            {data?.status === 1 && payment?.type === 'qr' && (
              <div className='flex flex-col items-center gap-8 text-center'>
                <h3 className='text-2xl font-bold tracking-tight'>{t('scanToPay')}</h3>
                <p className='flex items-center text-3xl font-bold'>{countdownDisplay}</p>
                <QRCodeCanvas
                  value={payment?.checkout_url || ''}
                  size={208}
                  imageSettings={{
                    src: `/payment/alipay_f2f.svg`,
                    width: 24,
                    height: 24,
                    excavate: true,
                  }}
                />
                <div className='flex gap-4'>
                  <Button asChild>
                    <Link href='/subscribe'>{t('productList')}</Link>
                  </Button>
                  <Button asChild variant='outline'>
                    <Link href='/order'>{t('orderList')}</Link>
                  </Button>
                </div>
              </div>
            )}

            {data?.status === 1 && payment?.type === 'stripe' && (
              <div className='flex flex-col items-center gap-8 text-center'>
                <h3 className='text-2xl font-bold tracking-tight'>{t('waitingForPayment')}</h3>
                <p className='flex items-center text-3xl font-bold'>{countdownDisplay}</p>
                {payment.stripe && <StripePayment {...payment.stripe} />}
                {/* <div className='flex gap-4'>
                <Button asChild>
                  <Link href='/subscribe'>{t('productList')}</Link>
                </Button>
                <Button asChild variant='outline'>
                  <Link href='/order'>{t('orderList')}</Link>
                </Button>
              </div> */}
              </div>
            )}

            {data?.status && [3, 4].includes(data?.status) && (
              <div className='flex flex-col items-center gap-8 text-center'>
                <h3 className='text-2xl font-bold tracking-tight'>{t('orderClosed')}</h3>
                <Icon icon='mdi:cancel' className='text-7xl text-red-500' />
                <div className='flex gap-4'>
                  <Button asChild>
                    <Link href='/subscribe'>{t('productList')}</Link>
                  </Button>
                  <Button asChild variant='outline'>
                    <Link href='/order'>{t('orderList')}</Link>
                  </Button>
                </div>
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </main>
  );
}
