<template>
  <div>
    <div
      class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom"
    >
      <h1 class="h2">Objects</h1>
    </div>
    <div
      v-for="obj in cache"
      v-bind:key="obj.metadata.uid"
    >{{ obj.kind }}/{{ obj.metadata.name }}@{{ obj.metadata.namespace }}: {{ obj.status.phase }}</div>
  </div>
</template>

<script>
export default {
  name: "Objects",
  data() {
    return {
      cache: {},
      es: undefined
    };
  },
  created() {
    this.setupStream();
  },
  destroyed() {
    if (this.es) {
      this.es.close()
    }
    this.es = undefined
  },
  methods: {
    setupStream() {
      this.es = new EventSource("/sse");
      this.es.onerror = function(err) {
        window.console.log("sse error", err);
      };
      var self = this;
      this.es.onmessage = function(event) {
        var msg = JSON.parse(event.data);
        var obj = msg.Content;
        switch (msg.Action) {
          case "delete":
            self.$delete(self.cache, obj.metadata.uid);
            break;
          case "add":
            self.$set(self.cache, obj.metadata.uid, obj);
            break;
          case "update":
            self.$set(self.cache, obj.metadata.uid, obj);
            break;
        }
      };
    }
  }
};
</script>
