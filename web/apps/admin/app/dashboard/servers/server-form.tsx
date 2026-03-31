'use client';

import { useNode } from '@/store/node';
import { zodResolver } from '@hookform/resolvers/zod';
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from '@workspace/ui/components/accordion';
import { Badge } from '@workspace/ui/components/badge';
import { Button } from '@workspace/ui/components/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@workspace/ui/components/dropdown-menu';
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@workspace/ui/components/form';
import { ScrollArea } from '@workspace/ui/components/scroll-area';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@workspace/ui/components/select';
import {
  Sheet,
  SheetContent,
  SheetFooter,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from '@workspace/ui/components/sheet';
import { Switch } from '@workspace/ui/components/switch';
import { EnhancedInput } from '@workspace/ui/custom-components/enhanced-input';
import { Icon } from '@workspace/ui/custom-components/icon';
import { cn } from '@workspace/ui/lib/utils';
import { useTranslations } from 'next-intl';
import { useEffect, useState } from 'react';
import { useForm, useWatch } from 'react-hook-form';
import { toast } from 'sonner';
import {
  FieldConfig,
  formSchema,
  getLabel,
  getProtocolDefaultConfig,
  PROTOCOL_FIELDS,
  protocols as PROTOCOLS,
} from './form-schema';

function DynamicField({
  field,
  control,
  form,
  protocolIndex,
  protocolData,
  t,
}: {
  field: FieldConfig;
  control: any;
  form: any;
  protocolIndex: number;
  protocolData: any;
  t: (key: string) => string;
}) {
  const fieldName = `protocols.${protocolIndex}.${field.name}` as const;

  if (field.condition && !field.condition(protocolData, {})) {
    return null;
  }

  const commonProps = {
    control,
    name: fieldName,
  };

  switch (field.type) {
    case 'input':
      return (
        <FormField
          {...commonProps}
          render={({ field: fieldProps }) => (
            <FormItem>
              <FormLabel>{t(field.label)}</FormLabel>
              <FormControl>
                <EnhancedInput
                  {...fieldProps}
                  type='text'
                  placeholder={
                    field.placeholder
                      ? typeof field.placeholder === 'function'
                        ? field.placeholder(t, protocolData)
                        : field.placeholder
                      : undefined
                  }
                  onValueChange={(v) => fieldProps.onChange(v)}
                  suffix={
                    field.generate ? (
                      field.generate.functions && field.generate.functions.length > 0 ? (
                        <DropdownMenu>
                          <DropdownMenuTrigger asChild>
                            <Button type='button' variant='ghost' size='sm'>
                              <Icon icon='mdi:key' className='h-4 w-4' />
                            </Button>
                          </DropdownMenuTrigger>
                          <DropdownMenuContent align='end'>
                            {field.generate.functions.map((genFunc, idx) => (
                              <DropdownMenuItem
                                key={idx}
                                onClick={async () => {
                                  const result = await genFunc.function();
                                  if (typeof result === 'string') {
                                    fieldProps.onChange(result);
                                  } else if (field.generate!.updateFields) {
                                    Object.entries(field.generate!.updateFields).forEach(
                                      ([fieldName, resultKey]) => {
                                        const fullFieldName = `protocols.${protocolIndex}.${fieldName}`;
                                        form.setValue(fullFieldName, (result as any)[resultKey]);
                                      },
                                    );
                                  } else {
                                    if (result.privateKey) {
                                      fieldProps.onChange(result.privateKey);
                                    }
                                  }
                                }}
                              >
                                {typeof genFunc.label === 'function'
                                  ? genFunc.label(t, protocolData)
                                  : genFunc.label}
                              </DropdownMenuItem>
                            ))}
                          </DropdownMenuContent>
                        </DropdownMenu>
                      ) : field.generate.function ? (
                        <Button
                          type='button'
                          variant='ghost'
                          size='sm'
                          onClick={async () => {
                            const result = await field.generate!.function!();
                            if (typeof result === 'string') {
                              fieldProps.onChange(result);
                            } else if (field.generate!.updateFields) {
                              Object.entries(field.generate!.updateFields).forEach(
                                ([fieldName, resultKey]) => {
                                  const fullFieldName = `protocols.${protocolIndex}.${fieldName}`;
                                  form.setValue(fullFieldName, (result as any)[resultKey]);
                                },
                              );
                            } else {
                              if (result.privateKey) {
                                fieldProps.onChange(result.privateKey);
                              }
                            }
                          }}
                        >
                          <Icon icon='mdi:key' className='h-4 w-4' />
                        </Button>
                      ) : null
                    ) : (
                      field.suffix
                    )
                  }
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
      );

    case 'number':
      return (
        <FormField
          {...commonProps}
          render={({ field: fieldProps }) => (
            <FormItem>
              <FormLabel>{t(field.label)}</FormLabel>
              <FormControl>
                <EnhancedInput
                  {...fieldProps}
                  type='number'
                  min={field.min}
                  max={field.max}
                  step={field.step || 1}
                  suffix={field.suffix}
                  placeholder={
                    field.placeholder
                      ? typeof field.placeholder === 'function'
                        ? field.placeholder(t, protocolData)
                        : field.placeholder
                      : undefined
                  }
                  onValueChange={(v) => fieldProps.onChange(v)}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
      );

    case 'select':
      if (!field.options || field.options.length <= 1) {
        return null;
      }

      return (
        <FormField
          {...commonProps}
          render={({ field: fieldProps }) => (
            <FormItem>
              <FormLabel>{t(field.label)}</FormLabel>
              <FormControl>
                <Select
                  value={fieldProps.value ?? field.defaultValue}
                  onValueChange={(v) => fieldProps.onChange(v)}
                >
                  <FormControl>
                    <SelectTrigger>
                      <SelectValue placeholder={t('please_select')} />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    {field.options?.map((option) => (
                      <SelectItem key={option} value={option}>
                        {getLabel(option)}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
      );

    case 'switch':
      return (
        <FormField
          {...commonProps}
          render={({ field: fieldProps }) => (
            <FormItem>
              <FormLabel>{t(field.label)}</FormLabel>
              <FormControl>
                <div className='pt-2'>
                  <Switch
                    checked={!!fieldProps.value}
                    onCheckedChange={(checked) => fieldProps.onChange(checked)}
                  />
                </div>
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
      );

    case 'textarea':
      return (
        <FormField
          {...commonProps}
          render={({ field: fieldProps }) => (
            <FormItem className='col-span-2'>
              <FormLabel>{t(field.label)}</FormLabel>
              <FormControl>
                <textarea
                  {...fieldProps}
                  value={fieldProps.value ?? ''}
                  className='border-input bg-background ring-offset-background placeholder:text-muted-foreground focus-visible:ring-ring flex min-h-[80px] w-full rounded-md border px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50'
                  placeholder={
                    field.placeholder
                      ? typeof field.placeholder === 'function'
                        ? field.placeholder(t, protocolData)
                        : field.placeholder
                      : undefined
                  }
                  onChange={(e) => fieldProps.onChange(e.target.value)}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
      );

    default:
      return null;
  }
}

function renderFieldsByGroup(
  fields: FieldConfig[],
  group: string,
  control: any,
  form: any,
  protocolIndex: number,
  protocolData: any,
  t: (key: string) => string,
) {
  const groupFields = fields.filter((field) => field.group === group);
  if (groupFields.length === 0) return null;

  return (
    <div className='grid grid-cols-2 gap-4'>
      {groupFields.map((field) => (
        <DynamicField
          key={field.name}
          field={field}
          control={control}
          form={form}
          protocolIndex={protocolIndex}
          protocolData={protocolData}
          t={t}
        />
      ))}
    </div>
  );
}

function renderGroupCard(
  title: string,
  fields: FieldConfig[],
  group: string,
  control: any,
  form: any,
  protocolIndex: number,
  protocolData: any,
  t: (key: string) => string,
) {
  const groupFields = fields.filter((field) => field.group === group);
  if (groupFields.length === 0) return null;

  const visibleFields = groupFields.filter(
    (field) => !field.condition || field.condition(protocolData, {}),
  );

  if (visibleFields.length === 0) return null;

  return (
    <div className='relative'>
      <fieldset className='border-border rounded-lg border'>
        <legend className='text-foreground bg-background ml-3 px-1 py-1 text-sm font-medium'>
          {t(title)}
        </legend>
        <div className='p-4 pt-2'>
          {renderFieldsByGroup(fields, group, control, form, protocolIndex, protocolData, t)}
        </div>
      </fieldset>
    </div>
  );
}

export default function ServerForm(props: {
  trigger: string;
  title: string;
  loading?: boolean;
  initialValues?: Partial<API.Server>;
  onSubmit: (values: Partial<API.Server>) => Promise<boolean> | boolean;
}) {
  const { trigger, title, loading, initialValues, onSubmit } = props;
  const t = useTranslations('servers');
  const [open, setOpen] = useState(false);
  const [accordionValue, setAccordionValue] = useState<string>();

  const { isProtocolUsedInNodes } = useNode();

  const form = useForm({
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: '',
      address: '',
      country: '',
      city: '',
      protocols: [] as any[],
      ...initialValues,
    },
  });
  const { control } = form;

  const protocolsValues = useWatch({ control, name: 'protocols' });

  useEffect(() => {
    if (initialValues) {
      form.reset({
        name: '',
        address: '',
        country: '',
        city: '',
        ...initialValues,
        protocols: PROTOCOLS.map((type) => {
          const existingProtocol = initialValues.protocols?.find((p) => p.type === type);
          const defaultConfig = getProtocolDefaultConfig(type);
          return existingProtocol ? { ...defaultConfig, ...existingProtocol } : defaultConfig;
        }),
      });
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [initialValues]);

  async function handleSubmit(values: Record<string, any>) {
    const filteredProtocols = (values?.protocols || []).filter((protocol: any) => {
      const port = Number(protocol?.port);
      return protocol && Number.isFinite(port) && port > 0 && port <= 65535;
    });

    const result = {
      name: values.name,
      country: values.country,
      city: values.city,
      address: values.address,
      protocols: filteredProtocols,
    };

    const ok = await onSubmit(result);
    if (ok) {
      form.reset();
      setOpen(false);
    }
  }

  return (
    <Sheet open={open} onOpenChange={setOpen}>
      <SheetTrigger asChild>
        <Button
          onClick={() => {
            if (!initialValues) {
              const full = PROTOCOLS.map((t) => getProtocolDefaultConfig(t));
              form.reset({
                name: '',
                address: '',
                country: '',
                city: '',
                protocols: full,
              });
            }
            setOpen(true);
          }}
        >
          {trigger}
        </Button>
      </SheetTrigger>
      <SheetContent className='w-[700px] max-w-full md:max-w-screen-md'>
        <SheetHeader>
          <SheetTitle>{title}</SheetTitle>
        </SheetHeader>
        <ScrollArea className='-mx-6 h-[calc(100dvh-48px-36px-36px-env(safe-area-inset-top))]'>
          <Form {...form}>
            <form className='grid grid-cols-1 gap-2 px-6 pt-4'>
              <div className='grid grid-cols-2 gap-2 md:grid-cols-4'>
                <FormField
                  control={control}
                  name='name'
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('name')}</FormLabel>
                      <FormControl>
                        <EnhancedInput {...field} onValueChange={(v) => field.onChange(v)} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={control}
                  name='address'
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('address')}</FormLabel>
                      <FormControl>
                        <EnhancedInput
                          {...field}
                          placeholder={t('address_placeholder')}
                          onValueChange={(v) => field.onChange(v)}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={control}
                  name='country'
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('country')}</FormLabel>
                      <FormControl>
                        <EnhancedInput {...field} onValueChange={(v) => field.onChange(v)} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={control}
                  name='city'
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('city')}</FormLabel>
                      <FormControl>
                        <EnhancedInput {...field} onValueChange={(v) => field.onChange(v)} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>
              <div className='my-3'>
                <h3 className='text-foreground text-sm font-semibold'>
                  {t('protocol_configurations')}
                </h3>
                <p className='text-muted-foreground mt-1 text-xs'>
                  {t('protocol_configurations_desc')}
                </p>
              </div>

              <Accordion
                type='single'
                collapsible
                className='w-full space-y-3'
                value={accordionValue}
                onValueChange={setAccordionValue}
              >
                {PROTOCOLS.map((type) => {
                  const i = Math.max(
                    0,
                    PROTOCOLS.findIndex((t) => t === type),
                  );
                  const current = (protocolsValues[i] || {}) as Record<string, any>;
                  const isEnabled = current?.enable;
                  const fields = PROTOCOL_FIELDS[type] || [];
                  return (
                    <AccordionItem key={type} value={type} className='mb-2 rounded-lg border'>
                      <AccordionTrigger className='px-4 py-3 hover:no-underline'>
                        <div className='flex w-full items-center justify-between'>
                          <div className='flex flex-col items-start gap-1'>
                            <div className='flex items-center gap-1'>
                              <span className='font-medium capitalize'>{type}</span>
                              {current.transport && (
                                <Badge variant='secondary' className='text-xs'>
                                  {current.transport.toUpperCase()}
                                </Badge>
                              )}
                              {current.security && current.security !== 'none' && (
                                <Badge variant='outline' className='text-xs'>
                                  {current.security.toUpperCase()}
                                </Badge>
                              )}
                              {current.port && <Badge className='text-xs'>{current.port}</Badge>}
                            </div>
                            <div className='flex items-center gap-1'>
                              <span
                                className={cn(
                                  'text-xs',
                                  isEnabled ? 'text-green-500' : 'text-muted-foreground',
                                )}
                              >
                                {isEnabled ? t('enabled') : t('disabled')}
                              </span>
                            </div>
                          </div>
                          <Switch
                            className='mr-2'
                            checked={!!isEnabled}
                            disabled={Boolean(
                              initialValues?.id &&
                                isProtocolUsedInNodes(initialValues?.id || 0, type) &&
                                isEnabled,
                            )}
                            onCheckedChange={(checked) => {
                              form.setValue(`protocols.${i}.enable`, checked);
                            }}
                            onClick={(e) => e.stopPropagation()}
                          />
                        </div>
                      </AccordionTrigger>
                      <AccordionContent className='px-4 pb-4 pt-0'>
                        <div className='-mx-4 space-y-4 rounded-b-lg border-t px-4 pt-4'>
                          {renderGroupCard('basic', fields, 'basic', control, form, i, current, t)}
                          {renderGroupCard('obfs', fields, 'obfs', control, form, i, current, t)}
                          {renderGroupCard(
                            'transport',
                            fields,
                            'transport',
                            control,
                            form,
                            i,
                            current,
                            t,
                          )}
                          {renderGroupCard(
                            'security',
                            fields,
                            'security',
                            control,
                            form,
                            i,
                            current,
                            t,
                          )}
                          {renderGroupCard(
                            'reality',
                            fields,
                            'reality',
                            control,
                            form,
                            i,
                            current,
                            t,
                          )}
                          {renderGroupCard(
                            'encryption',
                            fields,
                            'encryption',
                            control,
                            form,
                            i,
                            current,
                            t,
                          )}
                        </div>
                      </AccordionContent>
                    </AccordionItem>
                  );
                })}
              </Accordion>
            </form>
          </Form>
        </ScrollArea>
        <SheetFooter className='flex-row justify-end gap-2 pt-3'>
          <Button variant='outline' disabled={loading} onClick={() => setOpen(false)}>
            {t('cancel')}
          </Button>
          <Button
            disabled={loading}
            onClick={form.handleSubmit(handleSubmit, (errors) => {
              console.log(errors, form.getValues());
              const key = Object.keys(errors)[0] as keyof typeof errors;
              if (key) toast.error(String(errors[key]?.message));
              return false;
            })}
          >
            {loading && <Icon icon='mdi:loading' className='mr-2 animate-spin' />}
            {t('confirm')}
          </Button>
        </SheetFooter>
      </SheetContent>
    </Sheet>
  );
}
