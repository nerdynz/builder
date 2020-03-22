const path = require('path') // eslint-disable-line
const webpack = require('webpack') // eslint-disable-line
const projectRoot = path.resolve(__dirname, '../')

export default {
  mode: 'spa',
  /*
  ** Headers of the page
  */
  head: {
    title: process.env.npm_package_name || '',
    meta: [
      { charset: 'utf-8' },
      { name: 'viewport', content: 'width=device-width, initial-scale=1' },
      { hid: 'description', name: 'description', content: process.env.npm_package_description || '' }
    ],
    link: [
      { rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' }
    ]
  },
  /*
  ** Customize the progress-bar color
  */
  loading: { color: '#fff' },
  /*
  ** Global CSS
  */
  css: [
  ],
  /*
  ** Plugins to load before mounting the App
  */
  plugins: [
    '~/plugins/custom-controls',
    '~/plugins/custom-router',
    '~/plugins/axios-interceptors',
    '~/plugins/menu'
  ],
  /*
  ** Nuxt.js modules
  */
  modules: [
    'specky-service',
    // Doc: https://buefy.github.io/#/documentation
    // Doc: https://axios.nuxtjs.org/usage
    ['nuxt-buefy', {
      defaultIconPack: 'far',
      css: false,
      materialDesignIcons: false,
    }],
    '@nuxtjs/axios',
    '@nuxtjs/proxy',
    '@nuxtjs/pwa',
    '@nuxtjs/eslint-module'
  ],
  /*
  ** Customize the progress bar color
  */
  loading: { color: '#3B8070' },
  /*
  ** Build configuration
  */
  build: {
    /*
    ** Run ESLint on save
    */
    extend (config, ctx) {
      config.resolve.alias['~public'] = path.resolve(projectRoot, 'public')

  },
  axios: {
    browserBaseURL: '/'
  },
  proxy: [
    ['/api/', {
      target: 'http://localhost:5000/api/',
      changeOrigin: true,
      pathRewrite: {
        '^/api': ''
      }
    }],
    ['/fonts/', {
      target: 'http://localhost:5000/fonts/',
      changeOrigin: true,
      pathRewrite: {
        '^/fonts': ''
      }
    }],
    ['/attachments/', {
      target: 'http://localhost:5000/attachments/',
      changeOrigin: true,
      pathRewrite: {
        '^/attachments': ''
      }
    }],
    ['/firebase-messaging-sw.js', {
      target: 'http://localhost:5000/',
      changeOrigin: true
    }]
  ]
}
