<template>
<div class="w-100" style="height: 400px"></div>
</template>

<script>
import Vis from 'vis-network'

function color(node) {
  switch(node.phase) {
    case "Failed":
      return "#FB7E81"
    default:
      return "#97C2FC"
  }
}

function shape(node) {
  switch(node.type) {
    case "Pod":
      return "box"
    case "DAG":
      return "ellipse"
    default:
      return "box"
  }
}

export default {
  props: ['content', 'namespace', 'name'],
  data () {
    return {
      nodes: undefined,
      edges: undefined,
      options: {},
      network: undefined,
    }
  },
  mounted () {
    if (this.content.status) {
      this.update()
    }
  },
  methods: {
    update () {
      let wfNodes = this.content.status.nodes
      this.nodes = new Vis.DataSet([])
      this.edges = new Vis.DataSet([])
      this.$log("do update", wfNodes)
      Object.values(wfNodes).forEach( (node) => {
        this.nodes.add([{ id: node.id, label: node.displayName, shape: shape(node), color: color(node), type: node.type }])
      })
      Object.values(wfNodes).forEach( (node) => {
        (node.children || []).forEach( (child) => {
          this.edges.add([{ from: node.id, to: child, arrows: "to" }])
        })
      })
      this.network = new Vis.Network(this.$el, { nodes: this.nodes, edges: this.edges }, this.options)
      this.network.on("doubleClick", (ev) => {
        let node = this.nodes.get(ev.nodes[0])
        if (node && node.type == "Pod") {
          this.$router.push(`/workflow/${this.namespace}/${this.name}/pod/${ev.nodes[0]}`)
        }
      })
    }
  },
  watch: {
    content () {
      this.update()
    }
  }
}
</script>