import Vue from 'vue'
import App from './App.vue'
import router from './router'
import axios from 'axios'
import VueAxios from 'vue-axios'
import VueCodemirror from 'vue-codemirror'
import VueHighlightJS from 'vue-highlightjs'

import {BootstrapVue, IconsPlugin} from 'bootstrap-vue'

import 'bootstrap/dist/css/bootstrap.min.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'
import './assets/neon-glow.min.css'
import 'codemirror/lib/codemirror.css'
import 'highlight.js/styles/default.css'

Vue.use(BootstrapVue)
Vue.use(IconsPlugin)
Vue.use(VueAxios, axios)
Vue.use(VueCodemirror)
Vue.use(VueHighlightJS)

Vue.config.productionTip = false
Vue.prototype.$ajax= axios

new Vue({
    router,
    render: function (h) {
        return h(App)
    }
}).$mount('#app')
