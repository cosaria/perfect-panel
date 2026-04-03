"use client";

import { usePathname, useRouter } from "next/navigation";
import { useEffect } from "react";
import { oAuthLoginGetToken } from "@/services/user-api/sdk.gen";
import { getAllUrlParams, getRedirectUrl, setAuthorization } from "@/utils/common";

interface CertificationProps {
  platform: string;
  children: React.ReactNode;
}

export default function Certification({ platform, children }: CertificationProps) {
  const router = useRouter();
  const _pathname = usePathname();

  useEffect(() => {
    const searchParams = getAllUrlParams();
    oAuthLoginGetToken({
      body: {
        method: platform,
        callback: searchParams,
      },
    })
      .then(({ data }) => {
        const token = data?.token;
        if (!token) {
          throw new Error("Invalid token");
        }
        setAuthorization(token);
        router.replace(getRedirectUrl());
        router.refresh();
      })
      .catch((_error) => {
        router.replace("/auth");
      });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [platform, router.refresh, router.replace]);

  return children;
}
