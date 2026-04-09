"use client";

import { useQuery } from "@tanstack/react-query";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@workspace/ui/components/tabs";
import { useTranslations } from "@workspace/ui/components/i18n-provider";
import { Empty } from "@/components/empty";
import { queryDocumentList } from "@/services/user-api/sdk.gen";
import { DocumentButton } from "./document-button";

export default function Page() {
  const t = useTranslations("document");

  const { data } = useQuery({
    queryKey: ["queryDocumentList"],
    queryFn: async () => {
      const { data: response } = await queryDocumentList();
      const list = response?.list || [];
      return {
        tags: Array.from(
          new Set(list.reduce((acc: string[], item) => acc.concat(item.tags || []), [])),
        ),
        list,
      };
    },
  });
  const { tags, list: DocumentList } = data || { tags: [], list: [] };

  if (!DocumentList || DocumentList.length === 0) {
    return <Empty border />;
  }

  return (
    <div className="space-y-4">
      {DocumentList?.length > 0 && (
        <>
          <h2 className="flex items-center gap-1.5 font-semibold">{t("document")}</h2>
          <Tabs defaultValue="all">
            <TabsList className="h-full flex-wrap">
              <TabsTrigger value="all">{t("all")}</TabsTrigger>
              {tags?.map((item) => (
                <TabsTrigger key={item} value={item}>
                  {item}
                </TabsTrigger>
              ))}
            </TabsList>
            <TabsContent value="all">
              <DocumentButton items={DocumentList} />
            </TabsContent>
            {tags?.map((item) => (
              <TabsContent value={item} key={item}>
                <DocumentButton
                  items={DocumentList.filter((docs) => (item ? docs.tags?.includes(item) : true))}
                />
              </TabsContent>
            ))}
          </Tabs>
        </>
      )}
    </div>
  );
}
