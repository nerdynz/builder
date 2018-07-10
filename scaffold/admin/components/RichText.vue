<template>
  <froala :value="value" :config="config"  @input="input"></froala>
</template>
<script>

export default {
  // COMPONENT
  // ______________________________________
  components: {
  },
  props: {
    value: String,
    kind: {
      type: String,
      default: 'standard'
    }
  },
  computed: {
    config () {
      var config = {
        iconsTemplate: 'font_awesome_5r',
        paragraphFormat: {
          N: 'Normal',
          H1: 'Heading 1',
          H2: 'Heading 2',
          H3: 'Heading 3',
          H4: 'Heading 4',
          H5: 'Heading 5',
          H6: 'Heading 6',
          Code: 'Code'
        },
        pluginsEnabled: ['align', 'codeBeautifier', 'codeView', 'entities', 'file', 'image', 'imageManager', 'inlineStyle', 'lineBreaker', 'link', 'lists', 'paragraphFormat', 'paragraphStyle', 'quickInsert', 'quote', 'save', 'table', 'url', 'video', 'wordPaste'],
        toolbarButtons: ['undo', 'redo' | 'paragraphFormat', 'bold', 'italic', 'underline', 'strikeThrough', 'subscript', 'superscript', '|', 'fontFamily', 'fontSize', 'color', 'inlineStyle', 'paragraphStyle', '|', 'align', 'formatOL', 'formatUL', 'outdent', 'indent', 'quote', '-', 'insertLink', 'insertImage', 'insertVideo', 'embedly', 'insertFile', 'insertTable', '|', 'emoticons', 'specialCharacters', 'insertHR', '|', 'selectAll', 'clearFormatting', '|', 'print', 'spellChecker', 'help', 'html', '|']
      }
      if (this.kind === 'standard') {
        config.toolbarButtons = ['undo', 'redo', '|', 'paragraphFormat', 'bold', 'italic', 'underline', '|', 'align', 'formatOL', 'formatUL', 'outdent', 'indent', 'quote', '|', 'insertLink']
      }
      if (this.kind === 'mini') {
        config.toolbarButtons = ['bold', 'italic', 'underline', '|', 'align']
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
  .fr-box.fr-basic {
    .fr-element {
      padding: 2px;
    }
    .fr-toolbar {
      border-top: none;
      .fr-popup .fr-command.fr-btn.fr-dropdown i, .fr-popup .fr-command.fr-btn.fr-dropdown img, .fr-popup .fr-command.fr-btn.fr-dropdown span, .fr-toolbar .fr-command.fr-btn.fr-dropdown i, .fr-toolbar .fr-command.fr-btn.fr-dropdown img, .fr-toolbar .fr-command.fr-btn.fr-dropdown span {
        margin: 6px 12px 6px 6px!important;
      }
      .fr-command.fr-btn, .fr-popup .fr-command.fr-btn {
        width: 32px;
        height: 32px;
        i {
          margin: 6px 9px;
        }
      }
      .fr-separator.fr-vs {
        height: 28px;
      }
    }
  }
</style>
