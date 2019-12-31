<script>
import axios from "axios"

function getCookie(name) {
  var value = "; " + document.cookie
  var parts = value.split("; " + name + "=")
  if (parts.length == 2) return parts.pop().split(";").shift()
}

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
    if (window.argovue && window.argovue.api_base_url) {
      this.baseURL = window.argovue.api_base_url
    } else {
      this.baseURL = ''
    }
    this.$axios = axios.create({ baseURL: this.baseURL })
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
      let cookie = getCookie("auth-session")
      let ev = await this.$axios.post("/profile", { Cookie: cookie }, { headers: { 'Content-Type': 'application/json' } })
      window.console.log(ev)
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
      this.redirect('/logout')
    }
  }
};
</script>