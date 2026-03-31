'use client';

import { useLocale } from 'next-intl';
import { useTheme } from 'next-themes';
import { forwardRef, useEffect, useImperativeHandle } from 'react';
import Turnstile, { useTurnstile } from 'react-turnstile';

import useGlobalStore from '@/config/use-global';

export type TurnstileRef = {
  reset: () => void;
};

const CloudFlareTurnstile = forwardRef<
  TurnstileRef,
  {
    id?: string;
    value?: null | string;
    onChange: (value?: string) => void;
  }
>(function CloudFlareTurnstile({ id, value, onChange }, ref) {
  const { common } = useGlobalStore();
  const { verify } = common;
  const { resolvedTheme } = useTheme();
  const locale = useLocale();
  const turnstile = useTurnstile();

  useImperativeHandle(
    ref,
    () => ({
      reset: () => turnstile.reset(),
    }),
    [turnstile],
  );

  useEffect(() => {
    if (value === '') {
      turnstile.reset();
    }
  }, [turnstile, value]);

  return (
    verify.turnstile_site_key && (
      <Turnstile
        id={id}
        sitekey={verify.turnstile_site_key}
        theme={resolvedTheme as 'light' | 'dark'}
        language={locale.toLowerCase()}
        fixedSize
        onVerify={(token) => onChange(token)}
        // onError={() => {
        //   onChange();
        //   turnstile.reset();
        // }}
        onExpire={() => {
          onChange();
          turnstile.reset();
        }}
        onTimeout={() => {
          onChange();
          turnstile.reset();
        }}
      />
    )
  );
});

export default CloudFlareTurnstile;
