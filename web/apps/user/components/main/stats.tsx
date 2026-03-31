'use client';

import {
  NEXT_PUBLIC_HOME_LOCATION_COUNT,
  NEXT_PUBLIC_HOME_SERVER_COUNT,
  NEXT_PUBLIC_HOME_USER_COUNT,
} from '@/config/constants';
import { DotLottieReact } from '@lottiefiles/dotlottie-react';
import LocationsLittie from '@workspace/ui/lotties/locations.json';
import ServersLottie from '@workspace/ui/lotties/servers.json';
import UsersLottie from '@workspace/ui/lotties/users.json';
import { motion } from 'framer-motion';
import { useTranslations } from 'next-intl';
import { useEffect, useState } from 'react';

export function Stats() {
  const t = useTranslations('index');

  const list = [
    {
      name: t('users'),
      number: NEXT_PUBLIC_HOME_USER_COUNT,
      icon: <DotLottieReact className='size-24' data={UsersLottie} autoplay loop />,
    },
    {
      name: t('servers'),
      number: NEXT_PUBLIC_HOME_SERVER_COUNT,
      icon: <DotLottieReact className='size-24' data={ServersLottie} autoplay loop />,
    },
    {
      name: t('locations'),
      number: NEXT_PUBLIC_HOME_LOCATION_COUNT,
      icon: <DotLottieReact className='size-24' data={LocationsLittie} autoplay loop />,
    },
  ];
  return (
    <motion.section
      className='divide-muted z-10 grid w-full grid-cols-1 divide-y-2 rounded-lg sm:grid-cols-3 sm:divide-x-2 sm:divide-y-0'
      initial={{ opacity: 0, y: 50 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 1, ease: 'easeOut' }}
      viewport={{ once: true, amount: 0.8 }}
      whileInView={{ opacity: 1, y: 0 }}
    >
      {list.map((item, index) => (
        <motion.div
          className='mx-auto flex w-10/12 items-center justify-start px-4 py-4 sm:w-full sm:justify-center sm:py-6'
          key={item.name}
          initial={{ opacity: 0, scale: 0.8 }}
          whileInView={{ opacity: 1, scale: 1 }}
          transition={{ duration: 0.8, delay: index * 0.3, ease: 'easeOut' }}
          viewport={{ once: true, amount: 0.8 }}
        >
          <div className='flex w-full items-center sm:w-auto'>
            <div className='mr-4 flex h-20 w-20 items-center justify-center rounded-full'>
              {item.icon}
            </div>
            <div className='flex flex-col'>
              <p className='text-xl font-bold'>
                <CountUp end={item.number} duration={2000 + index * 500} />+
              </p>
              <p className='text-muted-foreground text-lg'>{item.name}</p>
            </div>
          </div>
        </motion.div>
      ))}
    </motion.section>
  );
}

function CountUp({ end, duration = 2000 }: { end: number; duration?: number }) {
  const [count, setCount] = useState(0);

  useEffect(() => {
    let startTime: number;
    let animationFrame: number;

    const easeOutQuad = (t: number) => t * (2 - t);

    const updateCount = (timestamp: number) => {
      if (!startTime) startTime = timestamp;
      const progress = timestamp - startTime;
      const easedProgress = easeOutQuad(Math.min(progress / duration, 1));
      const nextCount = Math.round(easedProgress * end);

      setCount(nextCount);

      if (progress < duration) {
        animationFrame = requestAnimationFrame(updateCount);
      }
    };

    animationFrame = requestAnimationFrame(updateCount);

    return () => cancelAnimationFrame(animationFrame);
  }, [end, duration]);

  return <>{count.toLocaleString()}</>;
}
