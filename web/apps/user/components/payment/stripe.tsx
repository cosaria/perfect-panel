import {
  CardCvcElement,
  CardExpiryElement,
  CardNumberElement,
  Elements,
  useElements,
  useStripe,
} from '@stripe/react-stripe-js';
import {
  loadStripe,
  PaymentIntentResult,
  StripeCardNumberElementOptions,
  StripeElementStyle,
} from '@stripe/stripe-js';
import { Button } from '@workspace/ui/components/button';
import { Input } from '@workspace/ui/components/input';
import { Label } from '@workspace/ui/components/label';
import { CheckCircle } from 'lucide-react';
import { useTranslations } from 'next-intl';
import { useTheme } from 'next-themes';
import { QRCodeCanvas } from 'qrcode.react';
import React, { useCallback, useEffect, useMemo, useState } from 'react';

interface StripePaymentProps {
  method: string;
  client_secret: string;
  publishable_key: string;
}

interface CardPaymentFormProps {
  clientSecret: string;
  onError: (message: string) => void;
}

const CardPaymentForm: React.FC<CardPaymentFormProps> = ({ clientSecret, onError }) => {
  const stripe = useStripe();
  const { theme, systemTheme } = useTheme();
  const elements = useElements();
  const [processing, setProcessing] = useState(false);
  const [succeeded, setSucceeded] = useState(false);
  const [errors, setErrors] = useState<{
    cardNumber?: string;
    cardExpiry?: string;
    cardCvc?: string;
    name?: string;
  }>({});
  const [cardholderName, setCardholderName] = useState('');
  const t = useTranslations('payment.stripe.card');

  const currentTheme = theme === 'system' ? systemTheme : theme;
  const elementStyle: StripeElementStyle = {
    base: {
      'fontSize': '16px',
      'color': currentTheme === 'dark' ? '#fff' : '#000',
      '::placeholder': {
        color: '#aab7c4',
      },
    },
    invalid: {
      color: '#EF4444',
      iconColor: '#EF4444',
    },
  };

  const elementOptions: StripeCardNumberElementOptions = {
    style: elementStyle,
    showIcon: true,
  };

  const handleChange = (event: any, field: keyof typeof errors) => {
    if (event.error) {
      setErrors((prev) => ({ ...prev, [field]: event.error.message }));
    } else {
      setErrors((prev) => ({ ...prev, [field]: undefined }));
    }
  };

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();

    if (!stripe || !elements) {
      onError(t('loading'));
      return;
    }

    if (!cardholderName.trim()) {
      setErrors((prev) => ({ ...prev, name: t('name_required') }));
      return;
    }

    setProcessing(true);

    const cardNumber = elements.getElement(CardNumberElement);
    const cardExpiry = elements.getElement(CardExpiryElement);
    const cardCvc = elements.getElement(CardCvcElement);

    if (!cardNumber || !cardExpiry || !cardCvc) {
      onError(t('element_error'));
      setProcessing(false);
      return;
    }

    const { error, paymentIntent } = await stripe.confirmCardPayment(clientSecret, {
      payment_method: {
        card: cardNumber,
        billing_details: {
          name: cardholderName,
        },
      },
    });

    if (error) {
      onError(error.message || t('payment_failed'));
      setProcessing(false);
    } else if (paymentIntent && paymentIntent.status === 'succeeded') {
      setSucceeded(true);
      setProcessing(false);
    } else {
      onError(t('processing'));
      setProcessing(false);
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      {succeeded ? (
        <div className='py-6 text-center'>
          <div className='mb-4 flex justify-center'>
            <CheckCircle className='h-12 w-12 text-green-500' />
          </div>
          <p className='text-xl font-medium'>{t('success_title')}</p>
          <p className='text-muted-foreground mt-2'>{t('success_message')}</p>
        </div>
      ) : (
        <>
          <div className='space-y-4'>
            {/* Cardholder Name */}
            <div className='space-y-1'>
              <Label htmlFor='cardholderName' className='text-sm font-medium'>
                {t('card_name')}
              </Label>
              <Input
                id='cardholderName'
                type='text'
                value={cardholderName}
                onChange={(e) => setCardholderName(e.target.value)}
                placeholder={t('name_placeholder')}
                className={errors.name ? 'border-destructive' : ''}
              />
              {errors.name && <p className='text-destructive text-xs'>{errors.name}</p>}
            </div>

            {/* Card Number */}
            <div className='space-y-1'>
              <Label htmlFor='cardNumber' className='text-sm font-medium'>
                {t('card_number')}
              </Label>
              <div className='relative'>
                <div
                  className={`focus-within:border-primary focus-within:ring-primary rounded-md border p-3 focus-within:ring-1 ${errors.cardNumber ? 'border-red-500' : ''}`}
                >
                  <CardNumberElement
                    id='cardNumber'
                    options={elementOptions}
                    onChange={(e) => handleChange(e, 'cardNumber')}
                  />
                </div>
              </div>
              {errors.cardNumber && <p className='text-destructive text-xs'>{errors.cardNumber}</p>}
            </div>

            <div className='grid grid-cols-2 gap-4'>
              {/* Expiry Date */}
              <div className='space-y-1'>
                <Label htmlFor='cardExpiry' className='text-sm font-medium'>
                  {t('expiry_date')}
                </Label>
                <div
                  className={`focus-within:border-primary focus-within:ring-primary rounded-md border p-3 focus-within:ring-1 ${errors.cardExpiry ? 'border-red-500' : ''}`}
                >
                  <CardExpiryElement
                    id='cardExpiry'
                    options={{ style: elementStyle }}
                    onChange={(e) => handleChange(e, 'cardExpiry')}
                  />
                </div>
                {errors.cardExpiry && (
                  <p className='text-destructive text-xs'>{errors.cardExpiry}</p>
                )}
              </div>

              {/* Security Code */}
              <div className='space-y-1'>
                <Label htmlFor='cardCvc' className='text-sm font-medium'>
                  {t('security_code')}
                </Label>
                <div
                  className={`focus-within:border-primary focus-within:ring-primary rounded-md border p-3 focus-within:ring-1 ${errors.cardCvc ? 'border-red-500' : ''}`}
                >
                  <CardCvcElement
                    id='cardCvc'
                    options={{ style: elementStyle }}
                    onChange={(e) => handleChange(e, 'cardCvc')}
                  />
                </div>
                {errors.cardCvc && <p className='text-destructive text-xs'>{errors.cardCvc}</p>}
              </div>
            </div>
          </div>
          <div className='mt-6 flex flex-col space-y-4'>
            <Button type='submit' disabled={processing || !stripe || !elements} className='w-full'>
              {processing ? t('processing_button') : t('pay_button')}
            </Button>
            <p className='text-muted-foreground text-center text-xs'>{t('secure_notice')}</p>
          </div>
        </>
      )}
    </form>
  );
};

