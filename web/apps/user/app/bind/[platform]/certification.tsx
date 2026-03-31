'use client';

import { bindOAuthCallback } from '@/services/user/user';
import { getAllUrlParams } from '@/utils/common';
import { usePathname, useRouter } from 'next/navigation';
import { useEffect } from 'react';

interface CertificationProps {
  platform: string;
  children: React.ReactNode;
}

export default function Certification({ platform, children }: CertificationProps) {
  const router = useRouter();
  const pathname = usePathname();

  useEffect(() => {
    const searchParams = getAllUrlParams();
    bindOAuthCallback({
      method: platform,
      callback: searchParams,
    })
      .then((res) => {
        router.replace('/profile');
        router.refresh();
      })
      .catch((error) => {
        router.replace('/auth');
      });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [pathname]);

  return children;
}
