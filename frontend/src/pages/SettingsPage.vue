<template>
  <div class="h-full overflow-y-auto p-8">
    <h1 class="text-2xl font-bold text-neutral-800 dark:text-neutral-100 mb-6">
      {{ t('settings.title') }}
    </h1>

    <!-- Language -->
    <section class="mb-8">
      <h2 class="text-lg font-medium text-neutral-800 dark:text-neutral-100 mb-3">
        {{ t('settings.language') }}
      </h2>
      <div class="bg-white dark:bg-neutral-800 rounded-xl ring-1 ring-black/5 dark:ring-white/10 p-4">
        <div class="flex items-center justify-between">
          <span class="text-neutral-700 dark:text-neutral-200">{{ t('settings.displayLanguage') }}</span>
          <div class="w-36">
            <n-select :value="appStore.currentLanguage" :options="languageOptions" @update:value="switchLanguage"
              class="w-40" size="small" />
          </div>
        </div>
      </div>
    </section>

    <!-- General -->
    <section class="mb-8">
      <h2 class="text-lg font-medium text-neutral-800 dark:text-neutral-100 mb-3">
        {{ t('settings.general') }}
      </h2>
      <div class="bg-white dark:bg-neutral-800 rounded-xl ring-1 ring-black/5 dark:ring-white/10 p-4">
        <div class="flex items-center justify-between gap-4">
          <span class="text-sm text-neutral-700 dark:text-neutral-200 shrink-0">{{ t('settings.saveLocation') }}</span>
          <n-radio-group :value="localMode" @update:value="onLocationModeChange">
            <n-radio value="desktop" class="mr-3">
              <span class="text-sm text-neutral-700 dark:text-neutral-200">{{ t('settings.locationDesktop') }}</span>
            </n-radio>
            <n-radio value="custom">
              <span class="text-sm text-neutral-700 dark:text-neutral-200">{{ t('settings.locationCustom') }}</span>
            </n-radio>
          </n-radio-group>
        </div>
        <n-collapse-transition :show="localMode === 'custom'">
          <div class="flex items-center gap-2 pt-6 pb-2">
            <n-input
              :value="customPath"
              :placeholder="t('settings.locationPlaceholder')"
              readonly
              class="flex-1"
              size="small"
            />
            <n-button size="small" type="primary" ghost @click="pickCustomPath">
              {{ t('settings.btnBrowse') }}
            </n-button>
          </div>
        </n-collapse-transition>

        <div class="flex items-center justify-between pt-3 border-t border-neutral-100 dark:border-neutral-700/50 mt-3">
          <span class="text-sm text-neutral-700 dark:text-neutral-200">{{ t('settings.showFloatButton') }}</span>
          <n-switch :value="appStore.showFloatButton" @update:value="val => { appStore.showFloatButton = val; appStore.saveConfig() }" />
        </div>
        <div class="flex items-center justify-between pt-3 border-t border-neutral-100 dark:border-neutral-700/50 mt-3">
          <span class="text-sm text-neutral-700 dark:text-neutral-200">{{ t('settings.openAfterBuild') }}</span>
          <n-switch :value="appStore.openAfterBuild" @update:value="val => { appStore.openAfterBuild = val; appStore.saveConfig() }" />
        </div>
      </div>
    </section>

    <!-- About -->
    <section class="mb-8">
      <h2 class="text-lg font-medium text-neutral-800 dark:text-neutral-100 mb-3">
        {{ t('settings.about') }}
      </h2>
      <div class="bg-white dark:bg-neutral-800 rounded-xl ring-1 ring-black/5 dark:ring-white/10 p-4">
        <div class="flex items-center justify-between gap-4">
          <div class="text-sm text-neutral-500 leading-relaxed">
            <strong class="text-neutral-700 dark:text-neutral-200">Foil</strong> · KaiZhou554 · v{{ appVersion }}<br>
            {{ t('settings.aboutDesc') }}
          </div>
          <GitHubButton class="shrink-0" />
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { NSelect, NRadio, NRadioGroup, NInput, NButton, NSwitch, NCollapseTransition } from 'naive-ui'
import { useAppStore } from '@/stores/appStore'
import { SelectOutputDir, GetAppVersion } from '../../wailsjs/go/main/App'
import GitHubButton from '@/components/GitHubButton.vue'

const { t, locale } = useI18n()
const appStore = useAppStore()

const appVersion = ref('')

onMounted(async () => {
  appVersion.value = await GetAppVersion()
})

const languageOptions = [
  { value: 'zh-CN', label: '中文' },
  { value: 'en', label: 'English' },
]

// ── Location state ──
const customPath = ref('')
const localMode = ref<'desktop' | 'custom'>(appStore.outputDir ? 'custom' : 'desktop')

function onLocationModeChange(val: string) {
  localMode.value = val as 'desktop' | 'custom'
  if (val === 'desktop') {
    customPath.value = ''
    appStore.outputDir = ''
    appStore.saveConfig()
  }
}

async function pickCustomPath() {
  const dir = await SelectOutputDir()
  if (dir) {
    customPath.value = dir
    appStore.outputDir = dir
    appStore.saveConfig()
  }
}

function switchLanguage(value: string) {
  const lang = value as 'zh-CN' | 'en'
  appStore.setLanguage(lang)
  locale.value = lang
  appStore.saveConfig()
}
</script>
