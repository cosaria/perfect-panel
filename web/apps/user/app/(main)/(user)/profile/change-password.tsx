'use client';

import { updateUserPassword } from '@/services/user/user';
import { zodResolver } from '@hookform/resolvers/zod';
import { Button } from '@workspace/ui/components/button';
import { Card, CardContent, CardHeader, CardTitle } from '@workspace/ui/components/card';
import { Form, FormControl, FormField, FormItem, FormMessage } from '@workspace/ui/components/form';
import { Input } from '@workspace/ui/components/input';
import { useTranslations } from 'next-intl';
import { useForm } from 'react-hook-form';
import { toast } from 'sonner';
import { z } from 'zod';

const FormSchema = z
  .object({
    password: z.string().min(6),
    repeat_password: z.string(),
  })
  .refine((data) => data.password === data.repeat_password, {
    message: 'passwordMismatch',
    path: ['repeat_password'],
  });

export default function ChangePassword() {
  const t = useTranslations('profile.accountSettings');
  const form = useForm<z.infer<typeof FormSchema>>({
    resolver: zodResolver(FormSchema),
  });

  async function onSubmit(data: z.infer<typeof FormSchema>) {
    await updateUserPassword({ password: data.password });
    toast.success(t('updateSuccess'));
    form.reset();
  }

  return (
    <Card className='min-w-80'>
      <CardHeader className='bg-muted/50'>
        <CardTitle className='flex items-center justify-between'>
          {t('accountSettings')}
          <Button type='submit' size='sm' form='password-form'>
            {t('updatePassword')}
          </Button>
        </CardTitle>
      </CardHeader>
      <CardContent className='p-6'>
        <Form {...form}>
          <form id='password-form' onSubmit={form.handleSubmit(onSubmit)} className='space-y-4'>
            <FormField
              control={form.control}
              name='password'
              render={({ field }) => (
                <FormItem>
                  <FormControl>
                    <Input type='password' placeholder={t('newPassword')} {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name='repeat_password'
              render={({ field }) => (
                <FormItem>
                  <FormControl>
                    <Input type='password' placeholder={t('repeatNewPassword')} {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
          </form>
        </Form>
      </CardContent>
    </Card>
  );
}
