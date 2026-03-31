"use client";

import { oAuthLoginGetToken } from "@/services/common/oauth";
import {
	getAllUrlParams,
	getRedirectUrl,
	setAuthorization,
} from "@/utils/common";
import { usePathname, useRouter } from "next/navigation";
import { useEffect } from "react";

interface CertificationProps {
	platform: string;
	children: React.ReactNode;
}

export default function Certification({
	platform,
	children,
}: CertificationProps) {
	const router = useRouter();
	const pathname = usePathname();

	useEffect(() => {
		const searchParams = getAllUrlParams();
		oAuthLoginGetToken({
			method: platform,
			callback: searchParams,
		})
			.then((res) => {
				const token = res?.data?.data?.token;
				if (!token) {
					throw new Error("Invalid token");
				}
				setAuthorization(token);
				router.replace(getRedirectUrl());
				router.refresh();
			})
			.catch((error) => {
				router.replace("/auth");
			});
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, [pathname]);

	return children;
}
