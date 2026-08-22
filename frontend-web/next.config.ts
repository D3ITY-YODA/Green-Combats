// next.config.ts

import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  // Enable React Strict Mode for better development experience and catching bugs early
  reactStrictMode: true,

  // Standalone output for Docker deployments
  output: "standalone",

  // Image optimization configuration
  images: {
    remotePatterns: [
      {
        protocol: "https",
        hostname: "**", // Allow images from any HTTPS source (adjust if specific CDN is used)
      },
      {
        protocol: "http",
        hostname: "localhost",
        port: "8080", // Allow images from local Go backend during development
        pathname: "/api/v1/**",
      },
    ],
    // Improve performance by using modern image formats
    formats: ["image/avif", "image/webp"],
  },

  // Security headers for production
  async headers() {
    return [
      {
        source: "/(.*)",
        headers: [
          {
            key: "X-DNS-Prefetch-Control",
            value: "on",
          },
          {
            key: "Strict-Transport-Security",
            value: "max-age=63072000; includeSubDomains; preload",
          },
          {
            key: "X-Frame-Options",
            value: "SAMEORIGIN",
          },
          {
            key: "X-Content-Type-Options",
            value: "nosniff",
          },
          {
            key: "Referrer-Policy",
            value: "origin-when-cross-origin",
          },
        ],
      },
    ];
  },

  // Performance optimizations
  experimental: {
    // Optimize imports for large libraries to reduce bundle size
    optimizePackageImports: ["lucide-react", "@tanstack/react-query", "maplibre-gl"],
  },

  // Redirects for legacy routes or marketing alignment (optional)
  async redirects() {
    return [
      {
        source: "/home",
        destination: "/",
        permanent: true,
      },
    ];
  },
};

export default nextConfig;
