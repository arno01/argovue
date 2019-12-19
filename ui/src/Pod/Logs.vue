<template>
<term :content="logs" :title="Logs"></term>
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
    let re = await this.$api.get(`/logs/pod/${this.namespace}/${this.name}/${this.container}`)
    this.logs = re.data
  }
}
</script>