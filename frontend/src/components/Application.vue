<template>
  <div>
    <v-app-bar color="primary" dense dark>
      <v-app-bar-nav-icon @click="drawer = true"></v-app-bar-nav-icon>

      <v-toolbar-title>{{ $t(this.$route.name) }}</v-toolbar-title>

      <v-spacer></v-spacer>

      <v-text-field
        v-model="search"
        prepend-inner-icon="mdi-magnify"
        :label="$t('application.search_label')"
        single-line
        hide-details
        clearable
      ></v-text-field>

      <v-spacer></v-spacer>

      <v-tooltip bottom style="margin-right: 12px;">
        <template v-slot:activator="{ on, attrs }">
          <v-btn
            v-bind="attrs"
            v-on="on"
            text
            v-on:click="$vuetify.theme.dark = !$vuetify.theme.dark"
          >
            <v-icon>mdi-theme-light-dark</v-icon>
          </v-btn>
        </template>
        <span>{{ $t("application.toggle_darkmode") }}</span>
      </v-tooltip>

      <v-btn-toggle
        v-model="$i18n.locale"
        @change="switchLanguageInCookie()"
        mandatory
        dense
        background-color="transparent"
        borderless
      >
        <v-btn value="de" icon>
          <img src="/de.svg" />
        </v-btn>
        <v-btn value="en" icon>
          <img src="/en.svg" />
        </v-btn>
      </v-btn-toggle>
    </v-app-bar>

    <v-navigation-drawer v-model="drawer" app temporary>
      <v-list-item>
        <v-list-item-content>
          <v-list-item-title class="title">
            {{ $t("application.title") }}
          </v-list-item-title>
          <v-list-item-subtitle>
            {{ $t("application.yourFS") }}
          </v-list-item-subtitle>
        </v-list-item-content>
      </v-list-item>

      <v-divider></v-divider>

      <v-list nav>
        <v-list-item
          v-for="item in items"
          :key="item.title"
          link
          :href="item.action"
          :target="item.target"
        >
          <v-list-item-icon>
            <v-icon>{{ item.icon }}</v-icon>
          </v-list-item-icon>

          <v-list-item-content>
            <v-list-item-title>{{ $t(item.title) }}</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>

    <router-view />
  </div>
</template>

<script>
import Vue from "vue";
import i18n from "../i18n";
export default {
  name: "Application",

  data: () => ({
    search: "",
    drawer: null,
    items: [
      {
        title: "application.get_exam",
        icon: "mdi-file-document",
        action: "exams",
      },
      {
        title: "application.hand_in",
        icon: "mdi-send",
        action:
          "mailto:pruefungsberichte@mathphys.info?subject=Neue Altklausur über die digitale Altklausurausleihe&body=Liebe Fachschaft,%0D%0A%0D%0Aich möchte euch eine neue Altklausur einreichen. Diese wurde im (Sommersemester/Wintersemester) XXXX für das Modul XXXX von der Lehrperson XXXX gestellt.%0D%0A%0D%0AViele Grüße,%0D%0A%0D%0AAnhang nicht vergessen, am liebsten als PDF oder TeX Datei",
      },
      {
        title: "application.github_project",
        icon: "mdi-github",
        action: "https://github.com/FachschaftMathPhysInfo/altklausur-ausleihe",
        target: "_blank",
      },
      {
        title: "application.dataprivacy",
        icon: "mdi-file-document-multiple",
        action: "privacy",
      },
      {
        title: "application.impress",
        icon: "mdi-information",
        action: "impress",
      },
    ],
    right: null,
  }),
  methods: {
    switchLanguageInCookie: () => {
      Vue.$cookies.set("language", i18n._vm.locale, "1y");
    },
  },
  mounted() {
    if (Vue.$cookies.get("language")) {
      // get language from cookie
      i18n._vm.locale = Vue.$cookies.get("language");
    } else {
      // set German as default language if none is set in cookie
      Vue.$cookies.set("language", "de", "1y");
    }

    // Dark theme if set in user preferences
    if (
      window.matchMedia &&
      window.matchMedia("(prefers-color-scheme: dark)").matches
    ) {
      // dark mode
      this.$vuetify.theme.dark = true;
    }

    window
      .matchMedia("(prefers-color-scheme: dark)")
      .addEventListener("change", (e) => {
        this.$vuetify.theme.dark = e.matches;
      });
  },
};
</script>
