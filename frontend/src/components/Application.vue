<template>
  <div>
    <v-app-bar color="primary" dense dark>
      <v-app-bar-nav-icon @click="drawer = true"></v-app-bar-nav-icon>

      <v-toolbar-title>{{ this.$route.name }}</v-toolbar-title>

      <v-spacer></v-spacer>
      <v-spacer></v-spacer>

      <v-text-field
        v-model="search"
        prepend-inner-icon="mdi-magnify"
        label="Filter Altklausuren, z.B. nach Prüfenden oder Veranstaltungen"
        single-line
        hide-details
        clearable
      ></v-text-field>
    </v-app-bar>

    <v-navigation-drawer v-model="drawer" app temporary>
      <v-list-item>
        <v-list-item-content>
          <v-list-item-title class="title">
            Altklausurausleihe
          </v-list-item-title>
          <v-list-item-subtitle>
            deiner Fachschaft MathPhysInfo
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
            <v-list-item-title>{{ item.title }}</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>

    <router-view />
  </div>
</template>

<script>
export default {
  name: "Application",

  data: () => ({
    search: "",
    drawer: null,
    items: [
      {
        title: "Altklausur ausleihen",
        icon: "mdi-file-document",
        action: "exams",
      },
      {
        title: "Altklausur einreichen",
        icon: "mdi-send",
        action:
          "mailto:pruefungsberichte@mathphys.info?subject=Neue Altklausur über die digitale Altklausurausleihe&body=Liebe Fachschaft,%0D%0A%0D%0Aich möchte euch eine neue Altklausur einreichen. Diese wurde im (Sommersemester/Wintersemester) XXXX für das Modul XXXX von der Lehrperson XXXX gestellt.%0D%0A%0D%0AViele Grüße,%0D%0A%0D%0AAnhang nicht vergessen, am liebsten als PDF oder TeX Datei",
      },
      {
        title: "GitHub Projekt",
        icon: "mdi-github",
        action: "https://github.com/FachschaftMathPhysInfo/altklausur-ausleihe",
        target: "_blank",
      },
      {
        title: "Datenschutz",
        icon: "mdi-file-document-multiple",
        action: "privacy",
      },
      { title: "Impressum", icon: "mdi-information", action: "impress" },
    ],
    right: null,
  }),
};
</script>
