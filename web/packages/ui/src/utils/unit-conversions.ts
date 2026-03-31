import { evaluate, format } from 'mathjs';

type ConversionType =
  | 'centsToDollars'
  | 'dollarsToCents'
  | 'bitsToMb'
  | 'mbToBits'
  | 'bytesToGb'
  | 'gbToBytes';

const conversionConfig: Record<ConversionType, { formula: string; precision: number }> = {
  centsToDollars: { formula: 'value / 100', precision: 2 },
  dollarsToCents: { formula: 'value * 100', precision: 0 },
  bitsToMb: { formula: 'value / 1024 / 1024', precision: 2 },
  mbToBits: { formula: 'value * 1024 * 1024', precision: 0 },
  bytesToGb: { formula: 'value / 1024 / 1024 / 1024', precision: 2 },
  gbToBytes: { formula: 'value * 1024 * 1024 * 1024', precision: 0 },
};

export function unitConversion(type: ConversionType, value?: number | string) {
  if (!value) return 0;

  const config = conversionConfig[type];
  if (!config) throw new Error('Invalid conversion type');

  const formula = config.formula.replace('value', `${value}`);
  const result = evaluate(formula);
  return Number(format(result, { notation: 'fixed', precision: config.precision }));
}

export function evaluateWithPrecision(expression: string) {
  const result = evaluate(expression);
  const formatted = format(result, { notation: 'fixed', precision: 2 });
  return Number(formatted);
}
