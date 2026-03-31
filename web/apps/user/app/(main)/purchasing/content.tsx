'use client';

import { SubscribeBilling } from '@/components/subscribe/billing';
import CouponInput from '@/components/subscribe/coupon-input';
import { SubscribeDetail } from '@/components/subscribe/detail';
import DurationSelector from '@/components/subscribe/duration-selector';
import PaymentMethods from '@/components/subscribe/payment-methods';
import useGlobalStore from '@/config/use-global';
import { prePurchaseOrder, purchase } from '@/services/user/portal';
import { useQuery } from '@tanstack/react-query';
import { Button } from '@workspace/ui/components/button';
import { Card, CardContent, CardHeader } from '@workspace/ui/components/card';
import { Separator } from '@workspace/ui/components/separator';
import { EnhancedInput } from '@workspace/ui/custom-components/enhanced-input';
import { Icon } from '@workspace/ui/custom-components/icon';
import { cn } from '@workspace/ui/lib/utils';
import { LoaderCircle } from 'lucide-react';
import { useTranslations } from 'next-intl';
import { useRouter } from 'next/navigation';
import { useCallback, useEffect, useState, useTransition } from 'react';

export default function Content({ subscription }: { subscription?: API.Subscribe }) {
  const t = useTranslations('subscribe');
  const { common } = useGlobalStore();
  const router = useRouter();
  const [params, setParams] = useState<API.PortalPurchaseRequest>({
    quantity: 1,
    subscribe_id: 0,
    payment: -1,
    coupon: '',
    auth_type: 'email',
    identifier: '',
    password: '',
  });
  const [loading, startTransition] = useTransition();
  const [isEmailValid, setIsEmailValid] = useState({
    valid: false,
    message: '',
  });

  const { data: order } = useQuery({
    enabled: !!subscription?.id && !!params.payment,
    queryKey: ['preCreateOrder', params.coupon, params.quantity, params.payment],
    queryFn: async () => {
      const { data } = await prePurchaseOrder({
        ...params,
        subscribe_id: subscription?.id as number,
      } as API.PrePurchaseOrderRequest);
      return data.data;
    },
  });

  useEffect(() => {
    if (subscription) {
      setParams((prev) => ({
        ...prev,
        quantity: 1,
        subscribe_id: subscription?.id,
      }));
    }
  }, [subscription]);

  const handleChange = useCallback((field: keyof typeof params, value: string | number) => {
    setParams((prev) => ({
      ...prev,
      [field]: value,
    }));
  }, []);

  const handleSubmit = useCallback(async () => {
    startTransition(async () => {
      try {
        const { data } = await purchase(params);
        const { order_no } = data.data!;
        if (order_no) {
          localStorage.setItem(
            order_no,
            JSON.stringify({
              auth_type: params.auth_type,
              identifier: params.identifier,
            }),
          );
          router.push(`/purchasing/order?order_no=${order_no}`);
        }
      } catch (error) {
        console.log(error);
      }
    });
  }, [params, router]);

  if (!subscription) {
    return <div className='p-6 text-center'>{t('subscriptionNotFound')}</div>;
  }

  return (
    <div className='mx-auto mt-8 flex max-w-4xl flex-col gap-8 md:grid md:grid-cols-2 md:flex-row'>
      <div className='flex flex-col gap-6'>
        <Card>
          <CardHeader>输入要用于 {common.site.site_name} 账户的电子邮件地址</CardHeader>
          <CardContent className='flex flex-col gap-2'>
            <div className='flex flex-col gap-2'>
              <EnhancedInput
                className={cn({
                  'border-destructive': !isEmailValid.valid && params.identifier !== '',
                })}
                placeholder='Email'
                type='email'
                value={params.identifier || ''}
                onValueChange={(value) => {
                  const email = value as string;
                  setParams((prev) => ({
                    ...prev,
                    identifier: email,
                  }));
                  const reg = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
                  if (!reg.test(email)) {
                    setIsEmailValid({
                      valid: false,
                      message: '请输入有效的邮箱地址',
                    });
                  } else if (common.auth.email.enable_domain_suffix) {
                    const domain = email.split('@')[1];
                    const isValid = common.auth.email?.domain_suffix_list
                      .split('\n')
                      .includes(domain || '');
                    if (!isValid) {
                      setIsEmailValid({
                        valid: false,
                        message: '邮箱域名不在白名单中',
                      });
                      return;
                    }
                  } else {
                    setIsEmailValid({
                      valid: true,
                      message: '',
                    });
                  }
                }}
                required
              />
              <p
                className={cn('text-muted-foreground text-xs', {
                  'text-destructive': !isEmailValid.valid && params.identifier !== '',
                })}
              >
                {isEmailValid.message || '请填写您的电子邮件地址。'}
              </p>
            </div>
            {params.identifier && isEmailValid.valid && (
              <div className='flex flex-col gap-2'>
                <EnhancedInput
                  placeholder='Password'
                  type='password'
                  value={params.password || ''}
                  onValueChange={(value) => handleChange('password', value)}
                />
                <p className='text-muted-foreground text-xs'>
                  如果您不填写密码，我们将会自动生成密码并发送到您的邮箱。
                </p>
              </div>
            )}
            {/* <div>
              <OAuthMethods />
            </div> */}
          </CardContent>
        </Card>
        <Card>
          <CardContent className='grid gap-3 p-6 text-sm'>
            <h2 className='text-xl font-semibold'>{subscription.name}</h2>
            <ul className='flex flex-grow flex-col gap-3'>
              {(() => {
                let parsedDescription;
                try {
                  parsedDescription = JSON.parse(subscription.description);
                } catch {
                  parsedDescription = { description: '', features: [] };
                }

                const { description, features } = parsedDescription;
                return (
                  <>
                    {description && <li className='text-muted-foreground'>{description}</li>}
                    {features?.map(
                      (
                        feature: {
                          icon: string;
                          label: string;
                          type: 'default' | 'success' | 'destructive';
                        },
                        index: number,
                      ) => (
                        <li
                          className={cn('flex items-center gap-1', {
                            'text-muted-foreground line-through': feature.type === 'destructive',
                          })}
                          key={index}
                        >
                          {feature.icon && (
                            <Icon
                              icon={feature.icon}
                              className={cn('text-primary size-5', {
                                'text-green-500': feature.type === 'success',
                                'text-destructive': feature.type === 'destructive',
                              })}
                            />
                          )}
                          {feature.label}
                        </li>
                      ),
                    )}
                  </>
                );
              })()}
            </ul>
            <SubscribeDetail
              subscribe={{
                ...subscription,
                quantity: params.quantity,
              }}
            />
            <Separator />
            <SubscribeBilling
              order={{
                ...order,
                quantity: params.quantity,
                unit_price: subscription?.unit_price,
              }}
            />
          </CardContent>
        </Card>
      </div>

      <div className='flex flex-col gap-6'>
        <Card>
          <CardContent className='p-6'>
            <div className='grid gap-6'>
              <DurationSelector
                quantity={params.quantity!}
                unitTime={subscription?.unit_time}
                discounts={subscription?.discount}
                onChange={(value) => handleChange('quantity', value)}
              />
              <CouponInput
                coupon={params.coupon}
                onChange={(value) => handleChange('coupon', value)}
              />
              <PaymentMethods
                balance={false}
                value={params.payment!}
                onChange={(value) => handleChange('payment', value)}
              />
            </div>
          </CardContent>
        </Card>

        <Button
          className='w-full'
          size='lg'
          disabled={!isEmailValid.valid || loading}
          onClick={handleSubmit}
        >
          {loading && <LoaderCircle className='mr-2 animate-spin' />}
          {t('buyNow')}
        </Button>
      </div>
    </div>
  );
}
