<template>
  <div class="h-full bg-gray-50 dark:bg-gray-900 flex items-center justify-center p-4">
      <div
        class="w-full max-w-lg p-10 text-center"
      >
        <!-- App icon / logo -->
        <img
          :src="AppIcon"
          alt=""
          class="w-20 h-20 mx-auto mb-6 rounded-2xl"
        />

        <!-- Slogan -->
        <h1 class="text-2xl font-bold text-gray-800 dark:text-gray-100 mb-2">
          {{ t('welcome.title') }}
        </h1>
        <p class="text-gray-500 dark:text-gray-400 mb-8">
          {{ t('welcome.slogan') }}
        </p>

        <!-- Language selector (dropdown) -->
        <div class="mb-8">
          <label class="block text-sm font-medium text-gray-600 dark:text-gray-300 mb-3">
            {{ t('welcome.selectLanguage') }}
          </label>
          <n-select
            :value="appStore.currentLanguage"
            :options="languageOptions"
            @update:value="selectLanguage"
            class="w-48 mx-auto"
            size="large"
          />
        </div>

        <!-- Get started button -->
        <BubbleButton @click="goToSetup">
          {{ t('welcome.getStarted') }}
        </BubbleButton>
      </div>
    </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { NSelect } from 'naive-ui'
import { useAppStore } from '@/stores/appStore'
import AppIcon from '@/assets/appicon_128px.webp'
import BubbleButton from '@/components/BubbleButton.vue'

const { t, locale } = useI18n()
const router = useRouter()
const appStore = useAppStore()

const languageOptions = [
  { value: 'zh-CN', label: '中文' },
  { value: 'en', label: 'English' },
]

function selectLanguage(value: string) {
  const lang = value as 'zh-CN' | 'en'
  appStore.setLanguage(lang)
  locale.value = lang
  appStore.saveConfig()
}

function goToSetup() {
  router.push('/setup')
}
</script>


