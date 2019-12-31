<template>
  <div>
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <h1 class="h2">{{namespace}}/{{kind}}/{{name}}</h1>
    </div>
    <div>
      <b-card no-body>
        <b-tabs card no-key-nav>
          <b-tab title="Pod" active>
            <jsoneditor :content="object"></jsoneditor>
          </b-tab>
          <b-tab v-for="container in containers" v-bind:key="container.name" :title="container.name" lazy>
            <logs :name="name" :namespace="namespace" :pod="pod" :container="container.name"></logs>
          </b-tab>
        </b-tabs>
      </b-card>
    </div>
  </div>
</template>

<script>
import SSE from '@/SSE/Object'
import JsonEditor from '@/JsonEditor'
import Logs from '@/Pod/Logs'

export default {
  props: ["namespace", "name"],
  extends: SSE,
  components: {
    jsoneditor: JsonEditor,
    logs: Logs,
  },
  data() {
    return {
      kind: "pods",
      containers: [],
    }
  },
  methods: {
    uri() {
      return `/k8s/pod/${this.namespace}/${this.name}`
    }
  },
  watch: {
    object (c) {
      this.containers = c.spec.containers
    }
  }
}
</script>
