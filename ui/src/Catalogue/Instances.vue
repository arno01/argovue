<template>
  <b-container>
    <b-row v-for="obj in orderedCache" v-bind:key="obj.metadata.uid">
      <b-col>
        <b-link :to="`/catalogue/${namespace}/${name}/instance/${obj.metadata.name}`">{{ obj.metadata.name }}</b-link>
      </b-col>
      <b-col cols=2>
        <b-dropdown variant="link" text="Control" toggle-class="p-0">
          <b-dropdown-item-button @click="del(obj.metadata.name)">Delete</b-dropdown-item-button>
        </b-dropdown>
      </b-col>
      <b-col cols=3 class="text-right">
        {{ formatTs(obj) }}
      </b-col>
    </b-row>
  </b-container>
</template>

<script>
import SSE from '@/SSE/Objects.vue'
import moment from 'moment'

export default {
  props: ['name', 'kind', 'namespace'],
  extends: SSE,
  data () {
    return {
    }
  },
  methods: {
    uri() {
      return `/catalogue/${this.namespace}/${this.name}/instances`
    },
    formatTs(obj) {
      return moment(obj.metadata.creationTimestamp).format("YYYY-MM-DD HH:mm:ss")
    },
    action: async function(instance, action) {
      let re = await this.$api.post(`/catalogue/${this.namespace}/${this.name}/instance/${instance}/action/${action}`)
      this.$bvToast.toast(`${re.data.action} ${re.data.status} ${re.data.message}`, {
        title: re.data.action,
        toaster: 'b-toaster-bottom-right',
        autoHideDelay: 3000,
        noCloseButton: true,
        variant: re.data.status == 'ok'? 'info' : 'error'
      })
    },
    del(instance) {
      this.action(instance, "delete")
    }
  },
}
</script>
