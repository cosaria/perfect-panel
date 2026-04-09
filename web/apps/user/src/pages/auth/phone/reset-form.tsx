import { zodResolver } from "@hookform/resolvers/zod";
import { Button } from "@workspace/ui/components/button";
import {
	Form,
	FormControl,
	FormField,
	FormItem,
	FormMessage,
} from "@workspace/ui/components/form";
import { useTranslations } from "@workspace/ui/components/i18n-provider";
import { Input } from "@workspace/ui/components/input";
import { AreaCodeSelect } from "@workspace/ui/custom-components/area-code-select";
import { Icon } from "@workspace/ui/custom-components/icon";
import { type Dispatch, type SetStateAction, useRef, useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";
import useGlobalStore from "@/config/use-global";
import SendCode from "../send-code";
import CloudFlareTurnstile, { type TurnstileRef } from "../turnstile";
import type { AuthView, PhoneAuthValues } from "./types";

export default function ResetForm({
	loading,
	onSubmit,
	initialValues,
	onSwitchForm,
}: {
	loading?: boolean;
	onSubmit: (data: PhoneAuthValues) => void;
	initialValues: PhoneAuthValues | undefined;
	setInitialValues: Dispatch<SetStateAction<PhoneAuthValues | undefined>>;
	onSwitchForm: Dispatch<SetStateAction<AuthView>>;
}) {
	const t = useTranslations("auth.reset");

	const { common } = useGlobalStore();
	const { verify, auth } = common;

	const [_targetDate, _setTargetDate] = useState<number>();

	const formSchema = z.object({
		telephone_area_code: z.string(),
		telephone: z.string(),
		password: z.string(),
		code: auth?.email?.enable_verify ? z.string() : z.string().nullish(),
		cf_token:
			verify.enable_register_verify && verify.turnstile_site_key
				? z.string()
				: z.string().nullish(),
	});
	const form = useForm<z.infer<typeof formSchema>>({
		resolver: zodResolver(formSchema),
		defaultValues: initialValues,
	});

	const turnstile = useRef<TurnstileRef>(null);
	const handleSubmit = form.handleSubmit((data) => {
		try {
			onSubmit({
				...data,
				code: data.code ?? undefined,
				cf_token: data.cf_token ?? undefined,
			});
		} catch (_error) {
			turnstile.current?.reset();
		}
	});

	return (
		<>
			<Form {...form}>
				<form onSubmit={handleSubmit} className="grid gap-6">
					<FormField
						control={form.control}
						name="telephone"
						render={({ field }) => (
							<FormItem>
								<FormControl>
									<div className="flex">
										<FormField
											control={form.control}
											name="telephone_area_code"
											render={({ field }) => (
												<FormItem>
													<FormControl>
														<AreaCodeSelect
															simple
															className="w-32 rounded-r-none border-r-0"
															placeholder="Area code..."
															value={field.value}
															onChange={(value) => {
																if (value.phone) {
																	form.setValue(
																		"telephone_area_code",
																		value.phone,
																	);
																}
															}}
														/>
													</FormControl>
													<FormMessage />
												</FormItem>
											)}
										/>
										<Input
											className="rounded-l-none"
											placeholder="Enter your telephone..."
											type="tel"
											{...field}
										/>
									</div>
								</FormControl>
								<FormMessage />
							</FormItem>
						)}
					/>
					<FormField
						control={form.control}
						name="code"
						render={({ field }) => (
							<FormItem>
								<FormControl>
									<div className="flex items-center gap-2">
										<Input
											placeholder="Enter code..."
											type="text"
											{...field}
											value={field.value as string}
										/>
										<SendCode
											type="phone"
											params={{
												...form.getValues(),
												type: 2,
											}}
										/>
									</div>
								</FormControl>
								<FormMessage />
							</FormItem>
						)}
					/>
					<FormField
						control={form.control}
						name="password"
						render={({ field }) => (
							<FormItem>
								<FormControl>
									<Input
										placeholder="Enter your new password..."
										type="password"
										{...field}
									/>
								</FormControl>
								<FormMessage />
							</FormItem>
						)}
					/>
					{verify.enable_reset_password_verify && (
						<FormField
							control={form.control}
							name="cf_token"
							render={({ field }) => (
								<FormItem>
									<FormControl>
										<CloudFlareTurnstile
											id="reset"
											{...field}
											ref={turnstile}
										/>
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>
					)}
					<Button type="submit" disabled={loading}>
						{loading && <Icon icon="mdi:loading" className="animate-spin" />}
						{t("title")}
					</Button>
				</form>
			</Form>
			<div className="mt-4 text-right text-sm">
				{t("existingAccount")}&nbsp;
				<Button
					variant="link"
					className="p-0"
					onClick={() => {
						onSwitchForm("login");
					}}
				>
					{t("switchToLogin")}
				</Button>
			</div>
		</>
	);
}
