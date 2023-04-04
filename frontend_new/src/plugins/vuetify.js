import { createVuetify } from 'vuetify'
import { aliases, mdi } from 'vuetify/iconsets/mdi'

export default createVuetify({
  icons: {
    defaultSet: 'mdi',
    aliases,
    sets: {
      mdi,
    }
  },
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
})
