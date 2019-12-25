import Vue from 'vue'
import VueRouter from 'vue-router'
import BootstrapVue from 'bootstrap-vue'

import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'

import App from '@/App'
import Front from '@/Front'
import Objects from '@/Objects'
import WorkflowObject from '@/Workflow/Object'
import CatalogueObject from '@/Catalogue/Object'
import PodObject from '@/Workflow/Pod/Object'
import API from '@/API'

Vue.use(VueRouter)
Vue.use(BootstrapVue)

Vue.prototype.$api = new Vue(API)
Vue.prototype.$log = window.console.log.bind(console)

function routeProps({params}) {
  return { namespace: params.namespace, name: params.name, kind: params.kind, pod: params.pod }
}

function workflowPods({params}) {
  return { namespace: params.namespace, name: params.name, pod: params.pod }
}

const router = new VueRouter({
  routes: [
    { path: '/', component: Front },
    { path: '/watch/:kind', component: Objects, props: routeProps },
    { path: '/workflows/:namespace/:name', component: WorkflowObject, props: routeProps },
    { path: '/workflow/:namespace/:name/pod/:pod', component: PodObject, props: workflowPods },
    { path: '/catalogue/:namespace/:name', component: CatalogueObject, props: routeProps },
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