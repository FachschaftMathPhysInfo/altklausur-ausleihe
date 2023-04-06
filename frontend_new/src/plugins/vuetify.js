// Styles
import '@mdi/font/css/materialdesignicons.css'
import 'vuetify/styles'
import { aliases, mdi } from 'vuetify/iconsets/mdi'

// Vuetify
import { createVuetify } from 'vuetify'

import 'vuetify/styles'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import { VDataTable } from 'vuetify/labs/VDataTable'

const mpiThemeLight = {
  dark: false,
  colors: {
    primary: '#990000',
    secondary: '#b0bec5',
    accent: '#8c9eff',
    error: '#b71c1c',
  }
}

const mpiThemeDark = {
  dark: true,
  colors: {
    primary: "#590d08",
  }
}

export default createVuetify({
  
  components: { ...components, VDataTable
  },
  directives,
  // https://vuetifyjs.com/en/introduction/why-vuetify/#feature-guides
  icons: {
    defaultSet: 'mdi',
    aliases,
    sets: {
      mdi,
    }
  },
  theme: {
    defaultTheme: 'mpiThemeLight',
    themes: {
      mpiThemeLight, mpiThemeDark
    }
  }
  
})



