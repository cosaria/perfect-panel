"use client";

import { useTranslations } from "@workspace/ui/components/i18n-provider";
import { type ReactNode, useState, useTransition } from "react";
import { toast } from "sonner";
import {
	VITE_DEFAULT_USER_EMAIL,
	VITE_DEFAULT_USER_PASSWORD,
} from "@/config/constants";
import {
	resetPassword,
	userLogin,
	userRegister,
} from "@/services/user-api/sdk.gen";
import { getRedirectUrl, setAuthorization } from "@/utils/common";
import { useRouter } from "@/utils/router";
import LoginForm from "./login-form";
import RegisterForm from "./register-form";
import ResetForm from "./reset-form";
import type { AuthView, EmailAuthValues } from "./types";

export default function EmailAuthForm() {
	const t = useTranslations("auth");
	const router = useRouter();
	const [type, setType] = useState<AuthView>("login");
	const [loading, startTransition] = useTransition();
	const [initialValues, setInitialValues] = useState<
		EmailAuthValues | undefined
	>({
		email: VITE_DEFAULT_USER_EMAIL,
		password: VITE_DEFAULT_USER_PASSWORD,
	});

	const handleFormSubmit = async (params: EmailAuthValues) => {
		const onLogin = async (token?: string) => {
			if (!token) return;
			setAuthorization(token);
			router.replace(getRedirectUrl());
			router.refresh();
		};
		startTransition(async () => {
			try {
				switch (type) {
					case "login": {
						if (!params.email || !params.password) {
							return;
						}

						const { data: login } = await userLogin({
							body: {
								identifier: params.identifier ?? params.email ?? "",
								email: params.email ?? "",
								password: params.password ?? "",
								cf_token: params.cf_token ?? "",
								IP: "",
								LoginType: "email",
								UserAgent: "",
							},
						});
						toast.success(t("login.success"));
						onLogin(login?.token);
						break;
					}
					case "register": {
						if (!params.email || !params.password) {
							return;
						}

						const { data: create } = await userRegister({
							body: {
								identifier: params.identifier ?? params.email ?? "",
								email: params.email ?? "",
								password: params.password ?? "",
								invite: params.invite ?? "",
								code: params.code ?? "",
								cf_token: params.cf_token ?? "",
								IP: "",
								LoginType: "email",
								UserAgent: "",
							},
						});
						toast.success(t("register.success"));
						onLogin(create?.token);
						break;
					}
					case "reset":
						if (!params.email || !params.password) {
							return;
						}

						await resetPassword({
							body: {
								identifier: params.identifier ?? params.email ?? "",
								email: params.email ?? "",
								password: params.password ?? "",
								code: params.code ?? "",
								cf_token: params.cf_token ?? "",
								IP: "",
								LoginType: "email",
								UserAgent: "",
							},
						});
						toast.success(t("reset.success"));
						setType("login");
						break;
				}
			} catch (_error) {
				/* empty */
			}
		});
	};

	let UserForm: ReactNode = null;
	switch (type) {
		case "login":
			UserForm = (
				<LoginForm
					loading={loading}
					onSubmit={handleFormSubmit}
					initialValues={initialValues}
					setInitialValues={setInitialValues}
					onSwitchForm={setType}
				/>
			);
			break;
		case "register":
			UserForm = (
				<RegisterForm
					loading={loading}
					onSubmit={handleFormSubmit}
					initialValues={initialValues}
					setInitialValues={setInitialValues}
					onSwitchForm={setType}
				/>
			);
			break;
		case "reset":
			UserForm = (
				<ResetForm
					loading={loading}
					onSubmit={handleFormSubmit}
					initialValues={initialValues}
					setInitialValues={setInitialValues}
					onSwitchForm={setType}
				/>
			);
			break;
	}

	return UserForm;
}
