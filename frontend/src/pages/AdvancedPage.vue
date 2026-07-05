<template>
  <div class="h-full overflow-y-auto p-6">
    <n-message-provider>
      <div class="mb-6">
        <h1 class="text-2xl font-bold text-neutral-800 dark:text-neutral-100">
          {{ t('advancedPage.header') }}
        </h1>
      </div>

      <div class="max-w-2xl space-y-5 pb-6">
        <!-- Project source -->
        <n-card size="small">
          <template #header>
            <div class="flex items-center justify-between">
              <span>{{ t('buildPage.sourceCard') }}</span>
              <n-tag v-if="!sourceStatus" type="default" size="small" round>
                {{ t('buildPage.statusNotSelected') }}
              </n-tag>
              <n-tag v-else type="success" size="small" round>
                <template #icon><n-icon :component="CheckmarkCircle24Regular" /></template>
                {{ t('buildPage.status' + sourceStatus) }}
              </n-tag>
            </div>
          </template>
          <n-tabs type="line" animated default-value="folder" @update:value="onTabChange">
            <n-tab-pane name="folder" :tab="t('buildPage.tabFolder')">
              <div class="flex items-center gap-2">
                <n-input v-model:value="projectPath" :placeholder="t('buildPage.placeholderFolder')" clearable @keydown.enter="onPathEnter" />
                <n-button @click="selectDir" type="primary" ghost>{{ t('buildPage.btnBrowse') }}</n-button>
              </div>
            </n-tab-pane>
            <n-tab-pane name="file" :tab="t('buildPage.tabFile')">
              <div class="flex items-center gap-2">
                <n-input v-model:value="filePath" :placeholder="t('buildPage.placeholderFile')" clearable />
                <n-button @click="selectFile" type="primary" ghost>{{ t('buildPage.btnSelectFile') }}</n-button>
              </div>
            </n-tab-pane>
          </n-tabs>
        </n-card>

        <!-- App name -->
        <n-card :title="t('buildPage.appNameCard')" size="small">
          <n-input v-model:value="appName" :placeholder="t('buildPage.placeholderAppName')" maxlength="30" />
        </n-card>

        <!-- Icon -->
        <n-card :title="t('buildPage.iconCard')" size="small">
          <div class="flex items-center gap-4">
            <img v-if="customIconData" :src="customIconData" class="w-16 h-16 rounded-xl object-cover shadow-sm border border-neutral-200 dark:border-neutral-700 shrink-0" />
            <div v-else class="w-16 h-16 rounded-xl bg-linear-to-br from-[#fff4dd] to-[#e16668] flex items-center justify-center text-white text-3xl font-bold shadow-sm shrink-0">{{ iconLetter }}</div>
            <div class="flex-1">
              <div class="text-neutral-500 dark:text-neutral-400 mb-2">{{ t('buildPage.iconHint') }}</div>
              <input ref="fileInputRef" type="file" accept="image/*" class="hidden" @change="onIconFileSelected" />
              <div class="flex items-center gap-2">
                <n-button size="small" @click="openIconPicker">{{ t('buildPage.btnSelectImage') }}</n-button>
                <n-button v-if="customIconData" size="small" tertiary @click="clearIcon">{{ t('buildPage.btnClear') }}</n-button>
              </div>
            </div>
          </div>
        </n-card>

        <!-- Package name & Version (always visible) -->
        <n-card size="small">
          <div>
            <div class="flex items-center justify-between mb-1.5">
              <span class="text-sm text-neutral-500 dark:text-neutral-400">{{ t('buildPage.pkgPreview') }}</span>
              <span class="text-sm text-neutral-400 font-mono">{{ previewPkgName }}</span>
            </div>
            <n-input-group>
              <n-input v-model:value="pkgSegment1" placeholder="com" :allow-input="onlyAllowPkg" :style="{ width: '33%' }" maxlength="20" />
              <n-input v-model:value="pkgSegment2" :placeholder="lastSeg2 || 'company'" :allow-input="onlyAllowPkg" :style="{ width: '33%' }" maxlength="20" />
              <n-input v-model:value="pkgSegment3" placeholder="auto" :allow-input="onlyAllowPkg" :style="{ width: '33%' }" maxlength="20" />
            </n-input-group>
          </div>
          <div class="border-t border-neutral-200 dark:border-neutral-700 pt-3">
            <label class="text-sm text-neutral-500 dark:text-neutral-400 block mb-1.5">{{ t('buildPage.versionLabel') }}</label>
            <n-input v-model:value="versionName" :placeholder="t('buildPage.versionPlaceholder')" :allow-input="onlyAllowVersion" maxlength="20" />
          </div>
        </n-card>

        <!-- Certificate -->
        <n-card size="small">
          <template #header>
            <span class="text-sm font-medium">{{ t('advancedPage.certificate') }}</span>
          </template>
          <n-radio-group :value="certMode" @update:value="onCertModeChange">
            <div class="space-y-3">
              <n-radio value="auto">
                <span class="text-sm text-neutral-700 dark:text-neutral-200">{{ t('advancedPage.certAuto') }}</span>
              </n-radio>
              <div>
                <n-radio value="custom">
                  <span class="text-sm text-neutral-700 dark:text-neutral-200">{{ t('advancedPage.certCustom') }}</span>
                </n-radio>
                <n-collapse-transition :show="certMode === 'custom'">
                  <div class="mt-2 ml-7 space-y-2">
                    <n-input v-model:value="certPath" :placeholder="t('advancedPage.certPathPlaceholder')" readonly>
                      <template #suffix>
                        <n-button text @click="pickCertFile">{{ t('advancedPage.btnBrowse') }}</n-button>
                      </template>
                    </n-input>
                    <n-input v-model:value="certPassword" :placeholder="t('advancedPage.certPasswordPlaceholder')" type="password" show-password-on="click" />
                    <n-input
                      v-model:value="certAlias"
                      :placeholder="t('advancedPage.aliasPlaceholder')"
                    />
                    <n-input
                      v-model:value="keyPassword"
                      :placeholder="t('advancedPage.keyPasswordPlaceholder')"
                      type="password"
                      show-password-on="click"
                    />
                    <n-divider />
                    <div class="flex items-center justify-end gap-2">
                      <n-icon :component="KeyMultiple20Filled" size="18" class="text-neutral-400" />
                      <span class="text-neutral-500 shrink-0">{{ t('advancedPage.remember') }}</span>
                      <n-select
                        v-model:value="rememberLevel"
                        :options="rememberOptions"
                        class="w-36"
                      />
                      <n-tooltip trigger="hover" placement="top" style="max-width: 300px">
                        <template #trigger>
                          <n-icon :component="Info16Regular" size="16" class="text-neutral-400 cursor-help" />
                        </template>
                        <span v-html="t('advancedPage.rememberTooltip')"></span>
                      </n-tooltip>
                    </div>
                  </div>
                </n-collapse-transition>
              </div>
            </div>
          </n-radio-group>
        </n-card>

        <!-- Build button -->
        <BuildButton :disabled="!canBuild || building" :building="building" :idle-text="t('buildPage.btnBuild')" :busy-text="t('buildPage.btnBuilding')" @click="buildAPK" class="max-w-sm" />
      </div>
    </n-message-provider>

    <!-- Float button -->
    <n-tooltip v-if="appStore.showFloatButton" trigger="hover" placement="left">
      <template #trigger>
        <n-float-button :right="24" :bottom="24">
          <n-icon>
            <BookPulse24Regular />
          </n-icon>
        </n-float-button>
      </template>
      <span class="text-xs whitespace-pre-wrap" v-text="buildLog || t('buildPage.floatPlaceholder')" />
    </n-tooltip>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { BuildAPK, GetIconPaths, SelectDirectory, SelectFile, SelectCertFile, PrepareFileInput, SaveCertInfo, LoadCertInfo, ListKeystoreAliases, SetCustomCert, OpenFolder } from '../../wailsjs/go/main/App'
