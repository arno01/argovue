<template>
  <div>
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <h1 class="h2">pvc/{{namespace}}/{{name}}</h1>
    </div>
    <div>
      <b-card no-body>
        <b-tabs card no-key-nav>
          <b-tab title="Pvc" active>
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
  props: ["namespace", "name"],
  extends: SSE,
  components: {
    jsoneditor: JsonEditor,
  },
  data() {
    return {
    }
  },
  methods: {
    uri() {
      return `/k8s/pvc/${this.namespace}/${this.name}`
    }
  },
}
</script>
