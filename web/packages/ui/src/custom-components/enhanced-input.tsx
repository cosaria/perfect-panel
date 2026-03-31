import { Input } from '@workspace/ui/components/input';
import { cn } from '@workspace/ui/lib/utils';
import { ChangeEvent, ReactNode, useEffect, useState } from 'react';

export interface EnhancedInputProps<T = string>
  extends Omit<React.InputHTMLAttributes<HTMLInputElement>, 'prefix' | 'value' | 'onChange'> {
  prefix?: string | ReactNode;
  suffix?: string | ReactNode;
  value?: T;
  formatInput?: (value: T) => string | number;
  formatOutput?: (value: string | number) => T;
  onValueChange?: (value: T) => void;
  onValueBlur?: (value: T) => void;
  min?: number;
  max?: number;
}

export function EnhancedInput<T = string>({
  suffix,
  prefix,
  formatInput,
  formatOutput,
  value: initialValue,
  className,
  onValueChange,
  onValueBlur,
  ...props
}: EnhancedInputProps<T>) {
  const getProcessedValue = (inputValue: unknown) => {
    if (inputValue === '' || inputValue === 0 || inputValue === '0') return '';
    const newValue = String(inputValue ?? '');
    return formatInput ? formatInput(inputValue as T) : newValue;
  };

  const [value, setValue] = useState<string | number>(() => getProcessedValue(initialValue));
  const [internalValue, setInternalValue] = useState<T | string | number>(initialValue ?? '');

  useEffect(() => {
    if (initialValue !== internalValue) {
      const newValue = getProcessedValue(initialValue);
      if (value !== newValue) {
        setValue(newValue);
        setInternalValue(initialValue ?? '');
      }
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [initialValue, formatInput]);

  const processValue = (inputValue: string | number): T => {
    let processedValue: number | string = inputValue?.toString().trim();

    if (processedValue === '0' && props.type === 'number') {
      return (formatOutput ? formatOutput(0) : 0) as T;
    }

    if (processedValue && props.type === 'number') processedValue = Number(processedValue);
    return formatOutput ? formatOutput(processedValue) : (processedValue as T);
  };

  const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
    let inputValue = e.target.value;

    if (props.type === 'number') {
      if (inputValue === '0') {
        setValue('');
        setInternalValue(0);
        onValueChange?.(processValue(0));
        return;
      }

      if (/^-?\d*\.?\d*$/.test(inputValue) || inputValue === '-' || inputValue === '.') {
        const numericValue = Number(inputValue);
        if (!isNaN(numericValue) && inputValue !== '-' && inputValue !== '.') {
          const min = Number.isFinite(props.min) ? props.min : -Infinity;
          const max = Number.isFinite(props.max) ? props.max : Infinity;
          const constrainedValue = Math.max(min!, Math.min(max!, numericValue));
          inputValue = String(constrainedValue);
          setInternalValue(constrainedValue);
        } else {
          setInternalValue(inputValue);
        }
        setValue(inputValue === '0' ? '' : inputValue);
      }
    } else {
      setValue(inputValue);
      setInternalValue(inputValue);
    }

    const outputValue = processValue(inputValue);
    onValueChange?.(outputValue);
  };

  const handleBlur = () => {
    if (props.type === 'number' && value) {
      if (value === '-' || value === '.') {
        setValue('');
        setInternalValue('');
        onValueBlur?.('' as T);
        return;
      }

      if (value === '0') {
        setValue('');
        onValueBlur?.(processValue(0));
        return;
      }
    }

    const outputValue = processValue(value);
    if ((initialValue || '') !== outputValue) {
      onValueBlur?.(outputValue);
    }
  };

  const renderPrefix = () => {
    return typeof prefix === 'string' ? (
      <div className='bg-muted relative mr-px flex h-9 items-center text-nowrap px-3'>{prefix}</div>
    ) : (
      prefix
    );
  };

  const renderSuffix = () => {
    return typeof suffix === 'string' ? (
      <div className='bg-muted relative ml-px flex h-9 items-center text-nowrap px-3'>{suffix}</div>
    ) : (
      suffix
    );
  };

  return (
    <div
      className={cn(
        'border-input flex w-full items-center overflow-hidden rounded-md border',
        className,
      )}
      suppressHydrationWarning
    >
      {renderPrefix()}
      <Input
        step={0.01}
        autoComplete='off'
        {...props}
        value={value}
        className='block rounded-none border-none'
        onChange={handleChange}
        onBlur={handleBlur}
      />
      {renderSuffix()}
    </div>
  );
}
