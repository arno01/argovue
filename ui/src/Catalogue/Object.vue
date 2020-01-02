<template>
  <div>
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <h1 class="h2">{{namespace}}/{{kind}}/{{name}}</h1>
    </div>
    <div>
      <b-card no-body>
        <b-tabs card no-key-nav>
          <b-tab title="Deploy" active>
            <deploy :object="object" :name="name" :namespace="namespace"></deploy>
          </b-tab>
          <b-tab title="Resources">
            <resources :name="name" :namespace="namespace" :kind="kind"></resources>
          </b-tab>
          <b-tab title="Instances">
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
import Deploy from '@/Catalogue/Deploy'
import Instances from '@/Catalogue/Instances'
import Resources from '@/Catalogue/Resources'

export default {
  props: ['namespace', 'name'],
  extends: SSE,
  components: {
    jsoneditor: JsonEditor,
    deploy: Deploy,
    instances: Instances,
    resources: Resources,
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
