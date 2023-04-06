import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import vuetify from './plugins/vuetify'
import { loadFonts } from './plugins/webfontloader'
import { createI18n } from 'vue-i18n'
// import translations
import de from "./locales/de.json";
import en from "./locales/en.json";
import { globalCookiesConfig } from 'vue3-cookies'

loadFonts()

// internationalization / translation
const i18n = createI18n({
    allowComposition: true, // you need to specify that!
    locale: process.env.VUE_APP_I18N_LOCALE || 'de',
    fallbackLocale: process.env.VUE_APP_I18N_FALLBACK_LOCALE || 'de',
    localeDir: 'locales',
    messages: { de, en },
  })

// cookies
globalCookiesConfig({
    expireTimes: "1y",
    secure: true,
});

createApp(App)
  .use(router)
  .use(vuetify)
  .use(i18n)
  .mount('#app')
