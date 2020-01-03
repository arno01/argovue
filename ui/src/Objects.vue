<template>
  <div>
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <h1 class="h2">{{kind}}</h1>
    </div>
    <b-container>
      <b-row v-for="obj in orderedCache" v-bind:key="obj.metadata.uid">
        <b-col>
          <b-link :to="`/${kind}/${obj.metadata.namespace}/${obj.metadata.name}`">{{obj.metadata.namespace}}/{{ obj.metadata.name }}</b-link>
        </b-col>
        <b-col cols=2 v-if="isGroup(obj)" class="text-right">
          {{ owner(obj) }}
        </b-col>
        <b-col cols=2 v-if="obj.status" class="text-right">
          {{ obj.status.phase }}
        </b-col>
        <b-col cols=3 class="text-right">
          {{ formatTs(obj) }}
        </b-col>
      </b-row>
    </b-container>
  </div>
</template>

<script>
import SSE from '@/SSE/Objects.vue'

function hex2a(hex) {
  var str = '';
  for (var i = 0; i < hex.length; i += 2) str += String.fromCharCode(parseInt(hex.substr(i, 2), 16));
  return str;
}

export default {
  props: ["kind"],
  extends: SSE,
  data() {
    return {
    }
  },
  methods: {
    owner(obj) {
      if (obj.metadata) {
        return hex2a(obj.metadata.labels['oidc.argovue.io/id']) || obj.metadata.labels['oidc.argovue.io/group'] || "unknown"
      }
    },
    isGroup (obj) {
      return obj.metadata && obj.metadata.labels && obj.metadata.labels['oidc.argovue.io/group']
    }
  },
};
</script>
