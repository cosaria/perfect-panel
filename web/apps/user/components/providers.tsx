'use client';

import useGlobalStore, { GlobalStore } from '@/config/use-global';
import { Logout } from '@/utils/common';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { ReactQueryStreamedHydration } from '@tanstack/react-query-next-experimental';
import { ThemeProvider as NextThemesProvider } from 'next-themes';
import React, { useEffect, useState } from 'react';
import Loading from './loading';

export default function Providers({
  children,
  common,
  user,
}: {
  children: React.ReactNode;
  common: Partial<GlobalStore['common']>;
  user: GlobalStore['user'];
}) {
  const [loading, setLoading] = useState(true);
  const [queryClient] = useState(
    () =>
      new QueryClient({
        defaultOptions: {
          queries: {
            staleTime: 5 * 1000,
            retry: false,
          },
        },
      }),
  );

  const { setCommon, setUser } = useGlobalStore();

  useEffect(() => {
    const initializeData = async () => {
      try {
        if (user) {
          setUser(user);
        } else {
          Logout();
        }
        setCommon(common);
      } finally {
        setTimeout(() => {
          setLoading(false);
        }, 1000);
      }
    };

    initializeData();
  }, [setUser, setCommon, user, common]);

  useEffect(() => {
    const searchParams = new URLSearchParams(location.search);
    const invite = searchParams.get('invite');
    if (invite) {
      localStorage.setItem('invite', invite);
    }
  }, []);

  return (
    <NextThemesProvider attribute='class' defaultTheme='system' enableSystem>
      <QueryClientProvider client={queryClient}>
        <ReactQueryStreamedHydration>
          <Loading loading={loading || queryClient.isMutating() > 0} />
          {children}
        </ReactQueryStreamedHydration>
      </QueryClientProvider>
    </NextThemesProvider>
  );
}
