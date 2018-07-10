<template>
  <aside class="menu app-sidebar animated" :class="{ slideInLeft: show, slideOutLeft: !show, 'push-down': !navbar.hidden}">
    <div class="logo-wrap u-mb">
      <img class="logo" src="~public/images/logo.png" alt="">
    </div>
    <picture-avatar v-show="navbar.hidden" :user="user"></picture-avatar>
    <ul class="menu-list u-mt" >
      <li v-if="checkAllowed(item)" v-for="(item, index) in menu" :key="index">
        <nuxt-link v-if="item.path" :to="item.path" :exact="true" :aria-expanded="(isExpanded(item) ? 'true' : 'false')" >
          <span class="icon is-small"><i :class="['far', item.icon]" ></i></span>
          <span class="menu-label-text">{{getLabel(item)}}</span>
          <span class="icon is-small is-angle" v-if="hasChildren(item)">
            <i class="fa fa-angle-down"></i>
          </span>
        </nuxt-link>
        <div v-else-if="item.label" class="menu-label">
          {{item.label}}
        </div>

        <a :aria-expanded="isExpanded(item)" v-else @click="toggle(index, item)">
          <span class="icon is-small"><i :class="['far', item.icon]"></i></span>
          <span class="menu-label-text">{{getLabel(item)}}</span>
          <span class="icon is-small is-angle" v-if="item.children && item.children.length">
            <i class="far fa-angle-down"></i>
          </span>
        </a>

        <b-collapse :open="hasChildren(item)" >
          <ul v-show="isExpanded(item)">
            <li v-for="(subItem, index) in item.children" v-if="subItem.path && checkAllowed(subItem)" :key="index">
              <nuxt-link :to="generatePath(item, subItem)" >
                <span class="menu-label-text">{{getLabel(subItem)}}</span>
              </nuxt-link>
            </li>
          </ul>
        </b-collapse>
      </li>
    </ul>
  </aside>
</template>

<script>
import PictureAvatar from '~/components/PictureAvatar'
import { mapGetters, mapActions } from 'vuex'

export default {
  components: {
    PictureAvatar
  },

  props: {
    show: Boolean
  },

  data () {
    return {
      isReady: false
    }
  },

  mounted () {
    let route = this.$route
    if (route.name) {
      this.isReady = true
      this.shouldExpandMatchItem(route)
    }
  },

  computed: {
    firstName: function () {
      if (this.user.Name) {
        return this.user.Name.split(' ')[0]
      }
    },
    ...mapGetters({
      navbar: 'app/navbar',
      user: 'auth/user'
      // userIsManager: 'userIsManager'
    }),
    menu () {
      return this.$menu.items
    }
  },

  methods: {
    ...mapActions({
      expandMenu: 'menu/EXPAND_MENU'
    }),

    checkAllowed (item) {
      let allowed = true
      if (item) {
        if (item.isRouteOnly) {
          return false
        }
        if (item.role) {
          allowed = item.role === this.user.Role
        }
        if (allowed && typeof (item.isAllowed) === 'function') {
          return item.isAllowed(this)
        }
        if (typeof (item.path) === 'undefined' && typeof (item.label) === 'undefined') {
          return false
        }
      }
      return allowed
    },
    getLabel (item) {
      if (item.title) {
        if (typeof (item.title) === 'function') {
          return item.title(this)
        }
        return item.title
      }
    },

    isExpanded (item) {
      return item.expanded
    },

    hasChildren (item) {
      let hasChildren = false

      // check if there are any childitems, children can specify isDisplayed = false to prevent being rendered
      if (Array.isArray(item.children)) {
        for (let i = 0; i < item.children.length; i++) {
          let childItem = item.children[i]
          if (hasChildren) {
            break // if we have children at any point then we know we need to show expand
          }
          hasChildren = true
          // console.log('childItem.isDisplayed', childItem.isDisplayed)
          if (childItem.routeOnly) {
            // console.log('hitting this a lot')
            hasChildren = false
          }
        }
      }
      // console.log(item.path, ' has children: ', hasChildren)
      return hasChildren
    },

    toggle (index, item) {
      this.expandMenu({
        index: index,
        expanded: !item.expanded
      })
    },

    shouldExpandMatchItem (route) {
      let matched = route.matched
      let lastMatched = matched[matched.length - 1]
      if (lastMatched) {
        let parent = lastMatched.parent || lastMatched
        const isParent = parent === lastMatched

        if (isParent) {
          const p = this.findParentFromMenu(route)
          if (p) {
            parent = p
          }
        }

        if ('expanded' in parent && !isParent) {
          this.expandMenu({
            item: parent,
            expanded: true
          })
        }
      }
    },

    generatePath (item, subItem) {
      return `${item.component ? item.path + '/' : ''}${subItem.path}`
    },

    findParentFromMenu (route) {
      let items = this.$menu.items
      for (let i = 0, l = items.length; i < l; i++) {
        let item = items[i]
        let k = item.children && item.children.length
        if (k) {
          for (let j = 0; j < k; j++) {
            if (item.children[j].name === route.name) {
              return item
            }
          }
        }
      }
    }
  },

  watch: {
    $route (route) {
      this.isReady = true
      this.shouldExpandMatchItem(route)
    }
  }

}
</script>

<style lang="scss">
@import "~public/scss/_variables";

.app-sidebar {
  position: fixed;
  top: 0;
  left: 0;
  bottom: 0;
  padding: 5px 0 50px;
  width: $sidebar-width;
  min-width: 45px;
  max-height: 100vh;
  // height: calc(100% - 50px);
  height: 100%;
  z-index: 501;
  background: $sidebar-color;
  box-shadow: 0 2px 3px rgba(17, 17, 17, 0.1), 0 0 0 1px rgba(17, 17, 17, 0.1);
  overflow-y: auto;
  overflow-x: hidden;
  color: invert($sidebar-color);
  // @include mobile() {
  //   transform: translate3d(-180px, 0, 0);
  // }

  .logo {
    padding: 1rem;
  }
  .icon {
    vertical-align: baseline;
    &.is-angle {
      position: absolute;
      right: 10px;
      transition: transform .377s ease;
    }
  }

  .menu-label {
    padding: 15px 5px 3px;
    color: rgba(252, 252, 255, 0.62);
  }

  .menu-list {
    a {
      border-radius: 0;
      &:hover {
        background: lighten($grey, 3);
      }
    }
    li a {
      &[aria-expanded="true"] {
        .is-angle {
          transform: rotate(180deg);
        }
      }
    }

    li a + ul {
      margin: 0 10px 0 15px;
    }

    .menu-label-text {
      padding-left: 5px;
    }
  }

  .current-client {
    min-height: 80px !important;
  }
}

.menu-list {
  a {
    color: invert($sidebar-color);
    font-weight: $weight-semibold;
    &.is-active {
      background-color: rgba(243, 243, 243, 0.1);
      color: invert($sidebar-color);
    }
  }
}

.logged-in-user {
  display: flex;
  align-items: center;
  justify-content:center;  
  margin-top: -1rem;
  padding: 1rem 0;

  .picture {
    width: 46px;
    margin-top: 4px;
    margin-right: 10px;
    img {
      border-radius: 40px;
    }
  }

  .logout {
    color: $white;
  }
  .first-name, .logout {
    margin: 0 0.5rem;
    font-weight: $weight-semibold;
    color: $white;
    a {
      color: $white;
    }
  }
}

</style>
