<script setup>
import { RouterView } from 'vue-router'
</script>

<template>
  <RouterView />
</template>

<style>
/* Reset global: El mapa y el login deben ocupar todo el espacio disponible */
body {
  margin: 0;
  padding: 0;
  overflow: hidden;
  font-family: system-ui, -apple-system, sans-serif;
}

#app {
  width: 100vw;
  height: 100vh;
}
</style>