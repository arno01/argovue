<template>
  <div>
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <h1 class="h2">{{namespace}}/{{instance}}</h1>
    </div>
    <div>
      <b-card no-body>
        <b-tabs card lazy>
          <b-tab title="Resources">
            <resources :name="name" :namespace="namespace" :instance="instance"></resources>
          </b-tab>
          <b-tab title="Release">
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
import Resources from '@/Catalogue/Instance/Resources'

export default {
  props: ['namespace', 'name', 'instance'],
  extends: SSE,
  components: {
    jsoneditor: JsonEditor,
    resources: Resources,
  },
  methods: {
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
