<template>
  <field :label="label" :for="name">
    <div v-if="type === 'date'" class="field has-addons">
      <p class="control u-fw">
        <b-datepicker ref="field" :value="dateVal" @input="dateChanged"
          :label="label"
          :disabled="disabled"
          placeholder="" class="title-input-child" 
          :date-formatter="dateFormatter"
          icon="calendar-o"
        >
        </b-datepicker>
      </p>
    </div>
    <div v-else-if="type === 'datetime'" class="field has-addons">
      <p class="control u-fw" >
        <b-datepicker ref="field" :value="value" @input="dateChanged"
          :label="label"
          :disabled="disabled"
          placeholder="" class="title-input-child" 
        >
        </b-datepicker>
      </p>
      <p class="control" v-if="!disabled">
        <a class="button" @click="openDatePicker"><i class="fa fa-calendar"></i></a>
      </p>
    </div>
    <div v-else-if="type === 'select'" class="field has-addons"> 
      <p class="control">
        <b-select :value="value" 
        :disabled="disabled"
        :placeholder="placeholder || '-- please select --'"
        @input="input">
          <option 
            v-for="(option, index) in options" 
            :key="index"
            :value="option[optionValue]">
              {{option[optionLabel]}}
          </option>
        </b-select>
      </p>
      <p v-if="showButtons" class="control">
        <button class="button" :class="klass" @click="decrement">
          <i class="fa fa-minus"></i>
        </button>
      </p>
      <p v-if="showButtons" class="control">
        <button class="button" :class="klass" @click="increment">
          <i class="fa fa-plus"></i>
        </button>
      </p>
    </div>
    <div v-else-if="type === 'number'" class="field has-addons numeric"> 
      <p class="control">
        <input
          ref="field"
          :id="name" 
          :name="name"
          type="number"
          :placeholder="placeholder"
          :disabled="disabled"
          :value="value"
          :class="klass"
          class="numeric-input"
          autocomplete="off"
          :maxlength="maxlength"
          @input="input"
          @change="change"
          @focus="focus"
          @blur="blur" 
        />
      </p>
      <p v-if="showButtons" class="control">
        <button class="button" :class="klass" @click="decrement">
          <i class="fa fa-minus"></i>
        </button>
      </p>
      <p v-if="showButtons" class="control">
        <button class="button" :class="klass" @click="increment">
          <i class="fa fa-plus"></i>
        </button>
      </p>
    </div>
    <p v-else-if="type === 'textarea'" class="control">
      <textarea 
        ref="field" 
        :id="name" 
        :name="name"
        :placeholder="placeholder"
        :disabled="disabled"
        :value="value"
        :class="klass"
        autocomplete="off"
        :maxlength="maxlength"
        @input="input"
        @change="change"
        @focus="focus"
        @blur="blur"
      >
      </textarea>
    </p>
    <p v-else class="control">
      <input 
        ref="field" 
        :id="name" 
        :name="name"
        :type="type"
        :placeholder="placeholder"
        :disabled="disabled"
        :value="value"
        :class="klass"
        autocomplete="off"
        :maxlength="maxlength"
        @input="input"
        @change="change"
        @focus="focus"
        @blur="blur" 
        @keydown="keydown"
      />
    </p>
  </field>
</template>

<script>
import {fmtDateTime} from '~/helpers/format'

export default {
  // COMPONENT
  // ______________________________________
  name: 'Control',
  inject: {
    $validator: '$validator'
  },
  components: {
  },
  props: {
    hint: {
      Type: null
    },
    name: {
      type: String,
      required: true
    },
    maxlength: String,
    showButtons: {
      type: Boolean,
      default: false
    },
    className: {
      type: String,
      default: ''
    },
    size: String,
    value: {
      required: true
    },
    label: String,
    placeholder: {
      type: String
    },
    type: {
      type: String,
      default: 'text'
    },
    disabled: Boolean,
    prefix: String,
    step: {
      type: Number,
      default: 1
    },
    options: Array,
    optionLabel: {
      type: String,
      default: 'label'
    },
    optionValue: {
      type: String,
      default: 'value'
    },
    tooltip: {
      type: Object
    }
    // errors: Array,
    // validation: {
    //   type: Object,
    //   default () {
    //     return {
    //       rules: {
    //       },
    //       errors: {
    //       }
    //     }
    //   }
    // }
  },
  computed: {
    hasValue () {
      return !!this.value
    },
    klass () {
      let className = this.className + ''

      if (this.type === 'textarea') {
        className += ' textarea '
      } else {
        className += ' input '
      }
      if (this.size) {
        className += ` is-${this.size.toLowerCase()}`
      }
      return className
    },
    dateVal () {
      let value = this.value
      if (this.type.indexOf('date') === 0) {
        if (value) {
          value = new Date(value) // str to date
        } else {
          value = new Date()
        }
        return value
      }
      return this.value
    }
  },
  methods: {
    dateFormatter (date) {
      return fmtDateTime(date, 'dddd, mmmm dS, yyyy')
    },
    placeChanged (a, b, c) {
      this.$emit('placechanged', a, b, c)
    },
    decrement () {
      var newVal = this.value - this.step
      this.$emit('input', parseInt(newVal))
    },
    increment () {
      var newVal = this.value + this.step
      this.$emit('input', parseInt(newVal))
    },
    openDatePicker () {
      this.$refs['field'].datepicker.open()
    },
    clearDatePicker () {
      this.$emit('input', '')
    },
    dateChanged (val) {
      let origDate = fmtDateTime(this.value, 'isoUtcDateTime')
      let newDate = fmtDateTime(val, 'isoUtcDateTime')
      if (newDate === origDate) {
        return // false positive, its just formatting
      }
      this.$emit('input', newDate)
    },
    labelClicked (ev) {
      this.$refs['field'].focus()
    },
    focus (ev) {
      this.isSelected = true
      this.$emit('focus', ev, ev.target.value)
    },
    blur (ev) {
      this.isSelected = false
      this.$emit('blur', ev, ev.target.value)
    },
    keydown (ev) {
      this.$emit('keydown', ev, ev.target.value)
    },
    change (valueOrEvent) {
      let value = valueOrEvent
      if (valueOrEvent && valueOrEvent.target) {
        value = valueOrEvent.target.value
      }

      if (this.type === 'number') {
        this.$emit('change', parseInt(value))
      } else {
        this.$emit('change', value)
      }
    },
    input (valueOrEvent) {
      let value = valueOrEvent
      if (valueOrEvent && valueOrEvent.target) {
        value = valueOrEvent.target.value
      }

      if (this.type === 'number') {
        this.$emit('input', parseInt(value))
      } else {
        this.$emit('input', value)
      }
    }
  },
  data () {
    return {
      isSelected: false
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
  .field {
    .numeric-input {
      width: 80px;
    }
  }
</style>
