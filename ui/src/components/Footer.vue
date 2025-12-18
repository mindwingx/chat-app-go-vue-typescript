<script setup lang="ts">
const currentYear = new Date().getFullYear()
const appTitle = import.meta.env.VITE_APP_NAME
const more = atob("R2l0SHViOiA8YSBocmVmPSJodHRwczovL2dpdGh1Yi5jb20vbWluZHdpbmd4IiB0YXJnZXQ9Il9ibGFuayI+TWluZFdpbmdYPC9hPg==")
</script>

<template>
  <footer
      class="footer text-center text-sm text-gray-400 py-8 fixed bottom-0 left-0 right-0 bg-gray-900 border-t border-gray-700">
     {{ appTitle }} <small>(dev)</small> - {{ currentYear }} | <span v-html="more"></span>
  </footer>
</template>