import { useAppStore } from '@/stores/appStore'
import CheckmarkCircle24Regular from '@vicons/fluent/es/CheckmarkCircle24Regular'
import KeyMultiple20Filled from '@vicons/fluent/es/KeyMultiple20Filled'
import Info16Regular from '@vicons/fluent/es/Info16Regular'
import BookPulse24Regular from '@vicons/fluent/es/BookPulse24Regular'
import { NInput, NButton, NTag, NIcon, NTabPane, NTabs, NCard, NInputGroup, NRadio, NRadioGroup, NSelect, NCollapseTransition, NDivider, NTooltip, NFloatButton, useMessage, NMessageProvider } from 'naive-ui'
import BuildButton from '@/components/BuildButton.vue'

const { t } = useI18n()
const message = useMessage()
const appStore = useAppStore()

// ── Source ──
const inputTab = ref('folder')
const projectPath = ref('')
const filePath = ref('')
const appName = ref('')
const building = ref(false)
const buildLog = ref('')

const sourceStatus = computed(() => {
  if (projectPath.value) return 'Folder'
  if (filePath.value) {
    const lower = filePath.value.toLowerCase()
    if (lower.endsWith('.zip')) return 'Zip'
    if (lower.endsWith('.html') || lower.endsWith('.htm')) return 'Html'
  }
  return null
})

