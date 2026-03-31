'use client';

import {
  ColumnDef,
  ColumnFiltersState,
  flexRender,
  getCoreRowModel,
  getFilteredRowModel,
  getPaginationRowModel,
  getSortedRowModel,
  SortingState,
  useReactTable,
  VisibilityState,
} from '@tanstack/react-table';
import { Alert, AlertDescription, AlertTitle } from '@workspace/ui/components/alert';
import { Button } from '@workspace/ui/components/button';
import { Checkbox } from '@workspace/ui/components/checkbox';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@workspace/ui/components/table';
import Empty from '@workspace/ui/custom-components/empty';
import { ColumnFilter, IParams } from '@workspace/ui/custom-components/pro-table/column-filter';
import { ColumnHeader } from '@workspace/ui/custom-components/pro-table/column-header';
import { ColumnToggle } from '@workspace/ui/custom-components/pro-table/column-toggle';
import { Pagination } from '@workspace/ui/custom-components/pro-table/pagination';
import { SortableRow } from '@workspace/ui/custom-components/pro-table/sortable-row';
import { ProTableWrapper } from '@workspace/ui/custom-components/pro-table/wrapper';
import { cn } from '@workspace/ui/lib/utils';
import { useSize } from 'ahooks';
import { GripVertical, ListRestart, Loader, RefreshCcw } from 'lucide-react';
import React, { Fragment, useEffect, useImperativeHandle, useRef, useState } from 'react';

export interface ProTableProps<TData, TValue> {
  columns: ColumnDef<TData, TValue>[];
  request: (
    pagination: {
      page: number;
      size: number;
    },
    filter: TValue,
  ) => Promise<{ list: TData[]; total: number }>;
  params?: IParams[];
  header?: {
    title?: React.ReactNode;
    toolbar?: React.ReactNode | React.ReactNode[];
    hidden?: boolean;
  };
  actions?: {
    render?: (row: TData) => React.ReactNode[];
    batchRender?: (rows: TData[]) => React.ReactNode[];
  };
  action?: React.Ref<ProTableActions | undefined>;
  texts?: Partial<{
    actions: string;
    asc: string;
    desc: string;
    hide: string;
    textRowsPerPage: string;
    textPageOf: (current: number, total: number) => string;
    selectedRowsText: (total: number) => string;
  }>;
  empty?: React.ReactNode;
  onSort?: (
    sourceId: string | number,
    targetId: string | number | null,
    items: TData[],
  ) => Promise<TData[]>;
  initialFilters?: Record<string, unknown>;
}

export interface ProTableActions {
  refresh: () => void;
  reset: () => void;
}

export function ProTable<
  TData extends Record<string, unknown> & { id?: string },
  TValue extends Record<string, unknown>,
