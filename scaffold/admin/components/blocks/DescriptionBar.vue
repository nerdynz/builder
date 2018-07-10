<template>
  <span class="description">
    <a @click="edit" class="desc-text">{{value}}</a>
    <a @click="removed" v-if="removeable || userIsDev" v-tooltip.top.hover="{content: 'Remove Block', class:'is-danger'}" class="desc-button is-danger"><i class="fa-trash fa"></i></a>
    <a @click="moveDown" v-tooltip.top.hover="{content: 'Move Block Down'}" class="desc-button"><i class="fa fa-chevron-down"></i></a>
    <a @click="moveUp" class="desc-button" v-tooltip.top.hover="{content: 'Move Block Up'}"><i class="fa fa-chevron-up"></i></a>
  </span>
</template>

<script>
// import {BLOCK_HOVERED} from 'vuex-store/mutation-types'
import {mapGetters} from 'vuex'

export default {
  // COMPONENT
  // ______________________________________
  name: 'DescriptionBar',
  components: {
  },
  props: {
    block: Object,
    value: String,
    removeable: Boolean
  },
  computed: {
    ...mapGetters({
      userIsDev: 'auth/userIsDev'
    })
  },
  methods: {
    edit () {
      this.$emit('edit')
    },
    removed () {
      this.$emit('removed')
    },
    moveUp () {
      this.$emit('move', 'up')
    },
    moveDown () {
      this.$emit('move', 'down')
    }
  },
  // watch: {},
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
