<template>
  <building-block :title="'Recent Work'" :block="block">
    <div class="columns">
      <div class="column is-full">
        <div class="custom-block">
          Recent Work
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
  name: 'OneColumn',
  meta: {
    name: 'Recent Work',
    svg: `
<svg  viewBox="0 0 960 200" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">
    <g id="Page-1" stroke="none" stroke-width="1" fill="none" fill-rule="evenodd">
        <g id="background">
            <use fill="#F7F7F7" fill-rule="evenodd" xlink:href="#path-1"></use>
            <rect stroke="#979797" stroke-width="2" x="1" y="1" width="958" height="198" rx="4"></rect>
        </g>
        <g id="content" transform="translate(30.000000, 22.000000)" fill="#D8D8D8">
            <rect id="Rectangle-2" x="154" y="0" width="746" height="15"></rect>
            <rect id="Rectangle-2" x="154" y="34" width="746" height="15"></rect>
            <rect id="Rectangle-2" x="154" y="68" width="746" height="15"></rect>
            <rect id="Rectangle-2" x="0" y="102" width="900" height="15"></rect>
            <rect id="Rectangle-2" x="0" y="136" width="900" height="15"></rect>
        </g>
        <text id="content" font-family="Helvetica-Bold, Helvetica" font-size="96" font-weight="bold" fill="#D8D8D8">
            <tspan x="30.640625" y="99">Aa</tspan>
        </text>
    </g>
</svg>
    `,
    category: 'Custom' // basic, interactive, form
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
