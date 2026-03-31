'use client';

import {
  createUserAuthMethod,
  deleteUserAuthMethod,
  updateUserAuthMethod,
} from '@/services/admin/user';
import { Badge } from '@workspace/ui/components/badge';
import { Button } from '@workspace/ui/components/button';
import { Card, CardContent, CardHeader, CardTitle } from '@workspace/ui/components/card';
import { EnhancedInput } from '@workspace/ui/custom-components/enhanced-input';
import { useTranslations } from 'next-intl';
import { useState } from 'react';
import { toast } from 'sonner';

export function AuthMethodsForm({ user, refetch }: { user: API.User; refetch: () => void }) {
  const t = useTranslations('user');

  const [emailChanges, setEmailChanges] = useState<Record<string, string>>({});

  const handleRemoveAuth = async (authType: string) => {
    await deleteUserAuthMethod({
      user_id: user.id,
      auth_type: authType,
    });
    toast.success(t('deleteSuccess'));
  };

  const handleUpdateEmail = async (email: string) => {
    await updateUserAuthMethod({
      user_id: user.id,
      auth_type: 'email',
      auth_identifier: email,
    });
    toast.success(t('updateSuccess'));
    refetch();
  };

  const handleCreateEmail = async (email: string) => {
    await createUserAuthMethod({
      user_id: user.id,
      auth_type: 'email',
      auth_identifier: email,
    });
    toast.success(t('createSuccess'));
    refetch();
  };

  const handleEmailChange = (authType: string, value: string) => {
    setEmailChanges((prev) => ({
      ...prev,
      [authType]: value,
    }));
  };

  const emailMethod = user.auth_methods.find((method) => method.auth_type === 'email');
  const otherMethods = user.auth_methods.filter((method) => method.auth_type !== 'email');

  const defaultEmailMethod = {
    auth_type: 'email',
    auth_identifier: '',
    verified: false,
    ...emailMethod,
  };

  const isEmailExists = !!emailMethod;
  const handleEmailAction = () => {
    const email = emailChanges['email'];
    if (isEmailExists) {
      handleUpdateEmail(email as string);
    } else {
      handleCreateEmail(email as string);
    }
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>{t('authMethodsTitle')}</CardTitle>
      </CardHeader>
      <CardContent className='space-y-6'>
        <div className='space-y-6'>
          <Card className='border-none shadow-none'>
            <CardContent className='space-y-3 p-0'>
              <div className='flex items-center justify-between'>
                <div className='font-medium uppercase'>email</div>
                <Badge variant={defaultEmailMethod.verified ? 'default' : 'destructive'}>
                  {defaultEmailMethod.verified ? t('verified') : t('unverified')}
                </Badge>
              </div>
              <div className='flex items-center gap-2'>
                <div className='flex-1'>
                  <EnhancedInput
                    value={defaultEmailMethod.auth_identifier}
                    placeholder={t('pleaseEnterEmail')}
                    onValueChange={(value) => handleEmailChange('email', value as string)}
                  />
                </div>
                <Button
                  onClick={handleEmailAction}
                  disabled={
                    !emailChanges['email'] ||
                    (isEmailExists && emailChanges['email'] === defaultEmailMethod.auth_identifier)
                  }
                >
                  {isEmailExists ? t('update') : t('add')}
                </Button>
              </div>
            </CardContent>
          </Card>

          {otherMethods.map((method) => (
            <Card key={method.auth_type} className='border-none shadow-none'>
              <CardContent className='space-y-3 p-0'>
                <div className='flex items-center justify-between'>
                  <div className='font-medium uppercase'>{method.auth_type}</div>
                  <Badge variant={method.verified ? 'default' : 'destructive'}>
                    {method.verified ? t('verified') : t('unverified')}
                  </Badge>
                </div>
                <div className='flex items-center gap-4'>
                  <div className='flex-1'>
                    <div className='text-muted-foreground text-sm'>{method.auth_identifier}</div>
                  </div>
                  <Button
                    variant='destructive'
                    size='sm'
                    onClick={() => handleRemoveAuth(method.auth_type)}
                  >
                    {t('remove')}
                  </Button>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      </CardContent>
    </Card>
  );
}
