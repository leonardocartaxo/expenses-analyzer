import path from 'node:path';
import type { NextConfig } from 'next';

const nextConfig: NextConfig = {
  output: 'standalone',
  outputFileTracingRoot: path.join(__dirname, '../..'),
  // Next 16.3.x + Node 24: standalone tracer omits @swc/helpers ESM files that
  // module-sync resolves at runtime under pnpm (upstream: vercel/next.js#97372).
  outputFileTracingIncludes: {
    '/**': ['../../node_modules/.pnpm/**/node_modules/@swc/helpers/esm/**/*'],
  },
  transpilePackages: ['@expenses/api-client'],
};

export default nextConfig;
