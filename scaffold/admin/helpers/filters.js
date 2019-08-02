function byUUID (array, uuid) {
  var findUUID = function (obj) {
    return obj.UUID === uuid
  }
  return array.find(findUUID)
}

function indexByUUID (array, uuid) {
  var findUUID = function (obj) {
    return obj.UUID === uuid
  }
  return array.findIndex(findUUID)
}

function indexByField (array, fieldName, value) {
  var find = function (obj) {
    return obj[fieldName] === value
  }
  return array.findIndex(find)
}

function byField (array, fieldName, value) {
  var find = function (obj) {
    return obj[fieldName] === value
  }
  return array.find(find)
}

function changeSortByUUID (records, fromUUID, toUUID) {
  let from = indexByUUID(records, fromUUID)
  let to = indexByUUID(records, toUUID)
  return changeSort(records, from, to)
}

function changeSort (records, from, to) {
  records = [
    ...records
  ]
  var newRecords = arrayMove(records, from, to)
  newRecords.forEach((record, index) => {
    record.SortPosition = (index + 1) * 50
  })
  return newRecords
}

export { changeSortByUUID, byUUID, indexByUUID, indexByField, byField, changeSort }

function arrayMove (arr, from, to) {
  if (to >= arr.length) {
    var k = to - arr.length + 1
    while (k--) {
      arr.push(undefined)
    }
  }
  arr.splice(to, 0, arr.splice(from, 1)[0])
  return arr // for testing
}
