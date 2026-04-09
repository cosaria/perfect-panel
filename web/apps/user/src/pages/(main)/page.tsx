"use client";

import { useEffect, useState } from "react";
import { GlobalMap } from "@/components/main/global-map";
import { Hero } from "@/components/main/hero";
import { ProductShowcase } from "@/components/main/product-showcase/index";
import { Stats } from "@/components/main/stats";
import { queryUserInfo } from "@/services/user-api/sdk.gen";
import { getAuthorization } from "@/utils/common";
import { useRouter } from "@/utils/router";

export default function Home() {
	const router = useRouter();
	const [ready, setReady] = useState(false);

	useEffect(() => {
		const auth = getAuthorization();
		if (auth) {
			queryUserInfo()
				.then(({ data }) => {
					if (data) {
						router.replace("/dashboard");
						return;
					}
					setReady(true);
				})
				.catch(() => {
					setReady(true);
				});
		} else {
			setReady(true);
		}
	}, [router]);

	if (!ready) return null;

	return (
		<main className="container space-y-16">
			<Hero />
			<Stats />
			<ProductShowcase />
			<GlobalMap />
		</main>
	);
}
