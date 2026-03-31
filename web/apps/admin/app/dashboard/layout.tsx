import { Header } from '@/components/header';
import { SidebarLeft } from '@/components/sidebar-left';
import { SidebarInset, SidebarProvider } from '@workspace/ui/components/sidebar';
import { cookies } from 'next/headers';

export default async function DashboardLayout({ children }: { children: React.ReactNode }) {
  const cookieStore = await cookies();
  const defaultOpen = cookieStore.get('sidebar_state')?.value === 'true';
  return (
    <SidebarProvider defaultOpen={defaultOpen}>
      <SidebarLeft />
      <SidebarInset className='relative flex-grow overflow-hidden'>
        <Header />
        <div className='h-[calc(100vh-56px)] flex-grow gap-4 overflow-auto p-4'>{children}</div>
      </SidebarInset>
    </SidebarProvider>
  );
}
