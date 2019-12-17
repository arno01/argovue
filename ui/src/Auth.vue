<script>
export default {
  data() {
    return {
      auth: "undefined",
      name: undefined
    };
  },
  methods: {
    isAuth() {
      return this.auth == "true"
    },
    isNot() {
      return this.auth == "false"
    },
    check: async function() {
      let ev = await this.$axios.get("/profile");
      if (ev.data && ev.data.name) {
        this.auth = "true"
        this.name = ev.data.name
      } else {
        this.auth = "false"
        this.name = undefined
      }
    },
    login() {
      window.location.href = this.$base+"/auth";
    },
    logout() {
      window.location = this.$base+"/logout";
    }
  }
};
</script>