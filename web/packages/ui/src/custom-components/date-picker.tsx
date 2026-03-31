'use client';

import { Button } from '@workspace/ui/components/button';
import { Calendar, CalendarProps } from '@workspace/ui/components/calendar';
import { Popover, PopoverContent, PopoverTrigger } from '@workspace/ui/components/popover';
import { cn } from '@workspace/ui/lib/utils';
import { intlFormat } from 'date-fns';
import { CalendarIcon } from 'lucide-react';
import * as React from 'react';

export function DatePicker({
  placeholder,
  value,
  onChange,
  ...props
}: CalendarProps & {
  placeholder?: string;
  value?: number;
  onChange?: (value?: number) => void;
}) {
  const [date, setDate] = React.useState<Date | undefined>(value ? new Date(value) : undefined);

  const handleSelect = (selectedDate: Date | undefined) => {
    setDate(selectedDate);
    if (onChange) {
      onChange(selectedDate?.getTime() || 0);
    }
  };

  return (
    <Popover>
      <PopoverTrigger asChild>
        <Button
          variant='outline'
          className={cn('w-full justify-between font-normal', !value && 'text-muted-foreground')}
        >
          {value ? intlFormat(value) : <span>{placeholder}</span>}
          <CalendarIcon className='size-4' />
        </Button>
      </PopoverTrigger>
      <PopoverContent className='w-auto p-0' align='start'>
        <Calendar {...props} mode='single' selected={date} onSelect={handleSelect} initialFocus />
      </PopoverContent>
    </Popover>
  );
}
