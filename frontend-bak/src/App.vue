<script setup lang="ts">
import { RouterView } from 'vue-router'
</script>

<template>
  <router-view v-slot="{ Component, route }">
    <template v-if="route.meta.keepAlive">
      <keep-alive>
        <component :is="Component" :key="route.fullPath" />
      </keep-alive>
    </template>
    <template v-else>
      <component :is="Component" :key="route.fullPath" />
    </template>
  </router-view>
</template>

<style>
#app {
  height: 100vh;
  margin: 0;
  padding: 0;
}
</style>