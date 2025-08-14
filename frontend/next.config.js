// next.config.js
/** @type {import('next').NextConfig} */
export default {
  output: 'export', // Enable static HTML export
  distDir: 'build', // Output directory for the static files
  trailingSlash: true,
  reactStrictMode: true,
  eslint: {
    ignoreDuringBuilds: true,
  },
  turbopack: {
  },
  typescript: {
    ignoreBuildErrors: true,
  },
  images: {
    unoptimized: true, // Required for static export
  },
  pageExtensions: ['js', 'jsx', 'ts', 'tsx'],
  // Configure for static export to work with Go server
  assetPrefix: '',
  basePath: '',
};
