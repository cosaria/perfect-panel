'use client';

import { Card, CardContent } from '@workspace/ui/components/card';
import { useTranslations } from 'next-intl';
import ConfigForm from './config-form';
import { ProtocolForm } from './protocol-form';

export default function SubscribePage() {
  const t = useTranslations('subscribe');

  return (
    <div className='space-y-4'>
      <h2 className='text-lg font-semibold'>{t('config.title')}</h2>
      <Card>
        <CardContent className='p-4'>
          <ConfigForm />
        </CardContent>
      </Card>

      <ProtocolForm />
    </div>
  );
}
