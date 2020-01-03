<template>
<div class="dag mx-auto" :style="`width: ${this.size.width}px; height: ${this.size.height}px`">
  <node v-for="id in nodes" v-bind:key="id" :node="getnode(id)" :name="name" :namespace="namespace"></node>
  <edge v-for="edge in edges" v-bind:key="edgekey(edge)" :edge="edge"></edge>
</div>
</template>

<script>
import Dagre from 'dagre'
import Edge from '@/Workflow/DAG/Edge'
import Node from '@/Workflow/DAG/Node'

const W = 152
const H = 42

function removeRetryNodes(nodes) {
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
}

export default {
  props: ['content', 'namespace', 'name'],
  data () {
    return {
      size: {},
      graph: {},
      nodes: [],
      edges: [],
    }
  },
  components: {
    node: Node,
    edge: Edge,
  },
  created () {
    this.graph = new Dagre.graphlib.Graph()
    if (this.content.status) {
      this.update()
    }
  },
  methods: {
    update () {
      let wfNodes = removeRetryNodes(JSON.parse(JSON.stringify(this.content.status.nodes)))
      this.graph.setGraph({})
      this.graph.setDefaultEdgeLabel(() => ({}))
      Object.values(wfNodes).forEach(
        (node) => {
          this.isVirtual(node) ?
          this.graph.setNode(node.id, {width: 1, height: 1, ...node}) :
          this.graph.setNode(node.id, {width: W, height: H, ...node})
        })
      Object.values(wfNodes).forEach(
        (node) => (node.children || []).forEach(
          childId => wfNodes[childId]? this.graph.setEdge(node.id, childId) : ''
          )
        )
      const onExitHandlerNodeId = Object.values(wfNodes).find(node => node.name === `${this.content.metadata.name}.onExit`)
      if (onExitHandlerNodeId) {
        this.getOutboundNodes(this.content.metadata.name).forEach(nodeId => this.graph.setEdge(nodeId, onExitHandlerNodeId))
      }
      Dagre.layout(this.graph)
      var edges = []
      this.graph.edges().forEach(edgeInfo => {
        const edge = this.graph.edge(edgeInfo)
        var lines = []
        if (edge.points.length > 1) {
          for (let i = 1; i < edge.points.length; i++) {
            const toNode = wfNodes[edgeInfo.w]
            lines.push({
              x1: edge.points[i - 1].x, y1: edge.points[i - 1].y,
              x2: edge.points[i].x, y2: edge.points[i].y,
              noArrow: this.isVirtual(toNode)
            })
          }
        }
        edges.push({from: edgeInfo.v, to: edgeInfo.w, lines});
      })
      this.size = this.getGraphSize(this.graph.nodes().map(id => this.graph.node(id)))
      this.nodes = this.graph.nodes()
      this.edges = edges
    },
    isVirtual(node) {
      return (node.type === 'StepGroup' || node.type === 'DAG' || node.type === 'TaskGroup') && !!node.boundaryID
    },
    getGraphSize(nodes) {
      let width = 0
      let height = 0
      nodes.forEach(node => {
        if (node) {
          width = Math.max(node.x + node.width / 2, width)
          height = Math.max(node.y + node.height / 2, height)
        }
      })
      return {width, height}
    },
    getOutboundNodes(nodeID) {
      const node = this.content.status.nodes[nodeID]
      if (node.type === 'Pod' || node.type === 'Skipped') {
        return [node.id]
      }
      let outbound = [];
      for (const outboundNodeID of node.outboundNodes || []) {
        const outNode = this.props.workflow.status.nodes[outboundNodeID]
        if (outNode.type === 'Pod') {
          outbound.push(outboundNodeID)
        } else {
          outbound = outbound.concat(this.getOutboundNodes(outboundNodeID))
        }
      }
      return outbound
    },
    getnode(id) {
      return this.graph.node(id)
    },
    edgekey(edge) {
      return `${edge.from}-${edge.to}`
    }
  },
  watch: {
    content () {
      this.update()
    }
  }
}
</script>

<style scoped>
.dag {
  position: relative;
}
</style>