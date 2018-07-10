export default function ({app, store}, inject) {
  logRoutes(app)
  // add this to a router, make it automatically set bread and current on change
  app.router.beforeEach((to, from, next) => {
    if (to.path === '/logout') {
      app.$toast.open({ message: 'You have been logged out', type: 'is-info' })
      app.router.replace({ path: 'login' })
      store.commit('auth/LOGOUT')
    }

    store.commit('app/SET_BUTTONS', [])
    store.commit('app/HIDE_SIDEBAR', false)
    store.commit('app/SET_IS_LOADING', true)
    store.commit('app/SET_SUBTITLE', '')
    store.commit('app/SET_SUBTITLE', null)
    if (from && from.path && from.name) {
      store.commit('app/SET_PREVIOUS_ROUTE', {
        name: from.name,
        path: from.path
      })
    }
    next()
  })
}

function logRoutes (app) {
  var routeInfo = '=== APPLICATION ROUTES ===\n'
  app.router.options.routes.forEach((route) => {
    routeInfo += `${route.name} => ${route.path}\n`
  })
  console.log(routeInfo)
}
