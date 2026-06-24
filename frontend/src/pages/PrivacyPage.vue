<template>
  <div>
    <!-- Sticky header — CSS scroll-driven animation, transparent bg -->
    <div class="sticky-header sticky top-0 z-2 pt-10 pb-4 px-8 rounded-tl-md
         bg-linear-to-b from-gray-50 from-60% to-transparent
         dark:from-gray-900 dark:to-transparent">
      <!-- Back button -->
      <button
        @click="$router.back()"
        class="back-btn flex items-center gap-2 text-gray-500 dark:text-gray-400
               hover:text-gray-700 dark:hover:text-gray-200 transition-colors"
      >
        <Icon size="20">
          <ChevronLeft20Regular />
        </Icon>
        <span class="text-sm">{{ t('settings.backToSettings') }}</span>
      </button>

      <!-- Title -->
      <h1
        class="header-title text-2xl font-bold text-gray-800 dark:text-gray-100 mt-1"
      >
        {{ t('settings.systemPermissions') }}
      </h1>
    </div>

    <!-- Content -->
    <div class="px-8 pt-8 pb-8">
      <PrivacySettings :app-name="t('app.placeholderAppName')" />
    </div>

    <!-- Bottom fade mask -->
    <div
      class="sticky bottom-0 h-6 z-10 pointer-events-none
             bg-linear-to-t from-gray-50 via-gray-50 to-transparent
             dark:from-gray-900 dark:via-gray-900"
    />
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { Icon } from '@vicons/utils'
import ChevronLeft20Regular from '@vicons/fluent/es/ChevronLeft20Regular'
import PrivacySettings from '@/components/PrivacySettings.vue'

const { t } = useI18n()
</script>

<style scoped>
/* Scroll-driven collapsing header */
@keyframes shrink-header {
  from {
    padding-top: 2.5rem;
    padding-bottom: 1rem;
  }
  to {
    padding-top: 1.5rem;
    padding-bottom: 0.5rem;
  }
}

@keyframes shrink-title {
  from {
    font-size: 1.5rem;
    line-height: 2rem;
  }
  to {
    font-size: 1.125rem;
    line-height: 1.75rem;
  }
}

@keyframes shrink-back {
  from {
    font-size: 0.875rem;
    opacity: 1;
  }
  to {
    font-size: 0.75rem;
    opacity: 0.6;
  }
}

.sticky-header {
  animation: shrink-header linear both;
  animation-timeline: scroll();
  animation-range: 0px 60px;
}

.header-title {
  animation: shrink-title linear both;
  animation-timeline: scroll();
  animation-range: 0px 60px;
}

.back-btn {
  animation: shrink-back linear both;
  animation-timeline: scroll();
  animation-range: 0px 60px;
}
</style>
