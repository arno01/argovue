<script>
import axios from "axios"

export default {
  data() {
    return {
      auth: "undefined",
      username: undefined,
      profile: {},
      baseURL: '',
      $axios: undefined,
    };
  },
  created () {
    if (window.kubevue && window.kubevue.api) {
      this.baseURL = window.kubevue.api
    } else {
      this.baseURL = ''
    }
    this.$axios = axios.create({ baseURL: this.base })
  },
  methods: {
    redirect(url) {
      window.location.href = this.baseURL+url
    },
    get(url) {
      return this.$axios.get(url)
    },
    post(url) {
      return this.$axios.post(url)
    },
    sse(url, onMessage) {
      let es = new EventSource(this.baseURL+url)
      es.onerror = (err) => this.$log("SSE", err)
      es.onmessage = onMessage
      return es
    },
    isAuth() {
      return this.auth == "true"
    },
    isNot() {
      return this.auth == "false"
    },
    verifyAuth: async function() {
      let ev = await this.$axios.get("/profile");
      if (ev.data && ev.data.name) {
        this.auth = "true"
        this.username = ev.data.name
        this.profile = ev.data
      } else {
        this.auth = "false"
        this.username = undefined
      }
    },
    login() {
      this.redirect('/auth')
    },
    logout() {
      this.redirect('/logoug')
    }
  }
};
</script>