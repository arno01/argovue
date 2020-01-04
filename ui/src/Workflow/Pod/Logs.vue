<template>
<term :lines="logs"></term>
</template>

<script>
import Term from '@/Term'

export default {
  props: ["namespace", "name", "pod", "container"],
  components: {
    term: Term
  },
  data() {
    return {
      logs: [],
      es: undefined,
    }
  },
  created: async function() {
    this.es = this.$api.sse(`/workflow/${this.namespace}/${this.name}/pod/${this.pod}/container/${this.container}/logs`,
      (event) => {
        this.logs.push(event.data)
      }
    )
  },
  destroyed () {
    if (this.es) {
      this.es.close()
    }
    this.es = undefined
  },
}
</script>