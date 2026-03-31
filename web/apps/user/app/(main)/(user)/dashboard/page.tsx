import Announcement from '@/components/announcement';
import { cookies } from 'next/headers';
import Content from './content';

export default async function Page() {
  return (
    <div className='flex min-h-[calc(100vh-64px-58px-32px-114px)] w-full flex-col gap-4 overflow-hidden'>
      <Announcement type='pinned' Authorization={(await cookies()).get('Authorization')?.value} />
      <Content />
    </div>
  );
}
