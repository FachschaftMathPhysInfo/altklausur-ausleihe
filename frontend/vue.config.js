module.exports = {
  "transpileDependencies": [
    "vuetify",
  ],

  devServer: {
    host: 'localhost',
  },

  pluginOptions: {
    i18n: {
      locale: 'de',
      fallbackLocale: 'de',
      localeDir: 'locales',
      enableInSFC: true,
      enableBridge: false
    }
  }
}
