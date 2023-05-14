import Vue from 'vue'

import AsyncComputed from 'vue-async-computed'
import Cookies from 'js-cookie'

import 'normalize.css/normalize.css'

import Element from 'element-ui'
import './styles/element-variables.scss'

import '@/styles/index.scss'

import App from './App'
import router from './router'
router.beforeEach((to, from, next) => {
  if (to.meta.title) {
    document.title = to.meta.title
  }
  next()
})

import './icons'
import * as filters from './filters'

Vue.use(Element, {
  size: Cookies.get('size') || 'mini'
})
Vue.use(AsyncComputed)

Object.keys(filters).forEach(key => {
  Vue.filter(key, filters[key])
})

Vue.config.productionTip = false

new Vue({
  el: '#app',
  router,
  render: h => h(App)
})
