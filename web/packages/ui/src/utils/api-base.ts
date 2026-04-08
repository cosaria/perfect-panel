export function getApiRoot(apiUrl?: string, siteUrl?: string): string {
  return apiUrl || siteUrl || "";
}

export function buildApiBaseUrl(
  apiUrl: string | undefined,
  siteUrl: string | undefined,
  prefix: string,
): string {
  const root = getApiRoot(apiUrl, siteUrl).replace(/\/+$/, "");
  return `${root}${prefix}`;
}
