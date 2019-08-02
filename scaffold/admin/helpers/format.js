import dateformat from 'dateformat'
import fromNow from 'time-from-now'

function fmtTimeAgo (date) {
  let d = new Date(date)
  d = fromNow(d)
  return d === 'now' ? d : d + ' ago'
}

function fmtDate (date) {
  return dateformat(date, 'dS mmm yyyy')
}

function fmtDateTime (date, format) {
  return dateformat(date, format)
}

function fmtTime (date) {
  return dateformat(date, 'h:MM TT')
}

function fmtCost (cost) {
  return parseFloat(cost).toFixed(2)
}

function titleCase (str) {
  return str.replace(/\w\S*/g, function (txt) { return txt.charAt(0).toUpperCase() + txt.substr(1).toLowerCase() })
}

function splitTitleCase (str) {
  return str.replace(/([a-z])([A-Z])/g, ' $1 $2')
}

function fmtCamel (str) {
  return str
    .replace(/\s(.)/g, function ($1) { return $1.toUpperCase() })
    .replace(/\s/g, '')
    .replace(/^(.)/, function ($1) { return $1.toLowerCase() })
}

function fmtKebab (str) {
  let result = str

  // Convert camelCase capitals to kebab-case.
  result = result.replace(/([a-z][A-Z])/g, function (match) {
    return match.substr(0, 1) + '-' + match.substr(1, 1).toLowerCase()
  })

  // Convert non-camelCase capitals to lowercase.
  result = result.toLowerCase()

  // Convert non-alphanumeric characters to hyphens
  result = result.replace(/[^-a-z0-9]+/g, '-')

  // Remove hyphens from both ends
  result = result.replace(/^-+/, '').replace(/-$/, '')

  return result
}

function fmtTextFromNum (num) {
  if (num === 1) {
    return ''
  }
  if (num === 2) {
    return 'Two'
  }
  if (num === 3) {
    return 'Three'
  }
  if (num === 4) {
    return 'Four'
  }
  if (num === 5) {
    return 'Five'
  }
  if (num === 6) {
    return 'Six'
  }
  if (num === 7) {
    return 'Seven'
  }
  if (num === 8) {
    return 'Eight'
  }
  if (num === 9) {
    return 'Nine'
  }
  if (num === 10) {
    return 'Ten'
  }
  if (num === 11) {
    return 'Eleven'
  }
  if (num === 12) {
    return 'Twelve'
  }
  if (num === 13) {
    return 'Thirteen'
  }
  if (num === 14) {
    return 'Forteen'
  }
  if (num === 15) {
    return 'Fifteen'
  }
  if (num === 16) {
    return 'Sixteen'
  }
  if (num === 17) {
    return 'Seventeen'
  }
  if (num === 18) {
    return 'Eighteen'
  }
  if (num === 19) {
    return 'Nineteen'
  }
  return 'Zero'
}

function fmtUpperCaseFirstLetter (str) {
  return str.charAt(0).toUpperCase() + str.slice(1)
}

// format suitable for golang
function fmtUtcDateTimez (val) {
  return dateformat(val, 'isoUtcDateTime')
}

function fmtAttachment (link) {
  if (!link) return ''
  let str = link
  if (link.indexOf('/attachments/') >= 0) {
    str = link.split('/attachments/')[1]
  }
  return `<a class="is-link" target="_blank" href="${link}">${str}</a>`
}

function truncate (str, length) {
  if (!length) length = 250
  let parts = str.split()
  parts = parts.slice(0, length)
  str = parts.join()
  return str
}

export {
  fmtTimeAgo,
  titleCase,
  splitTitleCase,
  fmtDate,
  fmtCost,
  fmtCamel,
  fmtKebab,
  fmtTextFromNum,
  fmtTime,
  fmtDateTime,
  fmtUtcDateTimez,
  fmtUpperCaseFirstLetter,
  truncate,
  fmtAttachment
}
