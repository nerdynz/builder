<template>
  <div class="application">
    <sidebar :show="sidebarIsOpen" />
    <main :class="{'main-body': true, 'sidebar-open': sidebarIsOpen, 'nav-showing': !navbar.hidden}">
      <levelbar />
      <nuxt />
    </main>
    <div class="footer-push" />
    <footer-bar />
  </div>
</template>
<script>
import { mapGetters, mapActions } from 'vuex'
import Levelbar from '~/components/layout/Levelbar'
import Sidebar from '~/components/layout/Sidebar'
import FooterBar from '~/components/layout/FooterBar'

export default {
  // COMPONENT
  // ______________________________________
  name: 'AdminLayout',
  middleware: 'authenticate',
  components: {
    Levelbar,
    Sidebar,
    FooterBar
    // Navbar
  },
  props: {},
  data () {
    return {
    }
  },
  computed: {
    ...mapGetters({
      navbar: 'app/navbar',
      sidebarIsOpen: 'app/sidebarIsOpen'
    })
  },
  watch: {},
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
  },
  methods: {
    updateToken (token) {
      this.$axios.get('/api/v1/updatetoken?token=' + token)
    },
    ...mapActions({
      toggleDevice: 'app/toggleDevice',
      toggleSidebar: 'app/toggleSidebar',
      toggleNavbar: 'app/toggleNavbar'
    })
  }
}
</script>

<style lang="scss">
  @import "~/assets/main.scss";

</style>
