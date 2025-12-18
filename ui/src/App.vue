<template>
  <div class="container max-w-[600px] mx-auto bg-gray-800 rounded-lg shadow h-[80vh] flex flex-col overflow-hidden mt-8">
   <router-view />
  </div>
  <Footer />

</template>

<script setup lang="ts">
import Footer from './components/Footer.vue'
</script>