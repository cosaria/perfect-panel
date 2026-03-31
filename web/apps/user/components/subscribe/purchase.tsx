'use client';

import CouponInput from '@/components/subscribe/coupon-input';
import DurationSelector from '@/components/subscribe/duration-selector';
import PaymentMethods from '@/components/subscribe/payment-methods';
import useGlobalStore from '@/config/use-global';
import { preCreateOrder, purchase } from '@/services/user/order';
import { useQuery } from '@tanstack/react-query';
import { Button } from '@workspace/ui/components/button';
import { Card, CardContent } from '@workspace/ui/components/card';
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@workspace/ui/components/dialog';
import { Separator } from '@workspace/ui/components/separator';
import { LoaderCircle } from 'lucide-react';
import { useTranslations } from 'next-intl';
import { useRouter } from 'next/navigation';
import { useCallback, useEffect, useRef, useState, useTransition } from 'react';
import { SubscribeBilling } from './billing';
import { SubscribeDetail } from './detail';

interface PurchaseProps {
  subscribe?: API.Subscribe;
  setSubscribe: (subscribe?: API.Subscribe) => void;
}

export default function Purchase({ subscribe, setSubscribe }: Readonly<PurchaseProps>) {
  const t = useTranslations('subscribe');
  const { getUserInfo } = useGlobalStore();
  const router = useRouter();
  const [params, setParams] = useState<Partial<API.PurchaseOrderRequest>>({
    quantity: 1,
    subscribe_id: 0,
    payment: -1,
    coupon: '',
  });
  const [loading, startTransition] = useTransition();
  const lastSuccessOrderRef = useRef<any>(null);

  const { data: order } = useQuery({
    enabled: !!subscribe?.id,
    queryKey: ['preCreateOrder', params],
    queryFn: async () => {
      try {
        const { data } = await preCreateOrder({
          ...params,
          subscribe_id: subscribe?.id as number,
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
        throw error;
      }
    },
  });

  useEffect(() => {
    if (subscribe) {
      setParams((prev) => ({
        ...prev,
        quantity: 1,
        subscribe_id: subscribe?.id,
      }));
    }
  }, [subscribe]);

  const handleChange = useCallback((field: keyof typeof params, value: string | number) => {
    setParams((prev) => ({
      ...prev,
      [field]: value,
    }));
  }, []);

  const handleSubmit = useCallback(async () => {
    startTransition(async () => {
      try {
        const response = await purchase(params as API.PurchaseOrderRequest);
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
    <Dialog
      open={!!subscribe?.id}
      onOpenChange={(open) => {
        if (!open) setSubscribe(undefined);
      }}
    >
      <DialogContent className='flex h-full max-w-screen-lg flex-col overflow-hidden border-none p-0 md:h-auto'>
        <DialogHeader className='p-6 pb-0'>
          <DialogTitle>{t('buySubscription')}</DialogTitle>
        </DialogHeader>
        <div className='grid w-full flex-grow gap-3 overflow-auto p-6 pt-0 lg:grid-cols-2'>
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
