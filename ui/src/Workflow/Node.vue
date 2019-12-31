<template>
<b-row>
  <b-col v-if="isPod()">
    <b-link :to="`/workflow/${namespace}/${name}/pod/${content.id}`">{{content.displayName}}</b-link>
  </b-col>
  <b-col v-else>
    {{content.displayName}}
  </b-col>
  <b-col cols=2>{{content.type}}</b-col>
  <b-col cols=2>{{content.phase}}</b-col>
  <b-col cols=1>{{duration}}s</b-col>
</b-row>
</template>

<script>
import moment from 'moment'
import 'moment-duration-format'

export default {
  props: ['content', 'name', 'namespace'],
  data () {
    return {
      duration: 0
    }
  },
  created () {
    setInterval(() => this.calcDuration(), 1000)
    this.calcDuration()
  },
  methods: {
    isPod () {
      return this.content.type == 'Pod'
    },
    calcDuration() {
      let start = moment(this.content.startedAt)
      let end = this.content.finishedAt? moment(this.content.finishedAt) : moment()
      let duration = moment.duration(end.unix() - start.unix(), 'seconds')
      this.duration = duration.format("h:mm:ss")
    },
  }
}
</script>