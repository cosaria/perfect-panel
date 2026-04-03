import type { z } from "zod";
import type { protocols } from "./constants";
import type { formSchema, protocolApiScheme } from "./schemas";

export type ProtocolConfig = z.infer<typeof protocolApiScheme>;
export type ServerFormValues = z.infer<typeof formSchema>;

type TranslationFn = (key: string) => string;
type ProtocolContext = Partial<ProtocolConfig>;
type GeneratedFieldValue = string | Record<string, string>;
type FieldDefaultValue = string | number | boolean | null;

export type FieldConfig = {
  name: string;
  type: "input" | "select" | "switch" | "number" | "textarea";
  label: string;
  placeholder?: string | ((t: TranslationFn, protocol: ProtocolContext) => string);
  options?: readonly string[];
  defaultValue?: FieldDefaultValue;
  min?: number;
  max?: number;
  step?: number;
  suffix?: string;
  generate?: {
    function?: () => Promise<GeneratedFieldValue> | GeneratedFieldValue;
    functions?: {
      label: string | ((t: TranslationFn, protocol: ProtocolContext) => string);
      function: () => Promise<GeneratedFieldValue> | GeneratedFieldValue;
    }[];
    updateFields?: Record<string, string>;
  };
  condition?: (protocol: ProtocolContext, values: Partial<ServerFormValues>) => boolean;
  group?: "basic" | "transport" | "security" | "reality" | "obfs" | "encryption";
  gridSpan?: 1 | 2;
};

export type ProtocolType = (typeof protocols)[number];
