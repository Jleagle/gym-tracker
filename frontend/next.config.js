module.exports = {
  compress: true,
  future: {
    webpack5: true,
  },
  env: {
    GYMTRACKER_PORT_BACKEND: process.env.GYMTRACKER_PORT_BACKEND,
  },
  async redirects() {
    return [
      {
        source: '/',
        destination: '/all',
        permanent: true,
      },
    ]
  },
}
