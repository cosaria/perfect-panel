'use client';

import {
  telephoneLogin,
  telephoneResetPassword,
  telephoneUserRegister,
} from '@/services/common/auth';
import { getRedirectUrl, setAuthorization } from '@/utils/common';
import { useTranslations } from 'next-intl';
import { useRouter } from 'next/navigation';
import { ReactNode, useState, useTransition } from 'react';
import { toast } from 'sonner';
import LoginForm from './login-form';
import RegisterForm from './register-form';
import ResetForm from './reset-form';

export default function PhoneAuthForm() {
  const t = useTranslations('auth');
  const router = useRouter();
  const [type, setType] = useState<'login' | 'register' | 'reset'>('login');
  const [loading, startTransition] = useTransition();
  const [initialValues, setInitialValues] = useState<API.TelephoneLoginRequest>({
    identifier: '',
    telephone: '',
    telephone_area_code: '1',
    password: '',
    telephone_code: '',
  });

  const handleFormSubmit = async (params: any) => {
    const onLogin = async (token?: string) => {
      if (!token) return;
      setAuthorization(token);
      router.replace(getRedirectUrl());
      router.refresh();
    };
    startTransition(async () => {
      try {
        switch (type) {
          case 'login': {
            const login = await telephoneLogin(params);
            toast.success(t('login.success'));
            onLogin(login.data.data?.token);
            break;
          }
          case 'register': {
            const create = await telephoneUserRegister(params);
            toast.success(t('register.success'));
            onLogin(create.data.data?.token);
            break;
          }
          case 'reset':
            await telephoneResetPassword(params);
            toast.success(t('reset.success'));
            setType('login');
            break;
        }
      } catch (error) {
        /* empty */
      }
    });
  };

  let UserForm: ReactNode = null;
  switch (type) {
    case 'login':
      UserForm = (
        <LoginForm
          loading={loading}
          onSubmit={handleFormSubmit}
          initialValues={initialValues}
          setInitialValues={setInitialValues}
          onSwitchForm={setType}
        />
      );
      break;
    case 'register':
      UserForm = (
        <RegisterForm
          loading={loading}
          onSubmit={handleFormSubmit}
          initialValues={initialValues}
          setInitialValues={setInitialValues}
          onSwitchForm={setType}
        />
      );
      break;
    case 'reset':
      UserForm = (
        <ResetForm
          loading={loading}
          onSubmit={handleFormSubmit}
          initialValues={initialValues}
          setInitialValues={setInitialValues}
          onSwitchForm={setType}
        />
      );
      break;
  }

  return UserForm;
}
