<template>
  <div class="h-full overflow-y-auto p-8">
    <h1 class="text-2xl font-bold text-gray-800 dark:text-gray-100 mb-4">
      {{ t('home.title') }}
    </h1>
    <p class="text-gray-500 dark:text-gray-400">
      {{ t('home.placeholder') }}
    </p>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
</script>
