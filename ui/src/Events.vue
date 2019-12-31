<template><div></div></template>

<script>
export default {
  data() {
    return {
      es: undefined
    };
  },
  created() {
    this.setupStream()
  },
  destroyed() {
    this.tearDown()
  },
  watch: {
    name () {
      this.tearDown()
      this.setupStream()
    }
  },
  methods: {
    tearDown() {
      if (this.es) {
        this.es.close()
      }
      this.object = {}
      this.es = undefined
    },
    setupStream() {
      this.es = this.$api.sse("/events", (event) => this.$log("event:", event))
      this.es.onerror = (err) => {
        this.$log("events error:", err)
        this.$api.logout()
      }
    }
  }
};
</script>