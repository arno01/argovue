<template>
<div class="w-100" style="height: 400px"></div>
</template>

<script>
import Vis from 'vis-network'

export default {
  props: ['content'],
  data () {
    return {
      nodes: undefined,
      edges: undefined,
      options: {},
      network: undefined,
    }
  },
  created () {
    this.nodes = new Vis.DataSet([])
    this.edges = new Vis.DataSet([])
  },
  mounted () {
    this.update()
    this.network = new Vis.Network(this.$el, { nodes: this.nodes, edges: this.edges }, this.options)
    this.network.fit()
  },
  methods: {
    update () {
      let wfNodes = this.content.status.nodes
      let nodeAlias = {}
      Object.values(wfNodes).forEach( (node) => {
        if (node.type == 'Retry') {
          if (node.children) {
            nodeAlias[node.id] = node.children[0]
          }
        } else {
          nodeAlias[node.id] = node.id
          this.nodes.add([{ id: node.id, label: node.displayName }])
        }
      })
      Object.values(wfNodes).forEach( (node) => {
        if (node.type != 'Retry' && node.children) {
          node.children.forEach( (child) => {
            this.edges.add([{ from: node.id, to: nodeAlias[child], arrows: "to" }])
          })
        }
      })
    },
  },
  watch: {
    content () {
    }
  }
}
</script>

<style>
div {
outline: none !important;
}
</style>