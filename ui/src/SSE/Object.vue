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
      this.es = this.$api.sse(`/watch/${this.namespace}/${this.kind}/${this.name}`, (event) => {
        let msg = JSON.parse(event.data);
        let obj = msg.Content;
        switch (msg.Action) {
          case "delete":
            this.$router.replace(`/watch/${this.namespace}/${this.kind}`)
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