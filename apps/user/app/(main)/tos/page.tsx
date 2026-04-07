"use client";

import { Markdown } from "@workspace/ui/custom-components/markdown";
import { useEffect, useState } from "react";
import { getTos } from "@/services/common-api/sdk.gen";

export default function Page() {
  const [content, setContent] = useState("");

  useEffect(() => {
    getTos()
      .then(({ data }) => {
        setContent(data?.tos_content || "");
      })
      .catch(() => {});
  }, []);

  return (
    <div className="container py-8">
      <Markdown>{content}</Markdown>
    </div>
  );
}
