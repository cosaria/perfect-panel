'use client';

import { createBatchSendEmailTask, getPreSendEmailCount } from '@/services/admin/marketing';
import { zodResolver } from '@hookform/resolvers/zod';
import { Button } from '@workspace/ui/components/button';
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@workspace/ui/components/form';
import { Input } from '@workspace/ui/components/input';
import { ScrollArea } from '@workspace/ui/components/scroll-area';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@workspace/ui/components/select';
import {
  Sheet,
  SheetContent,
  SheetFooter,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from '@workspace/ui/components/sheet';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@workspace/ui/components/tabs';
import { Textarea } from '@workspace/ui/components/textarea';
import { HTMLEditor } from '@workspace/ui/custom-components/editor';
import { EnhancedInput } from '@workspace/ui/custom-components/enhanced-input';
import { Icon } from '@workspace/ui/custom-components/icon';
import { useTranslations } from 'next-intl';
import { useEffect, useState } from 'react';
import { useForm } from 'react-hook-form';
import { toast } from 'sonner';
import { z } from 'zod';

export default function EmailBroadcastForm() {
  const t = useTranslations('marketing');

  // Define schema with internationalized error messages
  const emailBroadcastSchema = z.object({
    subject: z.string().min(1, t('subject') + ' ' + t('cannotBeEmpty')),
    content: z.string().min(1, t('content') + ' ' + t('cannotBeEmpty')),
    scope: z.number(),
    register_start_time: z.string().optional(),
    register_end_time: z.string().optional(),
    additional: z
      .string()
      .optional()
      .refine(
        (value) => {
          if (!value || value.trim() === '') return true;
          const emails = value.split('\n').filter((email) => email.trim() !== '');
          const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
          return emails.every((email) => emailRegex.test(email.trim()));
        },
        {
          message: t('pleaseEnterValidEmailAddresses'),
        },
      ),
    scheduled: z.string().optional(),
    interval: z.number().min(0.1, t('emailIntervalMinimum')).optional(),
    limit: z.number().min(1, t('dailyLimit')).optional(),
  });

  type EmailBroadcastFormData = z.infer<typeof emailBroadcastSchema>;
  const [open, setOpen] = useState(false);
  const [loading, setLoading] = useState(false);
  const [estimatedRecipients, setEstimatedRecipients] = useState<{
    users: number;
    additional: number;
    total: number;
  }>({ users: 0, additional: 0, total: 0 });

  const form = useForm<EmailBroadcastFormData>({
    resolver: zodResolver(emailBroadcastSchema),
    defaultValues: {
      subject: '',
      content: '',
      scope: 1, // ScopeAll
      register_start_time: '',
      register_end_time: '',
      additional: '',
      scheduled: '',
      interval: 1,
      limit: 1000,
    },
  });

  // Calculate recipient count
  const calculateRecipients = async () => {
    const formData = form.getValues();

    try {
      // Call API to get actual recipient count
      const scope = formData.scope || 1; // Default to ScopeAll

      // Convert dates to timestamps if they exist
      let register_start_time: number = 0;
      let register_end_time: number = 0;

      if (formData.register_start_time) {
        register_start_time = Math.floor(new Date(formData.register_start_time).getTime());
      }

      if (formData.register_end_time) {
        register_end_time = Math.floor(new Date(formData.register_end_time).getTime());
      }

      const response = await getPreSendEmailCount({
        scope,
        register_start_time,
        register_end_time,
      });

      const userCount = response.data?.data?.count || 0;

      // Calculate additional email count
      const additionalEmails = formData.additional || '';
      const additionalCount = additionalEmails
        .split('\n')
        .filter((email: string) => email.trim() !== '').length;

      const total = userCount + additionalCount;

      setEstimatedRecipients({
        users: userCount,
        additional: additionalCount,
        total,
      });
    } catch (error) {
      console.error('Failed to get recipient count:', error);
      // Set to 0 if API fails, don't use fallback simulation
      const additionalEmails = formData.additional || '';
      const additionalCount = additionalEmails
        .split('\n')
        .filter((email: string) => email.trim() !== '').length;

      setEstimatedRecipients({
        users: 0,
        additional: additionalCount,
        total: additionalCount,
      });
    }
  };

  // Listen to form changes
  const watchedValues = form.watch();

  // Use useEffect to respond to form changes, but only when sheet is open
  useEffect(() => {
    if (!open) return; // Only calculate when sheet is open

    const debounceTimer = setTimeout(() => {
      calculateRecipients();
    }, 500); // Add debounce to avoid too frequent API calls

    return () => clearTimeout(debounceTimer);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [
    open, // Add open dependency
    watchedValues.scope,
    watchedValues.register_start_time,
    watchedValues.register_end_time,
    watchedValues.additional,
  ]);

  const onSubmit = async (data: EmailBroadcastFormData) => {
    setLoading(true);
    try {
      // Validate scheduled send time
      let scheduled: number | undefined;
      if (data.scheduled && data.scheduled.trim() !== '') {
        const scheduledDate = new Date(data.scheduled);
        const now = new Date();
        if (scheduledDate <= now) {
          toast.error(t('scheduledSendTimeMustBeLater'));
          return;
        }
        scheduled = Math.floor(scheduledDate.getTime());
      }

      let register_start_time: number = 0;
      let register_end_time: number = 0;

      if (data.register_start_time) {
        register_start_time = Math.floor(new Date(data.register_start_time).getTime());
      }

      if (data.register_end_time) {
        register_end_time = Math.floor(new Date(data.register_end_time).getTime());
      }

      // Prepare API request data
      const requestData: API.CreateBatchSendEmailTaskRequest = {
        subject: data.subject,
        content: data.content,
        scope: data.scope,
        register_start_time,
        register_end_time,
        additional: data.additional || undefined,
        scheduled,
        interval: data.interval ? data.interval * 1000 : undefined, // Convert seconds to milliseconds
        limit: data.limit,
      };

      // Call API to create batch send email task
      await createBatchSendEmailTask(requestData);

      if (!data.scheduled || data.scheduled.trim() === '') {
        toast.success(t('emailBroadcastTaskCreatedSuccessfully'));
      } else {
        toast.success(t('emailAddedToScheduledQueue'));
      }

      form.reset();
      setOpen(false);
    } catch (error) {
      console.error('Email broadcast failed:', error);
      toast.error(t('sendFailed'));
    } finally {
      setLoading(false);
    }
  };

  return (
    <Sheet open={open} onOpenChange={setOpen}>
      <SheetTrigger asChild>
        <div className='flex cursor-pointer items-center justify-between transition-colors'>
          <div className='flex items-center gap-3'>
            <div className='bg-primary/10 flex h-10 w-10 items-center justify-center rounded-lg'>
              <Icon icon='mdi:email-send' className='text-primary h-5 w-5' />
            </div>
            <div className='flex-1'>
              <p className='font-medium'>{t('emailBroadcast')}</p>
              <p className='text-muted-foreground text-sm'>
                {t('createNewEmailBroadcastCampaign')}
              </p>
            </div>
          </div>
          <Icon icon='mdi:chevron-right' className='size-6' />
        </div>
      </SheetTrigger>
      <SheetContent className='w-[700px] max-w-full md:max-w-screen-lg'>
        <SheetHeader>
          <SheetTitle>{t('createBroadcast')}</SheetTitle>
        </SheetHeader>
        <ScrollArea className='-mx-6 h-[calc(100dvh-48px-36px-36px-env(safe-area-inset-top))] px-6'>
          <Form {...form}>
            <form
              id='broadcast-form'
              onSubmit={form.handleSubmit(onSubmit)}
              className='space-y-2 pt-4'
            >
              <Tabs defaultValue='content' className='space-y-2'>
                <TabsList className='grid w-full grid-cols-2'>
                  <TabsTrigger value='content'>{t('content')}</TabsTrigger>
                  <TabsTrigger value='settings'>{t('sendSettings')}</TabsTrigger>
                </TabsList>
                {/* Email Content Tab */}
                <TabsContent value='content' className='space-y-2'>
                  <FormField
                    control={form.control}
                    name='subject'
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>{t('subject')}</FormLabel>
                        <FormControl>
                          <Input
                            placeholder={`${t('pleaseEnter')} ${t('subject').toLowerCase()}`}
                            {...field}
                          />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={form.control}
                    name='content'
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>{t('content')}</FormLabel>
                        <FormControl>
                          <HTMLEditor
                            value={field.value}
                            onChange={(value) => {
                              form.setValue(field.name, value || '');
                            }}
                          />
                        </FormControl>
                        <FormDescription>{t('useMarkdownEditor')}</FormDescription>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </TabsContent>

                {/* Send Settings Tab */}
                <TabsContent value='settings' className='space-y-2'>
                  {/* Send scope and estimated recipients */}
                  <div className='grid grid-cols-2 items-center gap-4'>
                    <FormField
                      control={form.control}
                      name='scope'
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>{t('sendScope')}</FormLabel>
                          <Select
                            onValueChange={(value) => field.onChange(parseInt(value))}
                            value={field.value?.toString() || '1'}
                          >
                            <FormControl>
                              <SelectTrigger>
                                <SelectValue placeholder={t('selectSendScope')} />
                              </SelectTrigger>
                            </FormControl>
                            <SelectContent>
                              <SelectItem value='1'>{t('allUsers')}</SelectItem> {/* ScopeAll */}
                              <SelectItem value='2'>{t('subscribedUsersOnly')}</SelectItem>{' '}
                              {/* ScopeActive */}
                              <SelectItem value='3'>
                                {t('expiredSubscriptionUsersOnly')}
                              </SelectItem>{' '}
                              {/* ScopeExpired */}
                              <SelectItem value='4'>{t('noSubscriptionUsersOnly')}</SelectItem>{' '}
                              {/* ScopeNone */}
                              <SelectItem value='5'>{t('specificUsersOnly')}</SelectItem>{' '}
                              {/* ScopeSkip */}
                            </SelectContent>
                          </Select>
                          <FormDescription>{t('sendScopeDescription')}</FormDescription>
                        </FormItem>
                      )}
                    />

                    {/* Estimated recipients info */}
                    <div className='flex justify-end'>
                      <div className='border-l-primary bg-primary/10 border-l-4 px-4 py-3 text-sm'>
                        <span className='text-muted-foreground'>{t('estimatedRecipients')}: </span>
                        <span className='text-primary text-lg font-medium'>
                          {estimatedRecipients.total}
                        </span>
                        <span className='text-muted-foreground ml-2 text-xs'>
                          ({t('users')}: {estimatedRecipients.users}, {t('additional')}:{' '}
                          {estimatedRecipients.additional})
                        </span>
                      </div>
                    </div>
                  </div>

                  <div className='grid grid-cols-2 gap-4'>
                    <FormField
                      control={form.control}
                      name='register_start_time'
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>{t('registrationStartDate')}</FormLabel>
                          <FormControl>
                            <EnhancedInput
                              type='datetime-local'
                              step='1'
                              disabled={form.watch('scope') === 5} // ScopeSkip
                              value={field.value}
                              onValueChange={field.onChange}
                            />
                          </FormControl>
                          <FormDescription>{t('includeUsersRegisteredAfter')}</FormDescription>
                        </FormItem>
                      )}
                    />
                    <FormField
                      control={form.control}
                      name='register_end_time'
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>{t('registrationEndDate')}</FormLabel>
                          <FormControl>
                            <EnhancedInput
                              type='datetime-local'
                              step='1'
                              disabled={form.watch('scope') === 5} // ScopeSkip
                              value={field.value}
                              onValueChange={field.onChange}
                            />
                          </FormControl>
                          <FormDescription>{t('includeUsersRegisteredBefore')}</FormDescription>
                        </FormItem>
                      )}
                    />
                  </div>

                  {/* Additional recipients */}
                  <FormField
                    control={form.control}
                    name='additional'
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>{t('additionalRecipientEmails')}</FormLabel>
                        <FormControl>
                          <Textarea
                            placeholder={`${t('pleaseEnter')}${t('additionalRecipientEmails').toLowerCase()}，${t('onePerLine')}，for example:\nexample1@domain.com\nexample2@domain.com\nexample3@domain.com`}
                            className='min-h-[120px] font-mono text-sm'
                            {...field}
                          />
                        </FormControl>
                        <FormDescription>{t('additionalRecipientsDescription')}</FormDescription>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  {/* Send time settings */}
                  <FormField
                    control={form.control}
                    name='scheduled'
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>{t('scheduledSend')}</FormLabel>
                        <FormControl>
                          <EnhancedInput
                            type='datetime-local'
                            step='1'
                            placeholder={t('leaveEmptyForImmediateSend')}
                            value={field.value}
                            onValueChange={field.onChange}
                          />
                        </FormControl>
                        <FormDescription>{t('selectSendTime')}</FormDescription>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  {/* Send rate control */}
                  <div className='grid grid-cols-2 gap-4'>
                    <FormField
                      control={form.control}
                      name='interval'
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>{t('emailInterval')}</FormLabel>
                          <FormControl>
                            <Input
                              type='number'
                              min={1}
                              step={0.1}
                              placeholder='1'
                              {...field}
                              onChange={(e) => field.onChange(parseFloat(e.target.value) || 1)}
                            />
                          </FormControl>
                          <FormDescription>{t('intervalTimeBetweenEmails')}</FormDescription>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                    <FormField
                      control={form.control}
                      name='limit'
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>{t('dailySendLimit')}</FormLabel>
                          <FormControl>
                            <Input
                              type='number'
                              min={1}
                              step={1}
                              placeholder='1000'
                              {...field}
                              onChange={(e) => field.onChange(parseInt(e.target.value) || 1000)}
                            />
                          </FormControl>
                          <FormDescription>{t('maximumNumberPerDay')}</FormDescription>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                  </div>
                </TabsContent>
              </Tabs>
            </form>
          </Form>
        </ScrollArea>
        <SheetFooter className='flex flex-row items-center justify-end gap-2 pt-3'>
          <Button variant='outline' onClick={() => setOpen(false)}>
            {t('cancel')}
          </Button>
          <Button type='submit' form='broadcast-form' disabled={loading}>
            {loading && <Icon icon='mdi:loading' className='mr-2 h-4 w-4 animate-spin' />}
            {loading
              ? t('processing')
              : !form.watch('scheduled') || form.watch('scheduled')?.trim() === ''
                ? t('sendNow')
                : t('scheduleSend')}
          </Button>
        </SheetFooter>
      </SheetContent>
    </Sheet>
  );
}
