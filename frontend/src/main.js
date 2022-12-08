import Vue from 'vue';
import App from './App.vue';
import { createProvider } from './vue-apollo'
import VueCookies from 'vue-cookies'

Vue.config.productionTip = false;

Vue.use(VueCookies)

new Vue({
  apolloProvider: createProvider(),
  render: h => h(App)
}).$mount('#app')
