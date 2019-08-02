<template>
  <aside class="menu app-sidebar animated" :class="{ slideInLeft: show, slideOutLeft: !show, 'push-down': !navbar.hidden}">
    <div v-if="navbar.hidden" class="logo-wrap">
      <div>
        <logo />
      </div>
    </div>
    <picture-avatar v-show="navbar.hidden" :user="user" />
    <ul class="menu-list u-mt">
      <li v-for="(item, index) in menu" :key="index">
        <nuxt-link v-if="item.path" :to="item.path" :exact="true" :aria-expanded="(isExpanded(item) ? 'true' : 'false')">
          <span class="icon is-small"><i :class="['far', item.icon]" /></span>
          <span class="menu-label-text">&nbsp;{{ getLabel(item) }}</span>
          <span v-if="hasChildren(item)" class="icon is-small is-angle">
            <i class="far fa-angle-down" />
          </span>
        </nuxt-link>
        <div v-else-if="item.label" class="menu-label">
          {{ item.label }}
        </div>

        <a v-else :aria-expanded="isExpanded(item)" @click="toggle(index, item)">
          <span class="icon is-small"><i :class="['far', item.icon]" /></span>
          <span class="menu-label-text">&nbsp;{{ getLabel(item) }}</span>
          <span v-if="item.children && item.children.length" class="icon is-small is-angle">
            <i class="far fa-angle-down" />
          </span>
        </a>

        <b-collapse :open="hasChildren(item)">
          <ul v-show="isExpanded(item)">
            <li v-for="(subItem, childIndex) in filteredChildren(item)" :key="childIndex">
              <nuxt-link :to="generatePath(item, subItem)">
                <span class="menu-label-text">&nbsp;{{ getLabel(subItem) }}</span>
              </nuxt-link>
            </li>
          </ul>
        </b-collapse>
      </li>
    </ul>
  </aside>
</template>

<script>
import { mapGetters, mapActions } from 'vuex'
import PictureAvatar from '~/components/PictureAvatar'
import Logo from '~/components/Logo'

export default {
  components: {
    PictureAvatar,
    Logo
  },

  props: {
    show: Boolean
  },

  data () {
    return {
      isReady: false
    }
  },

  computed: {
    firstName () {
      if (this.user.Name) {
        return this.user.Name.split(' ')[0]
      }
      return ''
    },
    ...mapGetters({
      navbar: 'app/navbar',
      user: 'auth/user'
      // userIsManager: 'userIsManager'
    }),
    menu () {
      let menuItems = []
      Object.keys(this.$menu.items).forEach((key) => {
        let item = this.$menu.items[key]
        if (this.checkAllowed(item)) {
          menuItems.push(item)
        }
      })
      return menuItems
    }
  },

  watch: {
    $route (route) {
      this.isReady = true
      this.shouldExpandMatchItem(route)
    }
  },

  mounted () {
    let route = this.$route
    if (route.name) {
      this.isReady = true
      this.shouldExpandMatchItem(route)
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
    },

    filteredChildren (item) {
      let menuItems = []
      if (item.children) {
        Object.keys(item.children).forEach((key) => {
          let child = item.children[key]
          if (this.checkAllowed(child)) {
            menuItems.push(child)
          }
        })
      }
      return menuItems
    }
  }
}
</script>

<style lang="scss">
@import "~public/scss/variables";
@import "~public/scss/mixins";

</style>
