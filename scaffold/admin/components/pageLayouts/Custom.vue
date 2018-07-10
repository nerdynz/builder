<template>
  <div class="layout-example is-clearfix" @click="layoutClicked" :class="{'is-selected': value === selectval }">
    <div class="layout-image">
      <svg viewBox="0 0 640 750" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">
          <!-- Generator: Sketch 41 (35326) - http://www.bohemiancoding.com/sketch -->
          <title>Slice</title>
          <desc>Created with Sketch.</desc>
          <defs>
              <rect id="path-1" x="0" y="0" width="640" height="750"></rect>
              <mask id="mask-2" maskContentUnits="userSpaceOnUse" maskUnits="objectBoundingBox" x="0" y="0" width="640" height="750" fill="white">
                  <use xlink:href="#path-1"></use>
              </mask>
          </defs>
          <g id="Page-1" stroke="none" stroke-width="1" fill="none" fill-rule="evenodd">
              <use id="Rectangle" stroke="#E4E4E4" mask="url(#mask-2)" stroke-width="16" xlink:href="#path-1"></use>
              <text id="?" font-family="Helvetica" font-size="400" font-weight="normal" fill="#B6B6B6">
                  <tspan x="189.269531" y="513">?</tspan>
              </text>
          </g>
      </svg>
    </div>
    <div class="layout-label">{{selectval}}</div>
  </div>
</template>

<script>
export default {
  props: {
    value: {
      type: String,
      required: true
    },
    selectval: {
      type: String,
      default: 'Custom'
    }
  },
  methods: {
    layoutClicked () {
      this.$emit('input', this.selectval)
    }
  },
  data () {
    return {
    }
  }
}
</script>
