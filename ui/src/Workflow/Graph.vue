<template>
<div class="w-100" style="height: 600px"></div>
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
      Object.values(wfNodes).forEach( (node) => {
        this.nodes.add([{ id: node.id, label: node.displayName }])
        if (node.children) {
          node.children.forEach( (child) => {
            this.edges.add([{ from: node.id, to: child, arrows: "to" }])
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
