<template>
<b-container fluid>
  <b-row v-for="obj in orderedCache" v-bind:key="obj.metadata.uid">
    <b-col>
      <b-link :to="`/${kind}/${obj.metadata.namespace}/${obj.metadata.name}`">{{obj.metadata.namespace}}/{{ obj.metadata.name }}</b-link>
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
  extends: SSE,
  props: ['name', 'namespace'],
  data () {
    return {
      kind: "catalogue",
    }
  },
  methods: {
    uri() {
      return `/workflow/${this.namespace}/${this.name}/services`
    },
    formatTs(obj) {
      return moment(obj.metadata.creationTimestamp).format("YYYY-MM-DD HH:mm:ss")
    },
    del(instance) {
      return this._action(`/workflow/${this.namespace}/${this.name}/service/${instance}/action/delete`)
    },
  },
}
</script>


