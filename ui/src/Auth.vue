<script>
import axios from "axios";

export default {
  data() {
    return {
      auth: false,
      name: undefined
    };
  },
  methods: {
    isAuth() {
      return this.auth;
    },
    check: async function() {
      let ev = await axios.get("/profile");
      if (ev.data && ev.data.name) {
        this.auth = true;
        this.name = ev.data.name;
      } else {
        this.auth = false;
        this.name = undefined;
      }
    },
    login() {
      window.location.href = "/auth";
    },
    logout() {
      window.location = "/logout";
    }
  }
};
</script>