<template>
  <b-container>
    <b-row v-for="obj in orderedCache" v-bind:key="obj.metadata.uid">
      <b-col>
        <b-link :to="`/catalogue/${namespace}/${name}/instance/${obj.metadata.name}`">{{ obj.metadata.name }}</b-link>
      </b-col>
      <b-col cols=3>
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
      return moment(obj.metadata.creationTimestamp).format("YYYY-MM-DD HH:mm:ss1")
    }
  },
}
</script>