const StripePayment: React.FC<StripePaymentProps> = ({
  method,
  client_secret,
  publishable_key,
}) => {
  const stripePromise = useMemo(() => loadStripe(publishable_key), [publishable_key]);

  return (
    <Elements stripe={stripePromise}>
      <CheckoutForm method={method} client_secret={client_secret} />
    </Elements>
  );
};

const CheckoutForm: React.FC<Omit<StripePaymentProps, 'publishable_key'>> = ({
  client_secret,
  method,
}) => {
  const stripe = useStripe();
  const [errorMessage, setErrorMessage] = useState<string | null>(null);
  const [qrCodeUrl, setQrCodeUrl] = useState<string | null>(null);
  const [isSubmitted, setIsSubmitted] = useState(false);
  const t = useTranslations('payment.stripe');

  const handleError = useCallback((message: string) => {
    setErrorMessage(message);
    setIsSubmitted(false);
  }, []);

  const confirmPayment = useCallback(async (): Promise<PaymentIntentResult | null> => {
    if (!stripe) {
      handleError(t('card.loading'));
      return null;
    }

    if (method === 'alipay') {
      return await stripe.confirmAlipayPayment(
        client_secret,
        { return_url: window.location.href },
        { handleActions: false },
      );
    }
    if (method === 'wechat_pay') {
      return await stripe.confirmWechatPayPayment(
        client_secret,
        {
          payment_method_options: { wechat_pay: { client: 'web' } },
        },
        { handleActions: false },
      );
    }
    return null;
  }, [client_secret, method, stripe, handleError, t]);

  const autoSubmit = useCallback(async () => {
    if (isSubmitted || method === 'card') return;

    setIsSubmitted(true);

    try {
      const result = await confirmPayment();
      if (!result) return;

      const { error, paymentIntent } = result;
      if (error) return handleError(error.message!);

      if (paymentIntent?.status === 'requires_action') {
        const nextAction = paymentIntent.next_action as any;
        const qrUrl =
          method === 'alipay'
            ? nextAction?.alipay_handle_redirect?.url
            : nextAction?.wechat_pay_display_qr_code?.image_url_svg;

        setQrCodeUrl(qrUrl || null);
      }
    } catch (error) {
      handleError(t('error'));
    }
  }, [confirmPayment, isSubmitted, handleError, method, t]);

  useEffect(() => {
    autoSubmit();
  }, [autoSubmit]);

  return method === 'card' ? (
    <div className='min-w-80 text-left'>
      <CardPaymentForm clientSecret={client_secret} onError={handleError} />
    </div>
  ) : qrCodeUrl ? (
    <>
      <QRCodeCanvas
        value={qrCodeUrl}
        size={208}
        imageSettings={{
          src: `/payment/${method}.svg`,
          width: 24,
          height: 24,
          excavate: true,
        }}
      />
      <p className='text-muted-foreground mt-4 text-center'>{t(`qrcode.${method}`)}</p>
    </>
  ) : (
    errorMessage
  );
};

export default StripePayment;
