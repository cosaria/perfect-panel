'use client';

import { getNodeConfig } from '@/services/admin/system';
import { useQuery } from '@tanstack/react-query';
import { Button } from '@workspace/ui/components/button';
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@workspace/ui/components/dialog';
import { Input } from '@workspace/ui/components/input';
import { Label } from '@workspace/ui/components/label';
import { useTranslations } from 'next-intl';
import React, { useCallback, useEffect, useMemo, useState } from 'react';
import { toast } from 'sonner';

type Props = {
  server: API.Server;
};

export default function ServerInstall({ server }: Props) {
  const t = useTranslations('servers');
  const [open, setOpen] = useState(false);
  const [domain, setDomain] = useState('');

  const { data: cfgResp } = useQuery({
    queryKey: ['getNodeConfig'],
    queryFn: async () => {
      const { data } = await getNodeConfig();
      return data.data as API.NodeConfig | undefined;
    },
    enabled: open,
  });

  useEffect(() => {
    if (open) {
      const host = localStorage.getItem('API_HOST') ?? window.location.origin;
      setDomain(host);
    }
  }, [open]);

  const installCommand = useMemo(() => {
    const secret = cfgResp?.node_secret ?? '';
    return `wget -N https://raw.githubusercontent.com/perfect-panel/ppanel-node/master/scripts/install.sh && bash install.sh --api-host ${domain} --server-id ${server.id} --secret-key ${secret}`;
  }, [domain, server.id, cfgResp?.node_secret]);

  async function handleCopy() {
    try {
      if (navigator?.clipboard?.writeText) {
        await navigator.clipboard.writeText(installCommand);
      } else {
        // fallback for environments without clipboard API
        const el = document.createElement('textarea');
        el.value = installCommand;
        document.body.appendChild(el);
        el.select();
        document.execCommand('copy');
        document.body.removeChild(el);
      }
      toast.success(t('copied'));
      setOpen(false);
    } catch (error) {
      toast.error(t('copyFailed'));
    }
  }

  const onDomainChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    setDomain(e.target.value);
    localStorage.setItem('API_HOST', e.target.value);
  }, []);
  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button variant='secondary'>{t('connect')}</Button>
      </DialogTrigger>

      <DialogContent className='w-[720px] max-w-full md:max-w-screen-md'>
        <DialogHeader>
          <DialogTitle>{t('oneClickInstall')}</DialogTitle>
        </DialogHeader>

        <div className='space-y-4'>
          <div>
            <Label>{t('apiHost')}</Label>
            <div className='flex items-center gap-2'>
              <Input
                value={domain}
                placeholder={t('apiHostPlaceholder')}
                onChange={onDomainChange}
              />
            </div>
          </div>

          <div>
            <Label>{t('installCommand')}</Label>
            <div className='flex flex-col gap-2'>
              <textarea
                readOnly
                aria-label={t('installCommand')}
                value={installCommand}
                className='min-h-[88px] w-full rounded border p-2 font-mono text-sm'
              />
            </div>
          </div>
        </div>

        <DialogFooter className='flex-row justify-end gap-2 pt-3'>
          <Button variant='outline' onClick={() => setOpen(false)}>
            {t('close')}
          </Button>
          <Button onClick={handleCopy}>{t('copyAndClose')}</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
