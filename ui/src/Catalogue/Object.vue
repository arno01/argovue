<template>
  <div>
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <h1 class="h2">{{namespace}}/{{kind}}/{{name}}</h1>
    </div>
    <div>
      <control :object="object" :name="name" :namespace="namespace" style="margin-bottom: 20px"></control>
      <b-card no-body>
        <b-tabs card>
          <b-tab title="Instances" active>
            <instances :name="name" :namespace="namespace" :kind="kind"></instances>
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
import Control from '@/Catalogue/Control'
import Instances from '@/Catalogue/Instances'

export default {
  props: ['namespace', 'name'],
  extends: SSE,
  components: {
    jsoneditor: JsonEditor,
    control: Control,
    instances: Instances,
  },
  data() {
    return {
      kind: 'catalogue'
    }
  },
  methods: {
    uri() {
      return `/catalogue/${this.namespace}/${this.name}`
    },
  }
}
</script>
