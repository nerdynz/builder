import Vue from 'vue'
import Notification from 'vue-bulma-notification'
import Message from 'vue-bulma-message'

const NotificationComponent = Vue.extend(Notification)

const openNotification = (propsData = {
  title: '',
  message: '',
  type: '',
  direction: '',
  duration: 4500,
  container: '.notifications'
}) => {
  return new NotificationComponent({
    el: document.createElement('div'),
    propsData
  })
}

export function showNotification (titleOrObj, message, type, direction, duration, container) {
  let propObj = titleOrObj // allowing an object to be passed instead of each param
  if (typeof titleOrObj !== 'object') {
    // titleOrObj wasn't an object so we are using each param
    propObj = {
      title: titleOrObj,
      message,
      type,
      direction,
      duration,
      container
    }
  }
  openNotification(propObj)
}

const MessageComponent = Vue.extend(Message)

const openMessage = (propsData = {
  title: '',
  message: '',
  type: '',
  direction: '',
  duration: 1500,
  container: '.messages'
}) => {
  return new MessageComponent({
    el: document.createElement('div'),
    propsData
  })
}

export function showMessage (titleOrObj, message, type, callback, direction, duration, container) {
  let propObj = titleOrObj // allowing an object to be passed instead of each param
  if (typeof callback === 'function') { // dont close the notification, they need to action it.
    duration = 0
  }

  // are you sure is more of a yes/no question :)
  let testStr = (message + '').toLowerCase()
  let isProbablyYesNo = (testStr.indexOf('are you sure') >= 0) || (testStr.indexOf('want to continue') >= 0)

  if (typeof titleOrObj !== 'object') {
    // titleOrObj wasn't an object so we are using each param
    propObj = {
      title: titleOrObj,
      message,
      type,
      direction,
      duration: duration,
      container,
      showCloseButton: false,
      onConfirmCallback: callback,
      confirmButtonText: isProbablyYesNo ? 'Yes' : 'Ok',
      cancelButtonText: isProbablyYesNo ? 'No' : 'Cancel',
      showConfirmationButtons: (typeof callback === 'function')
    }
  }
  openMessage(propObj)
}
