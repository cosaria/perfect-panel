'use client';

import { getAvailablePaymentMethods } from '@/services/user/portal';
import { useQuery } from '@tanstack/react-query';
import { Label } from '@workspace/ui/components/label';
import { RadioGroup, RadioGroupItem } from '@workspace/ui/components/radio-group';
import { cn } from '@workspace/ui/lib/utils';
import { useTranslations } from 'next-intl';
import Image from 'next/image';
import React, { memo } from 'react';

interface PaymentMethodsProps {
  value: number;
  onChange: (value: number) => void;
  balance?: boolean;
}

const PaymentMethods: React.FC<PaymentMethodsProps> = ({ value, onChange, balance = true }) => {
  const t = useTranslations('subscribe');

  const { data } = useQuery({
    queryKey: ['getAvailablePaymentMethods', { balance }],
    queryFn: async () => {
      const { data } = await getAvailablePaymentMethods();
      const list = data.data?.list || [];
      const methods = balance ? list : list.filter((item) => item.id !== -1);
      const defaultMethod = methods.find((item) => item.id)?.id;
      if (defaultMethod) onChange(defaultMethod);
      return methods;
    },
  });
  return (
    <>
      <div className='font-semibold'>{t('paymentMethod')}</div>
      <RadioGroup
        className='grid grid-cols-2 gap-2 md:grid-cols-5'
        value={String(value)}
        onValueChange={(val) => {
          console.log(val);
          onChange(Number(val));
        }}
      >
        {data?.map((item) => (
          <div key={item.id} className='relative'>
            <RadioGroupItem
              value={String(item.id)}
              id={String(item.id)}
              className='absolute inset-0 z-10 h-full w-full cursor-pointer opacity-0'
            />
            <Label
              htmlFor={String(item.id)}
              className={cn(
                'border-muted bg-popover hover:bg-accent hover:text-accent-foreground flex flex-col items-center justify-between rounded-md border-2 py-2',
                String(value) === String(item.id) ? 'border-primary' : '',
              )}
            >
              <div className='mb-3 size-12'>
                <Image
                  src={item.icon || `/payment/balance.svg`}
                  width={48}
                  height={48}
                  alt={item.name}
                />
              </div>
              <span className='w-full overflow-hidden text-ellipsis whitespace-nowrap text-center'>
                {item.name}
              </span>
            </Label>
          </div>
        ))}
      </RadioGroup>
    </>
  );
};

export default memo(PaymentMethods);
