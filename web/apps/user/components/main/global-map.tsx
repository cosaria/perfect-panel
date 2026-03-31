'use client';

import { DotLottieReact } from '@lottiefiles/dotlottie-react';
import GlobalMapLottie from '@workspace/ui/lotties/global-map.json';
import { motion } from 'framer-motion';
import { useTranslations } from 'next-intl';

export function GlobalMap() {
  const t = useTranslations('index');
  return (
    <motion.section
      initial={{ opacity: 0 }}
      whileInView={{ opacity: 1 }}
      viewport={{ once: true }}
      transition={{ duration: 0.5 }}
    >
      <motion.h2
        initial={{ opacity: 0, y: -20 }}
        whileInView={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5 }}
        className='mb-2 text-center text-3xl font-bold'
      >
        {t('global_map_itle')}
      </motion.h2>
      <motion.p
        initial={{ opacity: 0, y: -20 }}
        whileInView={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5, delay: 0.1 }}
        className='text-muted-foreground mb-8 text-center text-lg'
      >
        {t('global_map_description')}
      </motion.p>
      <motion.div
        className='aspect-video w-full overflow-hidden'
        initial={{ scale: 0.9, opacity: 0 }}
        animate={{ scale: 1, opacity: 1 }}
        transition={{
          type: 'spring',
          stiffness: 100,
          damping: 15,
          delay: 0.4,
        }}
      >
        <DotLottieReact className='w-full scale-150' data={GlobalMapLottie} autoplay loop />
      </motion.div>
    </motion.section>
  );
}
