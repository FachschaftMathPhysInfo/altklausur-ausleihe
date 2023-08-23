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
    locale: 'de',
    fallbackLocale: 'de',
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

// TODOS:
// process.env.VARS einfügen
// SOLVED?: i18n $t in ExamList computed props zum laufen bekommen
// README schreiben mit commands, npm run dev etc 
// example https://www.newline.co/@kchan/building-a-graphql-application-with-vue-3-and-apollo--4679b402