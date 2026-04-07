declare module "monaco-themes/themes/Dracula.json" {
  interface MonacoTheme {
    base: string;
    inherit: boolean;
    rules: Array<{ token: string; foreground?: string; background?: string; fontStyle?: string }>;
    colors: Record<string, string>;
  }
  const theme: MonacoTheme;
  export default theme;
}
