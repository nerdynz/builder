<template>
  <div>
    <label for="" class="label thin-label">PDF Attachment</label>
    <b-field>
        <b-upload v-model="files" accept="application/pdf" @input="fileChanged">
            <a class="button is-primary" style="z-index:1;">
                <b-icon icon="upload"></b-icon>
                <span v-if="value">Change</span>
                <span v-else>Click to upload</span>
            </a>
        </b-upload>
        <div v-if="value">
            <span class="file-name">
                {{ filename }}
            </span>
        </div>
    </b-field>
  </div>
</template>

<script>
// ... imports

export default {
  // COMPONENT
  // ______________________________________
  name: 'FileUpload',
  components: {},
  props: {
    label: String,
    value: String
  },
  computed: {
    filename () {
      if (this.value.indexOf('/attachments/') >= 0) {
        return this.value.split('/attachments/')[1]
      }
      return this.value
    }
  },
  methods: {
    fileChanged (files) {
      let formData = new FormData() // eslint-disable-line
      formData.append('file', files[0])
      this.$http.post('https://cdn.nerdy.co.nz/upload/displayworks-signs/file', formData).then((resp) => {
        this.$emit('input', resp.data.link)
      })
    }
  },
  watch: {},
  data () {
    return {
      files: []
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

<style lang="scss">

</style>
