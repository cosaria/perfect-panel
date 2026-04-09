"use client";

import { default as _Empty } from "@workspace/ui/custom-components/empty";
import { useTranslations } from "@workspace/ui/components/i18n-provider";
import { useEffect, useState } from "react";

export function Empty() {
  const t = useTranslations("common");

  const [description, setDescription] = useState("");

  useEffect(() => {
    const random = Math.floor(Math.random() * 10);
    setDescription(t(`empty.${random}`));
  }, [t]);

  return <_Empty description={description} />;
}
