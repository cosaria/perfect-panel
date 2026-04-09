"use client";

import { Button } from "@workspace/ui/components/button";
import { useTranslations } from "@workspace/ui/components/i18n-provider";
import { useCountDown } from "ahooks";
import { useEffect, useState } from "react";
import useGlobalStore from "@/config/use-global";
import { sendEmailCode, sendSmsCode } from "@/services/common-api/sdk.gen";

interface SendCodeProps {
	type: "email" | "phone";
	params: {
		email?: string;
		type?: 1 | 2;
		telephone_area_code?: string;
		telephone?: string;
	};
}
export default function SendCode({ type, params }: SendCodeProps) {
	const t = useTranslations("auth");
	const { common } = useGlobalStore();
	const { verify_code_interval } = common.verify_code;
	const [targetDate, setTargetDate] = useState<number>();

	useEffect(() => {
		const storedEndTime = localStorage.getItem(`verify_code_${type}`);
		if (storedEndTime) {
			const endTime = parseInt(storedEndTime, 10);
			if (endTime > Date.now()) {
				setTargetDate(endTime);
			} else {
				localStorage.removeItem(`verify_code_${type}`);
			}
		}
	}, [type]);

	const [, { seconds }] = useCountDown({
		targetDate,
		onEnd: () => {
			setTargetDate(undefined);
			localStorage.removeItem(`verify_code_${type}`);
		},
	});

	const setCodeTimer = () => {
		const endTime = Date.now() + verify_code_interval * 1000;
		setTargetDate(endTime);
		localStorage.setItem(`verify_code_${type}`, endTime.toString());
	};

	const getEmailCode = async () => {
		if (params.email && params.type) {
			await sendEmailCode({
				body: {
					email: params.email,
					type: params.type,
				},
			});
			setCodeTimer();
		}
	};

	const getPhoneCode = async () => {
		if (params.telephone && params.telephone_area_code && params.type) {
			await sendSmsCode({
				body: {
					telephone: params.telephone,
					telephone_area_code: params.telephone_area_code,
					type: params.type,
				},
			});
			setCodeTimer();
		}
	};

	const handleSendCode = async () => {
		if (type === "email") {
			getEmailCode();
		} else {
			getPhoneCode();
		}
	};
	const disabled =
		seconds > 0 ||
		(type === "email"
			? !params.email
			: !params.telephone || !params.telephone_area_code);

	return (
		<Button type="button" onClick={handleSendCode} disabled={disabled}>
			{seconds > 0 ? `${seconds}s` : t("get")}
		</Button>
	);
}
