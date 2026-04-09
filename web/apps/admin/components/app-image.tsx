"use client";

import { forwardRef } from "react";

type AppImageProps = Omit<React.ImgHTMLAttributes<HTMLImageElement>, "src"> & {
  alt: string;
  fill?: boolean;
  loader?: unknown;
  placeholder?: string;
  priority?: boolean;
  quality?: number;
  sizes?: string;
  src: string | { src?: string };
  unoptimized?: boolean;
};

export default forwardRef<HTMLImageElement, AppImageProps>(function AppImage(
  { alt, fill = false, src, style, ...props },
  ref,
) {
  const resolvedSrc = typeof src === "string" ? src : (src?.src ?? "");

  return (
    // biome-ignore lint/performance/noImgElement: 管理端改为 Vite 后保留轻量图片组件即可。
    <img
      {...props}
      ref={ref}
      alt={alt}
      src={resolvedSrc}
      style={
        fill
          ? {
              ...style,
              height: "100%",
              inset: 0,
              objectFit: "cover",
              position: "absolute",
              width: "100%",
            }
          : style
      }
    />
  );
});
