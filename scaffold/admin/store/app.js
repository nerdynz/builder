export const state = () => ({
  isLoading: false,
  isDevelopment: (process.env.NODE_ENV === 'development'),
  sitename: 'nerdy',
  isElectron: typeof (window && window.process && window.process.type) !== 'undefined',
  device: {
    isMobile: false,
    isTablet: false
  },
  sidebar: {
    opened: true,
    hidden: false
  },
  navbar: {
    hidden: false
  },
  effect: {
    translate3d: true
  },
  buttons: [

  ],
  title: 'Generic Title',
  subtitle: '',
  breadcrumb: null,
  previousRoute: {
    name: '',
    path: ''
  },
  windowWidth: 1680,
  windowHeight: 1050,
  settings: {
  },
  fileQueue: {
  }
})

export const mutations = {
  BLOCK_HOVERED (state, isHovered) {
    state.blockIsHovered = isHovered
  },
  CURRENT_WINDOW_HEIGHT (state, h) {
    state.windowHeight = h
  },
  CURRENT_WINDOW_WIDTH (state, w) {
    state.windowWidth = w
  },
  SET_PREVIOUS_ROUTE (state, prevRoute) {
    if (prevRoute) {
      state.previousRoute = prevRoute
    } else {
      state.previousRoute = {
        name: '',
        path: ''
      }
    }
  },
  SET_IS_LOADING (state, isLoading) {
    state.isLoading = isLoading
  },

  TOGGLE_DEVICE (state, device) {
    state.device.isMobile = device === 'mobile'
    state.device.isTablet = device === 'tablet'
  },

  // SIDEBAR
  TOGGLE_SIDEBAR (state, opened) {
    state.sidebar.opened = opened
  },

  HIDE_SIDEBAR (state, isHidden) {
    state.sidebar.hidden = isHidden
  },

  // NAVBAR
  HIDE_NAVBAR (state, isHidden) {
    state.navbar.hidden = isHidden
  },

  SWITCH_EFFECT (state, effectItem) {
    for (let name in effectItem) {
      state.effect[name] = effectItem[name]
    }
  },

  SET_BUTTONS (state, buttonOptions) {
    state.buttons = buttonOptions
  },

  SET_TITLE (state, title) {
    state.title = title
  },

  SET_SUBTITLE (state, subtitle) {
    state.subtitle = subtitle
  },

  SET_SETTINGS (state, settings) {
    state.settings = settings
  },

  SET_BREADCRUMB (state, breadcrumb) {
    state.breadcrumb = breadcrumb
  },

  ADD_TO_UPLOADS (state, {ulid, data}) {
    state.fileQueue[ulid] = data
  },

  CLEAR_UPLOADS (state) {
    state.fileQueue = {
    }
  }
}

export const actions = {
  loadSettings ({ commit }) {
    this.$axios.get('/api/v1/settings/retrieve/1').then(({data}) => {
      commit('SET_SETTINGS', data)
    })
  },
  setButtons ({ commit }, buttons) {
    if (buttons) {
      commit('SET_BUTTONS', buttons)
    }
  },

  toggleSidebar ({ commit }, opened) {
    commit('TOGGLE_SIDEBAR', opened)
  },

  hideSidebar ({ commit }, isHidden) {
    commit('HIDE_SIDEBAR', isHidden)
  },

  toggleNavbar ({ commit }, opened) {
    commit('HIDE_NAVBAR', opened)
  },

  toggleDevice ({ commit }, device) {
    commit('TOGGLE_DEVICE', device)
  },

  setWindowHeight ({ commit }, h) {
    commit('SET_CURRENT_WINDOW_HEIGHT', h)
  },

  saveUploads (store, cb) {
    // delay slightly because we also delay updating the fileQueue so we dont do it too often as its expensive
    setTimeout(() => {
      let http = this.app.$axios
      Object.keys(store.state.fileQueue).forEach((key, index) => {
        let data = store.state.fileQueue[key]
        http.post('/api/v1/upload/crop', data).then((data) => {
          store.commit('CLEAR_UPLOADS')
          if (cb) cb(data)
        })
      })
    }, 200)
  }
}

export const getters = {
  settings (state) {
    return state.settings
  },
  breadcrumb (state) {
    return state.breadcrumb
  },
  current (state) {
    if (state.breadcrumb) {
      return state.breadcrumb[state.breadcrumb.length - 1]
    }
    return null
  },
  parent (state) {
    if (state.breadcrumb) {
      if (state.breadcrumb.length === 2) {
        // current is parent also
        return state.breadcrumb[state.breadcrumb.length - 1]
      }
      return state.breadcrumb[state.breadcrumb.length - 2]
    }
    return null
  },
  title (state) {
    if (state.breadcrumb && state.breadcrumb.length > 0) {
      return state.breadcrumb[state.breadcrumb.length - 1].title
    }
    if (state.title) {
      return state.title
    }
    return 'Generic Title'
  },
  subtitle (state) {
    return state.subtitle
  },
  navbar (state) {
    return state.navbar
  },
  sidebar (state) {
    return state.sidebar
  },
  sidebarIsOpen (state) {
    if (state.sidebar.hidden) {
      return false
    }
    return state.sidebar.opened
  },
  windowWidth (state) {
    return state.windowWidth
  },
  windowHeight (state) {
    return state.windowHeight
  },
  offsetWindowHeight (state) {
    return state.windowHeight - 150
  },
  leftButtons (state) {
    return state.buttons
      .filter(button => button.alignment === 'left' && (typeof (button.role) === 'undefined' || button.role === state.login.details.Person.Role))
  },
  rightButtons (state) {
    return state.buttons
      .filter(button => button.alignment === 'right' && (typeof (button.role) === 'undefined' || button.role === state.login.details.Person.Role))
  },
  centerButtons (state) {
    return state.buttons
      .filter(button => button.alignment === 'center' && (typeof (button.role) === 'undefined' || button.role === state.login.details.Person.Role))
  },
  specialButtons (state) {
    return state.buttons
      .filter(button => button.alignment === 'special' && (typeof (button.role) === 'undefined' || button.role === state.login.details.Person.Role))
  },
  isDev () {
    return (process.env.NODE_ENV === 'development')
  },
  sitename (state) {
    return state.sitename
  },
  logo (state) {
    let logo = state.settings.LogoPicture || ''
    if (logo && (process.env.NODE_ENV === 'development')) {
      logo = logo.replace('/attachments/', `https://cdn.nerdy.co.nz/attachments/${state.sitename}/`)
    }
    return logo
  }
}
