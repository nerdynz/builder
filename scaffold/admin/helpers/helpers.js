function getParameterByName (name, url) {
  if (!url) url = window.location.href
  name = name.replace(/[\[\]]/g, '\\$&') //   eslint-disable-line
  let regex = new RegExp('[?&]' + name + '(=([^&#]*)|&|#|$)')
  let results = regex.exec(url)
  if (!results) return null
  if (!results[2]) return ''
  return decodeURIComponent(results[2].replace(/\+/g, ' '))
}

function focusElement (selector, index) {
  // setTimeout(() => { // give it a mo to render etc.
  //   let els
  //   let targetEl
  //   // if (selector.indexOf('.') === 0) {
  //   //   // .something
  //   //   selector = selector.split('.')[1]
  //   //   els = document.getElementsByClassName(selector)
  //   // } else if (selector.indexOf('#') === 0) {
  //   //   // #something
  //   //   selector = selector.split('#')[1]
  //   //   els = document.getElementById(selector)
  //   // } else if (selector.indexOf('[') === 0) {
  //   //   // [something]
  //   //   selector = selector.split('[')[1]
  //   //   selector = selector.split(']')[0]
  //   //   els = document.getElementsByName(selector)
  //   // } else {
  //   //   els = document.getElementsByTag(selector)
  //   // }

  //   els = $(selector)
  //   if (els.length > 0) {
  //     if (index === -1) {
  //       targetEl = els[els.length - 1]
  //     } else if (index) {
  //       targetEl = els[index]
  //     } else {
  //       targetEl = els[0]
  //     }
  //   }
  //   let scrollTo = $(targetEl).offset().top + $('.scroller').scrollTop()
  //   $('.scroller').stop().animate({
  //     scrollTop: scrollTo
  //   }, 200)
  //   targetEl.focus()
  // }, 50)
}

export { getParameterByName, focusElement }
