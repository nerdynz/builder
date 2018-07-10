<template>
  <building-block :title="name" :block="block" :removeable="false">
    <div class="columns is-multiline">
      <div class="column is-half is-offset-6">
      </div>
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
        <div class="block-placeholder box u-s-mt" style="height: 300px;">
          Messenger Widget
        </div>
      </div>
    </div>
  </building-block>
</template>

<script>
import BuildingBlock from './BuildingBlock'

let meta = {
  name: 'Contact Messenger',
  svg: ``,
  category: 'Dev Only' // basic, interactive, form, custom, dev only
}

export default {
  // COMPONENT
  // ______________________________________
  name: meta.name.replace(/ /g, ''),
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
