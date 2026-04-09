import type { z } from "zod";
import type { protocols } from "./constants";
import type { formSchema, protocolApiScheme } from "./schemas";

export type ProtocolConfig = z.infer<typeof protocolApiScheme>;
export type ServerFormValues = z.infer<typeof formSchema>;

type TranslationFn = (key: string) => string;
/**
 * Loose protocol context for condition callbacks and UI field rendering.
 * Using Record<string, unknown> instead of Partial<ProtocolConfig> because
 * the discriminated union makes it impossible to access protocol-specific
 * fields (e.g. transport, security, obfs) without narrowing first,
 * which is impractical in generic form field conditions.
 */
/**
 * Loose protocol context for condition callbacks and UI field rendering.
 * All values accessed in condition callbacks are strings (cipher, transport,
 * security, obfs, cert_mode, encryption, encryption_rtt), so we use `string`
 * as the value type. This avoids narrowing issues from the discriminated union.
 */
type ProtocolContext = Record<string, string>;
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
  condition?: (protocol: ProtocolContext, values: Partial<ServerFormValues>) => unknown;
  group?: "basic" | "transport" | "security" | "reality" | "obfs" | "encryption";
  gridSpan?: 1 | 2;
};

export type ProtocolType = (typeof protocols)[number];
