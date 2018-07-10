<template>
  <div class="layout-example is-clearfix" @click="layoutClicked" :class="{'is-selected': value === selectVal }">
    <div class="layout-image">
      <svg viewBox="0 0 640 750"
           version="1.1"
           xmlns="http://www.w3.org/2000/svg"
           xmlns:xlink="http://www.w3.org/1999/xlink"
           style="background: #FFFFFF;">
        <defs>
          <rect id="path-1"
                x="0"
                y="0"
                width="640"
                height="750"></rect>
          <mask id="mask-2"
                maskContentUnits="userSpaceOnUse"
                maskUnits="objectBoundingBox"
                x="0"
                y="0"
                width="640"
                height="750"
                fill="white">
            <use xlink:href="#path-1"></use>
          </mask>
        </defs>
        <g id="Page-1"
           stroke="none"
           stroke-width="1"
           fill="none"
           fill-rule="evenodd">
          <use id="Rectangle"
               stroke="#E4E4E4"
               mask="url(#mask-2)"
               stroke-width="8"
               xlink:href="#path-1"></use>
          <text id="Lorem-Ipsum"
                font-family="Helvetica"
                font-size="24"
                font-weight="normal"
                fill="#B6B6B6">
            <tspan x="23"
                   y="405">Lorem Ipsum</tspan>
          </text>
          <text id="Lorem-ipsum-dolor-si"
                font-family="Helvetica"
                font-size="14"
                font-weight="normal"
                fill="#B6B6B6">
            <tspan x="24"
                   y="437.568126">Lorem ipsum dolor sit amet, consectetur adipiscing elit. Duis malesuada bibendum leo. In </tspan>
            <tspan x="24"
                   y="454.568126">venenatis ligula et leo interdum, id laoreet sapien luctus. Sed non massa lectus. In sed aliquet </tspan>
            <tspan x="24"
                   y="471.568126">velit. Aliquam erat volutpat. Nunc elementum justo id lorem ultrices, a viverra tellus mattis. Cras </tspan>
            <tspan x="24"
                   y="488.568126">faucibus interdum iaculis. Praesent eu velit id orci rutrum commodo. Nullam hendrerit diam in </tspan>
            <tspan x="24"
                   y="505.568126">dolor pulvinar malesuada. Fusce cursus mattis nibh et lobortis. Nullam ultricies feugiat enim a </tspan>
            <tspan x="24"
                   y="522.568126">accumsan. </tspan>
            <tspan x="24"
                   y="539.568126"></tspan>
            <tspan x="24"
                   y="556.568126">Pellentesque vestibulum in nibh sit amet vulputate. Donec nisi magna, tristique nec hendrerit a, </tspan>
            <tspan x="24"
                   y="573.568126">consectetur vitae odio. Maecenas vel nulla non ipsum finibus sagittis. Sed ac orci nibh. </tspan>
            <tspan x="24"
                   y="590.568126">Vestibulum congue a lectus vel congue. Curabitur vehicula lorem et turpis imperdiet, sed varius </tspan>
            <tspan x="24"
                   y="607.568126">mi scelerisque. Fusce varius nibh efficitur, tristique nunc ut, euismod quam. Curabitur luctus </tspan>
            <tspan x="24"
                   y="624.568126">porttitor eros vel semper.</tspan>
            <tspan x="24"
                   y="641.568126"></tspan>
            <tspan x="24"
                   y="658.568126">Duis hendrerit sed elit sit amet laoreet. Sed vel neque tristique, scelerisque mi sed, tincidunt </tspan>
            <tspan x="24"
                   y="675.568126">metus. Phasellus convallis condimentum finibus. Phasellus porta purus ac elementum mollis. </tspan>
            <tspan x="24"
                   y="692.568126">Nulla nec quam egestas, consectetur diam ut, dictum ligula. Donec mattis ipsum sed lectus </tspan>
            <tspan x="24"
                   y="709.568126">molestie efficitur eu et massa. Etiam et odio vel neque maximus cursus. Etiam magna sapien, </tspan>
            <tspan x="24"
                   y="726.568126">congue a hendrerit ut, efficitur at orci.</tspan>
          </text>
          <g id="image-placeholder_3x2"
             transform="translate(24.000000, 32.000000)">
            <g id="placeholder">
              <rect id="Rectangle-path"
                    fill="#E4E4E4"
                    x="0"
                    y="0"
                    width="591"
                    height="333"></rect>
              <g id="Group"
                 transform="translate(210.000000, 102.000000)"
                 fill="#FFFFFF">
                <path d="M0.0278866667,0.516666667 L0.0278866667,127.55 L171.972113,127.55 L171.972113,0.516666667 L0.0278866667,0.516666667 Z M158.112507,113.773667 L14.6951933,113.773667 L14.6951933,14.2917667 L158.112507,14.2917667 L158.112507,113.773667 L158.112507,113.773667 Z"
                      id="Shape"></path>
                <polygon id="Shape"
                         points="22.4084 103.0597 51.9689067 73.8802667 63.0442467 78.3807 97.9907333 43.4465333 111.746587 58.8977333 117.952087 55.3778 151.532707 103.0597"></polygon>
                <ellipse id="Oval"
                         cx="56.5576933"
                         cy="39.5886667"
                         rx="12.9100667"
                         ry="12.1236667"></ellipse>
              </g>
            </g>
          </g>
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
      selectVal: 'Hero Image'
    }
  }
}
</script>
