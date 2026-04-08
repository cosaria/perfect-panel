import { defaultPlugins, defineConfig } from "@hey-api/openapi-ts";

export default defineConfig({
  input: "../../../docs/openapi/common.json",
  output: {
    path: "./services/common-api",
    clean: true,
  },
  plugins: [
    ...defaultPlugins,
    "@hey-api/client-fetch",
    {
      name: "@hey-api/typescript",
      enums: "javascript",
    },
    {
      name: "@hey-api/sdk",
      auth: false,
    },
  ],
});
