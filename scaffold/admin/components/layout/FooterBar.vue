<template>
  <footer class="footer" :class="{'sidebar-open': sidebarIsOpen}">
    <div class="level">
      <div class="level-left">
        <div v-for="(item, i) in leftButtons" :key="i" class="level-item">
          <button class="button" :class="'is-' + item.kind" :type="item.kind" :disabled="item.isDisabled || !item.click" @click="(ev) => { item.click(ev) }">
            {{ item.text }}
          </button>
        </div>
      </div>
      <div v-for="(item, i) in centerButtons" :key="i" class="level-item">
        <button class="button" :class="'is-' + item.kind" :type="item.kind" :disabled="item.isDisabled || !item.click" @click="(ev) => { item.click(ev) }">
          {{ item.text }}
        </button>
      </div>
      <div class="level-right">
        <div v-for="(item, i) in rightButtons" :key="i" class="level-item">
          <button class="button" :class="'is-' + item.kind" :type="item.kind" :disabled="item.isDisabled || !item.click" @click="(ev) => { item.click(ev) }">
            {{ item.text }}
          </button>
        </div>
      </div>
    </div>
  </footer>
</template>

<script>
import { mapGetters, mapActions } from 'vuex'
export default {
  data () {
    return {
    }
  },
  computed: {
    ...mapGetters({
      leftButtons: 'app/leftButtons',
      centerButtons: 'app/centerButtons',
      rightButtons: 'app/rightButtons',
      navbar: 'app/navbar',
      sidebarIsOpen: 'app/sidebarIsOpen'
    })
  },
  methods: {
    ...mapActions([
      'toggleJobView'
    ])
  }
}
</script>
