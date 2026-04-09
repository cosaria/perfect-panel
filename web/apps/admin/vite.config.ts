import path from "node:path";
import react from "@vitejs/plugin-react";
import { defineConfig } from "vite";

const adminRoot = __dirname;
const uiRoot = path.resolve(adminRoot, "../../packages/ui/src");
const publicEnvKeys = [
  "VITE_DEFAULT_LANGUAGE",
  "VITE_SITE_URL",
  "VITE_API_URL",
  "VITE_ADMIN_PATH",
  "VITE_DEFAULT_USER_EMAIL",
  "VITE_DEFAULT_USER_PASSWORD",
] as const;

export default defineConfig({
  plugins: [react()],
  base: "/admin/",
  define: Object.fromEntries(
    publicEnvKeys.map((key) => [`import.meta.env.${key}`, JSON.stringify(process.env[key] ?? "")]),
  ),
  resolve: {
    alias: [
      {
        find: "@/",
        replacement: `${adminRoot}/`,
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
    port: 3000,
  },
  build: {
    outDir: "dist",
    emptyOutDir: true,
  },
});
