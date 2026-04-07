import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  transpilePackages: ["@workspace/ui"],
  output: "export",
  basePath: "/admin",
  images: {
    unoptimized: true,
  },
  webpack: (config) => {
    if (!config.resolve) config.resolve = {};

    // @hey-api/openapi-ts generates ESM imports with .js extensions
    // (e.g., './client.gen.js') but actual files are .ts.
    config.resolve.extensionAlias = {
      ".js": [".ts", ".tsx", ".js", ".jsx"],
    };

    // monaco-themes exports field doesn't include ./themes/*,
    // but the JSON files exist on disk. Bypass exports field check.
    const alias = (config.resolve.alias || {}) as Record<string, string>;
    alias["monaco-themes/themes"] = `${process.cwd()}/../../node_modules/monaco-themes/themes`;
    config.resolve.alias = alias;

    return config;
  },
};

export default nextConfig;
