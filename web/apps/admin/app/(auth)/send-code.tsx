'use client';

import { sendEmailCode, sendSmsCode } from '@/services/common/common';
import { Button } from '@workspace/ui/components/button';
import { useCountDown } from 'ahooks';
import { useTranslations } from 'next-intl';
import { useState } from 'react';

interface SendCodeProps {
  type: 'email' | 'phone';
  params: {
    email?: string;
    type?: 1 | 2;
    telephone_area_code?: string;
    telephone?: string;
  };
}
export default function SendCode({ type, params }: SendCodeProps) {
  const t = useTranslations('auth');
  const [targetDate, setTargetDate] = useState<number>();

  const [, { seconds }] = useCountDown({
    targetDate,
    onEnd: () => {
      setTargetDate(undefined);
    },
  });

  const getEmailCode = async () => {
    if (params.email && params.type) {
      await sendEmailCode({
        email: params.email,
        type: params.type,
      });
      setTargetDate(Date.now() + 60000);
    }
  };

  const getPhoneCode = async () => {
    if (params.telephone && params.telephone_area_code && params.type) {
      await sendSmsCode({
        telephone: params.telephone,
        telephone_area_code: params.telephone_area_code,
        type: params.type,
      });
      setTargetDate(Date.now() + 60000);
    }
  };

  const handleSendCode = async () => {
    if (type === 'email') {
      getEmailCode();
    } else {
      getPhoneCode();
    }
  };
  const disabled =
    seconds > 0 ||
    (type === 'email' ? !params.email : !params.telephone || !params.telephone_area_code);

  return (
    <Button type='button' onClick={handleSendCode} disabled={disabled}>
      {seconds > 0 ? `${seconds}s` : t('get')}
    </Button>
  );
}
