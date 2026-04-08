import "@workspace/ui/globals.css";
import type React from "react";
import ClientRoot from "@/components/client-root";

export const metadata = {
  title: {
    default: "PPanel",
    template: "%s | PPanel",
  },
  description: "",
  icons: {
    icon: [
      { url: "/favicon.ico", sizes: "48x48" },
      { url: "/favicon.svg", type: "image/svg+xml" },
    ],
    apple: "/apple-touch-icon.png",
  },
  manifest: "/site.webmanifest",
};

export const viewport = {
  themeColor: [
    { media: "(prefers-color-scheme: light)", color: "#FFFFFF" },
    { media: "(prefers-color-scheme: dark)", color: "#000000" },
  ],
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return <ClientRoot>{children}</ClientRoot>;
}
