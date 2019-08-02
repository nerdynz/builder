<template>
  <froala :value="value" :config="config" @input="input" />
</template>
<script>

export default {
  // COMPONENT
  // ______________________________________
  components: {
  },
  props: {
    value: {
      type: String,
      default: ''
    },
    kind: {
      type: String,
      default: 'standard'
    }
  },
  computed: {
    config () {
      let config = {
        iconsTemplate: 'font_awesome_5',
        paragraphFormat: {
          N: 'Normal',
          H1: 'Heading 1',
          H2: 'Heading 2',
          H3: 'Heading 3',
          H4: 'Heading 4',
          H5: 'Heading 5'
        },
        paragraphStyles: {
          // 'is-uppercase': 'Uppercase',
          // 'is-lowercase': 'Lowercase',

          'super-bold': 'Super Bold',
          'has-text-grey-dark': 'Dark Grey',
          'has-text-grey': 'Grey',
          'has-text-grey-light': 'Light Grey',
          'has-text-primary': 'EHS Red'
        },
        pluginsEnabled: ['align', 'codeBeautifier', 'codeView', 'entities', 'file', 'image', 'imageManager', 'inlineStyle', 'lineBreaker', 'link', 'lists', 'paragraphFormat', 'paragraphStyle', 'quickInsert', 'quote', 'save', 'table', 'url', 'video', 'wordPaste'],
        toolbarButtons: ['undo', 'redo' | 'paragraphFormat', 'paragraphStyle', 'bold', 'italic', 'underline', 'strikeThrough', 'subscript', 'superscript', '|', 'fontFamily', 'fontSize', 'color', 'inlineStyle', 'paragraphStyle', '|', 'align', 'formatOL', 'formatUL', 'outdent', 'indent', 'quote', '-', 'insertLink', 'insertImage', 'insertVideo', 'embedly', 'insertFile', 'insertTable', '|', 'emoticons', 'specialCharacters', 'insertHR', '|', 'selectAll', 'clearFormatting', '|', 'print', 'spellChecker', 'help', 'html', '|']
      }
      if (this.kind === 'standard') {
        config.toolbarButtons = ['undo', 'redo', '|', 'paragraphFormat', 'paragraphStyle', 'color', 'bold', 'italic', 'underline', '|', 'insertHR', 'align', 'formatOL', 'formatUL', 'outdent', 'indent', 'quote', '|', 'insertLink']
      } else if (this.kind === 'media') {
        config.toolbarButtons = ['undo', 'redo', '|', 'paragraphFormat', 'paragraphStyle', 'color', 'bold', 'italic', 'underline', '|', 'insertHR', 'align', 'formatOL', 'formatUL', 'outdent', 'indent', 'quote', '|', 'insertLink']
      } else if (this.kind === 'mini') {
        config.pluginsEnabled = ['link']
        config.toolbarButtons = ['bold', 'italic', 'underline', '|', 'align', 'insertLink', 'align']
      } else if (this.kind === 'micro') {
        config.pluginsEnabled = ['link', 'paragraphFormat']
        config.toolbarButtons = ['paragraphFormat', '|', 'bold', 'italic', 'underline', '|', 'align', 'insertLink', 'align']
      }
      config.toolbarInline = this.kind === 'inline' || this.kind === 'inline-simple'
      return config
    }
  },
  methods: {
    input (val) {
      this.$emit('input', val)
    }
  },
  watch: {},
  data () {
    return {
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
  .fr-dropdown-list a,
  .content {
    h1, h1 a{
      font-weight: 700 !important;
      strong {
        font-weight: 900 !important
      }
      font-size: 3.6rem !important;
      font-weight: 900 !important;
      margin: 0 !important;
    }
    h2, h2 a{
      font-weight: 700 !important;
      strong {
        font-weight: 900 !important
      }
      font-size: 2.6rem !important;
      margin: 0 !important;
    }
    h3, h3 a{
      font-weight: 700 !important;
      strong {
        font-weight: 900 !important
      }
      font-size: 2.1rem !important;
      margin: 0 !important;
    }
    h4, h4 a{
      strong {
        font-weight: 900 !important
      }
      font-size: 1.8rem !important;
      margin: 0 !important;
    }
    h5, h5 a{
      strong {
        font-weight: 900 !important
      }
      font-size: 1.5rem !important;
      margin: 0 !important;
    }
    p {
      font-size: 1.1rem;
      margin: 0 !important;
    }
    strong {
      color: inherit;
    }
    a {
      color: #dd2b26;
    }

    hr {
      background: #d0cfd4;
    }
  }
  .content h1:not(:last-child),.content h2:not(:last-child),.content h3:not(:last-child),.content h4:not(:last-child), .content h5:not(:last-child),.content h6:not(:last-child),
  .content p:not(:last-child), .content dl:not(:last-child), .content ol:not(:last-child), .content ul:not(:last-child), .content blockquote:not(:last-child), .content pre:not(:last-child), .content table:not(:last-child) {
    margin-bottom: 0;
  }

  .super-bold {
    font-weight: 900;
  }

</style>
