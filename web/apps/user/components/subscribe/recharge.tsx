'use client';

import useGlobalStore from '@/config/use-global';
import { recharge } from '@/services/user/order';
import { Button, ButtonProps } from '@workspace/ui/components/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@workspace/ui/components/dialog';
import { EnhancedInput } from '@workspace/ui/custom-components/enhanced-input';
import { unitConversion } from '@workspace/ui/utils';
import { LoaderCircle } from 'lucide-react';
import { useTranslations } from 'next-intl';
import { useRouter } from 'next/navigation';
import { useState, useTransition } from 'react';
import PaymentMethods from './payment-methods';

export default function Recharge(props: Readonly<ButtonProps>) {
  const t = useTranslations('subscribe');
  const { common } = useGlobalStore();
  const { currency } = common;

  const router = useRouter();
  const [open, setOpen] = useState<boolean>(false);
  const [loading, startTransition] = useTransition();

  const [params, setParams] = useState<API.RechargeOrderRequest>({
    amount: 0,
    payment: 1,
  });

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button {...props}>{t('recharge')}</Button>
      </DialogTrigger>
      <DialogContent className='flex h-full flex-col overflow-hidden md:h-auto'>
        <DialogHeader>
          <DialogTitle>{t('balanceRecharge')}</DialogTitle>
          <DialogDescription>{t('rechargeDescription')}</DialogDescription>
        </DialogHeader>
        <div className='flex flex-col justify-between text-sm'>
          <div className='grid gap-3'>
            <div className='font-semibold'>{t('rechargeAmount')}</div>
            <div className='flex'>
              <EnhancedInput
                type='number'
                placeholder={t('enterAmount')}
                min={0}
                value={params.amount}
                formatInput={(value) => unitConversion('centsToDollars', value)}
                formatOutput={(value) => unitConversion('dollarsToCents', value)}
                onValueChange={(value) => {
                  setParams((prev) => ({
                    ...prev,
                    amount: value as number,
                  }));
                }}
                prefix={currency.currency_symbol}
                suffix={currency.currency_unit}
              />
            </div>
            <PaymentMethods
              balance={false}
              value={params.payment}
              onChange={(value) => setParams({ ...params, payment: value })}
            />
          </div>
          <Button
            className='fixed bottom-0 left-0 w-full rounded-none md:relative md:mt-6'
            disabled={loading || !params.amount}
            onClick={() => {
              startTransition(async () => {
                try {
                  const response = await recharge(params);
                  const orderNo = response.data.data?.order_no;
                  if (orderNo) {
                    router.push(`/payment?order_no=${orderNo}`);
                    setOpen(false);
                  }
                } catch (error) {
                  /* empty */
                }
              });
            }}
          >
            {loading && <LoaderCircle className='mr-2 animate-spin' />}
            {t('rechargeNow')}
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}
