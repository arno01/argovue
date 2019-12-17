import Vue from 'vue'
import VueRouter from 'vue-router'
import BootstrapVue from 'bootstrap-vue'
import axios from "axios"

import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'

import App from '@/App.vue'
import Front from '@/Front.vue'
import Object from '@/Object.vue'
import Objects from '@/Objects.vue'
import Services from '@/Services.vue'
import Auth from '@/Auth.vue'

Vue.use(VueRouter)
Vue.use(BootstrapVue)

Vue.prototype.$auth = new Vue(Auth)
Vue.prototype.$base = window.kubevue.api
Vue.prototype.$axios = axios.create({ baseURL: window.kubevue.api })

function routeProps(route) {
  return {
    namespace: route.params.namespace,
    name: route.params.name,
    kind: route.params.kind
  }
}

const router = new VueRouter({
  routes: [
    { path: '/', component: Front },
    { path: '/watch/:namespace/services', component: Services, props: routeProps },
    { path: '/watch/:namespace/:kind/:name', component: Object, props: routeProps },
    { path: '/watch/:namespace/:kind', component: Objects, props: routeProps }
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