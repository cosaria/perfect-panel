'use client';

import { Badge } from '@workspace/ui/components/badge';
import { Label } from '@workspace/ui/components/label';
import { RadioGroup, RadioGroupItem } from '@workspace/ui/components/radio-group';
import { useTranslations } from 'next-intl';
import React, { useCallback } from 'react';

interface DurationSelectorProps {
  quantity: number;
  unitTime?: string;
  discounts?: Array<{ quantity: number; discount: number }>;
  onChange: (value: number) => void;
}

const DurationSelector: React.FC<DurationSelectorProps> = ({
  quantity,
  unitTime = 'Month',
  discounts = [],
  onChange,
}) => {
  const t = useTranslations('subscribe');
  const handleChange = useCallback(
    (value: string) => {
      onChange(Number(value));
    },
    [onChange],
  );

  const DurationOption: React.FC<{ value: string; label: string }> = ({ value, label }) => (
    <div className='relative'>
      <RadioGroupItem value={value} id={value} className='peer sr-only' />
      <Label
        htmlFor={value}
        className='border-muted bg-popover hover:bg-accent hover:text-accent-foreground peer-data-[state=checked]:border-primary relative flex h-full flex-col items-center justify-center gap-2 rounded-md border-2 p-2'
      >
        {label}
      </Label>
    </div>
  );

  const currentDiscount = discounts?.find((item) => item.quantity === quantity)?.discount;
  const discountPercentage = currentDiscount ? 100 - currentDiscount : 0;

  return (
    <>
      <div className='font-semibold'>{t('purchaseDuration')}</div>
      <RadioGroup
        value={String(quantity)}
        onValueChange={handleChange}
        className='flex flex-wrap gap-3'
      >
        {unitTime !== 'Minute' && <DurationOption value='1' label={`1 / ${t(unitTime)}`} />}
        {discounts?.map((item) => (
          <DurationOption
            key={item.quantity}
            value={String(item.quantity)}
            label={`${item.quantity} / ${t(unitTime)}`}
          />
        ))}
      </RadioGroup>
      <div className='flex items-center justify-between'>
        <span className='text-muted-foreground text-sm'>{t('discountInfo')}:</span>
        {discountPercentage > 0 ? (
          <Badge variant='destructive' className='h-6 text-sm'>
            -{discountPercentage}% {t('discount')}
          </Badge>
        ) : (
          <span className='text-muted-foreground h-6 text-sm'>--</span>
        )}
      </div>
    </>
  );
};

export default DurationSelector;
