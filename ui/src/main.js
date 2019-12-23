import Vue from 'vue'
import VueRouter from 'vue-router'
import BootstrapVue from 'bootstrap-vue'

import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'

import App from '@/App'
import Front from '@/Front'
import Object from '@/Object'
import Objects from '@/Objects'
import WorkflowObject from '@/Workflow/Object'
import PodObject from '@/Pod/Object'
import ServiceObject from '@/Service/Object'
import API from '@/API'

Vue.use(VueRouter)
Vue.use(BootstrapVue)

Vue.prototype.$api = new Vue(API)
Vue.prototype.$log = window.console.log.bind(console)

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
    { path: '/watch/:namespace/services/:name', component: ServiceObject, props: routeProps },
    { path: '/watch/:namespace/pods/:name', component: PodObject, props: routeProps },
    { path: '/watch/:namespace/workflows/:name', component: WorkflowObject, props: routeProps },
    { path: '/watch/:namespace/:kind/:name', component: Object, props: routeProps },
    { path: '/watch/:namespace/:kind', component: Objects, props: routeProps },
    { path: '/watch/:kind', component: Objects, props: routeProps }
  ]
})

router.beforeEach((to, from, next) => {
  if (router.app.$api.isAuth()) {
    next()
  } else {
    router.app.$api.verifyAuth()
    next()
  }
})

new Vue({
  router,
  render: h => h(App),
}).$mount('#app')