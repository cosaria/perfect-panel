import useGlobalStore from '@/config/use-global';
import { zodResolver } from '@hookform/resolvers/zod';
import { Button } from '@workspace/ui/components/button';
import { Form, FormControl, FormField, FormItem, FormMessage } from '@workspace/ui/components/form';
import { Input } from '@workspace/ui/components/input';
import { Icon } from '@workspace/ui/custom-components/icon';
import { Markdown } from '@workspace/ui/custom-components/markdown';
import { useTranslations } from 'next-intl';
import { Dispatch, SetStateAction, useRef } from 'react';
import { useForm } from 'react-hook-form';
import { z } from 'zod';
import SendCode from '../send-code';
import CloudFlareTurnstile, { TurnstileRef } from '../turnstile';

export default function RegisterForm({
  loading,
  onSubmit,
  initialValues,
  setInitialValues,
  onSwitchForm,
}: {
  loading?: boolean;
  onSubmit: (data: any) => void;
  initialValues: any;
  setInitialValues: Dispatch<SetStateAction<any>>;
  onSwitchForm: Dispatch<SetStateAction<'register' | 'reset' | 'login'>>;
}) {
  const t = useTranslations('auth.register');
  const { common } = useGlobalStore();
  const { verify, auth, invite } = common;

  const handleCheckUser = async (email: string) => {
    try {
      if (!auth.email.enable_domain_suffix) return true;
      const domain = email.split('@')[1];
      const isValid = auth.email?.domain_suffix_list.split('\n').includes(domain || '');
      return isValid;
    } catch (error) {
      console.log('Error checking user:', error);
      return false;
    }
  };

  const formSchema = z
    .object({
      email: z
        .string()
        .email(t('email'))
        .refine(handleCheckUser, {
          message: t('whitelist'),
        }),
      password: z.string(),
      repeat_password: z.string(),
      code: auth.email.enable_verify ? z.string() : z.string().nullish(),
      invite: invite.forced_invite ? z.string().min(1) : z.string().nullish(),
      cf_token:
        verify.enable_register_verify && verify.turnstile_site_key
          ? z.string()
          : z.string().nullish(),
    })
    .superRefine(({ password, repeat_password }, ctx) => {
      if (password !== repeat_password) {
        ctx.addIssue({
          code: z.ZodIssueCode.custom,
          message: t('passwordMismatch'),
          path: ['repeat_password'],
        });
      }
    });

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      ...initialValues,
      invite: localStorage.getItem('invite') || '',
    },
  });

  const turnstile = useRef<TurnstileRef>(null);
  const handleSubmit = form.handleSubmit((data) => {
    try {
      onSubmit(data);
    } catch (error) {
      turnstile.current?.reset();
    }
  });

  return (
    <>
      {auth.register.stop_register ? (
        <Markdown>{t('message')}</Markdown>
      ) : (
        <Form {...form}>
          <form onSubmit={handleSubmit} className='grid gap-6'>
            <FormField
              control={form.control}
              name='email'
              render={({ field }) => (
                <FormItem>
                  <FormControl>
                    <Input placeholder='Enter your email...' type='email' {...field} />
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
                  <FormControl>
                    <Input placeholder='Enter your password...' type='password' {...field} />
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
                    <Input
                      disabled={loading}
                      placeholder='Enter password again...'
                      type='password'
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            {auth.email.enable_verify && (
              <FormField
                control={form.control}
                name='code'
                render={({ field }) => (
                  <FormItem>
                    <FormControl>
                      <div className='flex items-center gap-2'>
                        <Input
                          disabled={loading}
                          placeholder='Enter code...'
                          type='text'
                          {...field}
                          value={field.value as string}
                        />
                        <SendCode
                          type='email'
                          params={{
                            ...form.getValues(),
                            type: 1,
                          }}
                        />
                      </div>
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            )}
            <FormField
              control={form.control}
              name='invite'
              render={({ field }) => (
                <FormItem>
                  <FormControl>
                    <Input
                      disabled={loading || !!localStorage.getItem('invite')}
                      placeholder={t('invite')}
                      {...field}
                      value={field.value || ''}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            {verify.enable_register_verify && (
              <FormField
                control={form.control}
                name='cf_token'
                render={({ field }) => (
                  <FormItem>
                    <FormControl>
                      <CloudFlareTurnstile id='register' {...field} ref={turnstile} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            )}
            <Button type='submit' disabled={loading}>
              {loading && <Icon icon='mdi:loading' className='animate-spin' />}
              {t('title')}
            </Button>
          </form>
        </Form>
      )}
      <div className='mt-4 text-right text-sm'>
        {t('existingAccount')}&nbsp;
        <Button
          variant='link'
          className='p-0'
          onClick={() => {
            setInitialValues(undefined);
            onSwitchForm('login');
          }}
        >
          {t('switchToLogin')}
        </Button>
      </div>
    </>
  );
}
