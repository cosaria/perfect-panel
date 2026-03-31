import rtlDetect from 'rtl-detect';

export function getLangDir(locale: string) {
  return rtlDetect.getLangDir(locale);
}
