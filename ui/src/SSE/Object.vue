<script>
export default {
  data() {
    return {
      object: {},
      name: undefined,
      kind: undefined,
      namespace: undefined,
      es: undefined
    };
  },
  created() {
    this.setupStream();
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
      this.es = new EventSource(`/watch/${this.namespace}/${this.kind}/${this.name}`)
      this.es.onerror = function(err) {
        window.console.log("sse error", err)
      };
      var self = this;
      this.es.onmessage = function(event) {
        var msg = JSON.parse(event.data);
        var obj = msg.Content;
        switch (msg.Action) {
          case "delete":
            // do something, go up?
            break
          case "add":
            self.$set(self, "object", obj)
            break
          case "update":
            self.$set(self, "object", obj)
            break
        }
      };
    }
  }
};
</script>