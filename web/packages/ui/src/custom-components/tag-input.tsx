import { Badge } from '@workspace/ui/components/badge';
import { Input } from '@workspace/ui/components/input';
import { cn } from '@workspace/ui/lib/utils';
import { X } from 'lucide-react';
import React, { useEffect, useRef, useState } from 'react';

interface TagInputProps {
  value?: string[];
  onChange?: (tags: string[]) => void;
  placeholder?: string;
  separator?: string;
  className?: string;
  options?: string[];
}

export function TagInput({
  value = [],
  onChange,
  placeholder,
  separator = ',',
  className,
  options = [],
}: TagInputProps) {
  const [inputValue, setInputValue] = useState('');
  const [tags, setTags] = useState<string[]>(value);
  const [open, setOpen] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    setTags(value.map((tag) => tag.trim()).filter((tag) => tag));
  }, [value]);

  function normalizeInput(input: string) {
    return input.replace(/，/g, ',');
  }

  function addTag(tagValue?: string) {
    let tagsToAdd: string[] = [];
    let shouldKeepOpen = false;

    if (tagValue) {
      if (!tags.includes(tagValue)) {
        tagsToAdd = [tagValue];
        shouldKeepOpen = true;
      }
    } else if (inputValue.trim()) {
      const normalizedInput = normalizeInput(inputValue);
      tagsToAdd = normalizedInput
        .split(separator)
        .map((tag) => tag.trim())
        .filter((tag) => tag && !tags.includes(tag));
    }

    if (tagsToAdd.length > 0) {
      const updatedTags = [...tags, ...tagsToAdd];
      updateTags(updatedTags);
    }
    setInputValue('');

    if (shouldKeepOpen && options.length > 0) {
      setTimeout(() => {
        setOpen(true);
      }, 10);
    } else {
      setOpen(false);
    }
  }

  function handleKeyDown(event: React.KeyboardEvent<HTMLInputElement>) {
    if (event.key === 'Enter' || event.key === separator || event.key === '，') {
      event.preventDefault();
      addTag();
    } else if (event.key === 'Backspace' && inputValue === '') {
      event.preventDefault();
      handleRemoveTag(tags.length - 1);
    } else if (event.key === 'Escape') {
      setOpen(false);
    }
  }

  function handleInputFocus() {
    if (options.length > 0) {
      setOpen(true);
    }
  }

  function handleInputBlur() {
    if (inputValue.trim()) addTag();
    setOpen(false);
  }

  function handleRemoveTag(index: number) {
    const newTags = tags.filter((_, i) => i !== index);
    updateTags(newTags);
  }

  function updateTags(newTags: string[]) {
    setTags(newTags);
    onChange?.(newTags);
  }

  const availableOptions = options
    .filter((option) => !tags.includes(option))
    .filter(
      (option) =>
        inputValue.trim() === '' || option.toLowerCase().includes(inputValue.toLowerCase()),
    );

  return (
    <div className={cn('relative', className)}>
      <div
        className={cn(
          'border-input focus-within:ring-primary flex min-h-9 w-full cursor-text flex-wrap items-center gap-2 rounded-md border bg-transparent p-2 shadow-sm transition-colors focus-within:ring-0',
        )}
        onClick={() => inputRef.current?.focus()}
      >
        {tags.map((tag, index) => (
          <Badge
            key={tag}
            variant='outline'
            className='border-primary bg-primary/10 flex items-center gap-1 px-1'
            onClick={(e) => e.stopPropagation()}
          >
            {tag}
            <X
              className='hover:text-destructive size-4 cursor-pointer rounded-sm'
              onClick={(e) => {
                e.stopPropagation();
                handleRemoveTag(index);
              }}
            />
          </Badge>
        ))}
        <div className='flex min-w-0 flex-1 items-center gap-2'>
          <Input
            ref={inputRef}
            className='h-full min-w-0 flex-1 border-none bg-transparent p-0 shadow-none !ring-0'
            value={inputValue}
            onChange={(e) => setInputValue(e.target.value)}
            onKeyDown={handleKeyDown}
            onFocus={handleInputFocus}
            onBlur={handleInputBlur}
            placeholder={placeholder}
          />

          {open && availableOptions.length > 0 && (
            <div className='bg-popover text-popover-foreground absolute left-0 top-full z-50 max-h-60 w-full overflow-auto rounded-md border shadow-md'>
              {availableOptions.map((option) => (
                <div
                  key={option}
                  className='hover:bg-accent hover:text-accent-foreground relative flex cursor-pointer select-none items-center px-2 py-1.5 text-sm'
                  onMouseDown={(e) => {
                    e.preventDefault();
                    addTag(option);
                    setTimeout(() => {
                      inputRef.current?.focus();
                    }, 10);
                  }}
                >
                  {option}
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

export default TagInput;
