'use client';

import CouponInput from '@/components/subscribe/coupon-input';
import DurationSelector from '@/components/subscribe/duration-selector';
import PaymentMethods from '@/components/subscribe/payment-methods';
import useGlobalStore from '@/config/use-global';
import { preCreateOrder, renewal } from '@/services/user/order';
import { useQuery } from '@tanstack/react-query';
import { Button } from '@workspace/ui/components/button';
import { Card, CardContent } from '@workspace/ui/components/card';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@workspace/ui/components/dialog';
import { Separator } from '@workspace/ui/components/separator';
import { LoaderCircle } from 'lucide-react';
import { useTranslations } from 'next-intl';
import { useRouter } from 'next/navigation';
import { useCallback, useEffect, useRef, useState, useTransition } from 'react';
import { SubscribeBilling } from './billing';
import { SubscribeDetail } from './detail';

interface RenewalProps {
  id: number;
  subscribe: API.Subscribe;
}

export default function Renewal({ id, subscribe }: Readonly<RenewalProps>) {
  const t = useTranslations('subscribe');
  const { getUserInfo } = useGlobalStore();
  const [open, setOpen] = useState<boolean>(false);
  const router = useRouter();
  const [params, setParams] = useState<Partial<API.RenewalOrderRequest>>({
    quantity: 1,
    payment: -1,
    coupon: '',
    user_subscribe_id: id,
  });
  const [loading, startTransition] = useTransition();
  const lastSuccessOrderRef = useRef<any>(null);

  const { data: order } = useQuery({
    enabled: !!subscribe.id && open,
    queryKey: ['preCreateOrder', params],
    queryFn: async () => {
      try {
        const { data } = await preCreateOrder({
          ...params,
          subscribe_id: subscribe.id,
        } as API.PurchaseOrderRequest);
        const result = data.data;
        if (result) {
          lastSuccessOrderRef.current = result;
        }
        return result;
      } catch (error) {
        if (lastSuccessOrderRef.current) {
          return lastSuccessOrderRef.current;
        }
      }
    },
  });

  useEffect(() => {
    if (subscribe.id && id) {
      setParams((prev) => ({
        ...prev,
        quantity: 1,
        subscribe_id: subscribe.id,
        user_subscribe_id: id,
      }));
    }
  }, [subscribe.id, id]);

  const handleChange = useCallback((field: keyof typeof params, value: string | number) => {
    setParams((prev) => ({
      ...prev,
      [field]: value,
    }));
  }, []);

  const handleSubmit = useCallback(async () => {
    startTransition(async () => {
      try {
        const response = await renewal(params as API.RenewalOrderRequest);
        const orderNo = response.data.data?.order_no;
        if (orderNo) {
          getUserInfo();
          router.push(`/payment?order_no=${orderNo}`);
        }
      } catch (error) {
        /* empty */
      }
    });
  }, [params, router, getUserInfo]);

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button size='sm'>{t('renew')}</Button>
      </DialogTrigger>
      <DialogContent className='flex h-full max-w-screen-lg flex-col overflow-hidden md:h-auto'>
        <DialogHeader>
          <DialogTitle>{t('renewSubscription')}</DialogTitle>
        </DialogHeader>
        <div className='grid w-full gap-3 lg:grid-cols-2'>
          <Card className='border-transparent shadow-none md:border-inherit md:shadow'>
            <CardContent className='grid gap-3 p-0 text-sm md:p-6'>
              <SubscribeDetail
                subscribe={{
                  ...subscribe,
                  quantity: params.quantity,
                }}
              />
              <Separator />
              <SubscribeBilling
                order={{
                  ...order,
                  quantity: params.quantity,
                  unit_price: subscribe?.unit_price,
                }}
              />
            </CardContent>
          </Card>
          <div className='flex flex-col justify-between text-sm'>
            <div className='mb-6 grid gap-3'>
              <DurationSelector
                quantity={params.quantity!}
                unitTime={subscribe?.unit_time}
                discounts={subscribe?.discount}
                onChange={(value) => {
                  handleChange('quantity', value);
                }}
              />
              <CouponInput
                coupon={params.coupon}
                onChange={(value) => handleChange('coupon', value)}
              />
              <PaymentMethods
                value={params.payment!}
                onChange={(value) => {
                  handleChange('payment', value);
                }}
              />
            </div>
            <Button
              className='fixed bottom-0 left-0 w-full rounded-none md:relative md:mt-6'
              disabled={loading}
              onClick={handleSubmit}
            >
              {loading && <LoaderCircle className='mr-2 animate-spin' />}
              {t('buyNow')}
            </Button>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}
