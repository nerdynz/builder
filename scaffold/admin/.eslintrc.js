module.exports = {
  root: true,
  env: {
    browser: true,
    node: true,
    jquery: true
  },
  parserOptions: {
    parser: 'babel-eslint'
  },
  extends: [
    '@nuxtjs',
    'plugin:nuxt/recommended',
    'standard'
  ],
  env: {
    'jquery': true
  },
  // add your custom rules here
  rules: {
    'prefer-const': 0,
    'no-console': 0
  }
}
