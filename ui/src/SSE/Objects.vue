<script>
function objKey(obj) {
  return obj.metadata.creationTimestamp
}

export default {
  data() {
    return {
      cache: {},
      kind: undefined,
      es: undefined
    };
  },
  created() {
    this.setupStream()
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
      this.es = this.$api.sse(`/watch/${this.kind}`, (event) => {
        var msg = JSON.parse(event.data)
        var obj = msg.Content
        switch (msg.Action) {
          case "delete":
            this.$delete(this.cache, obj.metadata.uid)
            break
          case "add":
            this.$set(this.cache, obj.metadata.uid, obj)
            break
          case "update":
            this.$set(this.cache, obj.metadata.uid, obj)
            break
        }
      })
    }
  }
}
</script>