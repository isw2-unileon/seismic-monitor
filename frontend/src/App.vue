<script setup>
import { RouterView } from 'vue-router'
</script>

<template>
  <RouterView />
</template>

<style>
/* Global reset: Map and login should occupy all available space */
body {
  margin: 0;
  padding: 0;
  font-family: system-ui, -apple-system, sans-serif;
  background-color: #1a1a2e; /* Add global background color to avoid white flashes */
}

#app {
  width: 100%;
  min-height: 100vh;
}
</style>