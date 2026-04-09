import path from "node:path";
import react from "@vitejs/plugin-react";
import { defineConfig } from "vite";

const userRoot = __dirname;
const uiRoot = path.resolve(userRoot, "../../packages/ui/src");
const publicEnvKeys = [
  "VITE_API_URL",
  "VITE_CDN_URL",
  "VITE_DEFAULT_LANGUAGE",
  "VITE_DEFAULT_USER_EMAIL",
  "VITE_DEFAULT_USER_PASSWORD",
  "VITE_DISCORD_LINK",
  "VITE_EMAIL",
  "VITE_GITHUB_LINK",
  "VITE_HOME_LOCATION_COUNT",
  "VITE_HOME_SERVER_COUNT",
  "VITE_HOME_USER_COUNT",
  "VITE_INSTAGRAM_LINK",
  "VITE_LINKEDIN_LINK",
  "VITE_SITE_URL",
  "VITE_TELEGRAM_LINK",
  "VITE_TWITTER_LINK",
] as const;

export default defineConfig({
  plugins: [react()],
  define: Object.fromEntries(
    publicEnvKeys.map((key) => [`import.meta.env.${key}`, JSON.stringify(process.env[key] ?? "")]),
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
