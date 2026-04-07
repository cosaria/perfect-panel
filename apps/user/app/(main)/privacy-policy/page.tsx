"use client";

import { Markdown } from "@workspace/ui/custom-components/markdown";
import { useEffect, useState } from "react";
import { getPrivacyPolicy } from "@/services/common-api/sdk.gen";

export default function Page() {
  const [content, setContent] = useState("");

  useEffect(() => {
    getPrivacyPolicy()
      .then(({ data }) => {
        setContent(data?.privacy_policy || "");
      })
      .catch(() => {});
  }, []);

  return (
    <div className="container py-8">
      <Markdown>{content}</Markdown>
    </div>
  );
}
