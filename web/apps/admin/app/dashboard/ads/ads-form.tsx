'use client';

import { zodResolver } from '@hookform/resolvers/zod';
import { Button } from '@workspace/ui/components/button';
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@workspace/ui/components/form';
import { RadioGroup, RadioGroupItem } from '@workspace/ui/components/radio-group';
import { ScrollArea } from '@workspace/ui/components/scroll-area';
import {
  Sheet,
  SheetContent,
  SheetFooter,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from '@workspace/ui/components/sheet';
import { EnhancedInput } from '@workspace/ui/custom-components/enhanced-input';
import { Icon } from '@workspace/ui/custom-components/icon';
import { useTranslations } from 'next-intl';
import { useEffect, useState } from 'react';
import { useForm } from 'react-hook-form';
import { z } from 'zod';

const formSchema = z.object({
  title: z.string(),
  type: z.enum(['image', 'video']),
  content: z.string(),
  description: z.string(),
  target_url: z.string().url(),
  start_time: z.number(),
  end_time: z.number(),
});

interface AdsFormProps<T> {
  onSubmit: (data: T) => Promise<boolean> | boolean;
  initialValues?: T;
  loading?: boolean;
  trigger: string;
  title: string;
}

export default function AdsForm<T extends Record<string, any>>({
  onSubmit,
  initialValues,
  loading,
  trigger,
  title,
}: AdsFormProps<T>) {
  const t = useTranslations('ads');
  const [open, setOpen] = useState(false);

  const form = useForm({
    resolver: zodResolver(formSchema),
    defaultValues: {
      ...initialValues,
    } as any,
  });

  useEffect(() => {
    form?.reset(initialValues);
  }, [form, initialValues]);

  const type = form.watch('type');
  const startTime = form.watch('start_time');

  const renderContentField = () => {
    return (
      <FormField
        control={form.control}
        name='content'
        render={({ field }) => (
          <FormItem>
            <FormLabel>{t('form.content')}</FormLabel>
            <FormControl>
              <EnhancedInput
                placeholder={
                  type === 'image'
                    ? 'https://example.com/image.jpg'
                    : 'https://example.com/video.mp4'
                }
                value={field.value}
                onValueChange={(value) => {
                  form.setValue('content', value);
                }}
              />
            </FormControl>
            <FormMessage />
          </FormItem>
        )}
      />
    );
  };

  async function handleSubmit(data: { [x: string]: any }) {
    const bool = await onSubmit(data as T);
    if (bool) setOpen(false);
  }

  return (
    <Sheet open={open} onOpenChange={setOpen}>
      <SheetTrigger asChild>
        <Button
          onClick={() => {
            form.reset();
            setOpen(true);
          }}
        >
          {trigger}
        </Button>
      </SheetTrigger>
      <SheetContent className='w-[500px] max-w-full md:max-w-screen-md'>
        <SheetHeader>
          <SheetTitle>{title}</SheetTitle>
        </SheetHeader>
        <ScrollArea className='-mx-6 h-[calc(100vh-48px-36px-36px-env(safe-area-inset-top))]'>
          <Form {...form}>
            <form onSubmit={form.handleSubmit(handleSubmit)} className='space-y-4 px-6 pt-4'>
              <FormField
                control={form.control}
                name='title'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('form.title')}</FormLabel>
                    <FormControl>
                      <EnhancedInput
                        placeholder={t('form.enterTitle')}
                        value={field.value}
                        onValueChange={(value) => {
                          form.setValue(field.name, value);
                        }}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name='type'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('form.type')}</FormLabel>
                    <FormControl>
                      <RadioGroup
                        defaultValue={field.value}
                        onValueChange={(value) => {
                          form.setValue(field.name, value);
                        }}
                        className='flex gap-4'
                      >
                        <FormItem className='flex items-center space-x-3 space-y-0'>
                          <FormControl>
                            <RadioGroupItem value='image' />
                          </FormControl>
                          <FormLabel className='font-normal'>{t('form.typeImage')}</FormLabel>
                        </FormItem>
                        <FormItem className='flex items-center space-x-3 space-y-0'>
                          <FormControl>
                            <RadioGroupItem value='video' />
                          </FormControl>
                          <FormLabel className='font-normal'>{t('form.typeVideo')}</FormLabel>
                        </FormItem>
                      </RadioGroup>
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              {renderContentField()}

              <FormField
                control={form.control}
                name='description'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('form.description')}</FormLabel>
                    <FormControl>
                      <EnhancedInput
                        placeholder={t('form.enterDescription')}
                        value={field.value}
                        onValueChange={(value) => {
                          form.setValue(field.name, value);
                        }}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name='target_url'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('form.targetUrl')}</FormLabel>
                    <FormControl>
                      <EnhancedInput
                        placeholder={t('form.enterTargetUrl')}
                        value={field.value}
                        onValueChange={(value) => {
                          form.setValue(field.name, value);
                        }}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name='start_time'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('form.startTime')}</FormLabel>
                    <FormControl>
                      <EnhancedInput
                        type='datetime-local'
                        step='1'
                        placeholder={t('form.enterStartTime')}
                        value={field.value ? new Date(field.value).toISOString().slice(0, 16) : ''}
                        min={Number(new Date().toISOString().slice(0, 16))}
                        onValueChange={(value) => {
                          const timestamp = value ? new Date(value).getTime() : 0;
                          form.setValue(field.name, timestamp);
                          const endTime = form.getValues('end_time');
                          if (endTime && timestamp > endTime) {
                            form.setValue('end_time', '');
                          }
                        }}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name='end_time'
                render={({ field }) => {
                  return (
                    <FormItem>
                      <FormLabel>{t('form.endTime')}</FormLabel>
                      <FormControl>
                        <EnhancedInput
                          type='datetime-local'
                          step='1'
                          placeholder={t('form.enterEndTime')}
                          value={
                            field.value ? new Date(field.value).toISOString().slice(0, 16) : ''
                          }
                          min={Number(
                            startTime
                              ? new Date(startTime).toISOString().slice(0, 16)
                              : new Date().toISOString().slice(0, 16),
                          )}
                          disabled={!startTime}
                          onValueChange={(value) => {
                            const timestamp = value ? new Date(value).getTime() : 0;
                            if (!startTime || timestamp < startTime) return;
                            form.setValue(field.name, timestamp);
                          }}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  );
                }}
              />
            </form>
          </Form>
        </ScrollArea>
        <SheetFooter className='flex-row justify-end gap-2 pt-3'>
          <Button
            variant='outline'
            disabled={loading}
            onClick={() => {
              setOpen(false);
            }}
          >
            {t('form.cancel')}
          </Button>
          <Button disabled={loading} onClick={form.handleSubmit(handleSubmit)}>
            {loading && <Icon icon='mdi:loading' className='mr-2 animate-spin' />}
            {t('form.confirm')}
          </Button>
        </SheetFooter>
      </SheetContent>
    </Sheet>
  );
}
