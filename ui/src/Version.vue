<template>
<div>
  <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
    <h1 class="h2">Versions</h1>
  </div>
  <b-container>
    <b-row>
      <b-col cols=3>App</b-col>
      <b-col md=auto>{{version.version}}</b-col>
      <b-col md=auto>{{version.builddate}}</b-col>
      <b-col md=auto>{{version.commit}}</b-col>
    </b-row>
    <b-row>
      <b-col cols=3>Kubernetes</b-col>
      <b-col md=auto>{{version.kubernetes.gitVersion}}</b-col>
      <b-col md=auto>{{version.kubernetes.buildDate}}</b-col>
    </b-row>
  </b-container>
</div>
</template>

<script>
export default {
  data () {
    return {
      version: {}
    }
  },
  created: async function() {
    var re = await this.$api.get("/version")
    this.$log("version", re.data)
    this.version = re.data
  }
}
</script>
