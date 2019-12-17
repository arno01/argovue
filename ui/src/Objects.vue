<template>
  <div>
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <h1 class="h2">{{namespace}}/{{kind}}</h1>
    </div>
    <div v-for="obj in orderedCache" v-bind:key="obj.metadata.uid">
      <b-link :to="`/watch/${obj.metadata.namespace}/${kind}/${obj.metadata.name}`">{{ obj.metadata.name }}</b-link>
    </div>
  </div>
</template>

<script>
import SSE from '@/SSE/Objects.vue'

export default {
  props: ["namespace", "kind"],
  extends: SSE,
  data() {
    return {
      cache: {},
      es: undefined
    };
  },
};
</script>
