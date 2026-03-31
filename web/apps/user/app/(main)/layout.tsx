import Footer from '@/components/footer';
import Header from '@/components/header';

export default async function MainLayout({ children }: { children: React.ReactNode }) {
  return (
    <>
      <Header />
      {children}
      <Footer />
    </>
  );
}
