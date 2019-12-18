<template>
  <div>
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <h1 class="h2">{{namespace}}/{{kind}}/{{name}}</h1>
    </div>
    <div>
      <b-card no-body>
        <b-tabs card>
          <b-tab title="Service" active>
            <jsoneditor :content="object"></jsoneditor>
          </b-tab>
          <b-tab title="Proxy">
            <a v-for="port in object.spec.ports" :key="port.port" target="_blank" :href="`/proxy/${namespace}/${name}/${port.port}`">
              {{ name }}:{{port.port}}
            </a>
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
  props: ['namespace', 'name'],
  extends: SSE,
  components: {
    jsoneditor: JsonEditor,
  },
  data() {
    return {
      kind: "services"
    }
  },
}
</script>
