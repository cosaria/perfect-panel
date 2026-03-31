"use client";

import { Table, TableBody, TableCell, TableRow } from "@workspace/ui/components/table";
import { useTranslations } from "next-intl";
import AppleForm from "./forms/apple-form";
import DeviceForm from "./forms/device-form";
import EmailSettingsForm from "./forms/email-settings-form";
import FacebookForm from "./forms/facebook-form";
import GithubForm from "./forms/github-form";
import GoogleForm from "./forms/google-form";
import PhoneSettingsForm from "./forms/phone-settings-form";
import TelegramForm from "./forms/telegram-form";

export default function Page() {
  const t = useTranslations("auth-control");

  const formSections = [
    {
      title: t("communicationMethods"),
      forms: [
        { key: "email-settings", component: EmailSettingsForm },
        { key: "phone-settings", component: PhoneSettingsForm },
      ],
    },
    {
      title: t("socialAuthMethods"),
      forms: [
        { key: "apple", component: AppleForm },
        { key: "google", component: GoogleForm },
        { key: "facebook", component: FacebookForm },
        { key: "github", component: GithubForm },
        { key: "telegram", component: TelegramForm },
      ],
    },
    {
      title: t("deviceAuthMethods"),
      forms: [{ key: "device", component: DeviceForm }],
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
