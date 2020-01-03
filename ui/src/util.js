function hex2a(hex) {
  var str = ''
  for (var i = 0; i < hex.length; i += 2) {
    str += String.fromCharCode(parseInt(hex.substr(i, 2), 16))
  }
  return str
}

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
  },
  phase (phase) {
    switch (phase) {
      case "Running":
        return "R"
      case "Pending":
        return "P"
      case "Succeeded":
        return "S"
      case "Failed":
        return "F"
      case "Unknown":
        return "U"
      default:
        return phase
    }
  },
  owner(obj) {
    if (obj && obj.metadata) {
      if (obj.metadata.labels['oidc.argovue.io/id']) {
        return hex2a(obj.metadata.labels['oidc.argovue.io/id'])
      } else if (obj.metadata.labels['oidc.argovue.io/group']) {
        return obj.metadata.labels['oidc.argovue.io/group']
      } else {
        return "unknown"
      }
    } else {
      return "undefined"
    }
  },
}