<script setup></script>

<template>
  <v-container>
    <div style="max-width: 760px; margin: 0px auto">
      <v-textarea
        readonly
        no-resize
        outlined
        rows="7"
        value="Fachschaften Mathematik, Physik, Informatik 
der Verfassten Studierendenschaft der Universität Heidelberg 
Im Neuenheimer Feld 205
69120 Heidelberg 
Tel. +49 6221 54-14999 
Fax +49 6221 54-14999
altklausuren@mathphys.stura.uni-heidelberg.de"
      >
      </v-textarea>

      {{ $t("impress.intro") }}
      <a
        href="https://www.stura.uni-heidelberg.de/impressum/"
        target="_blank"
        >{{ $t("impress.intro1") }}</a
      >{{ $t("impress.intro2") }}
      <a
        href="https://www.stura.uni-heidelberg.de/vs-strukturen/referatekonferenz-refkonf/vorsitz/"
        target="_blank"
        >{{ $t("impress.vorsitz") }}</a
      >
      {{ $t("impress.intro3") }}
    </div>
  </v-container>
</template>

<style scoped></style>
