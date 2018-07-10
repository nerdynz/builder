<template>
  <div class="layout-example is-clearfix" @click="layoutClicked" :class="{'is-selected': value === selectVal }">
    <div class="layout-image">
    <svg viewBox="0 0 640 750" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" style="background: #FFFFFF;">
        <defs>
            <rect id="path-1" x="0" y="0" width="640" height="750"></rect>
            <mask id="mask-2" maskContentUnits="userSpaceOnUse" maskUnits="objectBoundingBox" x="0" y="0" width="640" height="750" fill="white">
                <use xlink:href="#path-1"></use>
            </mask>
        </defs>
        <g id="Page-1" stroke="none" stroke-width="1" fill="none" fill-rule="evenodd">
            <use id="Rectangle" stroke="#E4E4E4" mask="url(#mask-2)" stroke-width="8" xlink:href="#path-1"></use>
            <text id="Lorem-Ipsum" font-family="Helvetica" font-size="24" font-weight="normal" fill="#B6B6B6">
                <tspan x="29" y="417">Lorem Ipsum</tspan>
            </text>
            <text id="Lorem-Ipsum" font-family="Helvetica" font-size="24" font-weight="normal" fill="#B6B6B6">
                <tspan x="29" y="417">Lorem Ipsum</tspan>
            </text>
            <text id="Lorem-ipsum-dolor-si" font-family="Helvetica" font-size="14" font-weight="normal" fill="#B6B6B6">
                <tspan x="28" y="446">Lorem ipsum dolor sit amet, </tspan>
                <tspan x="28" y="463">consectetur adipiscing elit. </tspan>
                <tspan x="28" y="480">Aenean porttitor semper </tspan>
                <tspan x="28" y="497">lacus. Vivamus fermentum </tspan>
                <tspan x="28" y="514">dolor et erat tincidunt </tspan>
                <tspan x="28" y="531">tempus. Ut tristique, purus </tspan>
                <tspan x="28" y="548">eget volutpat elementum, </tspan>
                <tspan x="28" y="565">quam purus luctus nisl, in </tspan>
                <tspan x="28" y="582">lacinia metus ipsum aliquet </tspan>
                <tspan x="28" y="599">mi. Praesent eleifend ante </tspan>
                <tspan x="28" y="616">mauris, quis varius tortor </tspan>
                <tspan x="28" y="633">dapibus sit amet. Curabitur </tspan>
                <tspan x="28" y="650">ac lacus sit amet nisi lobortis </tspan>
                <tspan x="28" y="667">blandit quis quis nibh. </tspan>
                <tspan x="28" y="684">Suspendisse convallis </tspan>
                <tspan x="28" y="701">cursus nisi, sit amet varius </tspan>
                <tspan x="28" y="718">neque tincidunt ultrices.</tspan>
            </text>
            <text id="Lorem-Ipsum" font-family="Helvetica" font-size="24" font-weight="normal" fill="#B6B6B6">
                <tspan x="237" y="417">Lorem Ipsum</tspan>
            </text>
            <text id="Lorem-ipsum-dolor-si" font-family="Helvetica" font-size="14" font-weight="normal" fill="#B6B6B6">
                <tspan x="236" y="446">Lorem ipsum dolor sit amet, </tspan>
                <tspan x="236" y="463">consectetur adipiscing elit. </tspan>
                <tspan x="236" y="480">Aenean porttitor semper </tspan>
                <tspan x="236" y="497">lacus. Vivamus fermentum </tspan>
                <tspan x="236" y="514">dolor et erat tincidunt </tspan>
                <tspan x="236" y="531">tempus. Ut tristique, purus </tspan>
                <tspan x="236" y="548">eget volutpat elementum, </tspan>
                <tspan x="236" y="565">quam purus luctus nisl, in </tspan>
                <tspan x="236" y="582">lacinia metus ipsum aliquet </tspan>
                <tspan x="236" y="599">mi. Praesent eleifend ante </tspan>
                <tspan x="236" y="616">mauris, quis varius tortor </tspan>
                <tspan x="236" y="633">dapibus sit amet. Curabitur </tspan>
                <tspan x="236" y="650">ac lacus sit amet nisi lobortis </tspan>
                <tspan x="236" y="667">blandit quis quis nibh. </tspan>
                <tspan x="236" y="684">Suspendisse convallis </tspan>
                <tspan x="236" y="701">cursus nisi, sit amet varius </tspan>
                <tspan x="236" y="718">neque tincidunt ultrices.</tspan>
            </text>
            <text id="Lorem-Ipsum" font-family="Helvetica" font-size="24" font-weight="normal" fill="#B6B6B6">
                <tspan x="440" y="417">Lorem Ipsum</tspan>
            </text>
            <text id="Lorem-ipsum-dolor-si" font-family="Helvetica" font-size="14" font-weight="normal" fill="#B6B6B6">
                <tspan x="439" y="446">Lorem ipsum dolor sit amet, </tspan>
                <tspan x="439" y="463">consectetur adipiscing elit. </tspan>
                <tspan x="439" y="480">Aenean porttitor semper </tspan>
                <tspan x="439" y="497">lacus. Vivamus fermentum </tspan>
                <tspan x="439" y="514">dolor et erat tincidunt </tspan>
                <tspan x="439" y="531">tempus. Ut tristique, purus </tspan>
                <tspan x="439" y="548">eget volutpat elementum, </tspan>
                <tspan x="439" y="565">quam purus luctus nisl, in </tspan>
                <tspan x="439" y="582">lacinia metus ipsum aliquet </tspan>
                <tspan x="439" y="599">mi. Praesent eleifend ante </tspan>
                <tspan x="439" y="616">mauris, quis varius tortor </tspan>
                <tspan x="439" y="633">dapibus sit amet. Curabitur </tspan>
                <tspan x="439" y="650">ac lacus sit amet nisi lobortis </tspan>
                <tspan x="439" y="667">blandit quis quis nibh. </tspan>
                <tspan x="439" y="684">Suspendisse convallis </tspan>
                <tspan x="439" y="701">cursus nisi, sit amet varius </tspan>
                <tspan x="439" y="718">neque tincidunt ultrices.</tspan>
            </text>
            <g id="image-placeholder_3x2" transform="translate(24.000000, 32.000000)">
                <g id="placeholder">
                    <rect id="Rectangle-path" fill="#E4E4E4" x="0" y="0" width="591" height="333"></rect>
                    <g id="Group" transform="translate(210.000000, 102.000000)" fill="#FFFFFF">
                        <path d="M0.0278866667,0.516666667 L0.0278866667,127.55 L171.972113,127.55 L171.972113,0.516666667 L0.0278866667,0.516666667 Z M158.112507,113.773667 L14.6951933,113.773667 L14.6951933,14.2917667 L158.112507,14.2917667 L158.112507,113.773667 L158.112507,113.773667 Z" id="Shape"></path>
                        <polygon id="Shape" points="22.4084 103.0597 51.9689067 73.8802667 63.0442467 78.3807 97.9907333 43.4465333 111.746587 58.8977333 117.952087 55.3778 151.532707 103.0597"></polygon>
                        <ellipse id="Oval" cx="56.5576933" cy="39.5886667" rx="12.9100667" ry="12.1236667"></ellipse>
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
      selectVal: 'Hero Image w\' Three Columns'
    }
  }
}
</script>
