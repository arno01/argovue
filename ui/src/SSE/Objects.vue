<script>
function objKey(obj) {
  return obj.metadata.creationTimestamp
}

export default {
  props: ["namespace"],
  data() {
    return {
      cache: {},
      kind: undefined,
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
     return Object.values(this.cache).sort( (a, b) => objKey(b).localeCompare(objKey(a)) )
    }
  },
  watch: {
    kind () {
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
      this.es = new EventSource(`/watch/${this.namespace}/${this.kind}`);
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