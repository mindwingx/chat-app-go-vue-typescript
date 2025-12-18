<template>
    <div class="p-6 flex flex-col gap-3 items-center justify-center flex-1">
        <div>
        <NotFound :size="250" class="mb-20"/>
        </div>
        <router-link class="text-center w-full max-w-[240px] bg-blue-500 text-white px-3 py-2 rounded-md cursor-pointer text-sm disabled:opacity-50 disabled:cursor-not-allowed" to="/">Go Back</router-link>
    </div>
</template>

<script setup lang="ts">
import NotFound from './../components/icons/NotFound.vue';
</script>