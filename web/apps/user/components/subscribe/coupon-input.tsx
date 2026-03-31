'use client';

import { Input } from '@workspace/ui/components/input';
import { useTranslations } from 'next-intl';
import React from 'react';

interface CouponInputProps {
  coupon?: string;
  onChange: (value: string) => void;
}

const CouponInput: React.FC<CouponInputProps> = ({ coupon, onChange }) => {
  const t = useTranslations('subscribe');

  return (
    <>
      <div className='font-semibold'>{t('coupon')}</div>
      <Input
        placeholder={t('enterCoupon')}
        value={coupon}
        onChange={(e) => onChange(e.target.value.trim())}
      />
    </>
  );
};

export default CouponInput;
