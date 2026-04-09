"use client";

import { Suspense, useEffect, useState } from "react";
import { getClientLocale } from "@/locales/client";
import { getSubscription } from "@/services/user-api/sdk.gen";
import type { Subscribe } from "@/services/user-api/types.gen";
import { useSearchParams } from "@/utils/router";
import Content from "./content";

function PurchasingContent() {
	const searchParams = useSearchParams();
	const id = searchParams.get("id");
	const [subscription, setSubscription] = useState<Subscribe | undefined>();

	useEffect(() => {
		const locale = getClientLocale();
		getSubscription({
			query: { language: locale },
		})
			.then(({ data }) => {
				const list = data?.list || [];
				const found = list.find((item) => item.id === Number(id));
				setSubscription(found);
			})
			.catch(() => {});
	}, [id]);

	return (
		<main className="container space-y-16">
			<Content subscription={subscription} />
		</main>
	);
}

export default function Page() {
	return (
		<Suspense>
			<PurchasingContent />
		</Suspense>
	);
}
