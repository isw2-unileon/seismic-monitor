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
  font-family: system-ui, -apple-system, sans-serif;
  background-color: #1a1a2e; /* Añadimos el color de fondo global para evitar flashes blancos */
}

#app {
  width: 100%;
  min-height: 100vh;
}
</style>