function onTabChange(name: string) { inputTab.value = name }
async function selectDir() {
  const dir = await SelectDirectory()
  if (dir) { projectPath.value = dir.replace(/^"(.*)"$/, '$1').trim(); filePath.value = '' }
}
function onPathEnter() { projectPath.value = projectPath.value.replace(/^"(.*)"$/, '$1').trim(); filePath.value = '' }
async function selectFile() {
  const f = await SelectFile()
  if (f) { filePath.value = f; projectPath.value = '' }
}

// ── Icon ──
const fileInputRef = ref<HTMLInputElement>()
const customIconData = ref<string | null>(null)
const iconLetter = computed(() => customIconData.value ? '' : (appName.value ? appName.value.charAt(0).toUpperCase() : '⌘'))
function openIconPicker() { fileInputRef.value?.click() }
function onIconFileSelected(e: Event) {
  const input = e.target as HTMLInputElement
  if (!input.files?.length) return
  const reader = new FileReader()
  reader.onload = () => { customIconData.value = reader.result as string }
  reader.readAsDataURL(input.files[0])
  input.value = ''
}
function clearIcon() { customIconData.value = null }

// ── Package ──
const pkgSegment1 = ref('com')
const pkgSegment2 = ref('')
const pkgSegment3 = ref('')
const lastSeg2 = ref('')
const versionName = ref('')
const previewPkgName = computed(() => `${pkgSegment1.value || 'com'}.${pkgSegment2.value || 'app'}.${pkgSegment3.value || 'app'}`)
function onlyAllowPkg(value: string) { return !value || /^[a-z][a-z0-9_.-]*$/.test(value) }
function onlyAllowVersion(value: string) {
  if (!value) return true
  if (!/^[\d.]*$/.test(value)) return false
  if (value.includes('..')) return false
  return (value.match(/\./g) || []).length <= 2
}

// ── Certificate ──
const certMode = ref<'auto' | 'custom'>('auto')
const certPath = ref('')
const certPassword = ref('')
const certAlias = ref('')
const keyPassword = ref('')
const rememberLevel = ref('off')
const rememberOptions = ref([
  { value: 'off', label: '' },
  { value: 'path', label: '' },
  { value: 'full', label: '' },
])
const aliasOptions = ref<{ label: string; value: string }[]>([])
const loadingAliases = ref(false)

// Set i18n labels for rememberOptions after t() is available
const setRememberLabels = () => {
  rememberOptions.value = [
    { value: 'off', label: t('advancedPage.rememberOff') },
    { value: 'path', label: t('advancedPage.rememberPath') },
    { value: 'full', label: t('advancedPage.rememberFull') },
  ]
}
setRememberLabels()

// On mount, restore remembered cert settings
onMounted(async () => {
  if (appStore.useCustomCert) {
    certMode.value = 'custom'
  }
  rememberLevel.value = appStore.rememberLevel || 'off'

  // Restore company name if enabled
  if (appStore.rememberCompany && appStore.companyName) {
    pkgSegment2.value = appStore.companyName
  }

  if (rememberLevel.value !== 'off') {
    const info = await LoadCertInfo()
    if (info && info.certPath) {
      certPath.value = info.certPath
      if (rememberLevel.value === 'full') {
        certPassword.value = info.certPassword
        certAlias.value = info.certAlias
        keyPassword.value = info.keyPassword
      }
      tryDetectAliases(info.certPath, certPassword.value)
    }
  }
})

function onCertModeChange(val: string) {
  certMode.value = val as 'auto' | 'custom'
  appStore.useCustomCert = val === 'custom'
  appStore.saveConfig()
}

async function persistCertInfo() {
  if (rememberLevel.value === 'off' || !certPath.value) return
  const savePwd = rememberLevel.value === 'full'
  try {
    await SaveCertInfo(
      certPath.value,
      savePwd ? certPassword.value : '',
      savePwd ? certAlias.value : '',
      savePwd ? keyPassword.value : '',
    )
  } catch (e) {
    console.warn('Failed to save cert info:', e)
  }
}

