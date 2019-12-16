<template>
  <div>
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <h1 class="h2">{{namespace}}/{{objects}}</h1>
    </div>
    <div v-for="obj in orderedCache" v-bind:key="obj.metadata.uid">
      {{ obj.metadata.name }} {{ obj.status? obj.status.phase : "" }}
    </div>
  </div>
</template>

<script>
function objKey(obj) {
  return obj.metadata.name
}

export default {
  props: ["namespace", "objects"],
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
    this.tearDown()
  },
  computed: {
    orderedCache: function() {
     return Object.values(this.cache).sort( (a, b) => objKey(a).localeCompare(objKey(b)) )
    }
  },
  watch: {
    objects () {
      this.tearDown()
      this.setupStream()
    }
  },
  methods: {
    tearDown() {
      if (this.es) {
        this.es.close()
      }
      this.cache = {}
      this.es = undefined
    },
    setupStream() {
      this.es = new EventSource(`/watch/${this.namespace}/${this.objects}`);
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