>({
  columns,
  request,
  params,
  header,
  actions,
  action,
  texts,
  empty,
  onSort,
  initialFilters,
}: ProTableProps<TData, TValue>) {
  const [sorting, setSorting] = useState<SortingState>([]);
  const [columnFilters, setColumnFilters] = useState<ColumnFiltersState>(() => {
    if (initialFilters) {
      return Object.entries(initialFilters).map(([id, value]) => ({
        id,
        value,
      })) as ColumnFiltersState;
    }
    return [];
  });
  const [columnVisibility, setColumnVisibility] = useState<VisibilityState>({});
  const [rowSelection, setRowSelection] = useState({});
  const [data, setData] = useState<TData[]>([]);
  const [rowCount, setRowCount] = useState<number>(0);
  const [pagination, setPagination] = useState({
    pageIndex: 0,
    pageSize: 10,
  });
  const loading = useRef(false);

  const table = useReactTable({
    data,
    columns: [
      ...(onSort
        ? [
            {
              id: 'sortable',
              header: (
                <GripVertical className='h-4 w-4 cursor-move text-gray-500 hover:text-gray-700' />
              ),
              enableSorting: false,
              enableHiding: false,
            },
          ]
        : []),
      ...(actions?.batchRender ? [createSelectColumn<TData, TValue>()] : []),
      ...columns.map(
        (column) =>
          ({
            enableSorting: false,
            ...column,
          }) as ColumnDef<TData, TValue>,
      ),
      ...(actions?.render
        ? ([
            {
              id: 'actions',
              header: texts?.actions,
              cell: ({ row }) => (
                <div className='flex items-center justify-end gap-2'>
                  {actions?.render?.(row.original).map((item, index) => (
                    <Fragment key={index}>{item}</Fragment>
                  ))}
                </div>
              ),
              enableSorting: false,
              enableHiding: false,
            },
          ] as ColumnDef<TData, TValue>[])
        : []),
    ] as ColumnDef<TData, TValue>[],
    onPaginationChange: setPagination,
    onSortingChange: setSorting,
    onColumnFiltersChange: setColumnFilters,
    getCoreRowModel: getCoreRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    getSortedRowModel: getSortedRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    onColumnVisibilityChange: setColumnVisibility,
    onRowSelectionChange: setRowSelection,
    state: {
      sorting,
      columnFilters,
      columnVisibility,
      rowSelection,
      pagination,
    },
    manualPagination: true,
    manualFiltering: true,
    rowCount: rowCount,
    manualSorting: true,
  });

  const fetchData = async () => {
    if (loading.current) return;
    loading.current = true;
    try {
      const response = await request(
        {
          page: pagination.pageIndex + 1,
          size: pagination.pageSize,
        },
        Object.fromEntries(columnFilters.map((item) => [item.id, item.value])) as TValue,
      );
      setData(response.list);
      setRowCount(response.total);
    } catch (error) {
      console.log('Fetch data error:', error);
    } finally {
      loading.current = false;
    }
  };
  const reset = async () => {
    table.resetSorting();
    table.resetColumnFilters();
    table.resetGlobalFilter(true);
    table.resetColumnVisibility();
    table.resetRowSelection();
    table.resetPagination();
  };
  const ref = useRef<HTMLDivElement>(null);
  const size = useSize(ref);

  useImperativeHandle(action, () => ({
    refresh: fetchData,
    reset,
  }));

  useEffect(() => {
    fetchData();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [pagination.pageIndex, pagination.pageSize, JSON.stringify(columnFilters)]);

  const selectedRows = table.getSelectedRowModel().flatRows.map((row) => row.original);
  const selectedCount = selectedRows.length;

  return (
    <div className='flex flex-col gap-4' ref={ref}>
      {!header?.hidden && (
        <div className='flex flex-wrap-reverse items-center justify-between gap-4'>
          <div>
            {params ? (
              <ColumnFilter
                table={table}
                params={params}
                filters={Object.fromEntries(columnFilters.map((item) => [item.id, item.value]))}
              />
            ) : (
              header?.title
            )}
          </div>
          <div className='flex flex-1 items-center justify-end gap-2'>
            <Button variant='outline' size='icon' onClick={fetchData}>
              <RefreshCcw />
            </Button>
            <ColumnToggle table={table} />
            <Button variant='outline' size='icon' onClick={reset}>
              <ListRestart />
            </Button>
            {header?.toolbar}
          </div>
        </div>
      )}

      {selectedCount > 0 && actions?.batchRender && (
        <Alert className='flex items-center justify-between'>
          <AlertTitle className='m-0'>
            {texts?.selectedRowsText?.(selectedCount) || `Selected ${selectedCount} rows`}
          </AlertTitle>
          <AlertDescription className='flex gap-2'>
            {actions.batchRender(selectedRows)}
          </AlertDescription>
        </Alert>
      )}

      <div
        className='relative w-auto overflow-x-auto rounded-md border'
        style={{
          width: size?.width,
        }}
      >
        <ProTableWrapper data={data} setData={setData} onSort={onSort}>
          <Table className='w-full'>
            <TableHeader>
              {table.getHeaderGroups().map((headerGroup) => (
                <TableRow key={headerGroup.id}>
                  {headerGroup.headers.map((header) => (
                    <TableHead
                      key={header.id}
                      className={cn('!z-auto', getTableHeaderClass(header.column.id))}
                    >
                      <ColumnHeader
                        header={header}
                        text={{
                          asc: texts?.asc,
                          desc: texts?.desc,
                          hide: texts?.hide,
                        }}
                      />
                    </TableHead>
                  ))}
                </TableRow>
              ))}
            </TableHeader>
            <TableBody>
              {table.getRowModel()?.rows?.length ? (
                onSort ? (
                  table.getRowModel().rows.map((row) => (
                    <SortableRow
                      key={row.original.id ? String(row.original.id) : String(row.index)}
                      id={row.original.id ? String(row.original.id) : String(row.index)}
                      data-state={row.getIsSelected() && 'selected'}
                      isSortable
                    >
                      {row
                        .getVisibleCells()
                        .filter((cell) => {
                          return cell.column.id !== 'sortable';
                        })
                        .map((cell) => (
                          <TableCell key={cell.id} className={getTableCellClass(cell.column.id)}>
                            {flexRender(cell.column.columnDef.cell, cell.getContext())}
                          </TableCell>
                        ))}
                    </SortableRow>
                  ))
                ) : (
                  table.getRowModel().rows.map((row) => (
                    <TableRow key={row.id} data-state={row.getIsSelected() && 'selected'}>
                      {row.getVisibleCells().map((cell) => (
                        <TableCell key={cell.id} className={getTableCellClass(cell.column.id)}>
                          {flexRender(cell.column.columnDef.cell, cell.getContext())}
                        </TableCell>
                      ))}
                    </TableRow>
                  ))
                )
              ) : (
                <TableRow>
                  <TableCell colSpan={columns.length + 2} className='py-24'>
                    {empty || <Empty />}
                  </TableCell>
                </TableRow>
              )}
            </TableBody>
          </Table>
        </ProTableWrapper>

        {loading.current && (
          <div className='bg-muted/80 absolute top-0 z-20 flex h-full w-full items-center justify-center'>
            <Loader className='h-4 w-4 animate-spin' />
          </div>
        )}
      </div>
      {rowCount > 0 && (
        <Pagination
          table={table}
          text={{
            textRowsPerPage: texts?.textRowsPerPage,
            textPageOf: texts?.textPageOf,
          }}
        />
      )}
    </div>
  );
}

function createSelectColumn<TData, TValue>(): ColumnDef<TData, TValue> {
  return {
    id: 'selected',
    header: ({ table }) => (
      <Checkbox
        checked={
          table.getIsAllPageRowsSelected() || (table.getIsSomePageRowsSelected() && 'indeterminate')
        }
        onCheckedChange={(value) => table.toggleAllPageRowsSelected(!!value)}
        aria-label='Select all'
      />
    ),
    cell: ({ row }) => (
      <Checkbox
        checked={row.getIsSelected()}
        onCheckedChange={(value) => row.toggleSelected(!!value)}
        aria-label='Select row'
      />
    ),
    enableSorting: false,
    enableHiding: false,
  };
}

function getTableHeaderClass(columnId: string) {
  if (['sortable', 'selected'].includes(columnId)) {
    return 'sticky left-0 z-10 bg-background shadow-[2px_0_5px_-2px_rgba(0,0,0,0.1)] [&:has([role=checkbox])]:pr-2';
  } else if (columnId === 'actions') {
    return 'sticky right-0 z-10 text-right bg-background shadow-[-2px_0_5px_-2px_rgba(0,0,0,0.1)]';
  }
  return 'truncate';
}

function getTableCellClass(columnId: string) {
  if (['sortable', 'selected'].includes(columnId)) {
    return 'sticky left-0 bg-background shadow-[2px_0_5px_-2px_rgba(0,0,0,0.1)]';
  } else if (columnId === 'actions') {
    return 'sticky right-0 bg-background shadow-[-2px_0_5px_-2px_rgba(0,0,0,0.1)]';
  }
  return 'truncate';
}
