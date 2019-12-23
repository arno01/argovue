<template>
<b-button-toolbar>
  <b-button-group size="sm" class="mr-1">
    <b-button @click="deploy()">Deploy</b-button>
  </b-button-group>
</b-button-toolbar>
</template>

<script>
export default {
  props: ['object', 'name', 'namespace'],
  data () {
    return {
      nodes: []
    }
  },
  methods: {
    status(status) {
      return this.object && this.object.status && this.object.status.phase == status
    },
    action: async function(action) {
      let re = await this.$api.post(`/catalogue/${this.namespace}/${this.name}/${action}`)
      this.$bvToast.toast(`${re.data.action} ${re.data.status} ${re.data.message}`, {
        title: re.data.action,
        toaster: 'b-toaster-bottom-right',
        autoHideDelay: 3000,
        noCloseButton: true,
        variant: re.data.status == 'ok'? 'info' : 'error'
      })
    },
    deploy () {
      this.action('deploy')
    },
  },
}
</script>
