import { titleCase } from '~/helpers/format'
import menu from '~/pages/routes'

export default class Menu {
  setMenuProps (store, to) {
    let currentRoute = to
    let bread = []
    if (currentRoute.path !== '/') {
      bread.push(this.items.index)
    }
    // untested for level 2....
    let goDeep = function (menu, parts, bread, level) {
      let routeName = parts.join('-')
      let keyarr = []
      for (let i = 0; i <= level; i++) {
        keyarr.push(parts[i])
      }
      let key = keyarr.join('-')
      let menuItem = { title: 'not-in-menu', name: 'not-in-menu', path: '/not-in-menu' }
      if (menu[routeName]) {
        menuItem = menu[routeName]
      } else if (menu[key]) {
        menuItem = menu[key]
      }
      let foundItem = true
      if (level > 0) {
        if (menu.children[key]) {
          menuItem = menu.children[key]
        } else {
          menuItem = menu // use previous one and retry match with higher level of parts. i.e. skip page-id // confusing but allows to just simply have one level of nested children
          foundItem = false
        }
      }

      if (foundItem) {
        bread.push(menuItem)
      }
      if (level < parts.length - 1 && menuItem.children) {
        goDeep(menuItem, parts, bread, (level + 1))
      }
    }
    if (currentRoute.name) {
      let parts = currentRoute.name.split('-')
      goDeep(this.items, parts, bread, 0)
    }
    this.bread = bread
    if (bread && bread.length > 0) {
      store.commit('app/SET_BREADCRUMB', this.bread)
    }
  }
  static install (Vue, { router, store }) {
    let m = new Menu()
    m.items = menu() // load from static menu js
    addPaths(m.items, router, 0)
    // 4. add an instance method
    Vue.prototype.$menu = m
    router.beforeEach((to, from, next) => {
      m.setMenuProps(store, to)
      next()
    })
  }
}

function addPaths (menuObj, router, level) {
  Object.keys(menuObj).forEach((key, index) => {
    let menuItem = menuObj[key] || {}
    let route = router.options.routes.find((route) => {
      return route.name === key
    })
    if (route && !menuItem.path) {
      menuItem.path = route.path
    }
    if (route && !menuItem.title) {
      menuItem.title = titleCase(route.name)
    }
    if (route && !menuItem.name) {
      menuItem.name = route.name
    }
    if (!menuItem.icon) {
      menuItem.icon = 'fa-circle-o'
    }
    if (menuItem.children) {
      addPaths(menuItem.children, router, (level + 1))
    }
  })
}
