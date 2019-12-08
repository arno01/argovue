import Vue from 'vue'
import VueRouter from 'vue-router'
import VueSSE from 'vue-sse'
import BootstrapVue from 'bootstrap-vue'

import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'

import App from '@/App.vue'
import Dashboard from '@/Dashboard.vue'
import Objects from '@/Objects.vue'
import Auth from '@/Auth.vue'

Vue.use(VueRouter)
Vue.use(VueSSE)
Vue.use(BootstrapVue)

Vue.prototype.$auth = new Vue(Auth)

const router = new VueRouter({
  routes: [
    { path: '/', component: Dashboard },
    { path: '/objects', component: Objects }
  ]
})

router.beforeEach((to, from, next) => {
  if (router.app.$auth.isAuth()) {
    next()
  } else {
    router.app.$auth.check()
    next()
  }
})

new Vue({
  router,
  render: h => h(App),
}).$mount('#app')