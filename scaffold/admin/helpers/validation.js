/**
 * Validator class
 * @author  Gustavo Ocanto <gustavoocanto@gmail.com>
 */

export default class Validator {
  /**
   * Create a new instance.
   *
   * @param {Object} data
   * @param {Object} rules
   * @param {Object} messages
   */
  constructor (data, rules, messages) {
    // contains the validation errors.
    this.errors = []

    // the data to be checked out.
    this.data = data

    // the rules required.
    this.rules = rules

    // Error messages container.
    this.messages = messages
  }

  /**
   * Create a new static instance.
   *
   * @param {Object} data
   * @param {Object} rules
   * @param {Object} messages
   */
  static make (data, rules, messages) {
    let validate = new Validator(data, rules, messages)
    return validate.handle()
  }

  static install (Vue, options) {
    // default messages
    let defaultMessages = {
      'url': 'The field format is invalid.',
      'numeric': 'The field must be a number.',
      'integer': 'The field must be an integer.',
      'required': 'Please complete this field',
      'alpha': 'The field may only contain letters.',
      'email': 'Please enter a valid email.',
      'dateISO': 'Please enter a valid date with a ISO format.',
      'alphaNum': 'The field may only contain letters and numbers.',
      'phone': 'Please enter a valid phone number.',
      'blank': 'The field is required and does not allow blank spaces.'
    }

    // console.log(Vue)

    Vue.prototype.$validation = {
      defaultMessages: options & options.messages ? options.messages : defaultMessages,
      make: function (data, rules, messages) {
        messages = messages || this.defaultMessages
        let validate = new Validator(data, rules, messages)
        return validate.handle()
      }
    }
  }
  /**
   * Walk through the validations rules.
   *
   * @returns {Object}
   */
  handle () {
    let methods = null
    // console.log(this.rules)
    for (let rule in this.rules) {
      methods = this.rules[rule].split(',')
      this.evaluate(methods, rule)
    }

    return this.errors
  }

  /**
   * Evaluate the input against rules.
   *
   * @param {Array} methods
   * @param {String} field
   * @returns {Array}
   */
  evaluate (methods, field) {
    let value = this.data[field]
    // We return if the rule required does not exist and the field value is empty.
    if (this.isNotRequired(methods, value)) {
      return
    }

    for (let method in methods) {
      // if the rule required exits and there was an error, the
      // stack errors method is called to keep track of rules
      // that did not pass the validation.
      if (this[methods[method]] && !this[methods[method]](value)) {
        this.stackErrors({
          key: field, // evaluated field.
          error: this.messages[methods[method]]
        })
      }
    }
  }

  /**
   * Check whether a field is not required and has an empty value.
   *
   * @param {Array} methods
   * @param {String} value
   * @return {Boolean}
   */
  isNotRequired (methods, value) {
    // if the rule required does not exist and the field value is empty.
    if ((methods.indexOf('required') === -1) && this.empty(value)) {
      return true
    }

    return false
  }

  /**
   * Keep errors tracked out.
   *
   * @param {Object} data
   * @return {Void}
   */
  stackErrors (data) {
    if (!this.errors.find((error) => error.key === data.key)) {
      let error = data.error
      if (error.indexOf('[[FIELDNAME]]') > 0) {
        error = error.split('[[FIELDNAME]]').join(splitTitleCase(data.key))
      }
      this.errors[data.key] = error
    }
  }

  /**
   * Check whether the field is blank.
   *
   * @param {String} value
   * @return {Boolean}
   */
  empty (value) {
    return value == null || value.length === 0 || (typeof (value) === 'string' && value.trim() === '')
  }

  /**
   * No blank fields.
   *
   * @param {String} value
   * @return {Boolean}
   */
  required (value) {
    if (typeof value === 'boolean') {
      return value
    }

    return !this.empty(value)
  }

  /**
   * Numeric rule.
   *
   * @param {String} value
   * @return {Boolean}
   */
  numeric (value) {
    return (/^-?(?:0$0(?=\d*\.)|[1-9]|0)\d*(\.\d+)?$/).test(value)
  }

  /**
   * Integer rule.
   *
   * @param {String} value
   * @return {Boolean}
   */
  integer (value) {
    return (/^(-?[1-9]\d*|0)$/).test(value)
  }

  /**
   * Digits rule.
   *
   * @param {String} value
   * @return {Boolean}
   */
  digits (value) {
    return (/^[\d() \.\:\-\+#]+$/).test(value) // eslint-disable-line no-useless-escape
  }

  /**
   * Alpha rule.
   *
   * @param {String} value
   * @return {Boolean}
   */
  alpha (value) {
    return (/^[a-zA-Z]+$/).test(value)
  }

  /**
   * Alpha num rule.
   *
   * @param {String} value
   * @return {Boolean}
   */
  alphaNum (value) {
    return !(/\W/).test(value)
  }

  /**
   * Emails rule.
   *
   * @param {String} value
   * @return {Boolean}
   */
  email (value) {
    return (/^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/).test(value) // eslint-disable-line no-useless-escape
  }

  /**
   * Url rule.
   *
   * @param {String} value
   * @return {Boolean}
   */
  url (value) {
    return (/^(https?|ftp|rmtp|mms):\/\/(([A-Z0-9][A-Z0-9_-]*)(\.[A-Z0-9][A-Z0-9_-]*)+)(:(\d+))?\/?/i).test(value)
  }

  /**
   * Length rule.
   *
   * @param {String} value
   * @return {Boolean}
   */
  length (value) {
    console.warn("this length check doesn't work, always returns true")
    return true
  // return value && value.length === +arg
  }

  /**
   * Check whether the values has a blank value.
   *
   * @param {String} value
   * @return {Boolean}
   */
  blank (value) {
    return value && value.trim() !== ''
  }

  /**
   * Check whether the values has IOS date format.
   *
   * @param {String} value
   * @return {Boolean}
   */
  dateISO (value) {
    return (/^\d{4}[\/\-](0?[1-9]|1[012])[\/\-](0?[1-9]|[12][0-9]|3[01])$/).test(value) // eslint-disable-line no-useless-escape
  }

  /**
   * Check whether a given value has a valid US phone number format.
   *
   * @param {String} value
   * @return {Boolean}
   */
  phone (value) {
    /*
    | ----------------------------------------------------------
    | VALID FORMATS
    | ----------------------------------------------------------
    |
    | 123-456-7890, (123) 456-7890, 123 456 7890
    | 123.456.7890, +91 (123) 456-7890
    |
    */

    // eslint-disable-next-line no-useless-escape
    return (/^[\+]?[(]?[0-9]{3}[)]?[-\s\.]?[0-9]{3,4}[-\s\.]?[0-9]{3,4}$/gmi).test(value)
  }
}

function splitTitleCase (str) {
  let returnValue = str + ''
  returnValue = returnValue.replace('_day', ' Day')
  returnValue = returnValue.replace('_month', ' Month')
  returnValue = returnValue.replace('_year', ' Year')
  returnValue = returnValue.replace(/([a-z])([A-Z])/g, '$1 $2')
  return returnValue
}
