<script setup>
import SeismicMap from './components/SeismicMap.vue'
</script>

<template>
  <main>
    <SeismicMap />
  </main>
</template>

<style>
/* Reset global para evitar márgenes blancos alrededor del mapa */
body {
  margin: 0;
  padding: 0;
  overflow: hidden;
}
</style>
