import { defineStore } from 'pinia'
import { ref } from 'vue'
import { GetConfig, SaveConfig } from '../../wailsjs/go/main/App'
import { config } from '../../wailsjs/go/models'
import { usePrivacyConfigStore } from './privacyConfigStore'

/**
 * App-level state store.
 * Configuration is persisted to TOML via Wails backend.
 * All keys use camelCase consistently — frontend store, Wails models, Go struct tags, TOML file.
 * No conversion anywhere.
 */
export const useAppStore = defineStore('app', () => {
  // ============== Config-backed state ==============
  const supportedLocales = ['zh-CN', 'en'] as const
  type Locale = (typeof supportedLocales)[number]

  const isFirstLaunch = ref<boolean>(true)
  const currentLanguage = ref<Locale>('zh-CN')
  const generalSettings = ref({
    autoStart: false,
    notifications: true,
    minimizeToTray: true,
  })

  /** Whether the config has been loaded from the backend yet. */
  const configLoaded = ref<boolean>(false)

  // ============== Config load / save ==============

  /** Load config from the Go backend. Called once on app startup. */
  async function loadConfig() {
    try {
      const cfg = await GetConfig()
      isFirstLaunch.value = cfg.firstLaunch
      currentLanguage.value = cfg.language as Locale
      generalSettings.value = {
        autoStart: cfg.general?.autoStart ?? false,
        notifications: cfg.general?.notifications ?? true,
        minimizeToTray: cfg.general?.minimizeToTray ?? true,
      }

      // Sync privacy config store from loaded config — keys match directly (camelCase)
      if (cfg.privacy) {
        const privacyStore = usePrivacyConfigStore()
        const opts = cfg.privacy.options as unknown as Record<string, boolean> | undefined
        const disp = cfg.privacy.display as unknown as Record<string, boolean> | undefined
        const meta = cfg.privacy.meta
        if (opts) {
          for (const key of Object.keys(opts)) {
            if (key in privacyStore.optionValues) {
              ;(privacyStore.optionValues as any)[key] = opts[key]
            }
          }
        }
        if (disp) {
          for (const key of Object.keys(disp)) {
            if (key in privacyStore.optionDisplay) {
              ;(privacyStore.optionDisplay as any)[key] = disp[key]
            }
          }
        }
        if (meta) {
          privacyStore.meta = {
            displayVersion: meta.displayVersion,
            updateTime: meta.updateTime,
            author: meta.author,
            desc: meta.desc,
          }
        }
      }

      configLoaded.value = true
    } catch (err) {
      console.error('Failed to load config from backend:', err)
      // Keep defaults; user can still use the app
      configLoaded.value = true
    }
  }

  /** Persist current in-memory state to the Go backend. */
  async function saveConfig() {
    try {
      const privacyStore = usePrivacyConfigStore()
      await SaveConfig(new config.Config({
        version: '0.1.0',
        firstLaunch: isFirstLaunch.value,
        language: currentLanguage.value,
        general: new config.GeneralConfig({
          autoStart: generalSettings.value.autoStart,
          notifications: generalSettings.value.notifications,
          minimizeToTray: generalSettings.value.minimizeToTray,
        }),
        privacy: new config.PrivacyConfig({
          options: new config.PrivacyFields(privacyStore.optionValues as any),
          display: new config.PrivacyFields(privacyStore.optionDisplay as any),
          meta: new config.PrivacyMeta(privacyStore.meta as any),
        }),
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

  /** Switch language. Does NOT auto-save — call saveConfig() explicitly. */
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
    generalSettings,
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
