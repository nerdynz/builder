<template>
  <div ref="blockOuter" class="building-block block-sort-pos" :style="{height: hoverHeight}" :uuid="block.UUID">
    <description-bar class="has-close" :removeable="removeable" :value="title" @edit="edit" @move="move" @removed="remove" :block="block">
    </description-bar>
    <div ref="blockInner" :class="{'block-placeholder': custom}">
      <slot>
      </slot>
    </div>
  </div>
</template>

<script>
import DescriptionBar from './DescriptionBar'
import {mapGetters} from 'vuex'

export default {
  // COMPONENT
  // ______________________________________
  name: 'BuildingBlock',
  components: {
    DescriptionBar
  },
  props: {
    title: String,
    block: Object,
    removeable: {
      type: Boolean,
      default: true
    },
    custom: Boolean
  },
  computed: {
    ...mapGetters({
      // blockIsHovered: 'blockIsHovered'
    })
  },
  methods: {
    move (direction) {
      this.$parent.$emit('move', {
        uuid: this.block.UUID,
        direction: direction
      })
    },
    remove () {
      this.$parent.$emit('removed', this.block.UUID)
    },
    edit () {
      this.$parent.$emit('edit', this.block.UUID)
    }
  },
  watch: {
    // blockIsHovered (isHovering) {
    //   // if (this.isHovering) {
    //   //   return
    //   // }

    //   // var paddingBottom = Number($(this.$refs.blockOuter).css('padding-bottom').replace('px', ''))
    //   if (isHovering) {
    //     // this.copyHtml = ''
    //     this.hoverHeight = this.$refs.blockOuter.scrollHeight + 2 + 'px'
    //   } else {
    //     // this.copyHtml = ''
    //     this.hoverHeight = 'initial'
    //   }
    //   // this.isHovering = isHovering
    // }
  },
  data () {
    return {
      copyHtml: '',
      // isHovering: false,
      hoverHeight: 'initial'
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
