/** @type {import('next').NextConfig} */
const nextConfig = {
  // Clean configuration for noVNC 1.5.0
  output: 'standalone',
  async rewrites() {
    return [
      { source: '/vm/:path*', destination: 'http://backend:8080/vm/:path*' }
    ];
  }
}

module.exports = nextConfig 