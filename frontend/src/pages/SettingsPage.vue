<template>
  <div class="h-full overflow-y-auto p-8">
    <h1 class="text-2xl font-bold text-gray-800 dark:text-gray-100 mb-6">
      {{ t('settings.title') }}
    </h1>

    <!-- Language Section -->
    <section class="mb-8">
      <h2 class="text-lg font-medium text-gray-800 dark:text-gray-100 mb-3">
        {{ t('settings.language') }}
      </h2>
      <div class="bg-white dark:bg-gray-800 rounded-xl ring-1 ring-black/5 dark:ring-white/10 p-4">
        <div class="flex items-center justify-between">
          <span class="text-gray-700 dark:text-gray-200">{{ t('settings.displayLanguage') }}</span>
          <div class="w-36">
            <n-select :value="appStore.currentLanguage" :options="languageOptions" @update:value="switchLanguage"
              class="w-40" size="small" />
          </div>
        </div>
      </div>
    </section>

    <!-- General settings -->
    <section class="mb-8">
      <h2 class="text-lg font-medium text-gray-800 dark:text-gray-100 mb-3">
        {{ t('settings.general') }}
      </h2>
      <div
        class="bg-white dark:bg-gray-800 rounded-xl ring-1 ring-black/5 dark:ring-white/10 divide-y divide-gray-100 dark:divide-gray-700/50">
        <div class="flex items-center justify-between px-4 py-3">
          <span class="text-gray-700 dark:text-gray-200">{{ t('settings.autoLaunch') }}</span>
          <n-switch :value="appStore.generalSettings.autoStart" @update:value="val => { appStore.generalSettings.autoStart = val; appStore.saveConfig() }" />
        </div>
        <div class="flex items-center justify-between px-4 py-3">
          <span class="text-gray-700 dark:text-gray-200">{{ t('settings.notifications') }}</span>
          <n-switch :value="appStore.generalSettings.notifications" @update:value="val => { appStore.generalSettings.notifications = val; appStore.saveConfig() }" />
        </div>
        <div class="flex items-center justify-between px-4 py-3">
          <span class="text-gray-700 dark:text-gray-200">{{ t('settings.minimizeToTray') }}</span>
          <n-switch :value="appStore.generalSettings.minimizeToTray" @update:value="val => { appStore.generalSettings.minimizeToTray = val; appStore.saveConfig() }" />
        </div>
      </div>
    </section>

    <!-- System Permissions -->
    <section class="mb-8">
      <h2 class="text-lg font-medium text-gray-800 dark:text-gray-100 mb-3">
        {{ t('settings.systemPermissions') }}
      </h2>
      <div class="bg-white dark:bg-gray-800 rounded-xl ring-1 ring-black/5 dark:ring-white/10 p-4">
        <button @click="goToPrivacy" class="w-full flex items-center justify-between text-left
                 text-gray-700 dark:text-gray-200 hover:bg-gray-50 dark:hover:bg-gray-700/50
                 rounded-full px-3 py-2 transition-colors">
          <div class="flex items-center gap-3">
            <span>{{ t('settings.managePermissions') }}</span>
          </div>
          <Icon class="text-gray-400">
            <ChevronRight20Regular />
          </Icon>
        </button>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { NSwitch, NSelect } from 'naive-ui'
import { Icon } from '@vicons/utils'
import ChevronRight20Regular from '@vicons/fluent/es/ChevronRight20Regular'
import { useAppStore } from '@/stores/appStore'

const { t, locale } = useI18n()
const router = useRouter()
const appStore = useAppStore()

const languageOptions = [
  { value: 'zh-CN', label: '中文' },
  { value: 'en', label: 'English' },
]

function switchLanguage(value: string) {
  const lang = value as 'zh-CN' | 'en'
  appStore.setLanguage(lang)
  locale.value = lang
  appStore.saveConfig()
}

function goToPrivacy() {
  router.push({ name: 'settings-privacy' })
}
</script>
