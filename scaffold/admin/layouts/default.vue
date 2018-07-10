<template>
  <div class="u-fh">
    <navbar v-show="!navbar.hidden"></navbar>
    <sidebar :show="sidebarIsOpen"></sidebar>
    <main :class="{'main-body': true, 'sidebar-open': sidebarIsOpen, 'nav-showing': !navbar.hidden}">
      <levelbar />
      <nuxt />
    </main>
    <footer-bar></footer-bar>
  </div>
</template>
<script>
import { mapGetters, mapActions } from 'vuex'
import Levelbar from '~/components/layout/Levelbar'
import Sidebar from '~/components/layout/Sidebar'
import Navbar from '~/components/layout/Navbar'
import FooterBar from '~/components/layout/FooterBar'

export default {
  // COMPONENT
  // ______________________________________
  name: 'AdminLayout',
  middleware: 'authenticate',
  components: {
    Levelbar,
    Sidebar,
    FooterBar,
    Navbar
  },
  props: {},
  computed: {
    ...mapGetters({
      navbar: 'app/navbar',
      sidebarIsOpen: 'app/sidebarIsOpen'
    })
  },
  methods: {
    ...mapActions({
      toggleDevice: 'app/toggleDevice',
      toggleSidebar: 'app/toggleSidebar',
      toggleNavbar: 'app/toggleNavbar'
    })
  },
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
    const { body } = document
    const WIDTH = 768
    const RATIO = 3
    const handler = () => {
      if (!document.hidden) {
        let rect = body.getBoundingClientRect()
        let isMobile = rect.width - RATIO < WIDTH
        this.toggleDevice(isMobile ? 'mobile' : 'other')
        this.toggleSidebar(!isMobile)
        this.toggleNavbar(!isMobile)
        this.$store.commit('app/CURRENT_WINDOW_WIDTH', window.innerWidth)
        this.$store.commit('app/CURRENT_WINDOW_HEIGHT', window.innerHeight)
      }
    }

    document.addEventListener('visibilitychange', handler)
    window.addEventListener('DOMContentLoaded', handler)
    window.addEventListener('resize', handler)
    setTimeout(() => {
      handler()
    }, 200)
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

<style lang="scss">
@import "~public/scss/admin.scss";
@import "~public/scss/admin-overrides.scss";
@import "~public/scss/content.scss";
@import "~public/scss/utility.scss";
@import 'vue-multiselect/dist/vue-multiselect.min.css';
.nuxt-progress {
  background-color: $green !important;
}

.main-body {
  @media (min-width: $tablet) {
    &.sidebar-open {
      padding-left: 250px;
    }
  }
  &.nav-showing {
    padding-top: 80px;
  }
}
</style>
