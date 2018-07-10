<template>
  <main class="simple main-body u-fh">
    <nuxt />
  </main>
</template>

<script>

export default {
  // COMPONENT
  // ______________________________________
  head () {
    return {
      title: 'Login - '
    }
  },
  components: {
  },
  props: {},
  computed: {},
  methods: {},
  watch: {},
  data () {
    return {
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
  },
  beforeUpdate () {
  },
  updated () {
  },
  beforeDestroy () {
  },
  destroyed () {
  }
}
</script>

<style lang="scss" >
  @import "~public/scss/admin.scss";
  @import "~public/scss/admin-overrides.scss";
  @import "~public/scss/utility.scss";
  @import "~public/scss/_variables";
  @import "~public/scss/utility";

  .simple.main-body {
    background: $sidebar-color;
    .logo {
      padding: 2rem 4rem;
    }
  }
</style>
