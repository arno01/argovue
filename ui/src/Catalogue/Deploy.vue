<template>
<b-form @submit="onSubmit">
  <b-form-group v-for="input in object.spec.input" :key="input.name" :label="input.caption" :label-for="input.name">
    <b-form-input :id="input.name" v-model="data[input.name]" type="text" placeholder=""></b-form-input>
  </b-form-group>
  <b-button type="submit" variant="primary">Deploy</b-button>
</b-form>
</template>

<script>
function serialize(data) {
  return Object.keys(data).map( (k) => ({ name: k, value: data[k] }))
}

export default {
  props: ['object', 'name', 'namespace'],
  data () {
    return {
      data: {}
    }
  },
  methods: {
    input () {
      return this.object && this.object.spec && this.object.spec.input? this.object.spec.input : []
    },
    onSubmit: async function(ev) {
      ev.preventDefault()
      let re = await this.$api.post2(`/catalogue/${this.namespace}/${this.name}/deploy`, serialize(this.data))
      if (re.data.status == 'ok') {
        this.data = {}
      }
      this.$bvToast.toast(`${re.data.action} ${re.data.status} ${re.data.message}`, {
        title: re.data.action,
        toaster: 'b-toaster-bottom-right',
        autoHideDelay: 3000,
        noCloseButton: true,
        variant: re.data.status == 'ok'? 'info' : 'error'
      })
    },
  },
}
</script>
