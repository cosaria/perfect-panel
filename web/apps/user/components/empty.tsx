'use client';

import { default as _Empty } from '@workspace/ui/custom-components/empty';
import { useTranslations } from 'next-intl';
import { useEffect, useState } from 'react';

interface EmptyProps {
  border?: boolean;
}
export function Empty({ border }: EmptyProps) {
  const t = useTranslations('common');

  const [description, setDescription] = useState('');

  useEffect(() => {
    const random = Math.floor(Math.random() * 10);
    setDescription(t(`empty.${random}`));
  }, [t]);

  return <_Empty description={description} border={border} />;
}
