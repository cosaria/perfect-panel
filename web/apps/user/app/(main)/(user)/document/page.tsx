'use client';

import { Empty } from '@/components/empty';
import { NEXT_PUBLIC_HIDDEN_TUTORIAL_DOCUMENT } from '@/config/constants';
import { queryDocumentList } from '@/services/user/document';
import { getTutorialList } from '@/utils/tutorial';
import { useQuery } from '@tanstack/react-query';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@workspace/ui/components/tabs';
import { useLocale, useTranslations } from 'next-intl';
import { DocumentButton } from './document-button';
import { TutorialButton } from './tutorial-button';

export default function Page() {
  const locale = useLocale();
  const t = useTranslations('document');

  const { data } = useQuery({
    queryKey: ['queryDocumentList'],
    queryFn: async () => {
      const response = await queryDocumentList();
      const list = response.data.data?.list || [];
      return {
        tags: Array.from(new Set(list.reduce((acc: string[], item) => acc.concat(item.tags), []))),
        list,
      };
    },
  });
  const { tags, list: DocumentList } = data || { tags: [], list: [] };

  const { data: TutorialList } = useQuery({
    queryKey: ['getTutorialList', locale],
    queryFn: async () => {
      const list = await getTutorialList();
      return list.get(locale);
    },
    enabled: NEXT_PUBLIC_HIDDEN_TUTORIAL_DOCUMENT !== 'true',
  });

  if (
    (!DocumentList || DocumentList.length === 0) &&
    (!TutorialList || TutorialList.length === 0)
  ) {
    return <Empty border />;
  }

  return (
    <div className='space-y-4'>
      {DocumentList?.length > 0 && (
        <>
          <h2 className='flex items-center gap-1.5 font-semibold'>{t('document')}</h2>
          <Tabs defaultValue='all'>
            <TabsList className='h-full flex-wrap'>
              <TabsTrigger value='all'>{t('all')}</TabsTrigger>
              {tags?.map((item) => (
                <TabsTrigger key={item} value={item}>
                  {item}
                </TabsTrigger>
              ))}
            </TabsList>
            <TabsContent value='all'>
              <DocumentButton items={DocumentList} />
            </TabsContent>
            {tags?.map((item) => (
              <TabsContent value={item} key={item}>
                <DocumentButton
                  items={DocumentList.filter((docs) => (item ? docs.tags.includes(item) : true))}
                />
              </TabsContent>
            ))}
          </Tabs>
        </>
      )}

      {TutorialList && TutorialList?.length > 0 && (
        <>
          <h2 className='flex items-center gap-1.5 font-semibold'>{t('tutorial')}</h2>
          <Tabs defaultValue={TutorialList?.[0]?.title}>
            <TabsList className='h-full flex-wrap'>
              {TutorialList?.map((tutorial) => (
                <TabsTrigger key={tutorial.title} value={tutorial.title}>
                  {tutorial.title}
                </TabsTrigger>
              ))}
            </TabsList>
            {TutorialList?.map((tutorial) => (
              <TabsContent key={tutorial.title} value={tutorial.title}>
                <TutorialButton
                  key={tutorial.path}
                  items={
                    tutorial.subItems && tutorial.subItems?.length > 0
                      ? tutorial.subItems
                      : [tutorial]
                  }
                />
              </TabsContent>
            ))}
          </Tabs>
        </>
      )}
    </div>
  );
}
