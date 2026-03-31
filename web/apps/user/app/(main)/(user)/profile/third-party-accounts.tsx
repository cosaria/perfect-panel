'use client';

import SendCode from '@/app/auth/send-code';
import useGlobalStore from '@/config/use-global';
import { bindOAuth, unbindOAuth, updateBindEmail, updateBindMobile } from '@/services/user/user';
import { zodResolver } from '@hookform/resolvers/zod';
import { Button } from '@workspace/ui/components/button';
import { Card, CardContent, CardHeader, CardTitle } from '@workspace/ui/components/card';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@workspace/ui/components/dialog';
import { Form, FormControl, FormField, FormItem, FormMessage } from '@workspace/ui/components/form';
import { Input } from '@workspace/ui/components/input';
import { AreaCodeSelect } from '@workspace/ui/custom-components/area-code-select';
import { Icon } from '@workspace/ui/custom-components/icon';
import { useTranslations } from 'next-intl';
import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { toast } from 'sonner';
import { z } from 'zod';

function MobileBindDialog({
  onSuccess,
  children,
}: {
  onSuccess: () => void;
  children: React.ReactNode;
}) {
  const t = useTranslations('profile.thirdParty');
  const { common } = useGlobalStore();
  const { enable_whitelist, whitelist } = common.auth.mobile;
  const [open, setOpen] = useState(false);

  const formSchema = z.object({
    area_code: z.string().min(1, 'Area code is required'),
    mobile: z.string().min(5, 'Phone number is required'),
    code: z.string().min(4, 'Verification code is required'),
  });

  type MobileBindFormValues = z.infer<typeof formSchema>;

  const form = useForm<MobileBindFormValues>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      area_code: '1',
      mobile: '',
      code: '',
    },
  });

  const onSubmit = async (values: MobileBindFormValues) => {
    try {
      await updateBindMobile(values);
      toast.success(t('bindSuccess'));
      onSuccess();
      setOpen(false);
    } catch (error) {
      toast.error(t('bindFailed'));
    }
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>{children}</DialogTrigger>
      <DialogContent className='sm:max-w-[425px]'>
        <DialogHeader>
          <DialogTitle>{t('bindMobile')}</DialogTitle>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className='space-y-4'>
            <FormField
              control={form.control}
              name='mobile'
              render={({ field }) => (
                <FormItem>
                  <FormControl>
                    <div className='flex'>
                      <FormField
                        control={form.control}
                        name='area_code'
                        render={({ field }) => (
                          <FormItem>
                            <FormControl>
                              <AreaCodeSelect
                                simple
                                className='w-32 rounded-r-none border-r-0'
                                placeholder='Area code...'
                                value={field.value}
                                whitelist={enable_whitelist ? whitelist : []}
                                onChange={(value) => {
                                  if (value.phone) {
                                    form.setValue(field.name, value.phone);
                                  }
                                }}
                              />
                            </FormControl>
                            <FormMessage />
                          </FormItem>
                        )}
                      />
                      <Input
                        className='rounded-l-none'
                        placeholder='Enter your telephone...'
                        type='tel'
                        {...field}
                      />
                    </div>
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name='code'
              render={({ field }) => (
                <FormItem>
                  <FormControl>
                    <div className='flex gap-2'>
                      <Input placeholder='Enter code...' type='text' {...field} />
                      <SendCode
                        type='phone'
                        params={{
                          telephone_area_code: form.getValues().area_code,
                          telephone: form.getValues().mobile,
                          type: 1,
                        }}
                      />
                    </div>
                  </FormControl>

                  <FormMessage />
                </FormItem>
              )}
            />
            <Button type='submit' className='w-full'>
              {t('confirm')}
            </Button>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}

export default function ThirdPartyAccounts() {
  const t = useTranslations('profile.thirdParty');
  const { user, getUserInfo, common } = useGlobalStore();
  const { oauth_methods } = common;

  const accounts = [
    {
      id: 'email',
      icon: 'logos:mailgun-icon',
      name: 'Email',
      type: 'Basic',
    },
    {
      id: 'mobile',
      icon: 'mdi:telephone',
      name: 'Mobile',
      type: 'Basic',
    },
    {
      id: 'telegram',
      icon: 'logos:telegram',
      name: 'Telegram',
      type: 'OAuth',
    },
    {
      id: 'apple',
      icon: 'uil:apple',
      name: 'Apple',
      type: 'OAuth',
    },
    {
      id: 'google',
      icon: 'logos:google',
      name: 'Google',
      type: 'OAuth',
    },
    {
      id: 'facebook',
      icon: 'logos:facebook',
      name: 'Facebook',
      type: 'OAuth',
    },
    {
      id: 'github',
      icon: 'uil:github',
      name: 'GitHub',
      type: 'OAuth',
    },
    {
      id: 'device',
      icon: 'mdi:devices',
      name: 'Device',
      type: 'OAuth',
    },
  ].filter((account) => oauth_methods?.includes(account.id));

  const [editValues, setEditValues] = useState<Record<string, any>>({});

  const handleBasicAccountUpdate = async (account: (typeof accounts)[0], value: string) => {
    if (account.id === 'email') {
      await updateBindEmail({ email: value });
      await getUserInfo();
      toast.success(t('updateSuccess'));
    }
  };

  const handleAccountAction = async (account: (typeof accounts)[number]) => {
    const isBound = user?.auth_methods?.find(
      (auth) => auth.auth_type === account.id,
    )?.auth_identifier;
    if (isBound) {
      await unbindOAuth({ method: account.id });
      await getUserInfo();
    } else {
      const res = await bindOAuth({
        method: account.id,
        redirect: `${window.location.origin}/bind/${account.id}`,
      });
      if (res.data?.data?.redirect) {
        window.location.href = res.data.data.redirect;
      }
    }
  };

  return (
    <>
      <Card>
        <CardHeader className='bg-muted/50'>
          <CardTitle>{t('title')}</CardTitle>
        </CardHeader>
        <CardContent className='p-6'>
          <div className='space-y-4'>
            {accounts.map((account) => {
              const method = user?.auth_methods?.find((auth) => auth.auth_type === account.id);
              const isEditing = account.id === 'email';
              const currentValue = method?.auth_identifier || editValues[account.id];
              let displayValue = '';

              switch (account.id) {
                case 'email':
                  displayValue = isEditing ? currentValue : method?.auth_identifier || '';
                  break;
                default:
                  displayValue = method?.auth_identifier || t(`${account.id}.description`);
              }

              return (
                <div key={account.id} className='flex w-full flex-col gap-2'>
                  <span className='flex gap-3 font-medium'>
                    <Icon icon={account.icon} className='size-6' />
                    {account.name}
                  </span>
                  <div className='flex items-center gap-2'>
                    <Input
                      value={displayValue}
                      disabled={!isEditing}
                      className='bg-muted flex-1 truncate'
                      onChange={(e) =>
                        isEditing &&
                        setEditValues((prev) => ({ ...prev, [account.id]: e.target.value }))
                      }
                      onKeyDown={(e) => {
                        if (e.key === 'Enter' && isEditing) {
                          handleBasicAccountUpdate(account, currentValue);
                        }
                      }}
                    />
                    {account.id === 'mobile' ? (
                      <MobileBindDialog onSuccess={getUserInfo}>
                        <Button
                          variant={method?.auth_identifier ? 'outline' : 'default'}
                          className='whitespace-nowrap'
                        >
                          {t(method?.auth_identifier ? 'update' : 'bind')}
                        </Button>
                      </MobileBindDialog>
                    ) : (
                      <Button
                        variant={method?.auth_identifier ? 'outline' : 'default'}
                        onClick={() =>
                          isEditing
                            ? handleBasicAccountUpdate(account, currentValue)
                            : handleAccountAction(account)
                        }
                        className='whitespace-nowrap'
                      >
                        {t(isEditing ? 'save' : method?.auth_identifier ? 'unbind' : 'bind')}
                      </Button>
                    )}
                  </div>
                </div>
              );
            })}
          </div>
        </CardContent>
      </Card>
    </>
  );
}
