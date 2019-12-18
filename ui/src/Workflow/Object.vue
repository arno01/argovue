<template>
  <div>
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <h1 class="h2">{{namespace}}/{{kind}}/{{name}}</h1>
    </div>
    <div>
      <b-card no-body>
        <b-tabs card>
          <b-tab title="Nodes" active>
            <nodes :content="object" :namespace="namespace"></nodes>
          </b-tab>
          <b-tab title="Workflow">
            <jsoneditor :content="object"></jsoneditor>
          </b-tab>
          <b-tab title="Control">
            <control :name="name" :namespace="namespace"></control>
          </b-tab>
        </b-tabs>
      </b-card>
    </div>
  </div>
</template>

<script>
import SSE from '@/SSE/Object'
import JsonEditor from '@/JsonEditor'
import Nodes from '@/Workflow/Nodes'
import Control from '@/Workflow/Control'

export default {
  props: ["namespace", "name"],
  extends: SSE,
  components: {
    jsoneditor: JsonEditor,
    nodes: Nodes,
    control: Control,
  },
  data() {
    return {
      kind: "workflows"
    }
  },
}
</script>
