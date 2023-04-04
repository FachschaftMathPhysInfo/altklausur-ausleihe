<script setup>
import { RouterView } from "vue-router";
</script>

<template>
  <v-app-bar :elevation="2"></v-app-bar>
  <v-btn>Tom</v-btn>
  <v-icon icon="mdi-home" />
  <RouterView />
</template>

<style scoped></style>
