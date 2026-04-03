"use client";

import { useRouter } from "next/navigation";
import { useTranslations } from "next-intl";
import { type ReactNode, useState, useTransition } from "react";
import { toast } from "sonner";
import {
  telephoneLogin,
  telephoneResetPassword,
  telephoneUserRegister,
} from "@/services/user-api/sdk.gen";
import { getRedirectUrl, setAuthorization } from "@/utils/common";
import LoginForm from "./login-form";
import RegisterForm from "./register-form";
import ResetForm from "./reset-form";
import type { AuthView, PhoneAuthValues } from "./types";

export default function PhoneAuthForm() {
  const t = useTranslations("auth");
  const router = useRouter();
  const [type, setType] = useState<AuthView>("login");
  const [loading, startTransition] = useTransition();
  const [initialValues, setInitialValues] = useState<PhoneAuthValues | undefined>({
    identifier: "",
    telephone: "",
    telephone_area_code: "1",
    password: "",
    telephone_code: "",
  });

  const handleFormSubmit = async (params: PhoneAuthValues) => {
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
            if (!params.telephone || !params.telephone_area_code) {
              return;
            }

            const { data: login } = await telephoneLogin({
              body: {
                identifier: params.identifier ?? params.telephone ?? "",
                telephone: params.telephone ?? "",
                telephone_area_code: params.telephone_area_code ?? "",
                telephone_code: params.telephone_code ?? "",
                password: params.password ?? "",
                cf_token: params.cf_token ?? "",
                IP: "",
                LoginType: "telephone",
                UserAgent: "",
              },
            });
            toast.success(t("login.success"));
            onLogin(login?.token);
            break;
          }
          case "register": {
            if (!params.telephone || !params.telephone_area_code || !params.password) {
              return;
            }

            const { data: create } = await telephoneUserRegister({
              body: {
                identifier: params.identifier ?? params.telephone ?? "",
                telephone: params.telephone ?? "",
                telephone_area_code: params.telephone_area_code ?? "",
                password: params.password ?? "",
                invite: params.invite ?? "",
                code: params.code ?? "",
                cf_token: params.cf_token ?? "",
                IP: "",
                LoginType: "telephone",
                UserAgent: "",
              },
            });
            toast.success(t("register.success"));
            onLogin(create?.token);
            break;
          }
          case "reset":
            if (!params.telephone || !params.telephone_area_code || !params.password) {
              return;
            }

            await telephoneResetPassword({
              body: {
                identifier: params.identifier ?? params.telephone ?? "",
                telephone: params.telephone ?? "",
                telephone_area_code: params.telephone_area_code ?? "",
                password: params.password ?? "",
                code: params.code ?? "",
                cf_token: params.cf_token ?? "",
                IP: "",
                LoginType: "telephone",
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
