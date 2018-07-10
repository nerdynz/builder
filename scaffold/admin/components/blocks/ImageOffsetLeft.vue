<template>
  <building-block :title="name" :block="block">
    <div class="columns">
      <div class="column is-3">
        <image-upload :width="480" :height="360" output="png" v-model="block.Picture"></image-upload>
      </div>
      <div class="column is-9">
        <div class="content is-clearfix">
          <rich-text v-model="block.HTML"></rich-text>
        </div>
      </div>
    </div>
  </building-block>
</template>

<script>
import BuildingBlock from './BuildingBlock'

let meta = {
  name: 'Image Offset Left',
  svg: `
<svg viewBox="0 0 960 200" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">
    <g id="Page-1" stroke="none" stroke-width="1" fill="none" fill-rule="evenodd">
        <g id="background">
            <use fill="#F7F7F7" fill-rule="evenodd" xlink:href="#path-1"></use>
            <rect stroke="#979797" stroke-width="2" x="1" y="1" width="958" height="198" rx="4"></rect>
        </g>
        <g id="Group" transform="translate(256.000000, 21.000000)" fill="#D8D8D8">
            <g id="content" transform="translate(11.000000, 0.000000)">
                <rect id="Rectangle-2" x="115.157778" y="0" width="557.842222" height="15.8940397"></rect>
                <rect id="Rectangle-2" x="115.157778" y="36.0264901" width="557.842222" height="15.8940397"></rect>
                <rect id="Rectangle-2" x="115.157778" y="72.0529801" width="557.842222" height="15.8940397"></rect>
                <rect id="Rectangle-2" x="0" y="108.07947" width="673" height="15.8940397"></rect>
                <rect id="Rectangle-2" x="0" y="144.10596" width="673" height="15.8940397"></rect>
            </g>
            <text id="content" font-family="Helvetica-Bold, Helvetica" font-size="72" font-weight="bold">
                <tspan x="10.9804688" y="74.5993789">Aa</tspan>
            </text>
        </g>
        <g id="content" transform="translate(20.000000, 23.000000)" fill-rule="nonzero" fill="#000000">
            <path d="M205.142727,0 L3.5989952,0 C1.61234985,0 0,1.55478261 0,3.47826087 L0,156.521739 C0,158.445217 1.61234985,160 3.5989952,160 L205.142727,160 C207.129372,160 208.741722,158.445217 208.741722,156.521739 L208.741722,3.47826087 C208.741722,1.55478261 207.129372,0 205.142727,0 Z M201.543731,153.043478 L7.19799041,153.043478 L7.19799041,6.95652174 L201.543731,6.95652174 L201.543731,153.043478 Z" id="Shape"></path>
            <path d="M56.6887417,77.3509934 C67.4985343,77.3509934 76.2913907,68.558137 76.2913907,57.7518643 C76.2913907,46.9385518 67.4985343,38.1456954 56.6887417,38.1456954 C45.8789491,38.1456954 37.0860927,46.9385518 37.0860927,57.7483444 C37.0860927,68.558137 45.8789491,77.3509934 56.6887417,77.3509934 Z M56.6887417,45.1856124 C63.6160201,45.1856124 69.2514737,50.824586 69.2514737,57.7483444 C69.2514737,64.6721028 63.6160201,70.3110763 56.6887417,70.3110763 C49.7614634,70.3110763 44.1260098,64.6756227 44.1260098,57.7518643 C44.1260098,50.8281059 49.7614634,45.1856124 56.6887417,45.1856124 Z" id="Shape"></path>
            <path d="M25.8832514,138.807947 C26.7325954,138.807947 27.5891988,138.52461 28.2788371,137.947569 L87.4897305,88.3220687 L124.882646,123.915459 C126.301849,125.266495 128.595803,125.266495 130.015007,123.915459 C131.43421,122.564424 131.43421,120.380653 130.015007,119.029618 L112.567157,102.41983 L145.891206,67.6799061 L186.76498,103.349314 C188.242258,104.638154 190.539842,104.541404 191.893711,103.135084 C193.24758,101.728763 193.149579,99.5415367 191.668671,98.2526973 L148.112568,60.2440283 C147.401151,59.6255236 146.453806,59.3179989 145.499201,59.3387309 C144.537337,59.3801949 143.629919,59.784469 142.980207,60.4617143 L107.431167,97.5270772 L90.2156167,81.1384302 C88.8581181,79.8495908 86.6875723,79.7839395 85.2502208,80.9863956 L23.484036,132.757658 C21.9777207,134.018855 21.8325337,136.202626 23.1573652,137.636589 C23.8760409,138.414039 24.8778313,138.807947 25.8832514,138.807947 Z" id="Shape"></path>
        </g>
    </g>
</svg>
  `,
  category: 'Images' // basic, interactive, form
}

export default {
  // COMPONENT
  // ______________________________________
  name: 'ImageOffsetLeft',
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
