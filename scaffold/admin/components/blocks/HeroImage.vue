<template>
  <building-block title="Hero Image" :block="block">
    <div class="columns">
      <div class="column is-full">
        <div class="content is-clearfix">
          <image-upload name="HeroImage" :width="1920" :height="840" output="jpeg" v-model="block.Picture"></image-upload>
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
    name: 'Hero Image',
    svg: `
<svg viewBox="0 0 960 400" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">
    <defs>
        <rect id="path-1" x="0" y="0" width="960" height="400" rx="4"></rect>
    </defs>
    <g id="Page-1" stroke="none" stroke-width="1" fill="none" fill-rule="evenodd">
        <g id="background">
            <use fill="#F7F7F7" fill-rule="evenodd" xlink:href="#path-1"></use>
            <rect stroke="#979797" stroke-width="2" x="1" y="1" width="958" height="398" rx="4"></rect>
        </g>
        <g id="content" transform="translate(61.000000, 37.000000)" fill-rule="nonzero" fill="#000000">
            <path d="M823.551724,0 L14.4482759,0 C6.47282759,0 0,3.20673913 0,7.17391304 L0,322.826087 C0,326.793261 6.47282759,330 14.4482759,330 L823.551724,330 C831.527172,330 838,326.793261 838,322.826087 L838,7.17391304 C838,3.20673913 831.527172,0 823.551724,0 Z M820,315.652174 L18,315.652174 L18,14.3478261 L820,14.3478261 L820,315.652174 Z" id="Shape"></path>
            <path d="M312.5,159 C336.487879,159 356,140.833543 356,118.507272 C356,96.1664572 336.487879,78 312.5,78 C288.512121,78 269,96.1664572 269,118.5 C269,140.833543 288.512121,159 312.5,159 Z M312.5,92.5448016 C327.872239,92.5448016 340.377806,104.195188 340.377806,118.5 C340.377806,132.804812 327.872239,144.455198 312.5,144.455198 C297.127761,144.455198 284.622194,132.812085 284.622194,118.507272 C284.622194,104.20246 297.127761,92.5448016 312.5,92.5448016 Z" id="Shape"></path>
            <path d="M244.089323,286 C245.981224,286 247.889294,285.415286 249.425453,284.224466 L381.316799,181.813976 L464.608937,255.26687 C467.770189,258.054957 472.879938,258.054957 476.041191,255.26687 C479.202443,252.478783 479.202443,247.972208 476.041191,245.184121 L437.176377,210.907055 L511.405179,139.215433 L602.450876,212.825201 C605.741489,215.484936 610.859323,215.285278 613.875045,212.383101 C616.890767,209.480924 616.672471,204.967218 613.373772,202.307483 L516.353227,123.87026 C514.768558,122.593872 512.658361,121.959244 510.531994,122.002028 C508.389457,122.087596 506.368196,122.921882 504.920973,124.319491 L425.736038,200.810045 L387.388668,166.989341 C384.364861,164.329606 379.530004,164.194123 376.328326,166.675592 L238.745108,273.51422 C235.389814,276.11691 235.066412,280.623485 238.017454,283.582707 C239.618293,285.187105 241.849766,286 244.089323,286 Z" id="Shape"></path>
        </g>
    </g>
</svg>
    `,
    category: 'Images' // basic, interactive, form
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
