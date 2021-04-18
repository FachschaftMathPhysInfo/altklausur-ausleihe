<template>
  <v-tooltip
    :left="position == 'left'"
    :right="position == 'right'"
    :bottom="position == 'bottom'"
    :top="position == 'top'"
  >
    <template v-slot:activator="{ on }">
      <v-icon v-on="on" @click="clicked" :color="color">{{ icon }}</v-icon>
    </template>
    <span>{{ text }}</span>
  </v-tooltip>
</template>

<script>
export default {
  name: "tooltipped-icon",
  props: {
    icon: String,
    text: String,
    color: String,
    position: {
      type: String,
      default: "bottom",
    },
  },
  methods: {
    clicked() {
      this.$emit("clicked");
    },
  },
};
</script>
