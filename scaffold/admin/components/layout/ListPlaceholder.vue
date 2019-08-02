<template>
  <div class="list-placeholder rel-overlay">
    <div v-if="showing">
      <slot />
    </div>
    <div v-else-if="loading" class="u-rel">
      <b-loading :active="true" :can-cancel="true" />
    </div>
    <div v-else class="no-data" />
  </div>
</template>

<script>
// ... imports

export default {
  // COMPONENT
  // ______________________________________
  name: 'ListPlaceholder',
  components: {},
  props: {
    showing: Boolean,
    loading: Boolean
  },
  data () {
    return {
      hello: true
    }
  },
  computed: {},
  watch: {},
  methods: {
    add () {
      this.$emit('create')
    }
  }
}
</script>
