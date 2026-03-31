'use client';

import useGlobalStore from '@/config/use-global';
import { updateUserBasicInfo } from '@/services/admin/user';
import { zodResolver } from '@hookform/resolvers/zod';
import { Button } from '@workspace/ui/components/button';
import { Card, CardContent, CardHeader, CardTitle } from '@workspace/ui/components/card';
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@workspace/ui/components/form';
import { Input } from '@workspace/ui/components/input';
import { Switch } from '@workspace/ui/components/switch';
import { EnhancedInput } from '@workspace/ui/custom-components/enhanced-input';
import { UploadImage } from '@workspace/ui/custom-components/upload-image';
import { unitConversion } from '@workspace/ui/utils';
import { useTranslations } from 'next-intl';
import { useForm } from 'react-hook-form';
import { toast } from 'sonner';
import * as z from 'zod';

const basicInfoSchema = z.object({
  avatar: z.string().optional(),
  balance: z.number().optional(),
  commission: z.number().optional(),
  gift_amount: z.number().optional(),
  refer_code: z.string().optional(),
  referer_id: z.number().optional(),
  referral_percentage: z.number().optional(),
  only_first_purchase: z.boolean().optional(),
  is_admin: z.boolean().optional(),
  password: z.string().optional(),
  enable: z.boolean(),
});

type BasicInfoValues = z.infer<typeof basicInfoSchema>;

export function BasicInfoForm({ user, refetch }: { user: API.User; refetch: () => void }) {
  const t = useTranslations('user');

  const { common } = useGlobalStore();
  const { currency } = common;

  const form = useForm<BasicInfoValues>({
    resolver: zodResolver(basicInfoSchema),
    defaultValues: {
      avatar: user.avatar,
      balance: user.balance,
      commission: user.commission,
      gift_amount: user.gift_amount,
      refer_code: user.refer_code,
      referer_id: user.referer_id,
      referral_percentage: user.referral_percentage,
      only_first_purchase: user.only_first_purchase,
      is_admin: user.is_admin,
      enable: user.enable,
    },
  });

  async function onSubmit(data: BasicInfoValues) {
    await updateUserBasicInfo({
      user_id: user.id,
      telegram: user.telegram,
      ...data,
    } as API.UpdateUserBasiceInfoRequest);
    toast.success(t('updateSuccess'));
    refetch();
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)}>
        <Card>
          <CardHeader className='flex flex-row items-center justify-between'>
            <CardTitle>{t('basicInfoTitle')}</CardTitle>
            <Button type='submit' size='sm'>
              {t('save')}
            </Button>
          </CardHeader>
          <CardContent className='space-y-4'>
            <FormField
              control={form.control}
              name='enable'
              render={({ field }) => (
                <FormItem className='flex items-center justify-between space-x-2'>
                  <FormLabel>{t('accountEnable')}</FormLabel>
                  <FormControl>
                    <Switch checked={field.value} onCheckedChange={field.onChange} />
                  </FormControl>
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name='is_admin'
              render={({ field }) => (
                <FormItem className='flex items-center justify-between space-x-2'>
                  <FormLabel>{t('administrator')}</FormLabel>
                  <FormControl>
                    <Switch checked={field.value} onCheckedChange={field.onChange} />
                  </FormControl>
                </FormItem>
              )}
            />
            <div className='grid grid-cols-3 gap-4'>
              <FormField
                control={form.control}
                name='balance'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('balance')}</FormLabel>
                    <FormControl>
                      <EnhancedInput
                        type='number'
                        min={0}
                        value={field.value}
                        prefix={currency?.currency_symbol ?? '$'}
                        formatInput={(value) => unitConversion('centsToDollars', value)}
                        formatOutput={(value) => unitConversion('dollarsToCents', value)}
                        onValueChange={(value) => {
                          form.setValue(field.name, value as number);
                        }}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name='commission'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('commission')}</FormLabel>
                    <FormControl>
                      <EnhancedInput
                        type='number'
                        min={0}
                        value={field.value}
                        prefix={currency?.currency_symbol ?? '$'}
                        formatInput={(value) => unitConversion('centsToDollars', value)}
                        formatOutput={(value) => unitConversion('dollarsToCents', value)}
                        onValueChange={(value) => {
                          form.setValue(field.name, value as number);
                        }}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name='gift_amount'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('giftAmount')}</FormLabel>
                    <FormControl>
                      <EnhancedInput
                        type='number'
                        min={0}
                        value={field.value}
                        prefix={currency?.currency_symbol ?? '$'}
                        formatInput={(value) => unitConversion('centsToDollars', value)}
                        formatOutput={(value) => unitConversion('dollarsToCents', value)}
                        onValueChange={(value) => {
                          form.setValue(field.name, value as number);
                        }}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            <div className='grid grid-cols-2 gap-4'>
              <FormField
                control={form.control}
                name='refer_code'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('referralCode')}</FormLabel>
                    <FormControl>
                      <EnhancedInput
                        value={field.value}
                        onValueChange={(value) => {
                          form.setValue(field.name, value as string);
                        }}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name='referer_id'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('referrerUserId')}</FormLabel>
                    <FormControl>
                      <EnhancedInput
                        type='number'
                        value={field.value}
                        onValueChange={(value) => {
                          form.setValue(field.name, value as number);
                        }}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            <FormField
              control={form.control}
              name='referral_percentage'
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t('referralPercentage')}</FormLabel>
                  <FormControl>
                    <EnhancedInput
                      type='number'
                      min={0}
                      max={100}
                      value={field.value}
                      suffix='%'
                      onValueChange={(value) => {
                        form.setValue(field.name, value as number);
                      }}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name='only_first_purchase'
              render={({ field }) => (
                <FormItem className='flex items-center justify-between space-x-2'>
                  <FormLabel>{t('onlyFirstPurchase')}</FormLabel>
                  <FormControl>
                    <Switch checked={field.value} onCheckedChange={field.onChange} />
                  </FormControl>
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name='avatar'
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t('avatar')}</FormLabel>
                  <FormControl>
                    <EnhancedInput
                      value={field.value}
                      onValueChange={(value) => {
                        form.setValue(field.name, value as string);
                      }}
                      suffix={
                        <UploadImage
                          className='bg-muted h-9 rounded-none border-none px-2'
                          returnType='base64'
                          onChange={(value) => form.setValue('avatar', value as string)}
                        />
                      }
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name='password'
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t('password')}</FormLabel>
                  <FormControl>
                    <Input type='password' placeholder={t('passwordPlaceholder')} {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
          </CardContent>
        </Card>
      </form>
    </Form>
  );
}
