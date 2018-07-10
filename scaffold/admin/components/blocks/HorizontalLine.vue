<template>
  <building-block title="Horizontal Line" :block="block">
    <div class="columns">
      <div class="column is-full">
        <div class="content is-clearfix">
         <hr style="background-color: #666;">
        </div>
      </div>
    </div>
  </building-block>
</template>

<script>
import BuildingBlock from './BuildingBlock'

export default {
  // COMPONENT
  // ______________________________________
  name: 'HorizontalLine',
  meta: {
    name: 'Horizontal Line',
    svg: `
<svg viewBox="0 0 960 76" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">
    <defs>
        <rect id="path-1" x="0" y="0" width="960" height="75" rx="4"></rect>
    </defs>
    <g id="Page-1" stroke="none" stroke-width="1" fill="none" fill-rule="evenodd">
        <g id="background">
            <use fill="#F7F7F7" fill-rule="evenodd" xlink:href="#path-1"></use>
            <rect stroke="#979797" stroke-width="2" x="1" y="1" width="958" height="73" rx="4"></rect>
        </g>
        <rect id="content" fill="#D8D8D8" x="25" y="30" width="903" height="10"></rect>
    </g>
</svg>
    `,
    category: 'Spacing' // basic, interactive, form
  },
  components: {
    BuildingBlock
  },
  props: {
    block: Object
  },
  computed: {
  },
  methods: {
  },
  // watch: {},
  data () {
    return {
      isInit: true
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
