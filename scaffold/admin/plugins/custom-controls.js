import Vue from 'vue'
// import VeeValidate from 'vee-validate'
import VueDragDrop from 'vue-drag-drop'
import Multiselect from 'vue-multiselect'
import VueQuill from 'vue-quill'
import { fmtDate, fmtDateTime } from '~/helpers/format'
import { focusElement } from '~/helpers/helpers'

Vue.component('multiselect', Multiselect)
Vue.use(VueDragDrop)
Vue.use(VueQuill)

Vue.mixin({
  methods: {
    fmtDate,
    fmtDateTime,
    focusElement
  }
})
