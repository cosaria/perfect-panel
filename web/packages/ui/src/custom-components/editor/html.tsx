'use client';

import {
  MonacoEditor,
  MonacoEditorProps,
} from '@workspace/ui/custom-components/editor/monaco-editor';
import { useEffect, useRef } from 'react';

export function HTMLEditor(props: MonacoEditorProps) {
  return (
    <MonacoEditor
      title='HTML Editor'
      description='Support HTML'
      {...props}
      language='markdown'
      render={(value) => <HTMLPreview value={value} />}
    />
  );
}

interface HTMLPreviewProps {
  value?: string;
}

function HTMLPreview({ value }: HTMLPreviewProps) {
  const iframeRef = useRef<HTMLIFrameElement>(null);

  useEffect(() => {
    const iframeDocument = iframeRef.current?.contentDocument;
    if (iframeDocument) {
      iframeDocument.open();
      iframeDocument.write(value || '');
      iframeDocument.close();
    }
  }, [value]);

  return <iframe ref={iframeRef} title='HTML Preview' className='h-full w-full border-0' />;
}
