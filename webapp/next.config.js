/** @type {import('next').NextConfig} */
const nextConfig = {
  output: 'standalone',
  env: {
    PROBE_API_URL: process.env.PROBE_API_URL || 'http://localhost:8080',
  },
}

module.exports = nextConfig
