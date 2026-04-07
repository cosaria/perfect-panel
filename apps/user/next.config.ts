import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  transpilePackages: ["@workspace/ui"],
  output: "export",
  images: {
    unoptimized: true,
  },
  webpack: (config) => {
    if (!config.resolve) config.resolve = {};
    config.resolve.extensionAlias = {
      ".js": [".ts", ".tsx", ".js", ".jsx"],
    };
    const alias = (config.resolve.alias || {}) as Record<string, string>;
    alias["monaco-themes/themes"] = `${process.cwd()}/../../node_modules/monaco-themes/themes`;
    config.resolve.alias = alias;
    return config;
  },
};

export default nextConfig;