// Auto-save cert info when any field changes
watch([certPath, certPassword, certAlias, keyPassword, rememberLevel], () => {
  appStore.rememberLevel = rememberLevel.value
  appStore.saveConfig()
  persistCertInfo()
})

// Auto-save company name when enabled
watch(pkgSegment2, (val) => {
  if (appStore.rememberCompany) {
    appStore.companyName = val
    appStore.saveConfig()
  }
})

async function tryDetectAliases(path: string, pass: string) {
  if (!path || !pass) return
  loadingAliases.value = true
  try {
    const aliases = await ListKeystoreAliases(path, pass)
    aliasOptions.value = aliases.map(a => ({ label: a, value: a }))
    if (aliases.length === 1 && !certAlias.value) {
      certAlias.value = aliases[0]
    }
  } catch {
    // keytool not available or failed — alias can be typed manually
  } finally {
    loadingAliases.value = false
  }
}

async function pickCertFile() {
  const file = await SelectCertFile()
  if (file) {
    certPath.value = file
    // Auto-detect aliases if password is already filled
    if (certPassword.value) {
      tryDetectAliases(file, certPassword.value)
    }
  }
}

// ── Build ──
const canBuild = computed(() => {
  if (building.value) return false
  const hasSource = inputTab.value === 'folder' ? !!projectPath.value : !!filePath.value
  if (!hasSource || !appName.value) return false
  // Alias is required when using custom cert with a selected file
  if (certMode.value === 'custom' && certPath.value && !certAlias.value) return false
  return true
})

async function buildAPK() {
  if (!canBuild.value || building.value) return
  building.value = true
  try {
    let projectDir = ''
    if (inputTab.value === 'folder') projectDir = projectPath.value
    else projectDir = await PrepareFileInput(filePath.value)

    const seg1 = pkgSegment1.value || 'com'
    const seg2 = pkgSegment2.value || 'app'
    let customPkg = `${seg1}.${seg2}`
    if (pkgSegment3.value) {
      customPkg += `.${pkgSegment3.value}`
    } else {
      // Empty third segment → let backend auto-generate
      customPkg = ''
    }

    const iconPaths = await GetIconPaths()
    const icons: Record<string, string> = {}
    if (customIconData.value) {
      for (const iconPath of iconPaths) {
        const size = extractIconSize(iconPath)
        if (size > 0) {
          const blob = await resizeImage(customIconData.value, size, iconPath.includes('foreground'))
          icons[iconPath] = await blobToBase64(blob)
        }
      }
    } else {
      const letter = appName.value.charAt(0).toUpperCase()
      const colors = ['#4F46E5', '#7C3AED', '#EC4899', '#F59E0B', '#10B981', '#3B82F6', '#EF4444']
      const colorIndex = appName.value.charCodeAt(0) % colors.length
      for (const iconPath of iconPaths) {
        const size = extractIconSize(iconPath)
        if (size > 0) {
          const blob = await generateTextIcon(letter, colors[colorIndex], size, iconPath.includes('foreground'))
          icons[iconPath] = await blobToBase64(blob)
        }
      }
    }

    // Set custom cert if configured
    if (certMode.value === 'custom' && certPath.value) {
      await SetCustomCert(certPath.value, certPassword.value, certAlias.value, keyPassword.value)
    } else {
      await SetCustomCert('', '', '', '')
    }

    const res = await BuildAPK(projectDir, appName.value, customPkg, versionName.value.replace(/\.+$/, ''), icons)
    buildLog.value += res.Log || ''
    message.success(t('buildPage.successTitle') + '\n' + res.APKPath, { duration: 4000, keepAliveOnHover: true })
    if (appStore.openAfterBuild) {
      const dir = res.APKPath.substring(0, res.APKPath.lastIndexOf('\\'))
      if (dir) OpenFolder(dir)
    }
  } catch (err: any) {
    const msg = String(err?.message || err || '')
    let errHint: string
    if (msg.includes('index.html')) {
      errHint = t('buildPage.errorNoIndex')
    } else if (msg.includes('key alias')) {
      errHint = t('advancedPage.errorAliasRequired')
    } else {
      errHint = t('buildPage.errorGeneric')
    }
    message.error(t('buildPage.failTitle') + '\n' + errHint, { duration: 5000, keepAliveOnHover: true })
  } finally {
    building.value = false
  }
}

