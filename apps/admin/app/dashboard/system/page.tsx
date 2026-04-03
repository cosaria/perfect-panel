"use client";

import { Table, TableBody, TableCell, TableRow } from "@workspace/ui/components/table";
import { useTranslations } from "next-intl";
import CurrencyForm from "./basic-settings/currency-form";
import PrivacyPolicyForm from "./basic-settings/privacy-policy-form";
import SiteForm from "./basic-settings/site-form";
import TosForm from "./basic-settings/tos-form";
import LogCleanupForm from "./log-cleanup/log-cleanup-form";
import InviteForm from "./user-security/invite-form";
import RegisterForm from "./user-security/register-form";
import VerifyCodeForm from "./user-security/verify-code-form";
import VerifyForm from "./user-security/verify-form";

export default function Page() {
  const t = useTranslations("system");

  const formSections = [
    {
      title: t("basicSettings"),
      forms: [
        { key: "site", component: SiteForm },
        { key: "currency", component: CurrencyForm },
        { key: "tos", component: TosForm },
        { key: "privacy-policy", component: PrivacyPolicyForm },
      ],
    },
    {
      title: t("userSecuritySettings"),
      forms: [
        { key: "register", component: RegisterForm },
        { key: "invite", component: InviteForm },
        { key: "verify", component: VerifyForm },
        { key: "verify-code", component: VerifyCodeForm },
      ],
    },
    {
      title: t("logSettings"),
      forms: [{ key: "log-cleanup", component: LogCleanupForm }],
    },
  ];

  return (
    <div className="space-y-8">
      {formSections.map((section) => (
        <div key={section.title}>
          <h2 className="mb-4 text-lg font-semibold">{section.title}</h2>
          <Table>
            <TableBody>
              {section.forms.map((form) => {
                const FormComponent = form.component;
                return (
                  <TableRow key={form.key}>
                    <TableCell>
                      <FormComponent />
                    </TableCell>
                  </TableRow>
                );
              })}
            </TableBody>
          </Table>
        </div>
      ))}
    </div>
  );
}
