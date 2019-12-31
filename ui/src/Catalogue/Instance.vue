<template>
  <div>
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <h1 class="h2">{{namespace}}/{{instance}}</h1>
    </div>
    <div>
      <b-card no-body>
        <b-tabs card>
          <b-tab title="Proxy" active>
            <a v-for="port in object.spec.ports" :key="port.port" target="_blank" :href="proxy_uri(port.port)">
              {{ instance }}:{{ port.port }}
            </a>
          </b-tab>
          <b-tab title="Service">
            <jsoneditor :content="object"></jsoneditor>
          </b-tab>
        </b-tabs>
      </b-card>
    </div>
  </div>
</template>

<script>
import SSE from '@/SSE/Object'
import JsonEditor from '@/JsonEditor'

export default {
  props: ['namespace', 'name', 'instance'],
  extends: SSE,
  components: {
    jsoneditor: JsonEditor,
  },
  methods: {
    proxy_uri(port) {
      return this.$api.uri(`/proxy/${this.namespace}/${this.instance}/${port}`)
    },
    uri() {
      return `/catalogue/${this.namespace}/${this.name}/instance/${this.instance}`
    },
  },
  data() {
    return {
      kind: "services",
      object: {
        spec: {
          ports: [],
        }
      }
    }
  },
}
</script>
