import createPersistedState from 'vuex-persistedstate'

export const state = () => ({
})

export const mutations = {
}

export const plugins = [
  createPersistedState({
    key: 'nerdy-vuex',
    paths: ['auth']
  })
]
