import { useEffect, useRef, useState } from "react";
import { useLocation } from "react-router-dom";

type RouterTopLoaderProps = {
  color?: string;
  height?: number;
  shadow?: string;
  showSpinner?: boolean;
  zIndex?: number;
};

export default function RouterTopLoader({
  color = "hsl(221 83% 53%)",
  height = 3,
  shadow = "0 0 10px rgba(59, 130, 246, 0.45)",
  showSpinner = false,
  zIndex = 60,
}: RouterTopLoaderProps) {
  const location = useLocation();
  const isFirstRender = useRef(true);
  const [isVisible, setIsVisible] = useState(false);
  const [progress, setProgress] = useState(0);
  const routeSignature = `${location.pathname}${location.search}${location.hash}`;

  useEffect(() => {
    if (routeSignature.length === 0) {
      return;
    }

    if (isFirstRender.current) {
      isFirstRender.current = false;
      return;
    }

    setIsVisible(true);
    setProgress(0.14);

    const animationFrame = requestAnimationFrame(() => {
      setProgress(0.38);
    });
    const nearCompleteTimer = window.setTimeout(() => {
      setProgress(0.78);
    }, 180);
    const completeTimer = window.setTimeout(() => {
      setProgress(1);
    }, 360);
    const hideTimer = window.setTimeout(() => {
      setIsVisible(false);
    }, 560);
    const resetTimer = window.setTimeout(() => {
      setProgress(0);
    }, 760);

    return () => {
      cancelAnimationFrame(animationFrame);
      window.clearTimeout(nearCompleteTimer);
      window.clearTimeout(completeTimer);
      window.clearTimeout(hideTimer);
      window.clearTimeout(resetTimer);
    };
  }, [routeSignature]);

  return (
    <>
      <div
        aria-hidden="true"
        className="pointer-events-none fixed inset-x-0 top-0"
        style={{ zIndex }}
      >
        <div
          style={{
            background: color,
            boxShadow: isVisible ? shadow : "none",
            height,
            opacity: isVisible ? 1 : 0,
            transform: `scaleX(${progress})`,
            transformOrigin: "left center",
            transition: "transform 220ms ease, opacity 220ms ease",
            width: "100%",
          }}
        />
      </div>
      {showSpinner ? null : null}
    </>
  );
}
