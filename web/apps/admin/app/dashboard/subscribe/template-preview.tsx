'use client';

import { previewSubscribeTemplate } from '@/services/admin/application';
import { useQuery } from '@tanstack/react-query';
import { Button } from '@workspace/ui/components/button';
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from '@workspace/ui/components/sheet';
import { MonacoEditor } from '@workspace/ui/custom-components/editor/monaco-editor';
import { Icon } from '@workspace/ui/custom-components/icon';
import { useTranslations } from 'next-intl';
import { useState } from 'react';

interface TemplatePreviewProps {
  applicationId: number;
  output_format?: string;
}

export function TemplatePreview({ applicationId, output_format }: TemplatePreviewProps) {
  const t = useTranslations('subscribe.templatePreview');
  const [isOpen, setIsOpen] = useState(false);

  const { data, isLoading, error } = useQuery({
    queryKey: ['previewSubscribeTemplate', applicationId],
    queryFn: () => previewSubscribeTemplate({ id: applicationId }, { skipErrorHandler: true }),
    enabled: isOpen && !!applicationId,
    retry: false,
  });

  const originalContent = data?.data?.data?.template || '';
  const errorMessage = (error as any)?.data?.msg || error?.message || t('failed');

  const getDecodedContent = () => {
    if (output_format === 'base64' && originalContent) {
      try {
        return atob(originalContent);
      } catch {
        return t('base64.decodeError');
      }
    }
    return '';
  };

  const getDisplayContent = () => {
    if (error) return errorMessage;
    if (!originalContent) return '';
    switch (output_format) {
      case 'base64': {
        const decoded = getDecodedContent();
        return `${t('base64.originalContent')}:\n${originalContent}\n\n${t('base64.decodedContent')}:\n${decoded}`;
      }
      default:
        return originalContent;
    }
  };
  const mapLanguage = (fmt?: string) => {
    switch (fmt) {
      case 'json':
        return 'json';
      case 'yaml':
        return 'yaml';
      case 'base64':
        return 'ini';
      case 'plain':
        return 'ini';
      case 'conf':
        return 'ini';
      default:
        return 'ini';
    }
  };

  return (
    <Sheet open={isOpen} onOpenChange={setIsOpen}>
      <SheetTrigger asChild>
        <Button variant='ghost'>
          <Icon icon='mdi:eye' className='h-4 w-4' />
          {t('preview')}
        </Button>
      </SheetTrigger>
      <SheetHeader>
        <SheetTitle></SheetTitle>
      </SheetHeader>
      <SheetContent className='w-[800px] max-w-[90vw] pt-10 md:max-w-screen-md'>
        {isLoading ? (
          <div className='flex items-center justify-center'>
            <Icon icon='mdi:loading' className='h-6 w-6 animate-spin' />
            <span className='ml-2'>{t('loading')}</span>
          </div>
        ) : (
          <MonacoEditor
            title={t('title')}
            value={getDisplayContent()}
            language={mapLanguage(output_format)}
            showLineNumbers
            readOnly
          />
        )}
      </SheetContent>
    </Sheet>
  );
}
