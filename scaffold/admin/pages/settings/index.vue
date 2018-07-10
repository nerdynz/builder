<template>
  <div>
    <form v-if="record">
      <div class="columns is-multiline">
        <div class="column is-6">
          <div class="field">
            <control label="Top bar Email" name="EmailAddress" type="email" v-model="record.EmailAddress" v-validate="'required'"></control>
          </div>
        </div>
        <div class="column is-6">
        </div>
        <div class="column is-6">
          <div class="field">
            <control label="Top bar Phone" name="Phone" type="phone" v-model="record.Phone" v-validate="'required'"></control>
          </div>
        </div>
        <div class="column is-6">
        </div>
        <div class="column is-6">
          <div class="field">
            <control label="Tagline" v-model="record.Tagline" name="Tagline" type="text" v-validate="'required'"></control>
          </div>
        </div>
        <div class="column is-6">
        </div>
        <div class="column is-6">
          <div class="field">
            <label class="label thin-label">Tagline Image</label>
            <image-upload label="" ratio="16:5" size="1920,600" v-model="record.TaglineImage"></image-upload>
          </div>
        </div>
      </div>
    </form>
  </div>
</template>

<script>
import {mapActions} from 'vuex'
// import { Tabs, TabPane } from 'vue-bulma-tabs'

export default {
  inject: {
    $validator: '$validator'
  },
  components: {
  },

  computed: {
    buttons () {
      return [
        {text: 'Save', alignment: 'left', kind: 'success', click: this.save, isDisabled: this.$validator.errors.any()}
      ]
    }
  },
  methods: {
    ...mapActions({
      setButtons: 'app/setButtons'
    }),

    loadRecord () {
      this.$service.retrieve('settings', 1).then((newRecord) => {
        this.record = newRecord
      })
    },

    save (goBack) {
      if (this.errors.any()) {
        this.$toast.open({title: 'Failed to save!', message: 'Please double check the fields highlighted.', type: 'is-danger'})
        return
      }
      this.$service.save('settings', this.record)
        .then((newRecord) => {
          this.$toast.open({message: `Settings Saved`, type: 'is-success'})
        })
    }
  },

  watch: {
    'errors.items' () {
      console.log('asdf')
      this.setButtons(this.buttons)
    }
  },

  data () {
    return {
      record: null
    }
  },

  // lifecycle method
  created () {
    this.loadRecord()
    this.setButtons(this.buttons)
  }
}
</script>
