import { Button } from "@workspace/ui/components/button";
import { Label } from "@workspace/ui/components/label";
import { Switch } from "@workspace/ui/components/switch";
import { Textarea } from "@workspace/ui/components/textarea";
import { Combobox } from "@workspace/ui/custom-components/combobox";
import {
  EnhancedInput,
  type EnhancedInputProps,
} from "@workspace/ui/custom-components/enhanced-input";
import { cn } from "@workspace/ui/lib/utils";
import { CircleMinusIcon, CirclePlusIcon } from "lucide-react";
import { useEffect, useRef, useState } from "react";

type InputRecord = Record<string, unknown>;

const getScalarFieldValue = (fieldValue: unknown): string | number => {
  return typeof fieldValue === "string" || typeof fieldValue === "number" ? fieldValue : "";
};

interface FieldConfig extends Omit<EnhancedInputProps<string | number>, "type"> {
  name: string;
  type: "text" | "number" | "select" | "time" | "boolean" | "textarea";
  options?: { label: string; value: string }[];
  // optional per-item visibility function: returns true to show the field for the given item
  visible?: (item: InputRecord) => boolean;
}

interface ObjectInputProps<T> {
  value: T;
  onChange: (value: T) => void;
  fields: FieldConfig[];
  className?: string;
}

export function ObjectInput<T extends InputRecord>({
  value,
  onChange,
  fields,
  className,
}: ObjectInputProps<T>) {
  const [internalState, setInternalState] = useState<T>(value);

  useEffect(() => {
    setInternalState(value);
  }, [value]);

  const updateField = (key: keyof T, fieldValue: string | number | boolean) => {
    const updatedInternalState = { ...internalState, [key]: fieldValue };
    setInternalState(updatedInternalState);
    onChange(updatedInternalState);
  };
  const renderField = (field: FieldConfig) => {
    // if visible callback exists and returns false for current item, don't render
    if (field.visible && !field.visible(internalState)) return null;
    switch (field.type) {
      case "select":
        return (
          field.options && (
            <Combobox<string, false>
              placeholder={field.placeholder}
              options={field.options}
              value={internalState[field.name] as string | undefined}
              onChange={(fieldValue) => updateField(field.name, fieldValue)}
            />
          )
        );
      case "boolean":
        return (
          <div className="flex h-full items-center space-x-2">
            <Switch
              checked={internalState[field.name] as boolean}
              onCheckedChange={(fieldValue) => updateField(field.name, fieldValue)}
            />
            {field.placeholder && <Label>{field.placeholder}</Label>}
          </div>
        );
      case "textarea":
        return (
          <div className="w-full space-y-2">
            {field.prefix && <Label className="text-sm font-medium">{field.prefix}</Label>}
            <Textarea
              value={getScalarFieldValue(internalState[field.name])}
              onChange={(e) => updateField(field.name, e.target.value)}
              placeholder={field.placeholder}
              className="min-h-32"
            />
          </div>
        );
      default:
        return (
          <EnhancedInput<string | number>
            {...field}
            value={getScalarFieldValue(internalState[field.name])}
            onValueChange={(fieldValue) => updateField(field.name, fieldValue)}
          />
        );
    }
  };
  return (
    <div className={cn("flex flex-1 flex-wrap gap-4", className)}>
      {fields.map((field) => {
        const node = renderField(field);
        if (node === null) return null; // don't render wrapper if field hidden
        return (
          <div key={field.name} className={cn("flex-1", field.className)}>
            {node}
          </div>
        );
      })}
    </div>
  );
}
interface ArrayInputProps<T> {
  value?: T[];
  onChange: (value: T[]) => void;
  fields: FieldConfig[];
  isReverse?: boolean;
  className?: string;
}

export function ArrayInput<T extends InputRecord>({
  value = [],
  onChange,
  fields,
  isReverse = false,
  className,
}: ArrayInputProps<T>) {
  const itemKeyMap = useRef(new WeakMap<T, string>());
  const nextItemKey = useRef(0);

  const getItemKey = (item: T) => {
    const existingKey = itemKeyMap.current.get(item);
    if (existingKey) {
      return existingKey;
    }

    const newKey = `dynamic-input-${nextItemKey.current++}`;
    itemKeyMap.current.set(item, newKey);
    return newKey;
  };

  const initializeDefaultItem = (): T =>
    fields.reduce((acc, field) => {
      acc[field.name as keyof T] = undefined as T[keyof T];
      return acc;
    }, {} as T);

  const [displayItems, setDisplayItems] = useState<T[]>(() => {
    return value.length > 0 ? value : [initializeDefaultItem()];
  });

  const isItemModified = (item: T): boolean =>
    fields.some((field) => {
      const val = item[field.name];
      return val !== undefined && val !== null && val !== "";
    });

  const handleItemChange = (index: number, updatedItem: T) => {
    const newDisplayItems = [...displayItems];
    newDisplayItems[index] = updatedItem;
    setDisplayItems(newDisplayItems);

    const modifiedItems = newDisplayItems.filter(isItemModified);
    onChange(modifiedItems);
  };

  const createField = () => {
    if (isReverse) {
      setDisplayItems([initializeDefaultItem(), ...displayItems]);
    } else {
      setDisplayItems([...displayItems, initializeDefaultItem()]);
    }
  };

  const deleteField = (index: number) => {
    const newDisplayItems = displayItems.filter((_, i) => i !== index);
    setDisplayItems(newDisplayItems);

    const modifiedItems = newDisplayItems.filter(isItemModified);
    onChange(modifiedItems);
  };

  useEffect(() => {
    if (value.length > 0) {
      setDisplayItems(value);
    }
  }, [value]);

  return (
    <div className="flex flex-col gap-4">
      {displayItems.map((item, index) => (
        <div key={getItemKey(item)} className="flex items-center gap-4">
          <ObjectInput
            value={item}
            onChange={(updatedItem) => handleItemChange(index, updatedItem)}
            fields={fields}
            className={className}
          />
          <div className="flex min-w-20 items-center">
            {displayItems.length > 1 && (
              <Button
                variant="ghost"
                size="icon"
                type="button"
                className="text-destructive p-0 text-lg"
                onClick={() => deleteField(index)}
              >
                <CircleMinusIcon />
              </Button>
            )}
            {(isReverse ? index === 0 : index === displayItems.length - 1) && (
              <Button
                variant="ghost"
                size="icon"
                type="button"
                className="text-primary p-0 text-lg"
                onClick={createField}
              >
                <CirclePlusIcon />
              </Button>
            )}
          </div>
        </div>
      ))}
    </div>
  );
}
