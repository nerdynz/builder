const path = require('path')
const webpack = require('webpack')
const projectRoot = path.resolve(__dirname, '../')

module.exports = {
  /*
  ** Headers of the page
  */
  head: {
    titleTemplate: '%s [[[SITENAME]]]',
    meta: [
      { charset: 'utf-8' },
      { name: 'viewport', content: 'width=device-width, initial-scale=1' },
      { hid: 'description', name: 'description', content: '[[[SITENAME]]] Admin' }
    ],
    link: [
      { rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' }
    ]
  },
  mode: 'spa',
  modules: [
    '@nuxtjs/axios',
    '@nuxtjs/proxy'
  ],
  router: {
    // mode: 'hash'
    base: '/admin'
  },
  plugins: [
    '~/plugins/buefy',
    '~/plugins/custom-controls',
    '~/plugins/custom-router',
    '~/plugins/axios-interceptors',
    '~/plugins/service',
    '~/plugins/menu'
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
      config.resolve.alias['bulma'] = path.resolve(projectRoot, 'public/scss/bulma')
      config.resolve.alias['public'] = path.resolve(projectRoot, 'public')
      if (ctx.isDev && ctx.isClient) {
        config.module.rules.push({
          enforce: 'pre',
          test: /\.(js|vue)$/,
          loader: 'eslint-loader',
          exclude: /(node_modules)/
        })
      }
    },
    plugins: [
      new webpack.ProvidePlugin({
        $: 'jquery',
        jquery: 'jquery',
        'window.jQuery': 'jquery',
        jQuery: 'jquery'
      })
    ]
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
    }]
  ]
}
