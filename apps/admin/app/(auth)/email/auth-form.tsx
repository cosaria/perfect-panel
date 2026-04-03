"use client";

import { useRouter } from "next/navigation";
import { useTranslations } from "next-intl";
import { type ReactNode, useState, useTransition } from "react";
import { toast } from "sonner";
import {
  NEXT_PUBLIC_DEFAULT_USER_EMAIL,
  NEXT_PUBLIC_DEFAULT_USER_PASSWORD,
} from "@/config/constants";
import { resetPassword, userLogin, userRegister } from "@/services/user-api/sdk.gen";
import type {
  UserLoginRequest,
  UserRegisterRequest,
  ResetPasswordRequest,
} from "@/services/user-api/types.gen";
import { getRedirectUrl, setAuthorization } from "@/utils/common";
import LoginForm from "./login-form";
import RegisterForm from "./register-form";
import ResetForm from "./reset-form";
import type { AuthView, EmailAuthValues } from "./types";

export default function EmailAuthForm() {
  const t = useTranslations("auth");
  const router = useRouter();
  const [type, setType] = useState<AuthView>("login");
  const [loading, startTransition] = useTransition();
  const [initialValues, setInitialValues] = useState<EmailAuthValues | undefined>({
    email: NEXT_PUBLIC_DEFAULT_USER_EMAIL,
    password: NEXT_PUBLIC_DEFAULT_USER_PASSWORD,
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
            const { data: login } = await userLogin({ body: params as UserLoginRequest });
            toast.success(t("login.success"));
            onLogin(login?.token);
            break;
          }
          case "register": {
            const { data: created } = await userRegister({ body: params as UserRegisterRequest });
            toast.success(t("register.success"));
            onLogin(created?.token);
            break;
          }
          case "reset":
            await resetPassword({ body: params as ResetPasswordRequest });
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

  return (
    <>
      <div className="mb-11 text-center">
        <h1 className="mb-3 text-2xl font-bold">{t(`${type || "check"}.title`)}</h1>
        <div className="text-muted-foreground font-medium">
          {t(`${type || "check"}.description`)}
        </div>
      </div>
      {UserForm}
    </>
  );
}
