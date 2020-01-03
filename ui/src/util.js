export default {
  removeRetryNodes (nodes) {
  var toReplace = {}
  var toRemove = {}
  Object.values(nodes).forEach( (node) => {
    if (node.type == 'Retry' && node.children && node.children.length > 0) {
      var last = node.children.length-1
      toRemove[node.id] = true
      node.children.forEach( (nodeId, i) => i != last? toRemove[nodeId] = true : '')
      toReplace[node.id] = node.children[last]
    }
  })
  var re = {}
  Object.values(nodes).forEach( (node) => {
    if (!toRemove[node.id]) {
      if (node.children) {
        var children = node.children.
          map( (nodeId) => toReplace[nodeId]? toReplace[nodeId] : nodeId )
        node.children = children
      }
      re[node.id] = node
    }
  })
  return re
  },
  deepCopy (thing) {
    return JSON.parse(JSON.stringify(thing))
  }
}