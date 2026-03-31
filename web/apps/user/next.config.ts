import type { NextConfig } from 'next';
import createNextIntlPlugin from 'next-intl/plugin';

const withNextIntl = createNextIntlPlugin('./locales/request.ts');

const nextConfig: NextConfig = {
  transpilePackages: ['@workspace/ui'],
  output: 'standalone',
  images: {
    remotePatterns: [
      {
        protocol: 'https',
        hostname: '**.**.**',
      },
      {
        protocol: 'https',
        hostname: '**.**',
      },
    ],
  },
};

export default withNextIntl(nextConfig);
