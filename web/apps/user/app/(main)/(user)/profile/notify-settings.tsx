'use client';

import useGlobalStore from '@/config/use-global';
import { updateUserNotify } from '@/services/user/user';
import { zodResolver } from '@hookform/resolvers/zod';
import { Button } from '@workspace/ui/components/button';
import { Card, CardContent, CardHeader, CardTitle } from '@workspace/ui/components/card';
import { Form, FormControl, FormField, FormItem, FormLabel } from '@workspace/ui/components/form';
import { Switch } from '@workspace/ui/components/switch';
import { useTranslations } from 'next-intl';
import { useForm } from 'react-hook-form';
import { toast } from 'sonner';
import { z } from 'zod';

const FormSchema = z.object({
  enable_balance_notify: z.boolean(),
  enable_login_notify: z.boolean(),
  enable_subscribe_notify: z.boolean(),
  enable_trade_notify: z.boolean(),
});

export default function NotifySettings() {
  const t = useTranslations('profile');
  const { user, getUserInfo } = useGlobalStore();
  const form = useForm<z.infer<typeof FormSchema>>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      enable_balance_notify: user?.enable_balance_notify ?? false,
      enable_login_notify: user?.enable_login_notify ?? false,
      enable_subscribe_notify: user?.enable_subscribe_notify ?? false,
      enable_trade_notify: user?.enable_trade_notify ?? false,
    },
  });

  async function onSubmit(data: z.infer<typeof FormSchema>) {
    await updateUserNotify(data);
    toast.success(t('notify.updateSuccess'));
    await getUserInfo();
  }

  return (
    <Card className='min-w-80'>
      <CardHeader className='bg-muted/50'>
        <CardTitle className='flex items-center justify-between'>
          {t('notify.notificationSettings')}
          <Button type='submit' size='sm' form='notify-form'>
            {t('notify.save')}
          </Button>
        </CardTitle>
      </CardHeader>
      <CardContent className='grid gap-6 p-6'>
        <Form {...form}>
          <form id='notify-form' onSubmit={form.handleSubmit(onSubmit)} className='space-y-4'>
            <div className='space-y-4'>
              {[
                { name: 'enable_balance_notify', label: 'balanceChange' },
                { name: 'enable_login_notify', label: 'login' },
                { name: 'enable_subscribe_notify', label: 'subscribe' },
                { name: 'enable_trade_notify', label: 'finance' },
              ].map(({ name, label }) => (
                <FormField
                  key={name}
                  control={form.control}
                  name={name as any}
                  render={({ field }) => (
                    <FormItem className='flex items-center justify-between space-x-4'>
                      <FormLabel className='text-muted-foreground'>
                        {t(`notify.${label}`)}
                      </FormLabel>
                      <FormControl>
                        <Switch checked={field.value} onCheckedChange={field.onChange} />
                      </FormControl>
                    </FormItem>
                  )}
                />
              ))}
            </div>
          </form>
        </Form>
      </CardContent>
    </Card>
  );
}
