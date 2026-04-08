import path from "node:path";
import react from "@vitejs/plugin-react";
import { defineConfig } from "vite";

const userRoot = __dirname;
const uiRoot = path.resolve(userRoot, "../../packages/ui/src");
const publicEnvKeys = [
  "NEXT_PUBLIC_API_URL",
  "NEXT_PUBLIC_CDN_URL",
  "NEXT_PUBLIC_DEFAULT_LANGUAGE",
  "NEXT_PUBLIC_DEFAULT_USER_EMAIL",
  "NEXT_PUBLIC_DEFAULT_USER_PASSWORD",
  "NEXT_PUBLIC_DISCORD_LINK",
  "NEXT_PUBLIC_EMAIL",
  "NEXT_PUBLIC_GITHUB_LINK",
  "NEXT_PUBLIC_HOME_LOCATION_COUNT",
  "NEXT_PUBLIC_HOME_SERVER_COUNT",
  "NEXT_PUBLIC_HOME_USER_COUNT",
  "NEXT_PUBLIC_INSTAGRAM_LINK",
  "NEXT_PUBLIC_LINKEDIN_LINK",
  "NEXT_PUBLIC_SITE_URL",
  "NEXT_PUBLIC_TELEGRAM_LINK",
  "NEXT_PUBLIC_TWITTER_LINK",
] as const;

export default defineConfig({
  plugins: [react()],
  define: Object.fromEntries(
    publicEnvKeys.map((key) => [`process.env.${key}`, JSON.stringify(process.env[key] ?? "")]),
  ),
  resolve: {
    alias: [
      {
        find: "@/",
        replacement: `${userRoot}/`,
      },
      {
        find: "@workspace/ui/globals.css",
        replacement: path.resolve(uiRoot, "styles/globals.css"),
      },
      {
        find: /^@workspace\/ui\/utils$/,
        replacement: path.resolve(uiRoot, "utils/index.ts"),
      },
      {
        find: /^@workspace\/ui\/components\/(.*)$/,
        replacement: `${path.resolve(uiRoot, "components")}/$1`,
      },
      {
        find: /^@workspace\/ui\/custom-components\/(.*)$/,
        replacement: `${path.resolve(uiRoot, "custom-components")}/$1`,
      },
      {
        find: /^@workspace\/ui\/hooks\/(.*)$/,
        replacement: `${path.resolve(uiRoot, "hooks")}/$1`,
      },
      {
        find: /^@workspace\/ui\/lib\/(.*)$/,
        replacement: `${path.resolve(uiRoot, "lib")}/$1`,
      },
      {
        find: /^@workspace\/ui\/lotties\/(.*)$/,
        replacement: `${path.resolve(uiRoot, "lotties")}/$1`,
      },
      {
        find: /^@workspace\/ui\/utils\/(.*)$/,
        replacement: `${path.resolve(uiRoot, "utils")}/$1`,
      },
      {
        find: "next/legacy/image",
        replacement: path.resolve(userRoot, "src/compat/next-image.tsx"),
      },
      {
        find: "next/image",
        replacement: path.resolve(userRoot, "src/compat/next-image.tsx"),
      },
      {
        find: "next/link",
        replacement: path.resolve(userRoot, "src/compat/next-link.tsx"),
      },
      {
        find: "next/navigation",
        replacement: path.resolve(userRoot, "src/compat/next-navigation.ts"),
      },
    ],
  },
  server: {
    host: "0.0.0.0",
    port: 3001,
  },
  build: {
    outDir: "dist",
    emptyOutDir: true,
  },
});
