<template>
  <b-container>
    <b-row v-for="obj in orderedCache" v-bind:key="obj.metadata.uid">
      <b-col>
        <b-link :to="`/k8s/${obj.kind.toLowerCase()}/${obj.metadata.namespace}/${obj.metadata.name}`">
          {{ obj.kind.toLowerCase() }}/{{ obj.metadata.namespace }}/{{ obj.metadata.name }}
        </b-link>
      </b-col>
      <b-col cols=2 v-if="obj.status" class="text-right">
        {{ obj.status.phase }}
      </b-col>
      <b-col cols=3 class="text-right">
        {{ formatTs(obj) }}
      </b-col>
    </b-row>
  </b-container>
</template>

<script>
import SSE from '@/SSE/Objects.vue'

export default {
  props: ['name', 'kind', 'namespace'],
  extends: SSE,
  data () {
    return {
    }
  },
  methods: {
    uri() {
      return `/catalogue/${this.namespace}/${this.name}/resources`
    },
  },
}
</script>