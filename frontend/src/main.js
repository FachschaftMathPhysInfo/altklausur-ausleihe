import Vue from 'vue';
import App from './App.vue';
import vuetify from './plugins/vuetify';
import router from './router';
import { createProvider } from './vue-apollo'
import i18n from './i18n'
import VueCookies from 'vue-cookies'

Vue.config.productionTip = false;

Vue.use(VueCookies)

new Vue({
  vuetify,
  router,
  apolloProvider: createProvider(),
  i18n,
  render: h => h(App)
}).$mount('#app')
