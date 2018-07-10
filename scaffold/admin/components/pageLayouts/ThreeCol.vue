<template>
  <div class="layout-example is-clearfix" @click="layoutClicked" :class="{'is-selected': value === selectVal }">
    <div class="layout-image">
<svg viewBox="0 0 640 750" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" style="background: #FFFFFF;">
    <!-- Generator: Sketch 41 (35326) - http://www.bohemiancoding.com/sketch -->
    <title>Wireframes</title>
    <desc>Created with Sketch.</desc>
    <defs>
        <rect id="path-1" x="0" y="0" width="640" height="750"></rect>
        <mask id="mask-2" maskContentUnits="userSpaceOnUse" maskUnits="objectBoundingBox" x="0" y="0" width="640" height="750" fill="white">
            <use xlink:href="#path-1"></use>
        </mask>
    </defs>
    <g id="Page-1" stroke="none" stroke-width="1" fill="none" fill-rule="evenodd">
        <use id="Rectangle" stroke="#E4E4E4" mask="url(#mask-2)" stroke-width="8" xlink:href="#path-1"></use>
        <text id="Lorem-Ipsum" font-family="Helvetica" font-size="24" font-weight="normal" fill="#B6B6B6">
            <tspan x="26" y="40">Lorem Ipsum</tspan>
            <tspan x="234" y="40">Lorem Ipsum</tspan>
            <tspan x="436" y="40">Lorem Ipsum</tspan>
        </text>
        <text id="Lorem-ipsum-dolor-si" font-family="Helvetica" font-size="14" font-weight="normal" fill="#B6B6B6">
            <tspan x="28" y="63.0992994">Lorem ipsum dolor sit amet, </tspan>
            <tspan x="28" y="80.0992994">consectetur adipiscing elit. </tspan>
            <tspan x="28" y="97.0992994">Duis malesuada bibendum </tspan>
            <tspan x="28" y="114.099299">leo. In venenatis ligula et leo </tspan>
            <tspan x="28" y="131.099299">interdum, id laoreet sapien </tspan>
            <tspan x="28" y="148.099299">luctus. Sed non massa </tspan>
            <tspan x="28" y="165.099299">lectus. In sed aliquet velit. </tspan>
            <tspan x="28" y="182.099299">Aliquam erat volutpat. Nunc </tspan>
            <tspan x="28" y="199.099299">elementum justo id lorem </tspan>
            <tspan x="28" y="216.099299">ultrices, a viverra tellus </tspan>
            <tspan x="28" y="233.099299">mattis. Cras faucibus </tspan>
            <tspan x="28" y="250.099299">interdum iaculis. Praesent eu </tspan>
            <tspan x="28" y="267.099299">velit id orci rutrum commodo. </tspan>
            <tspan x="28" y="284.099299">Nullam hendrerit diam in </tspan>
            <tspan x="28" y="301.099299">dolor pulvinar malesuada. </tspan>
            <tspan x="28" y="318.099299">Maecenas vel nulla non </tspan>
            <tspan x="28" y="335.099299">ipsum finibus sagittis. Sed ac </tspan>
            <tspan x="28" y="352.099299">orci nibh. Vestibulum congue </tspan>
            <tspan x="28" y="369.099299">a lectus vel congue. </tspan>
            <tspan x="28" y="386.099299">Curabitur vehicula lorem et </tspan>
            <tspan x="28" y="403.099299">turpis imperdiet, sed varius </tspan>
            <tspan x="28" y="420.099299">mi scelerisque. Fusce varius </tspan>
            <tspan x="28" y="437.099299">nibh efficitur, tristique nunc </tspan>
            <tspan x="28" y="454.099299">ut, euismod quam. Curabitur </tspan>
            <tspan x="28" y="471.099299">luctus porttitor eros vel </tspan>
            <tspan x="28" y="488.099299">semper.</tspan>
            <tspan x="28" y="505.099299"></tspan>
            <tspan x="28" y="522.099299">Duis hendrerit sed elit sit </tspan>
            <tspan x="28" y="539.099299">amet laoreet. Sed vel neque </tspan>
            <tspan x="28" y="556.099299">tristique, scelerisque mi sed, </tspan>
            <tspan x="28" y="573.099299">tincidunt metus. Phasellus </tspan>
            <tspan x="28" y="590.099299">convallis condimentum </tspan>
            <tspan x="28" y="607.099299">finibus. Phasellus porta </tspan>
            <tspan x="28" y="624.099299">purus ac elementum mollis. </tspan>
            <tspan x="28" y="641.099299">Nulla nec quam egestas, </tspan>
            <tspan x="28" y="658.099299">consectetur diam ut, dictum </tspan>
            <tspan x="28" y="675.099299">ligula. Donec mattis ipsum </tspan>
            <tspan x="28" y="692.099299">sed lectus molestie efficitur </tspan>
            <tspan x="28" y="709.099299">eu et massa. </tspan>
        </text>
        <text id="Lorem-ipsum-dolor-si" font-family="Helvetica" font-size="14" font-weight="normal" fill="#B6B6B6">
            <tspan x="229.961798" y="63.0992994">Lorem ipsum dolor sit amet, </tspan>
            <tspan x="229.961798" y="80.0992994">consectetur adipiscing elit. </tspan>
            <tspan x="229.961798" y="97.0992994">Duis malesuada bibendum </tspan>
            <tspan x="229.961798" y="114.099299">leo. In venenatis ligula et leo </tspan>
            <tspan x="229.961798" y="131.099299">interdum, id laoreet sapien </tspan>
            <tspan x="229.961798" y="148.099299">luctus. Sed non massa </tspan>
            <tspan x="229.961798" y="165.099299">lectus. In sed aliquet velit. </tspan>
            <tspan x="229.961798" y="182.099299">Aliquam erat volutpat. Nunc </tspan>
            <tspan x="229.961798" y="199.099299">elementum justo id lorem </tspan>
            <tspan x="229.961798" y="216.099299">ultrices, a viverra tellus </tspan>
            <tspan x="229.961798" y="233.099299">mattis. Cras faucibus </tspan>
            <tspan x="229.961798" y="250.099299">interdum iaculis. Praesent eu </tspan>
            <tspan x="229.961798" y="267.099299">velit id orci rutrum commodo. </tspan>
            <tspan x="229.961798" y="284.099299">Nullam hendrerit diam in </tspan>
            <tspan x="229.961798" y="301.099299">dolor pulvinar malesuada. </tspan>
            <tspan x="229.961798" y="318.099299">Maecenas vel nulla non </tspan>
            <tspan x="229.961798" y="335.099299">ipsum finibus sagittis. Sed ac </tspan>
            <tspan x="229.961798" y="352.099299">orci nibh. Vestibulum congue </tspan>
            <tspan x="229.961798" y="369.099299">a lectus vel congue. </tspan>
            <tspan x="229.961798" y="386.099299">Curabitur vehicula lorem et </tspan>
            <tspan x="229.961798" y="403.099299">turpis imperdiet, sed varius </tspan>
            <tspan x="229.961798" y="420.099299">mi scelerisque. Fusce varius </tspan>
            <tspan x="229.961798" y="437.099299">nibh efficitur, tristique nunc </tspan>
            <tspan x="229.961798" y="454.099299">ut, euismod quam. Curabitur </tspan>
            <tspan x="229.961798" y="471.099299">luctus porttitor eros vel </tspan>
            <tspan x="229.961798" y="488.099299">semper.</tspan>
            <tspan x="229.961798" y="505.099299"></tspan>
            <tspan x="229.961798" y="522.099299">Duis hendrerit sed elit sit </tspan>
            <tspan x="229.961798" y="539.099299">amet laoreet. Sed vel neque </tspan>
            <tspan x="229.961798" y="556.099299">tristique, scelerisque mi sed, </tspan>
            <tspan x="229.961798" y="573.099299">tincidunt metus. Phasellus </tspan>
            <tspan x="229.961798" y="590.099299">convallis condimentum </tspan>
            <tspan x="229.961798" y="607.099299">finibus. Phasellus porta </tspan>
            <tspan x="229.961798" y="624.099299">purus ac elementum mollis. </tspan>
            <tspan x="229.961798" y="641.099299">Nulla nec quam egestas, </tspan>
            <tspan x="229.961798" y="658.099299">consectetur diam ut, dictum </tspan>
            <tspan x="229.961798" y="675.099299">ligula. Donec mattis ipsum </tspan>
            <tspan x="229.961798" y="692.099299">sed lectus molestie efficitur </tspan>
            <tspan x="229.961798" y="709.099299">eu et massa. </tspan>
        </text>
        <text id="Lorem-ipsum-dolor-si" font-family="Helvetica" font-size="14" font-weight="normal" fill="#B6B6B6">
            <tspan x="431.923596" y="63">Lorem ipsum dolor sit amet, </tspan>
            <tspan x="431.923596" y="80">consectetur adipiscing elit. </tspan>
            <tspan x="431.923596" y="97">Duis malesuada bibendum </tspan>
            <tspan x="431.923596" y="114">leo. In venenatis ligula et leo </tspan>
            <tspan x="431.923596" y="131">interdum, id laoreet sapien </tspan>
            <tspan x="431.923596" y="148">luctus. Sed non massa </tspan>
            <tspan x="431.923596" y="165">lectus. In sed aliquet velit. </tspan>
            <tspan x="431.923596" y="182">Aliquam erat volutpat. Nunc </tspan>
            <tspan x="431.923596" y="199">elementum justo id lorem </tspan>
            <tspan x="431.923596" y="216">ultrices, a viverra tellus </tspan>
            <tspan x="431.923596" y="233">mattis. Cras faucibus </tspan>
            <tspan x="431.923596" y="250">interdum iaculis. Praesent eu </tspan>
            <tspan x="431.923596" y="267">velit id orci rutrum commodo. </tspan>
            <tspan x="431.923596" y="284">Nullam hendrerit diam in </tspan>
            <tspan x="431.923596" y="301">dolor pulvinar malesuada. </tspan>
            <tspan x="431.923596" y="318">Maecenas vel nulla non </tspan>
            <tspan x="431.923596" y="335">ipsum finibus sagittis. Sed ac </tspan>
            <tspan x="431.923596" y="352">orci nibh. Vestibulum congue </tspan>
            <tspan x="431.923596" y="369">a lectus vel congue. </tspan>
            <tspan x="431.923596" y="386">Curabitur vehicula lorem et </tspan>
            <tspan x="431.923596" y="403">turpis imperdiet, sed varius </tspan>
            <tspan x="431.923596" y="420">mi scelerisque. Fusce varius </tspan>
            <tspan x="431.923596" y="437">nibh efficitur, tristique nunc </tspan>
            <tspan x="431.923596" y="454">ut, euismod quam. Curabitur </tspan>
            <tspan x="431.923596" y="471">luctus porttitor eros vel </tspan>
            <tspan x="431.923596" y="488">semper.</tspan>
            <tspan x="431.923596" y="505"></tspan>
            <tspan x="431.923596" y="522">Duis hendrerit sed elit sit </tspan>
            <tspan x="431.923596" y="539">amet laoreet. Sed vel neque </tspan>
            <tspan x="431.923596" y="556">tristique, scelerisque mi sed, </tspan>
            <tspan x="431.923596" y="573">tincidunt metus. Phasellus </tspan>
            <tspan x="431.923596" y="590">convallis condimentum </tspan>
            <tspan x="431.923596" y="607">finibus. Phasellus porta </tspan>
            <tspan x="431.923596" y="624">purus ac elementum mollis. </tspan>
            <tspan x="431.923596" y="641">Nulla nec quam egestas, </tspan>
            <tspan x="431.923596" y="658">consectetur diam ut, dictum </tspan>
            <tspan x="431.923596" y="675">ligula. Donec mattis ipsum </tspan>
            <tspan x="431.923596" y="692">sed lectus molestie efficitur </tspan>
            <tspan x="431.923596" y="709">eu et massa. </tspan>
        </text>
    </g>
</svg>
    </div>
    <div class="layout-label">{{selectVal}}</div>
  </div>
</template>

<script>
export default {
  props: {
    value: {
      type: String,
      required: true
    }
  },
  methods: {
    layoutClicked () {
      this.$emit('input', this.selectVal)
    }
  },
  data () {
    return {
      selectVal: 'Three Columns'
    }
  }
}
</script>
