import Vue from 'vue'
import Menu from '~/helpers/menu.js'
export default function ({app}) {
  Vue.use(Menu, {
    router: app.router,
    store: app.store
  })
}
