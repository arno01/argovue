<template>
<term :content="logs"></term>
</template>

<script>
import Term from '@/Term'

export default {
  props: ["namespace", "name", "container"],
  components: {
    term: Term
  },
  data() {
    return {
      logs: undefined
    }
  },
  created: async function() {
    let re = await this.$api.get(`/k8s/pod/${this.namespace}/${this.name}/container/${this.container}/logs`)
    this.logs = re.data
  }
}
</script>