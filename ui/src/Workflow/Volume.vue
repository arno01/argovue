<template>
<b-row>
  <b-col>{{namespace}}/{{content.name}}</b-col>
</b-row>
</template>

<script>
export default {
  props: ['content', 'namespace'],
  data () {
    return {
    }
  },
  methods: {
  }
}
</script>