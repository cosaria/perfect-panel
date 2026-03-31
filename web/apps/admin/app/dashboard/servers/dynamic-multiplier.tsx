'use client';

import { getNodeMultiplier, setNodeMultiplier } from '@/services/admin/system';
import { useQuery } from '@tanstack/react-query';
import { Button } from '@workspace/ui/components/button';
import { Card, CardContent } from '@workspace/ui/components/card';
import { ScrollArea } from '@workspace/ui/components/scroll-area';
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetFooter,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from '@workspace/ui/components/sheet';
import { ArrayInput } from '@workspace/ui/custom-components/dynamic-Inputs';
import { Icon } from '@workspace/ui/custom-components/icon';
import { useTranslations } from 'next-intl';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';

export default function DynamicMultiplier() {
  const t = useTranslations('servers');
  const [open, setOpen] = useState(false);
  const [timeSlots, setTimeSlots] = useState<API.TimePeriod[]>([]);

  const { data: periodsResp, refetch: refetchPeriods } = useQuery({
    queryKey: ['getNodeMultiplier'],
    queryFn: async () => {
      const { data } = await getNodeMultiplier();
      return (data.data?.periods || []) as API.TimePeriod[];
    },
    enabled: open,
  });

  useEffect(() => {
    if (periodsResp) {
      setTimeSlots(periodsResp);
    }
  }, [periodsResp]);

  async function savePeriods() {
    await setNodeMultiplier({ periods: timeSlots });
    await refetchPeriods();
    toast.success(t('server_config.saveSuccess'));
    setOpen(false);
  }

  return (
    <Sheet open={open} onOpenChange={setOpen}>
      <SheetTrigger asChild>
        <Card>
          <CardContent className='p-4'>
            <div className='flex cursor-pointer items-center justify-between'>
              <div className='flex items-center gap-3'>
                <div className='bg-primary/10 flex h-10 w-10 items-center justify-center rounded-lg'>
                  <Icon icon='mdi:clock-time-eight' className='text-primary h-5 w-5' />
                </div>
                <div className='flex-1'>
                  <p className='font-medium'>{t('server_config.dynamic_multiplier')}</p>
                  <p className='text-muted-foreground truncate text-sm'>
                    {t('server_config.dynamic_multiplier_desc')}
                  </p>
                </div>
              </div>
              <Icon icon='mdi:chevron-right' className='size-6' />
            </div>
          </CardContent>
        </Card>
      </SheetTrigger>

      <SheetContent className='w-[600px] max-w-full md:max-w-screen-md'>
        <SheetHeader>
          <SheetTitle>{t('server_config.dynamic_multiplier')}</SheetTitle>
          <SheetDescription>{t('server_config.dynamic_multiplier_desc')}</SheetDescription>
        </SheetHeader>
        <ScrollArea className='-mx-6 h-[calc(100dvh-48px-36px-60px-env(safe-area-inset-top))] px-6'>
          <div className='space-y-4 pt-4'>
            <ArrayInput<API.TimePeriod>
              fields={[
                {
                  name: 'start_time',
                  prefix: t('server_config.fields.start_time'),
                  type: 'time',
                  step: '1',
                },
                {
                  name: 'end_time',
                  prefix: t('server_config.fields.end_time'),
                  type: 'time',
                  step: '1',
                },
                {
                  name: 'multiplier',
                  prefix: t('server_config.fields.multiplier'),
                  type: 'number',
                  placeholder: '0',
                },
              ]}
              value={timeSlots}
              onChange={setTimeSlots}
            />
          </div>
        </ScrollArea>

        <SheetFooter className='flex-row justify-between pt-3'>
          <Button variant='outline' onClick={() => setTimeSlots(periodsResp || [])}>
            {t('server_config.fields.reset')}
          </Button>
          <div className='flex gap-2'>
            <Button variant='outline' onClick={() => setOpen(false)}>
              {t('actions.cancel')}
            </Button>
            <Button onClick={savePeriods}>{t('actions.save')}</Button>
          </div>
        </SheetFooter>
      </SheetContent>
    </Sheet>
  );
}
