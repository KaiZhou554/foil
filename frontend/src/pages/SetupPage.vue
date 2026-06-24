<template>
  <div class="h-full bg-gray-50 dark:bg-gray-900 overflow-hidden">
      <SetupCard
        :app-name="t('app.defaultAppName')"
        theme-color="oklch(0.716 0.191 2.477)"
        @finish="handleFinish"
      >
        <PrivacySettings :app-name="t('app.placeholderAppName')" />

        <p class="text-gray-500/40 dark:text-gray-400/40 mt-6 text-right text-sm">
          {{ t('app.versionPrefix') }}{{ store.meta.displayVersion }}
        </p>
      </SetupCard>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import SetupCard from '@/components/SetupCard.vue'
import PrivacySettings from '@/components/PrivacySettings.vue'
import { usePrivacyConfigStore } from '@/stores/privacyConfigStore'
import { useAppStore } from '@/stores/appStore'

const { t } = useI18n()
const router = useRouter()
const store = usePrivacyConfigStore()
const appStore = useAppStore()

function handleFinish() {
  appStore.completeOnboarding()
  router.push('/main/home')
}

// ── Ctrl+Enter keyboard shortcut ──
function onKeydown(e: KeyboardEvent) {
  if ((e.ctrlKey || e.metaKey) && e.key === 'Enter') {
    e.preventDefault()
    handleFinish()
  }
}

onMounted(() => document.addEventListener('keydown', onKeydown))
onUnmounted(() => document.removeEventListener('keydown', onKeydown))
</script>
