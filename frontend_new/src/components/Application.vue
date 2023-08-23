<script setup>
import { RouterView } from "vue-router";
import { ref, onMounted } from "vue";
import { useTheme } from "vuetify";
import { useCookies } from "vue3-cookies";
import { useI18n } from "vue-i18n";

const search = ref("");
const drawer = ref(null);
const items = [
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
];

//Darkmode
const theme = useTheme();
function toggleTheme() {
  return (theme.global.name.value = theme.global.current.value.dark
    ? "mpiThemeLight"
    : "mpiThemeDark");
}

//Language cookies
const { cookies } = useCookies();
const i18n = useI18n();

function switchLanguageInCookie() {
  cookies.set("language", i18n.locale.value, "1y");
}

onMounted(() => {
  if (cookies.get("language")) {
    // get language from cookie
    i18n.locale.value = cookies.get("language");
  } else {
    // set German as default language if none is set in cookie
    cookies.set("language", "de", "1y");
  }

  // Dark theme if set in user preferences
  if (
    window.matchMedia &&
    window.matchMedia("(prefers-color-scheme: dark)").matches
  ) {
    // dark mode
    theme.global.name.value = "mpiThemeDark";
  }

  window
    .matchMedia("(prefers-color-scheme: dark)")
    .addEventListener("change", (e) => {
      theme.global.name.value = theme.global.current.value.dark
        ? "mpiThemeLight"
        : "mpiThemeDark";
    });
});
</script>

<template>
  <div>
    <v-app-bar color="primary" density="compact" dark>
      <v-app-bar-nav-icon @click="drawer = true"></v-app-bar-nav-icon>

      <v-toolbar-title>{{ $t($route.name) }}</v-toolbar-title>

      <v-spacer></v-spacer>

      <v-text-field
        v-model="search"
        prepend-inner-icon="mdi-magnify"
        :label="$t('application.search_label')"
        hide-details
        clearable
        density="compact"
      ></v-text-field>

      <v-spacer></v-spacer>

      <v-tooltip bottom style="margin-right: 12px">
        <template v-slot:activator="{ props }">
          <v-btn v-bind="props" text @click="toggleTheme">
            <v-icon>mdi-theme-light-dark</v-icon>
          </v-btn>
        </template>
        <span>{{ $t("application.toggle_darkmode") }}</span>
      </v-tooltip>

      <v-btn-toggle
        v-model="$i18n.locale"
        @update:model-value="switchLanguageInCookie()"
        mandatory
        density="comfortable"
        color="transparent"
        :border="10"
      >
        <v-btn value="de" icon>
          <img src="/de.svg" />
        </v-btn>
        <v-btn value="en" icon>
          <img src="/en.svg" />
        </v-btn>
      </v-btn-toggle>
    </v-app-bar>

    <v-navigation-drawer v-model="drawer" :order="-1" temporary>
      <v-list-item>
        <v-list-item-title class="title">
          {{ $t("application.title") }}
        </v-list-item-title>
        <v-list-item-subtitle>
          {{ $t("application.yourFS") }}
        </v-list-item-subtitle>
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
          <template v-slot:prepend>
            <v-icon>{{ item.icon }}</v-icon>
          </template>

          <v-list-item-title>{{ $t(item.title) }}</v-list-item-title>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>

    <RouterView />
  </div>
</template>

<style scoped>
.v-btn {
  display: block;
  padding: 2px;
  background: transparent;
}

.v-toolbar {
  position: relative !important;
}
</style>
