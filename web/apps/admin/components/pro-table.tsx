'use client';

import { ProTable as _ProTable, ProTableProps } from '@workspace/ui/custom-components/pro-table';
import { useTranslations } from 'next-intl';
import { Empty } from './empty';
export { type ProTableActions } from '@workspace/ui/custom-components/pro-table';

export function ProTable<
  TData extends Record<string, unknown>,
  TValue extends Record<string, unknown>,
>(props: ProTableProps<TData, TValue>) {
  const t = useTranslations('common.table');
  return (
    <_ProTable
      {...props}
      texts={{
        actions: t('actions'),
        asc: t('asc'),
        desc: t('desc'),
        hide: t('hide'),
        selectedRowsText(total) {
          return t('selectedItems', { total });
        },
        textPageOf(current, total) {
          return t('pageInfo', { current, total });
        },
        textRowsPerPage: t('rowsPerPage'),
      }}
      empty={<Empty />}
    />
  );
}
