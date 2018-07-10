import Vue from 'vue'
import VeeValidate from 'vee-validate'
import Tooltip from 'vue-directive-tooltip'
import 'vue-directive-tooltip/css/index.css'
import Control from '~/components/Control'
import Field from '~/components/Field'
import Multiselect from 'vue-multiselect'
import ImageUpload from '~/components/ImageUpload'
import RichText from '~/components/RichText'
import VueFroala from '~/components/froala-editor'
import Sortable from 'sortablejs'

Vue.use(VeeValidate, {
  inject: false
})
Vue.component('control', Control)
Vue.component('field', Field)
Vue.component('multiselect', Multiselect)
Vue.component('image-upload', ImageUpload)
Vue.component('rich-text', RichText)
Vue.use(require('vue-prevent-parent-scroll'))
Vue.use(Tooltip)
Vue.use(VueFroala, {
  defaultConfig: {
    key: 'kjc1tcqE5idF4ct=='
  }
})
Vue.directive('sortable', {
  inserted: function (el, binding) {
    new Sortable(el, binding.value || {}) // eslint-disable-line no-new
  }
})
require('~/helpers/imageupload')
