import { defineStore } from 'pinia'
import { ref } from 'vue'
import { GetConfig, SaveConfig } from '../../wailsjs/go/main/App'
import { config } from '../../wailsjs/go/models'

/**
 * App-level state store.
 * Configuration is persisted to TOML via Wails backend.
 */
export const useAppStore = defineStore('app', () => {
  // ============== Config-backed state ==============
  const supportedLocales = ['zh-CN', 'en'] as const
  type Locale = (typeof supportedLocales)[number]

  const isFirstLaunch = ref<boolean>(true)
  const currentLanguage = ref<Locale>('zh-CN')
  const outputDir = ref<string>('')
  const showFloatButton = ref<boolean>(false)
  const openAfterBuild = ref<boolean>(true)
  const useCustomCert = ref<boolean>(false)
  const rememberLevel = ref<string>('off')

  /** Whether the config has been loaded from the backend yet. */
  const configLoaded = ref<boolean>(false)

  // ============== Config load / save ==============

  /** Load config from the Go backend. Called once on app startup. */
  async function loadConfig() {
    try {
      const cfg = await GetConfig()
      isFirstLaunch.value = cfg.firstLaunch
      currentLanguage.value = cfg.language as Locale
      outputDir.value = cfg.outputDir || ''
      showFloatButton.value = cfg.showFloatButton ?? false
      openAfterBuild.value = cfg.openAfterBuild ?? true
      useCustomCert.value = cfg.useCustomCert ?? false
      rememberLevel.value = cfg.rememberLevel || 'off'
      configLoaded.value = true
    } catch (err) {
      console.error('Failed to load config from backend:', err)
      configLoaded.value = true
    }
  }

  /** Persist current in-memory state to the Go backend. */
  async function saveConfig() {
    try {
      await SaveConfig(new config.Config({
        version: '0.1.0',
        firstLaunch: isFirstLaunch.value,
        language: currentLanguage.value,
        outputDir: outputDir.value,
        showFloatButton: showFloatButton.value,
        openAfterBuild: openAfterBuild.value,
        useCustomCert: useCustomCert.value,
        rememberLevel: rememberLevel.value,
      }))
    } catch (err) {
      console.error('Failed to save config to backend:', err)
      throw err
    }
  }

  // ============== Onboarding ==============

  /** Mark onboarding as complete and persist. */
  async function completeOnboarding() {
    isFirstLaunch.value = false
    await saveConfig()
  }

  // ============== Language ==============

  function setLanguage(locale: Locale) {
    currentLanguage.value = locale
  }

  // ============== Window state (for TitleBar) ==============
  const isMaximized = ref<boolean>(false)

  function setMaximized(val: boolean) {
    isMaximized.value = val
  }

  // ============== Sidebar state ==============
  const sidebarOpen = ref<boolean>(false)
  const sidebarAnimating = ref<boolean>(false)

  function toggleSidebar() {
    sidebarAnimating.value = true
    sidebarOpen.value = !sidebarOpen.value
    setTimeout(() => {
      sidebarAnimating.value = false
    }, 300)
  }

  function setSidebarOpen(val: boolean) {
    sidebarOpen.value = val
  }

  return {
    // state
    isFirstLaunch,
    currentLanguage,
    outputDir,
    showFloatButton,
    openAfterBuild,
    useCustomCert,
    rememberLevel,
    configLoaded,
    isMaximized,
    sidebarOpen,
    sidebarAnimating,
    // actions
    loadConfig,
    saveConfig,
    completeOnboarding,
    setLanguage,
    setMaximized,
    toggleSidebar,
    setSidebarOpen,
  }
})
