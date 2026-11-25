<template>
  <app-shell />
</template>

<script setup lang="ts">
import AppShell from "#/widgets/layout/AppShell.vue"
import { useAppInit } from "#/features/session/model/useAppInit"

const { init } = useAppInit()
init()
</script>

<style>
/* App.vue now only acts as a root container */
</style>
