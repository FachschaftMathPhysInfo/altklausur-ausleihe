import { createApp } from 'vue'
import App from '../App.vue'

// Vuetify
import 'vuetify/styles'
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'

const vuetify = createVuetify({
  components,
  directives,
})

createApp(App).use(vuetify, {
  theme: {
    themes: {
      light: {
        primary: '#990000',
        secondary: '#b0bec5',
        accent: '#8c9eff',
        error: '#b71c1c',
      },
      dark: {
        primary: "#590d08",
      },
    },
  },
}).mount('#app')