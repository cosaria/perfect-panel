'use client';

import {
  MonacoEditor,
  MonacoEditorProps,
} from '@workspace/ui/custom-components/editor/monaco-editor';
import { Markdown } from '@workspace/ui/custom-components/markdown';

export function MarkdownEditor(props: MonacoEditorProps) {
  return (
    <MonacoEditor
      title='Markdown Editor'
      description='Support markdwon and html syntax'
      {...props}
      language='markdown'
      render={(value) => <Markdown>{value || ''}</Markdown>}
    />
  );
}
