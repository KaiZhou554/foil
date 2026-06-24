<template>
  <n-config-provider :theme-overrides="themeOverrides" :theme="theme">
    <div class="h-screen flex flex-col bg-white/20 dark:bg-neutral-900/20">
      <!-- TitleBar — fixed above all page transitions -->
      <TitleBar
        :show-sidebar-toggle="showSidebarToggle"
        :back-button="showBackButton"
        :icon-animating="appStore.sidebarAnimating"
        :is-maximized="appStore.isMaximized"
        :app-icon="AppIcon"
        :app-name="t('app.defaultAppName')"
        @toggle-sidebar="appStore.toggleSidebar()"
        @back="handleBack"
        @minimize="handleMinimize"
        @maximize="handleMaximize"
        @close="handleClose"
      />

      <!-- Page content area — hidden until config is loaded to avoid flash -->
      <div v-if="appStore.configLoaded" class="flex-1 overflow-hidden">
        <router-view v-slot="{ Component, route }">
          <transition :name="transitionName" mode="out-in">
            <component :is="Component" :key="route.matched[0]?.path ?? route.path" />
          </transition>
        </router-view>
      </div>
      <!-- Loading placeholder while config loads -->
      <div v-else class="flex-1 overflow-hidden flex items-center justify-center">
        <div class="w-6 h-6 border-2 border-gray-300 border-t-pink-400 rounded-full animate-spin" />
      </div>
    </div>
  </n-config-provider>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { NConfigProvider, darkTheme, useOsTheme  } from 'naive-ui'
import { WindowMinimise, Quit } from '../wailsjs/runtime/runtime'
import AppIcon from '@/assets/appicon_128px.webp'
import TitleBar from '@/components/TitleBar.vue'
import { useAppStore } from '@/stores/appStore'
import { usePrivacyConfigStore } from '@/stores/privacyConfigStore'
import { useWindowState } from '@/composables/useWindowState'

const router = useRouter()
const route = useRoute()
const { t, locale } = useI18n()
const appStore = useAppStore()
const osTheme = useOsTheme()
const theme = computed(() => (osTheme.value === 'dark' ? darkTheme : null))

// ── Window state synchronisation ──
const {
  isMaximized,
  startWindowStateSync,
  stopWindowStateSync,
  maximize,
  unmaximize,
} = useWindowState()

// Keep the store in sync with the composable's reactive ref.
// TitleBar reads appStore.isMaximized via :is-maximized prop.
watch(isMaximized, (val) => appStore.setMaximized(val))

/**
 * Global Naive UI theme overrides.
 */
const themeOverrides = {
  Switch: {
    loadingColor: '#ce597eFF',
    railColorActive: '#ce597eFF',
  },
  common: {
    primaryColor: '#ce597eFF',
    primaryColorHover: '#da6a8cFF',
    primaryColorPressed: '#9e4260FF',
    primaryColorSuppl: '#da6a8cFF',
  },
}

// ── TitleBar conditional props ──
const showSidebarToggle = computed(() =>
  route.path.startsWith('/main') && !appStore.isFirstLaunch,
)
const showBackButton = computed(() => route.path === '/setup')

// ── TitleBar event handlers ──
function handleBack() {
  router.back()
}
function handleMinimize() {
  WindowMinimise()
}
async function handleMaximize() {
  if (appStore.isMaximized) {
    await unmaximize()
  } else {
    await maximize()
  }
}
function handleClose() {
  Quit()
}

// ── Top-level page transition ──
const transitionName = ref('fade')

/** First-segment path → first-segment path → CSS transition name */
const topTransitionMap: Record<string, Record<string, string>> = {
  '/welcome': { '/setup': 'slide-left' },
  '/setup':   { '/welcome': 'slide-right', '/main': 'fade' },
}

watch(
  () => router.currentRoute.value.matched[0]?.path,
  (to, from) => {
    const fromPath = (from as string | undefined) ?? null
    const toPath = (to as string | undefined) ?? null
    if (!fromPath || !toPath) {
      transitionName.value = 'fade'
      return
    }
    transitionName.value = topTransitionMap[fromPath]?.[toPath] ?? 'fade'
  },
)

const privacyStore = usePrivacyConfigStore()

// When isFirstLaunch flips from default (true) to loaded (false), redirect to home.
// Using a watch avoids the async-timing issues of router.replace inside onMounted.
watch(
  () => appStore.isFirstLaunch,
  (val) => {
    if (!val) {
      const p = route.path
      if (p === '/welcome' || p === '/setup') {
        router.replace('/main/home')
      }
    }
  },
)

// Auto-save config when privacy settings change (debounced 400ms)
let privacySaveTimer: ReturnType<typeof setTimeout> | null = null
function schedulePrivacySave() {
  if (!appStore.configLoaded) return
  if (privacySaveTimer) clearTimeout(privacySaveTimer)
  privacySaveTimer = setTimeout(() => appStore.saveConfig(), 400)
}
watch(() => privacyStore.optionValues, schedulePrivacySave, { deep: true })
watch(() => privacyStore.optionDisplay, schedulePrivacySave, { deep: true })

onMounted(async () => {
  // Start polling window state (also does an immediate first sync)
  startWindowStateSync()

  // Load persisted config from Go backend
  await appStore.loadConfig()

  // Sync i18n locale with config language
  locale.value = appStore.currentLanguage

  // Navigate to correct initial page (router-view is hidden until now, so no flash)
  if (appStore.isFirstLaunch) {
    if (route.path !== '/welcome') router.replace('/welcome')
  } else {
    router.replace('/main/home')
  }
})

onUnmounted(() => {
  stopWindowStateSync()
})
</script>
