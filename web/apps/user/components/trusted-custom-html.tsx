"use client";

import DOMPurify from "dompurify";

type TrustedCustomHtmlProps = {
  html: string;
};

export function TrustedCustomHtml({ html }: TrustedCustomHtmlProps) {
  const sanitizedHtml = DOMPurify.sanitize(html, {
    USE_PROFILES: { html: true },
  });

  if (!sanitizedHtml) {
    return null;
  }

  // biome-ignore lint/security/noDangerouslySetInnerHtml: custom_html is sanitized and only rendered through this audited wrapper.
  return <div id="custom_html" dangerouslySetInnerHTML={{ __html: sanitizedHtml }} />;
}
