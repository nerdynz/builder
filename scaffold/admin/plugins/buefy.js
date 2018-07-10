import Vue from 'vue'
import * as Buefy from 'buefy/src/index.js'

export default function ({app}, inject) {
  Vue.use(Buefy.default, {
    defaultIconPack: 'far'
  })
  inject('dialog', Buefy.Dialog)
  inject('loading', Buefy.LoadingProgrammatic)
  inject('modal', Buefy.ModalProgrammatic)
  inject('snackbar', Buefy.Snackbar)
  inject('toast', Buefy.Toast)
}
