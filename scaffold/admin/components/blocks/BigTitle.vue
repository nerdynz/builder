<template>
  <building-block title="Big Title" :block="block">
    <div class="columns big-title-block">
      <div class="column is-full">
        <div class="content is-clearfix">
          <control name="HTML" type="text" v-model="block.HTML"></control>
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
  meta: {
    name: 'Big Title',
    svg: `
<svg  viewBox="0 0 960 128" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">
    <defs>
        <rect id="path-1" x="0" y="0" width="960" height="128" rx="4"></rect>
    </defs>
    <g id="Page-1" stroke="none" stroke-width="1" fill="none" fill-rule="evenodd">
        <g id="background">
            <use fill="#F7F7F7" fill-rule="evenodd" xlink:href="#path-1"></use>
            <rect stroke="#979797" stroke-width="2" x="1" y="1" width="958" height="126" rx="4"></rect>
        </g>
        <g id="Group" transform="translate(15.000000, 19.000000)" fill="#D8D8D8">
            <g id="content" transform="translate(169.000000, 8.000000)">
                <rect id="Rectangle-2" x="0" y="0" width="746" height="57"></rect>
            </g>
            <text id="content" font-family="Helvetica-Bold, Helvetica" font-size="64" font-weight="bold">
                <tspan x="35.59375" y="62">Aa</tspan>
            </text>
        </g>
    </g>
</svg>
    `,
    category: 'Titles' // basic, interactive, form
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
<style lang="scss">
  @import "~public/scss/_variables";
  .big-title-block input#HTML {
    font-size: $size-1;
    padding: 2px 5px;
    height: initial
  }
</style>

