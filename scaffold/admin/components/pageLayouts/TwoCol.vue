<template>
  <div class="layout-example is-clearfix" @click="layoutClicked" :class="{'is-selected': value === selectVal }">
    <div class="layout-image">
<svg  viewBox="0 0 640 750" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" style="background: #FFFFFF;">
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
            <tspan x="340" y="40">Lorem Ipsum</tspan>
        </text>
        <text id="Lorem-ipsum-dolor-si" font-family="Helvetica" font-size="14" font-weight="normal" fill="#B6B6B6">
            <tspan x="28" y="67.5681256">Lorem ipsum dolor sit amet, consectetur </tspan>
            <tspan x="28" y="84.5681256">adipiscing elit. Duis malesuada bibendum </tspan>
            <tspan x="28" y="101.568126">leo. In venenatis ligula et leo interdum, id </tspan>
            <tspan x="28" y="118.568126">laoreet sapien luctus. Sed non massa lectus. </tspan>
            <tspan x="28" y="135.568126">In sed aliquet velit. Aliquam erat volutpat. </tspan>
            <tspan x="28" y="152.568126">Nunc elementum justo id lorem ultrices, a </tspan>
            <tspan x="28" y="169.568126">viverra tellus mattis. Cras faucibus interdum </tspan>
            <tspan x="28" y="186.568126">iaculis. Praesent eu velit id orci rutrum </tspan>
            <tspan x="28" y="203.568126">commodo. Nullam hendrerit diam in dolor </tspan>
            <tspan x="28" y="220.568126">pulvinar malesuada. Maecenas vel nulla non </tspan>
            <tspan x="28" y="237.568126">ipsum finibus sagittis. Sed ac orci nibh. </tspan>
            <tspan x="28" y="254.568126">Vestibulum congue a lectus vel congue. </tspan>
            <tspan x="28" y="271.568126">Curabitur vehicula lorem et turpis imperdiet, </tspan>
            <tspan x="28" y="288.568126">sed varius mi scelerisque. Fusce varius nibh </tspan>
            <tspan x="28" y="305.568126">efficitur, tristique nunc ut, euismod quam. </tspan>
            <tspan x="28" y="322.568126">Curabitur luctus porttitor eros vel semper.</tspan>
            <tspan x="28" y="339.568126"></tspan>
            <tspan x="28" y="356.568126">Duis hendrerit sed elit sit amet laoreet. Sed </tspan>
            <tspan x="28" y="373.568126">vel neque tristique, scelerisque mi sed, </tspan>
            <tspan x="28" y="390.568126">tincidunt metus. Phasellus convallis </tspan>
            <tspan x="28" y="407.568126">condimentum finibus. Phasellus porta purus </tspan>
            <tspan x="28" y="424.568126">ac elementum mollis. Nulla nec quam </tspan>
            <tspan x="28" y="441.568126">egestas, consectetur diam ut, dictum ligula. </tspan>
            <tspan x="28" y="458.568126">Donec mattis ipsum sed lectus molestie </tspan>
            <tspan x="28" y="475.568126">efficitur eu et massa. Etiam et odio vel neque </tspan>
            <tspan x="28" y="492.568126">maximus cursus. Etiam magna sapien, </tspan>
            <tspan x="28" y="509.568126">congue a hendrerit ut, efficitur at orci.</tspan>
            <tspan x="28" y="526.568126"></tspan>
            <tspan x="28" y="543.568126">Fusce cursus mattis nibh et lobortis. Nullam </tspan>
            <tspan x="28" y="560.568126">ultricies feugiat enim a accumsan. </tspan>
            <tspan x="28" y="577.568126">Pellentesque vestibulum in nibh sit amet </tspan>
            <tspan x="28" y="594.568126">vulputate. Donec nisi magna, tristique nec </tspan>
            <tspan x="28" y="611.568126">hendrerit a, consectetur vitae odio. Fusce </tspan>
            <tspan x="28" y="628.568126">placerat velit felis, nec mattis urna pharetra </tspan>
            <tspan x="28" y="645.568126">in. Sed eu turpis placerat, sodales elit a, </tspan>
            <tspan x="28" y="662.568126">suscipit justo. Aliquam luctus tincidunt </tspan>
            <tspan x="28" y="679.568126">rutrum. Maecenas ut odio non mauris </tspan>
            <tspan x="28" y="696.568126">placerat lacinia sit amet sed ipsum. Donec et </tspan>
            <tspan x="28" y="713.568126">nisl quis turpis semper vulputate. </tspan>
        </text>
        <text id="Lorem-ipsum-dolor-si" font-family="Helvetica" font-size="14" font-weight="normal" fill="#B6B6B6">
            <tspan x="338" y="67.5681256">Lorem ipsum dolor sit amet, consectetur </tspan>
            <tspan x="338" y="84.5681256">adipiscing elit. Duis malesuada bibendum </tspan>
            <tspan x="338" y="101.568126">leo. In venenatis ligula et leo interdum, id </tspan>
            <tspan x="338" y="118.568126">laoreet sapien luctus. Sed non massa lectus. </tspan>
            <tspan x="338" y="135.568126">In sed aliquet velit. Aliquam erat volutpat. </tspan>
            <tspan x="338" y="152.568126">Nunc elementum justo id lorem ultrices, a </tspan>
            <tspan x="338" y="169.568126">viverra tellus mattis. Cras faucibus interdum </tspan>
            <tspan x="338" y="186.568126">iaculis. Praesent eu velit id orci rutrum </tspan>
            <tspan x="338" y="203.568126">commodo. Nullam hendrerit diam in dolor </tspan>
            <tspan x="338" y="220.568126">pulvinar malesuada. Maecenas vel nulla non </tspan>
            <tspan x="338" y="237.568126">ipsum finibus sagittis. Sed ac orci nibh. </tspan>
            <tspan x="338" y="254.568126">Vestibulum congue a lectus vel congue. </tspan>
            <tspan x="338" y="271.568126">Curabitur vehicula lorem et turpis imperdiet, </tspan>
            <tspan x="338" y="288.568126">sed varius mi scelerisque. Fusce varius nibh </tspan>
            <tspan x="338" y="305.568126">efficitur, tristique nunc ut, euismod quam. </tspan>
            <tspan x="338" y="322.568126">Curabitur luctus porttitor eros vel semper.</tspan>
            <tspan x="338" y="339.568126"></tspan>
            <tspan x="338" y="356.568126">Duis hendrerit sed elit sit amet laoreet. Sed </tspan>
            <tspan x="338" y="373.568126">vel neque tristique, scelerisque mi sed, </tspan>
            <tspan x="338" y="390.568126">tincidunt metus. Phasellus convallis </tspan>
            <tspan x="338" y="407.568126">condimentum finibus. Phasellus porta purus </tspan>
            <tspan x="338" y="424.568126">ac elementum mollis. Nulla nec quam </tspan>
            <tspan x="338" y="441.568126">egestas, consectetur diam ut, dictum ligula. </tspan>
            <tspan x="338" y="458.568126">Donec mattis ipsum sed lectus molestie </tspan>
            <tspan x="338" y="475.568126">efficitur eu et massa. Etiam et odio vel neque </tspan>
            <tspan x="338" y="492.568126">maximus cursus. Etiam magna sapien, </tspan>
            <tspan x="338" y="509.568126">congue a hendrerit ut, efficitur at orci.</tspan>
            <tspan x="338" y="526.568126"></tspan>
            <tspan x="338" y="543.568126">Fusce cursus mattis nibh et lobortis. Nullam </tspan>
            <tspan x="338" y="560.568126">ultricies feugiat enim a accumsan. </tspan>
            <tspan x="338" y="577.568126">Pellentesque vestibulum in nibh sit amet </tspan>
            <tspan x="338" y="594.568126">vulputate. Donec nisi magna, tristique nec </tspan>
            <tspan x="338" y="611.568126">hendrerit a, consectetur vitae odio. Fusce </tspan>
            <tspan x="338" y="628.568126">placerat velit felis, nec mattis urna pharetra </tspan>
            <tspan x="338" y="645.568126">in. Sed eu turpis placerat, sodales elit a, </tspan>
            <tspan x="338" y="662.568126">suscipit justo. Aliquam luctus tincidunt </tspan>
            <tspan x="338" y="679.568126">rutrum. Maecenas ut odio non mauris </tspan>
            <tspan x="338" y="696.568126">placerat lacinia sit amet sed ipsum. Donec et </tspan>
            <tspan x="338" y="713.568126">nisl quis turpis semper vulputate. </tspan>
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
      selectVal: 'Two Columns'
    }
  }
}
</script>
