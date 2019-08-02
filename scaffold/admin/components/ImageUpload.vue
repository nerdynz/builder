<template>
  <div>
    <div class="image-upload rel-overlay" :style="{'max-width': maxWidth ? maxWidth : '100%', opacity: loading ? 0 : 1 }">
      <b-loading :active="loading" />
      <div v-show="!hasImage" class="is-overlay uploader-msg is-vertical-centered" @click="clickUpload">
        <div class="has-text-centered">
          <i class="far fa-cloud-upload" />
          Click or Drag <br> to upload
        </div>
      </div>
      <div class="cropper-overlay is-overlay is-fixed" />
      <div ref="cropper" class="cropbox" style="width: 100%;" :style="{opacity: loading ? 0 : 1 }">
        <input ref="clicker" type="file" name="thumb" required="required">
      </div>
    </div>
  </div>
</template>

<script>

export default {
  // COMPONENT
  // ______________________________________
  name: 'ImageUpload',
  components: {},
  props: {
    value: {
      type: String,
      default: ''
    },
    width: {
      type: Number,
      default: 640
    },
    height: {
      type: Number,
      default: 480
    },
    maxWidth: {
      type: String,
      default: ''
    },
    forceJpeg: {
      type: Boolean,
      default: false
    },
    url: {
      type: String,
      default: '/api/v1/upload/crop'
    }
  },
  data () {
    return {
      isTransparent: false,
      currentULID: '',
      meta: {
      },
      loading: true,
      hasImage: false
    }
  },
  computed: {
    existing () {
      return this.value && this.value.indexOf('ncdn_') >= 0
    },
    image () {
      if (this.existing) {
        let ogExt = ''
        if (this.value.indexOf('.png') >= 0) {
          ogExt = '.png'
        } else {
          ogExt = '.jpg'
        }
        // ncdn_kdkyr29392_39749.png
        let split = this.value.split('_')
        if (split.length === 3) {
          return split[1] + ogExt
        }
      }
      return this.value
    },
    ext () {
      if (this.forceJpeg) {
        return '.jpg'
      }
      return this.isTransparent ? '.png' : '.jpg'
    }
  },
  watch: {},
  // LIFECYCLE METHODS
  // ______________________________________
  beforeCreate () {
  },
  created () {
  },
  beforeMount () {
  },
  mounted () {
    this.load()
  },
  beforeUpdate () {
  },
  updated () {
  },
  beforeDestroy () {
  },
  destroyed () {
  },
  methods: {
    clickUpload () {
      this.$refs.clicker.click()
    },
    load () {
      this.currentULID = this.$service.ULID()
      if (this.existing) {
        let split = this.value.split('_')
        if (split.length === 3) {
          this.currentULID = split[1]
        }
        this.$service.retrieve('imageMeta', this.currentULID).then((data) => {
          this.meta = data
          this.meta.originalName = '/attachments/' + this.meta.originalName
          this.init()
        }).catch((err) => {
          console.error(err) // eslint-disable-line
          this.meta.originalName = this.value
          this.init()
        })
      } else {
        this.meta.originalName = this.value
        this.init()
      }
    },
    init () {
      this.loading = true
      this.isTransparent = (this.value.split('.').pop() + '') === 'png'
      setTimeout(() => {
        $(this.$refs.cropper).ImageUpload({
          width: this.width,
          height: this.height,
          resize: true,
          smaller: true,
          saveOriginal: true,
          save: false,
          buttonEdit: true,
          meta: this.meta,
          background: 'white',
          // url: this.url ? this.url : '/api/v1/upload/image',
          onAfterInitImage: this.afterInit,
          onSave: this.save,
          onAfterSelectImage: this.imageSelected,
          onAfterRemoveImage: this.remove,
          onLoadFailed: this.remove
        })
      }, 1)
    },
    afterInit () {
      this.loading = false
      if (this.value !== '') {
        this.hasImage = true
      }
    },
    imageSelected (val) {
      this.hasImage = true
    },
    save (data) {
      if (this.existing) {
        data.isExisting = true
        data.oldFilename = this.value
        data.original = '' // dont need it, save bandwidth
      }
      data.uniqueID = this.currentULID // dont keep making a new ULID because of filequeue
      data.originalExt = data.name.split('.').pop().toLowerCase()
      if (data.originalExt.toLowerCase() === 'png') {
        this.isTransparent = true
      }
      data.originalName = `${data.uniqueID}.${data.originalExt}`
      data.name = `ncdn_${data.uniqueID}_${this.$service.ULID()}${this.ext}` // this.ext is dynamic based on isTransparent or forceJpeg
      data.ext = this.ext.replace('.', '')
      // data.data = ''
      this.$axios.post(this.url, data)
      // this.$store.commit('app/ADD_TO_UPLOADS', {ulid: data.currentULID, data: data})
      this.$emit('input', '/attachments/' + data.name)
    },
    remove () {
      this.hasImage = false
      this.loading = false
      this.$emit('input', '')
    }
  }
}
</script>

<style lang="css">
.cropper-overlay {
  background: #161616;
  z-index: 1500;
  display: none;
}
.image-upload {
  opacity: 0;
}
.image-upload-transparent-check {
  text-transform: uppercase;
  font-size: 12px;
  margin-top: -4px;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol";
}
.uploader-msg {
  display: flex;
  justify-content: center;
  align-self: center;
  flex-direction: column;
  z-index: 20;
  background: #DDDDDD;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif; cursor: pointer;
}
.cropbox { border: 1px solid #ccc; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif; }
.cropbox { background-color: #eeeeee; text-align: center; position: relative;  display: inline-block;  }
.cropbox.done:after { content: '';}
.cropbox.done:before { content: '';}
.cropbox.notAnImage { background-color: #f2dede; border-color: #ebccd1; }
.cropbox.notAnImage:after { content: 'The selected file is not an image!'; color: #a94442; }
.cropbox.notAnImage:before { content: ''; color: #ebccd1; }
.cropbox.alert-danger { background-color: #f2dede; }
.cropbox.alert-danger:after { content: ''; }
.cropbox.smalltext:before,
.cropbox.smalltext:after { font-size: 20px; }
.cropbox > span { font-size: 30px; color: #bbbbbb; position: absolute; top: 35%; left: 0; width: 100%; text-align: center; z-index:0;}
.cropbox > span.loader { display: none; }
.cropbox > input[type=file] { position: absolute; top: 0; left: 0; bottom: 0; right: 0; width: 100%; opacity: 0; cursor: pointer; z-index:2; height: 100% /* IE HACK*/ }
.cropbox > input[type=text] { display: none; }
.cropbox .progress { bottom: 10px; left: 10px; right: 10px; display: none;  }
.cropbox .cropWrapper { overflow: hidden; position: absolute; top:0; bottom: 0; left: 0; right: 0; z-index: 10; background-color: #eeeeee;  text-align: left;}
.cropbox img { z-index: 5; position: relative;  max-width: initial;}
.cropbox img.ghost { opacity: .2; z-index: auto; float:left /* HACK for not scrolling*/; }
.cropbox img.main { cursor: move; }
.cropbox .final img.main  { cursor: auto; }
.cropbox img.preview { width: 100%; }
.cropbox .tools { position: absolute; top: 10px; right: 10px; z-index: 999; display: flex; }
.cropbox .tools .button { width: 32px; height: 30px }
.cropbox .tools .saving { width: 100px; height: 30px }
.cropbox .download { position: absolute; bottom: 10px; left: 10px; z-index: 999; display: inline-block; }
.cropbox .download > * { margin: 0 0 0 5px; }
</style>
