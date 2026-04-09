"use client";

import { useEffect } from "react";
import { bindOAuthCallback } from "@/services/user-api/sdk.gen";
import { getAllUrlParams } from "@/utils/common";
import { usePathname, useRouter } from "@/utils/router";

interface CertificationProps {
	platform: string;
	children: React.ReactNode;
}

export default function Certification({
	platform,
	children,
}: CertificationProps) {
	const router = useRouter();
	const _pathname = usePathname();

	useEffect(() => {
		const searchParams = getAllUrlParams();
		bindOAuthCallback({
			body: {
				method: platform,
				callback: searchParams,
			},
		})
			.then((_res) => {
				router.replace("/profile");
				router.refresh();
			})
			.catch((_error) => {
				router.replace("/auth");
			});
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, [platform, router.refresh, router.replace]);

	return children;
}
