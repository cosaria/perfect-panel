import HyperText from '@workspace/ui/components/hyper-text';
import { OrbitingCircles } from '@workspace/ui/components/orbiting-circles';
import { Icon } from '@workspace/ui/custom-components/icon';
import { getTranslations } from 'next-intl/server';
import Certification from './certification';

export async function generateStaticParams() {
  return [
    {
      platform: 'telegram',
    },
    {
      platform: 'apple',
    },
    {
      platform: 'facebook',
    },
    {
      platform: 'google',
    },
    {
      platform: 'github',
    },
  ];
}

export default async function Page({
  params,
}: {
  params: Promise<{
    platform: string;
  }>;
}) {
  const { platform } = await params;
  const t = await getTranslations('auth');
  return (
    <Certification platform={platform}>
      <div className='bg-background relative flex h-screen w-full flex-col items-center justify-center overflow-hidden'>
        <div className='pointer-events-none flex animate-pulse flex-col items-center whitespace-pre-wrap bg-gradient-to-r from-blue-500 via-indigo-500 to-violet-500 bg-clip-text text-center font-black tracking-tight text-transparent dark:from-blue-400 dark:via-indigo-300 dark:to-violet-400'>
          <HyperText className='text-xl uppercase md:text-2xl'>{platform}</HyperText>
          <HyperText className='text-lg md:text-xl'>{t('authenticating')}</HyperText>
        </div>

        <OrbitingCircles iconSize={40} speed={0.8}>
          <Icon icon='logos:telegram' className='size-12' />
          <Icon icon='uil:apple' className='size-12' />
          <Icon icon='logos:google-icon' className='size-12' />
          <Icon icon='logos:facebook' className='size-12' />
          <Icon icon='uil:github' className='size-12' />
        </OrbitingCircles>
        <OrbitingCircles iconSize={30} radius={100} reverse speed={0.4}>
          <Icon icon='logos:telegram' className='size-10' />
          <Icon icon='uil:apple' className='size-10' />
          <Icon icon='logos:google-icon' className='size-10' />
          <Icon icon='logos:facebook' className='size-10' />
          <Icon icon='uil:github' className='size-10' />
        </OrbitingCircles>
      </div>
    </Certification>
  );
}
