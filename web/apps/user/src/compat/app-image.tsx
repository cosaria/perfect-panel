import { forwardRef } from "react";

type ImageProps = Omit<React.ImgHTMLAttributes<HTMLImageElement>, "src"> & {
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

export default forwardRef<HTMLImageElement, ImageProps>(function Image(
  { alt, fill = false, src, style, ...props },
  ref,
) {
  const resolvedSrc = typeof src === "string" ? src : (src?.src ?? "");

  return (
    // biome-ignore lint/performance/noImgElement: Vite 应用这里使用轻量图片封装。
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
