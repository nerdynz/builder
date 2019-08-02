<template>
  <div class="level app-levelbar u-p -u-mb">
    <div class="level-left">
      <div class="level-item">
        <h3 class="subtitle is-5">
          <strong>{{ getTitle(current) }} <span v-if="subtitle"> - {{ subtitle }}</span></strong>
        </h3>
      </div>
    </div>
    <div class="level-right is-hidden-mobile">
      <nav class="breadcrumb" aria-label="breadcrumbs">
        <ul>
          <li v-for="(item, index) in $menu.bread" :key="index" :class="{'is-active': index === $menu.bread.length - 1}">
            <nuxt-link :to="item.path" :exact="true">
              {{ getTitle(item) }}
            </nuxt-link>
          </li>
        </ul>
      </nav>
    </div>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'

export default {
  // COMPONENT
  // ______________________________________
  name: 'Levelbar',
  components: {},
  props: {},
  head () {
    return {
      title: this.actualTitle
    }
  },
  data () {
    return {
    }
  },
  computed: {
    ...mapGetters({
      current: 'app/current',
      title: 'app/title',
      subtitle: 'app/subtitle'
    }),
    actualTitle () {
      let title = typeof (this.title) === 'function' ? this.title(this) : this.title
      if (title) {
        return title + ' - '
      }
      return ''
    }
  },
  watch: {
  },
  methods: {
    getTitle (item) {
      if (!item) {
        return 'Generic Title'
      }
      return typeof (item.title) === 'function' ? item.title(this) : item.title
    }
  }
}
</script>

<style lang="scss">

</style>
