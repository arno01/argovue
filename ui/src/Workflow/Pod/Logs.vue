<template>
<term :content="logs" :title="Logs"></term>
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
      logs: undefined
    }
  },
  created: async function() {
    let re = await this.$api.get(`/workflow/${this.namespace}/${this.name}/pod/${this.pod}/container/${this.container}/logs`)
    this.logs = re.data
  }
}
</script>