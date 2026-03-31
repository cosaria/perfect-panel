'use client';

import { getVersion, restartSystem } from '@/services/admin/tool';
import { formatDate } from '@/utils/common';
import { useQuery } from '@tanstack/react-query';
import {
  AlertDialog,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from '@workspace/ui/components/alert-dialog';
import { Badge } from '@workspace/ui/components/badge';
import { Button } from '@workspace/ui/components/button';
import { Card, CardContent, CardHeader, CardTitle } from '@workspace/ui/components/card';
import { Icon } from '@workspace/ui/custom-components/icon';
import { useTranslations } from 'next-intl';
import Link from 'next/link';
import { useState } from 'react';
import packageJson from '../../../../package.json';
import SystemLogsDialog from './system-logs-dialog';

export default function SystemVersionCard() {
  const t = useTranslations('tool');
  const [openRestart, setOpenRestart] = useState(false);
  const [isRestarting, setIsRestarting] = useState(false);

  const { data: versionInfo } = useQuery({
    queryKey: ['getVersionInfo'],
    queryFn: async () => {
      try {
        const [webResponse, serverResponse, systemResponse] = await Promise.all([
          fetch(
            'https://data.jsdelivr.com/v1/packages/gh/perfect-panel/ppanel-web/resolved?specifier=latest',
          ),
          fetch(
            'https://data.jsdelivr.com/v1/packages/gh/perfect-panel/server/resolved?specifier=latest',
          ),
          getVersion(),
        ]);

        const webData = webResponse.ok ? await webResponse.json() : null;
        const serverData = serverResponse.ok ? await serverResponse.json() : null;
        const systemData = systemResponse.data.data;

        const rawVersion = (systemData?.version || '').replace(' Develop', '').trim();
        const timeMatch = rawVersion.match(/\(([^)]+)\)/);
        const timestamp = timeMatch ? timeMatch[1] : '';
        const versionWithoutTime = rawVersion.replace(/\([^)]*\)/, '').trim();

        const isDevelopment = !/^[Vv]?\d+\.\d+\.\d+(-[a-zA-Z]+(\.\d+)?)?$/.test(versionWithoutTime);

        let displayVersion = versionWithoutTime;
        if (
          !isDevelopment &&
          !versionWithoutTime.startsWith('V') &&
          !versionWithoutTime.startsWith('v')
        ) {
          displayVersion = `V${versionWithoutTime}`;
        }
        const lastUpdated = formatDate(new Date(timestamp || Date.now())) || '';

        const systemInfo = {
          isRelease: !isDevelopment,
          version: displayVersion,
          lastUpdated,
        };

        const latestReleases = {
          web: webData
            ? {
                version: webData.version,
                url: `https://github.com/perfect-panel/ppanel-web/releases/tag/v${webData.version}`,
              }
            : null,
          server: serverData
            ? {
                version: serverData.version,
                url: `https://github.com/perfect-panel/server/releases/tag/v${serverData.version}`,
              }
            : null,
        };

        const hasNewVersion =
          latestReleases.web &&
          packageJson.version !== latestReleases.web.version.replace(/^v/, '');

        const hasServerNewVersion =
          latestReleases.server &&
          systemInfo.version &&
          systemInfo.version.replace(/^V/, '') !== latestReleases.server.version.replace(/^v/, '');

        return {
          systemInfo,
          latestReleases,
          hasNewVersion,
          hasServerNewVersion,
        };
      } catch (error) {
        console.error('Failed to fetch version info:', error);
        return {
          systemInfo: { isRelease: true, version: 'V1.0.0', lastUpdated: '' },
          latestReleases: { web: null, server: null },
          hasNewVersion: false,
          hasServerNewVersion: false,
        };
      }
    },
    staleTime: 0,
    retry: 1,
    retryDelay: 10000,
    initialData: {
      systemInfo: { isRelease: true, version: 'V1.0.0', lastUpdated: '' },
      latestReleases: { web: null, server: null },
      hasNewVersion: false,
      hasServerNewVersion: false,
    },
  });

  const { systemInfo, latestReleases, hasNewVersion, hasServerNewVersion } = versionInfo;

  return (
    <Card className='p-3'>
      <CardHeader className='mb-2 p-0'>
        <CardTitle className='flex items-center justify-between'>
          {t('systemServices')}
          <div className='flex items-center space-x-2'>
            <SystemLogsDialog variant='outline' size='sm' />
            <AlertDialog open={openRestart} onOpenChange={setOpenRestart}>
              <AlertDialogTrigger asChild>
                <Button variant='destructive' size='sm'>
                  {t('systemReboot')}
                </Button>
              </AlertDialogTrigger>
              <AlertDialogContent>
                <AlertDialogHeader>
                  <AlertDialogTitle>{t('confirmSystemReboot')}</AlertDialogTitle>
                  <AlertDialogDescription>{t('rebootDescription')}</AlertDialogDescription>
                </AlertDialogHeader>
                <AlertDialogFooter>
                  <AlertDialogCancel>{t('cancel')}</AlertDialogCancel>
                  <Button
                    onClick={async () => {
                      setIsRestarting(true);
                      await restartSystem();
                      await new Promise((resolve) => setTimeout(resolve, 5000));
                      setIsRestarting(false);
                      setOpenRestart(false);
                    }}
                    disabled={isRestarting}
                  >
                    {isRestarting && <Icon icon='mdi:loading' className='mr-2 animate-spin' />}
                    {isRestarting ? t('rebooting') : t('confirmReboot')}
                  </Button>
                </AlertDialogFooter>
              </AlertDialogContent>
            </AlertDialog>
          </div>
        </CardTitle>
      </CardHeader>
      <CardContent className='space-y-3 p-0'>
        <div className='flex flex-1 items-center justify-between'>
          <div className='flex items-center'>
            <Icon icon='mdi:web' className='mr-2 h-4 w-4 text-green-600' />
            <span className='text-sm font-medium'>{t('webVersion')}</span>
          </div>
          <div className='flex items-center space-x-2'>
            <Badge>V{packageJson.version}</Badge>
            {hasNewVersion && (
              <Link
                href={
                  latestReleases?.web?.url || 'https://github.com/perfect-panel/ppanel-web/releases'
                }
                target='_blank'
                rel='noopener noreferrer'
                className='flex items-center space-x-1'
              >
                <Badge variant='destructive' className='animate-pulse px-2 py-0.5 text-xs'>
                  {t('newVersionAvailable')}
                  <Icon icon='mdi:open-in-new' />
                </Badge>
              </Link>
            )}
          </div>
        </div>
        <div className='flex flex-1 items-center justify-between'>
          <div className='flex items-center'>
            <Icon icon='mdi:server' className='mr-2 h-4 w-4 text-blue-600' />
            <span className='text-sm font-medium'>{t('serverVersion')}</span>
          </div>
          <div className='flex items-center space-x-2'>
            <Badge variant={!systemInfo?.isRelease ? 'destructive' : 'default'}>
              {systemInfo?.version || 'V1.0.0'}
            </Badge>
            {hasServerNewVersion && (
              <Link
                href={
                  latestReleases?.server?.url || 'https://github.com/perfect-panel/server/releases'
                }
                target='_blank'
                rel='noopener noreferrer'
                className='flex items-center space-x-1'
              >
                <Badge variant='destructive' className='animate-pulse px-2 py-0.5 text-xs'>
                  {t('newVersionAvailable')}
                  <Icon icon='mdi:open-in-new' />
                </Badge>
              </Link>
            )}
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
