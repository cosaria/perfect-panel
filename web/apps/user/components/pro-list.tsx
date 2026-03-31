'use client';

import { ProList as _ProList, ProListProps } from '@workspace/ui/custom-components/pro-list';
import { useTranslations } from 'next-intl';
export { type ProListActions, type ProListProps } from '@workspace/ui/custom-components/pro-list';

export function ProList<TData, TValue extends Record<string, unknown>>(
  props: ProListProps<TData, TValue>,
) {
  const t = useTranslations('common.list');
  return (
    <_ProList
      {...props}
      texts={{
        selectedRowsText(total) {
          return t('selectedItems', { total });
        },
        textPageOf(current, total) {
          return t('pageInfo', { current, total });
        },
        textRowsPerPage: t('rowsPerPage'),
      }}
    />
  );
}
