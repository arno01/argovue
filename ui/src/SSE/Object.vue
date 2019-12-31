<script>
export default {
  data() {
    return {
      object: {},
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
    uri() {
      return `/watch/${this.namespace}/${this.kind}/${this.name}`
    },
    parent() {
      return `/watch/${this.namespace}/${this.kind}`
    },
    setupStream() {
      this.es = this.$api.sse(this.uri(), (event) => {
        let msg = JSON.parse(event.data);
        let obj = msg.Content;
        switch (msg.Action) {
          case "delete":
            this.$router.replace(this.parent())
            break
          case "add":
            this.$set(this, "object", obj)
            break
          case "update":
            this.$set(this, "object", obj)
            break
        }
      })
    }
  }
};
</script>