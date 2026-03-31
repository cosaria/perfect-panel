'use client';

import { getSystemLog } from '@/services/admin/tool';
import { useQuery } from '@tanstack/react-query';
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from '@workspace/ui/components/accordion';
import { Button } from '@workspace/ui/components/button';
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@workspace/ui/components/dialog';
import { ScrollArea } from '@workspace/ui/components/scroll-area';
import { Icon } from '@workspace/ui/custom-components/icon';
import { useTranslations } from 'next-intl';
import { useState } from 'react';

const getLogLevelColor = (level: string) => {
  const colorMap: { [key: string]: string } = {
    INFO: 'bg-blue-100 text-blue-800 hover:bg-blue-200',
    WARN: 'bg-yellow-100 text-yellow-800 hover:bg-yellow-200',
    ERROR: 'bg-red-100 text-red-800 hover:bg-red-200',
  };
  return colorMap[level] || 'bg-gray-100 text-gray-800 hover:bg-gray-200';
};

interface SystemLogsDialogProps {
  trigger?: React.ReactNode;
  variant?: 'default' | 'outline' | 'ghost' | 'secondary';
  size?: 'sm' | 'default' | 'lg';
}

export default function SystemLogsDialog({
  trigger,
  variant = 'outline',
  size = 'sm',
}: SystemLogsDialogProps) {
  const t = useTranslations('tool');
  const [open, setOpen] = useState(false);

  const {
    data: logs,
    refetch,
    isLoading,
  } = useQuery({
    queryKey: ['getSystemLog'],
    queryFn: async () => {
      const { data } = await getSystemLog();
      return data.data?.list || [];
    },
    enabled: open,
  });

  const defaultTrigger = (
    <Button variant={variant} size={size}>
      {t('systemLogs')}
    </Button>
  );

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>{trigger || defaultTrigger}</DialogTrigger>
      <DialogContent className='max-w-4xl'>
        <DialogHeader>
          <DialogTitle>{t('systemLogs')}</DialogTitle>
        </DialogHeader>
        <ScrollArea className='bg-muted/30 h-[60vh] max-h-[80vh] min-h-[400px] w-full rounded-lg border p-1'>
          {isLoading ? (
            <div className='flex h-full items-center justify-center'>
              <Icon icon='uil:loading' className='text-primary h-8 w-8 animate-spin' />
            </div>
          ) : (
            <Accordion type='single' collapsible className='w-full'>
              {logs?.map((log: any, index: number) => (
                <AccordionItem key={index} value={`item-${index}`} className='px-4'>
                  <AccordionTrigger className='hover:no-underline'>
                    <div className='flex w-full flex-col items-start space-y-2 sm:flex-row sm:items-center sm:space-x-4 sm:space-y-0'>
                      <span className='text-xs font-medium sm:text-sm'>{log.timestamp}</span>
                    </div>
                  </AccordionTrigger>
                  <AccordionContent className='px-2'>
                    {Object.entries(log).map(([key, value]) => (
                      <div
                        key={key}
                        className='grid grid-cols-1 gap-2 text-xs sm:grid-cols-2 sm:text-sm'
                      >
                        <span className='font-medium'>{key}:</span>
                        <span className='break-all'>{value as string}</span>
                      </div>
                    ))}
                  </AccordionContent>
                </AccordionItem>
              ))}
            </Accordion>
          )}
        </ScrollArea>
        <DialogFooter>
          <Button
            variant='outline'
            onClick={() => {
              refetch();
            }}
          >
            <Icon icon='uil:refresh' className={`h-5 w-5 ${isLoading ? 'animate-spin' : ''}`} />
            <span>{t('refreshLogs')}</span>
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
