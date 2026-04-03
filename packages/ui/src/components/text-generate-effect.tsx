"use client";
import { cn } from "@workspace/ui/lib/utils";
import { motion, stagger, useAnimate } from "framer-motion";
import { useEffect } from "react";

export const TextGenerateEffect = ({
  words,
  className,
  filter = true,
  duration = 0.5,
}: {
  words: string;
  className?: string;
  filter?: boolean;
  duration?: number;
}) => {
  const [scope, animate] = useAnimate();
  const wordsArray = words.split(" ");
  useEffect(() => {
    animate(
      "span",
      {
        opacity: 1,
        filter: filter ? "blur(0px)" : "none",
      },
      {
        duration: duration ? duration : 1,
        delay: stagger(0.2),
      },
    );
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [animate, duration, filter]);

  const renderWords = () => {
    const keyedWords = wordsArray.reduce<Array<{ key: string; word: string }>>((result, word) => {
      const occurrence = result.filter((item) => item.word === word).length + 1;
      result.push({ key: `${word}-${occurrence}`, word });
      return result;
    }, []);

    return (
      <motion.div ref={scope}>
        {keyedWords.map(({ key, word }) => {
          return (
            <motion.span
              key={key}
              className="text-black opacity-0 dark:text-white"
              style={{
                filter: filter ? "blur(10px)" : "none",
              }}
            >
              {word}{" "}
            </motion.span>
          );
        })}
      </motion.div>
    );
  };

  return (
    <div className={cn("font-bold", className)}>
      <div className="mt-4">
        <div className="text-2xl leading-snug tracking-wide text-black dark:text-white">
          {renderWords()}
        </div>
      </div>
    </div>
  );
};
