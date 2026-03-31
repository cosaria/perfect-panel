"use client";

import { useEffect, useRef } from "react";

interface RawHtmlProps {
  html: string;
  className?: string;
  id?: string;
}

export function RawHtml({ html, className, id }: RawHtmlProps) {
  const ref = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const container = ref.current;

    if (!container) {
      return;
    }

    container.replaceChildren();

    if (!html) {
      return;
    }

    const range = document.createRange();
    range.selectNode(container);
    container.append(range.createContextualFragment(html));
  }, [html]);

  return <div id={id} ref={ref} className={className} />;
}
