<template>
  <building-block :title="name" :block="block">
    <div class="columns">
      <div class="column is-half">
        <div class="content is-clearfix">
          <rich-text v-model="block.HTML"></rich-text>
        </div>
      </div>
      <div class="column is-half">
        <div class="content is-clearfix">
          <rich-text v-model="block.HTMLTwo"></rich-text>
        </div>
      </div>
    </div>
  </building-block>
</template>

<script>
import BuildingBlock from './BuildingBlock'

let meta = {
  name: 'Two Column',
  svg: `
<svg viewBox="0 0 960 200" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">
    <g id="Page-1" stroke="none" stroke-width="1" fill="none" fill-rule="evenodd">
        <g id="background">
            <use fill="#F7F7F7" fill-rule="evenodd" xlink:href="#path-1"></use>
            <rect stroke="#979797" stroke-width="2" x="1" y="1" width="958" height="198" rx="4"></rect>
        </g>
        <g id="Group" transform="translate(11.000000, 22.000000)" fill="#D8D8D8">
            <g id="content" transform="translate(21.000000, 0.000000)">
                <rect id="Rectangle-2" x="108.632813" y="0" width="324.367187" height="15"></rect>
                <rect id="Rectangle-2" x="108.632813" y="34" width="324.367187" height="15"></rect>
                <rect id="Rectangle-2" x="108.632813" y="68" width="324.367187" height="15"></rect>
                <rect id="Rectangle-2" x="0" y="102" width="433" height="15"></rect>
                <rect id="Rectangle-2" x="0" y="136" width="433" height="15"></rect>
            </g>
            <text id="content" font-family="Helvetica-Bold, Helvetica" font-size="72" font-weight="bold">
                <tspan x="27.9804688" y="69">Aa</tspan>
            </text>
        </g>
        <g id="Group" transform="translate(478.000000, 22.000000)" fill="#D8D8D8">
            <g id="content" transform="translate(21.000000, 0.000000)">
                <rect id="Rectangle-2" x="108.632813" y="0" width="324.367187" height="15"></rect>
                <rect id="Rectangle-2" x="108.632813" y="34" width="324.367187" height="15"></rect>
                <rect id="Rectangle-2" x="108.632813" y="68" width="324.367187" height="15"></rect>
                <rect id="Rectangle-2" x="0" y="102" width="433" height="15"></rect>
                <rect id="Rectangle-2" x="0" y="136" width="433" height="15"></rect>
            </g>
            <text id="content" font-family="Helvetica-Bold, Helvetica" font-size="72" font-weight="bold">
                <tspan x="27.9804688" y="69">Aa</tspan>
            </text>
        </g>
    </g>
</svg>
  `,
  category: 'Basic' // basic, interactive, form
}

export default {
  // COMPONENT
  // ______________________________________
  name: 'TwoColumn',
  meta: meta,
  components: {
    BuildingBlock
  },
  props: {
    block: Object
  },
  computed: {
    name () {
      return meta.name
    }
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
