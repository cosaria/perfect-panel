'use client';

import useGlobalStore from '@/config/use-global';
import { DotLottieReact } from '@lottiefiles/dotlottie-react';
import { HoverBorderGradient } from '@workspace/ui/components/hover-border-gradient';
import { TextGenerateEffect } from '@workspace/ui/components/text-generate-effect';
import NetworkSecurityLottie from '@workspace/ui/lotties/network-security.json';
import { motion } from 'framer-motion';
import { useTranslations } from 'next-intl';
import Link from 'next/link';

export function Hero() {
  const t = useTranslations('index');
  const { common, user } = useGlobalStore();
  const { site } = common;

  return (
    <motion.div
      initial={{ opacity: 0, y: -50 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ type: 'spring', stiffness: 100, damping: 20 }}
      viewport={{ once: true, amount: 0.2 }}
      className='grid gap-8 pt-16 sm:grid-cols-2'
    >
      <motion.div
        initial={{ opacity: 0, y: 50 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ type: 'spring', stiffness: 80, damping: 15, delay: 0.3 }}
        viewport={{ once: true, amount: 0.3 }}
        className='flex flex-col items-start justify-center'
      >
        <h1 className='my-6 text-4xl font-bold lg:text-6xl'>
          {t('welcome')} {site.site_name}
        </h1>
        {site.site_desc && (
          <TextGenerateEffect
            words={site.site_desc}
            className='*:text-muted-foreground mb-8 max-w-xl'
          />
        )}
        <Link href={user ? '/dashboard' : '/auth'}>
          <HoverBorderGradient
            containerClassName='rounded-full'
            as='button'
            className='m-0.5 flex items-center space-x-2 text-white'
          >
            {t('started')}
          </HoverBorderGradient>
        </Link>
      </motion.div>
      <motion.div
        initial={{ opacity: 0, y: 50 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ type: 'spring', stiffness: 80, damping: 15, delay: 0.5 }}
        viewport={{ once: true, amount: 0.3 }}
        className='flex w-full'
      >
        <DotLottieReact data={NetworkSecurityLottie} autoplay loop />
      </motion.div>
    </motion.div>
  );
}