// ── Icon generation helpers ──

function extractIconSize(path: string): number {
  const sizes: Record<string, number> = { mdpi: 48, hdpi: 72, xhdpi: 96, xxhdpi: 144, xxxhdpi: 192 }
  for (const [key, sz] of Object.entries(sizes)) {
    if (path.includes(key)) {
      return path.includes('foreground') ? Math.round(sz * (108 / 48)) : sz
    }
  }
  return 0
}

function generateTextIcon(letter: string, bgColor: string, size: number, isForeground: boolean): Promise<Blob> {
  return new Promise((resolve, reject) => {
    const canvas = document.createElement('canvas')
    canvas.width = size
    canvas.height = size
    const ctx = canvas.getContext('2d')!
    const half = size / 2
    const radius = size * 0.22

    ctx.clearRect(0, 0, size, size)

    if (!isForeground) {
      ctx.fillStyle = bgColor
      ctx.beginPath()
      ctx.moveTo(radius, 0)
      ctx.lineTo(size - radius, 0)
      ctx.quadraticCurveTo(size, 0, size, radius)
      ctx.lineTo(size, size - radius)
      ctx.quadraticCurveTo(size, size, size - radius, size)
      ctx.lineTo(radius, size)
      ctx.quadraticCurveTo(0, size, 0, size - radius)
      ctx.lineTo(0, radius)
      ctx.quadraticCurveTo(0, 0, radius, 0)
      ctx.closePath()
      ctx.fill()
    } else {
      ctx.fillStyle = bgColor
      ctx.beginPath()
      ctx.moveTo(radius, 0)
      ctx.lineTo(size - radius, 0)
      ctx.quadraticCurveTo(size, 0, size, radius)
      ctx.lineTo(size, size - radius)
      ctx.quadraticCurveTo(size, size, size - radius, size)
      ctx.lineTo(radius, size)
      ctx.quadraticCurveTo(0, size, 0, size - radius)
      ctx.lineTo(0, radius)
      ctx.quadraticCurveTo(0, 0, radius, 0)
      ctx.closePath()
      ctx.fill()
    }

    ctx.fillStyle = '#FFFFFF'
    const fontSize = Math.round(size * 0.5)
    ctx.font = `bold ${fontSize}px system-ui, -apple-system, sans-serif`
    ctx.textAlign = 'center'
    ctx.textBaseline = 'middle'
    ctx.fillText(letter, half, half + 1)

    canvas.toBlob((blob) => {
      if (blob) resolve(blob)
      else reject(new Error('WebP encoding failed'))
    }, 'image/webp', 0.9)
  })
}

function resizeImage(dataUrl: string, size: number, _isForeground: boolean): Promise<Blob> {
  return new Promise((resolve, reject) => {
    const img = new Image()
    img.onload = () => {
      const canvas = document.createElement('canvas')
      canvas.width = size
      canvas.height = size
      const ctx = canvas.getContext('2d')!
      ctx.clearRect(0, 0, size, size)

      const radius = size * 0.22
      ctx.beginPath()
      ctx.moveTo(radius, 0)
      ctx.lineTo(size - radius, 0)
      ctx.quadraticCurveTo(size, 0, size, radius)
      ctx.lineTo(size, size - radius)
      ctx.quadraticCurveTo(size, size, size - radius, size)
      ctx.lineTo(radius, size)
      ctx.quadraticCurveTo(0, size, 0, size - radius)
      ctx.lineTo(0, radius)
      ctx.quadraticCurveTo(0, 0, radius, 0)
      ctx.closePath()
      ctx.clip()

      const min = Math.min(img.width, img.height)
      const sx = (img.width - min) / 2
      const sy = (img.height - min) / 2
      ctx.drawImage(img, sx, sy, min, min, 0, 0, size, size)
      canvas.toBlob((blob) => {
        if (blob) resolve(blob)
        else reject(new Error('WebP encoding failed'))
      }, 'image/webp', 0.9)
    }
    img.onerror = () => reject(new Error('Failed to load image'))
    img.src = dataUrl
  })
}

function blobToBase64(blob: Blob): Promise<string> {
  return new Promise((res, rej) => {
    const r = new FileReader()
    r.onload = () => { const s = r.result as string; res(s.substring(s.indexOf(',') + 1)) }
    r.onerror = () => rej(new Error('Failed'))
    r.readAsDataURL(blob)
  })
}
</script>
