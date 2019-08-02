<template>
  <section class="app-navbar ">
    <nav class="navbar is-flex-mobile">
      <div class="navbar-start is-flex-mobile" style="flex-grow:0;padding: 1rem 0.5rem;">
        <a class="navbar-item is-hidden-1230" @click="toggleSidebar(!sidebar.opened)">
          <i class="fa fa-bars" aria-hidden="true"></i>
        </a>
      </div>
      <div class="navbar-end is-flex-mobile">
        <div class="navbar-item">
          <picture-avatar :user="user"></picture-avatar>
        </div>
        <div class="navbar-item logout">
          <router-link to="/logout" ><b-icon icon="sign-out" size="large"></b-icon></router-link>
        </div>
      </div>
    </nav>
  </section>
</template>

<script>
import { mapGetters, mapActions } from 'vuex'
import PictureAvatar from '~/components/PictureAvatar'
export default {

  components: {
    PictureAvatar
  },

  props: {
    show: Boolean
  },

  computed: {
    ...mapGetters({
      sidebar: 'app/sidebar',
      user: 'auth/user',
      title: 'app/title'
    })
  },

  methods: mapActions({
    toggleSidebar: 'app/toggleSidebar'
  })

}
</script>

<style lang="scss">
@import "~public/scss/_variables";
@import '~bulma/utilities/all';

.app-navbar {
  position: fixed;
  min-width: 100%;
  z-index: 1024;
  // box-shadow: 0 2px 3px rgba(17, 17, 17, 0.1), 0 0 0 1px rgba(17, 17, 17, 0.1);
  background-color: $sidebar-color;

  .container {
    margin: auto 10px;
  }

  .logo-wrap {
    padding: 0;
    img {
      max-width:300px;
    }
  }

  .navbar {
    height: 60px;
    box-shadow: none;

    .navbar-item {
      padding: 0 0.5rem;
    }
  }

  .navbar-end {
    align-items: center;
    flex: 1;
    justify-content: flex-end;
    overflow: hidden;
    white-space: nowrap;

    .picture {
      width: 46px;
      margin-top: 4px;
      margin-right: 10px;
      img {
        border-radius: 40px;
      }
    }

    .first-name, .logout {
      color: $dark;
      a {
        color: $dark;
      }
    }
    // padding-right: 40px;
  }
}

.hero-brand {
  .vue {
    margin-left: 10px;
    color: #36AC70;
  }
  .admin {
    color: #28374B;
  }
}

.logo-wrap {
  display: flex;
  align-items: center;
  justify-content: center;
}

.top-title {
  color: $white;
  font-weight: $weight-bold;
  font-size: 16px;
}
</style>
