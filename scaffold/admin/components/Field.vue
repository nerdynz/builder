<template>
  <div :class="{'is-danger': errors.has(name) }">
    <label v-if="label" :for="name" class="label thin-label">
    <span v-if="isRequired" class="required-star">*</span>
      {{label}}
    </label>
    <div class="u-rel">
      <span v-if="errors.has(name)" class="validation-message">
        {{errors.first(name)}}
      </span>
    </div>
    <div class="field" :class="{'has-addons': hasAddons}">
      <slot>
      </slot>
    </div>
  </div>
</template>

<script>
// ... imports

export default {
  // COMPONENT
  // ______________________________________
  name: 'Field',
  inject: {
    $validator: '$validator'
  },
  components: {},
  props: {
    for: {
      type: String,
      required: true
    },
    label: String,
    validator: Object,
    hasAddons: Boolean
  },
  computed: {
    name () {
      return this.for
    }
  },
  methods: {
    checkRequired (name) {
      if (name && typeof (this.$validator.fields) !== 'undefined') {
        var item = this.$validator.fields.items.find((item) => {
          return item.name === name
        })
        if (item && item.rules.hasOwnProperty('required')) {
          this.isRequired = true
          return
        }

        // if (!rule) {
        //   return false
        // }
        // rule = rule.toLowerCase()
        // return rule.indexOf('required') >= 0
      }
      this.isRequired = false
    }
  },
  watch: {
  },
  data () {
    return {
      isRequired: false
    }
  },

  // LIFECYCLE METHODS
  // ______________________________________
  beforeCreate () {
  },
  created () {
  },
  beforeMount () {
  },
  mounted () {
    let self = this
    this.$nextTick(() => {
      self.checkRequired(this.name)
    })
  },
  beforeUpdate () {
  },
  updated () {
    this.checkRequired(this.name)
  },
  beforeDestroy () {
  },
  destroyed () {
  }
}
</script>

<style lang="scss">

</style>
