<template>
<b-container fluid>
  <volume v-for="v in volumes" v-bind:key="v.name" :content="v" :namespace="namespace" />
</b-container>
</template>

<script>
import Volume from '@/Workflow/Volume'

export default {
  props: ['content', 'namespace'],
  components: {
    volume: Volume,
  },
  data () {
    return {
      volumes: []
    }
  },
  methods: {
  },
  watch: {
    content (c) {
      let volumes = JSON.parse(JSON.stringify(c.status.persistentVolumeClaims)) // deep copy
      this.$set(this, "volumes", volumes)
    }
  }
}
</script>


