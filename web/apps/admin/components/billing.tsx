import { Avatar, AvatarFallback, AvatarImage } from '@workspace/ui/components/avatar';
import { Card, CardDescription, CardHeader, CardTitle } from '@workspace/ui/components/card';
import { getTranslations } from 'next-intl/server';
import Link from 'next/link';

interface BillingProps {
  type: 'dashboard' | 'payment';
}

interface ItemType {
  logo: string;
  title: string;
  description: string;
  expiryDate: string;
  href: string;
}

async function getBillingURL() {
  try {
    const response = await fetch(
      'https://api.github.com/repos/perfect-panel/ppanel-assets/commits',
    );
    const json = await response.json();
    const version = json[0]?.sha || 'latest';
    const url = new URL('https://cdn.jsdmirror.com/gh/perfect-panel/ppanel-assets');
    url.pathname += `@${version}/billing/index.json`;
    return url.toString();
  } catch (error) {
    return 'https://cdn.jsdmirror.com/gh/perfect-panel/ppanel-assets/billing/index.json';
  }
}

export default async function Billing({ type }: BillingProps) {
  const t = await getTranslations('common.billing');
  let list: ItemType[] = [];

  try {
    const url = await getBillingURL();
    const response = await fetch(url, {
      headers: {
        Accept: 'application/json',
      },
    });
    const data = await response.json();
    const now = Date.now();

    list = Array.isArray(data[type])
      ? data[type].filter((item: { expiryDate: string }) => {
          const expiryDate = Date.parse(item.expiryDate);
          return !isNaN(expiryDate) && expiryDate > now;
        })
      : [];
  } catch (error) {
    console.log(error);
    return null;
  }

  if (!list?.length) return null;

  return (
    <>
      <h1 className='text mt-2 font-bold'>
        <span>{t('title')}</span>
        <span className='text-muted-foreground ml-2 text-xs'>{t('description')}</span>
      </h1>
      <div className='grid gap-3 md:grid-cols-3 lg:grid-cols-6'>
        {list.map((item, index) => (
          <Link href={item.href} target='_blank' key={index}>
            <Card className='h-full cursor-pointer'>
              <CardHeader className='flex flex-row gap-2 p-3'>
                <Avatar>
                  <AvatarImage src={item.logo} />
                  <AvatarFallback>{item.title}</AvatarFallback>
                </Avatar>
                <div>
                  <CardTitle>{item.title}</CardTitle>
                  <CardDescription className='mt-2'>{item.description}</CardDescription>
                </div>
              </CardHeader>
            </Card>
          </Link>
        ))}
      </div>
    </>
  );
}
