import Vue from 'vue'
import VueRouter from 'vue-router'
import BootstrapVue from 'bootstrap-vue'
import axios from "axios"

import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'

import App from '@/App'
import Front from '@/Front'
import Object from '@/Object'
import Objects from '@/Objects'
import Services from '@/Services'
import WorkflowObject from '@/Workflow/Object'
import PodObject from '@/Pod/Object'
import Auth from '@/Auth'

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
    { path: '/watch/:namespace/pods/:name', component: PodObject, props: routeProps },
    { path: '/watch/:namespace/workflows/:name', component: WorkflowObject, props: routeProps },
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