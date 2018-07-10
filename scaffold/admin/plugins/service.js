import Vue from 'vue'
import Service from '~/helpers/service.js'

export default function ({app}) {
  Vue.use(Service, {
    http: app.$axios,
    store: app.store
  })
}
