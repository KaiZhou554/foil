<template>
  <div class="h-full overflow-y-auto p-6">
      <!-- Header -->
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-neutral-800 dark:text-neutral-100">
        {{ t('buildPage.header') }}
      </h1>
    </div>

    <div class="max-w-2xl space-y-5 pb-6">
      <!-- HTML 项目来源 (Tabs) -->
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
            <div class="space-y-3">
              <div class="flex items-center gap-2">
                <n-input
                  v-model:value="projectPath"
                  :placeholder="t('buildPage.placeholderFolder')"
                  clearable
                  @keydown.enter="onPathEnter"
                />
                <n-button @click="selectDir" type="primary" ghost>
                  {{ t('buildPage.btnBrowse') }}
                </n-button>
              </div>
            </div>
          </n-tab-pane>

          <n-tab-pane name="file" :tab="t('buildPage.tabFile')">
            <div class="space-y-3">
              <div class="flex items-center gap-2">
                <n-input
                  v-model:value="filePath"
                  :placeholder="t('buildPage.placeholderFile')"
                  clearable
                />
                <n-button @click="selectFile" type="primary" ghost>
                  {{ t('buildPage.btnSelectFile') }}
                </n-button>
              </div>
            </div>
          </n-tab-pane>
        </n-tabs>
      </n-card>

      <!-- 应用名称 -->
      <n-card :title="t('buildPage.appNameCard')" size="small">
        <n-input
          v-model:value="appName"
          :placeholder="t('buildPage.placeholderAppName')"
          maxlength="30"
          :input-props="{ onInput: (e: Event) => { const t = (e.target as HTMLInputElement)?.value; if (t) previewPkg(); } }"
        />
      </n-card>

      <!-- 图标上传 -->
      <n-card :title="t('buildPage.iconCard')" size="small">
        <div class="flex items-center gap-4">
          <img
            v-if="customIconData"
            :src="customIconData"
            class="w-16 h-16 rounded-xl object-cover shadow-sm border border-neutral-200 dark:border-neutral-700 shrink-0"
          />
          <div
            v-else
            class="w-16 h-16 rounded-xl bg-taffy-500/50 flex items-center justify-center text-white text-3xl font-bold shadow-sm shrink-0"
          >
            {{ iconLetter }}
          </div>
          <div class="flex-1">
            <div class="text-neutral-500 dark:text-neutral-400 mb-2">
              {{ t('buildPage.iconHint') }}
            </div>
            <input
              ref="fileInputRef"
              type="file"
              accept="image/*"
              class="hidden"
              @change="onIconFileSelected"
            />
            <div class="flex items-center gap-2">
              <n-button size="small" @click="openIconPicker">{{ t('buildPage.btnSelectImage') }}</n-button>
              <n-button v-if="customIconData" size="small" tertiary @click="clearIcon">
                {{ t('buildPage.btnClear') }}
              </n-button>
            </div>
          </div>
        </div>
      </n-card>

      <!-- 高级选项 -->
      <n-collapse>
        <n-collapse-item :title="t('buildPage.advanced')">
          <div class="space-y-4">
            <!-- Package name -->
            <div>
              <div class="flex items-center justify-between mb-1.5">
                <span class="text-sm text-neutral-500 dark:text-neutral-400">{{ t('buildPage.pkgPreview') }}</span>
                <span class="text-sm text-neutral-400 font-mono">{{ previewPkgName }}</span>
              </div>
              <n-input-group>
                <n-input
                  v-model:value="pkgSegment1"
                  placeholder="com"
                  :allow-input="onlyAllowPkg"
                  :style="{ width: '33%' }"
                  maxlength="20"
                />
                <n-input
                  v-model:value="pkgSegment2"
                  :placeholder="lastSeg2 || 'company'"
                  :allow-input="onlyAllowPkg"
                  :style="{ width: '33%' }"
                  maxlength="20"
                />
                <n-input
                  v-model:value="pkgSegment3"
                  placeholder="auto"
                  :allow-input="onlyAllowPkg"
                  :style="{ width: '33%' }"
                  maxlength="20"
                />
              </n-input-group>
            </div>

            <!-- Version -->
            <div class="border-t border-neutral-200 dark:border-neutral-700 pt-3">
              <label class="text-sm text-neutral-500 dark:text-neutral-400 block mb-1.5">{{ t('buildPage.versionLabel') }}</label>
              <n-input
                v-model:value="versionName"
                :placeholder="t('buildPage.versionPlaceholder')"
                :allow-input="onlyAllowVersion"
                maxlength="20"
              />
            </div>
          </div>
        </n-collapse-item>
      </n-collapse>

      <!-- 生成按钮 -->
      <BuildButton
        :disabled="!canBuild || building"
        :building="building"
        :idle-text="t('buildPage.btnBuild')"
        :busy-text="t('buildPage.btnBuilding')"
        @click="buildAPK"
        class="max-w-sm"
      />

    </div>

    <!-- Float button -->
    <n-popover v-if="appStore.showFloatButton" trigger="hover" placement="left">
      <template #trigger>
        <n-float-button :right="24" :bottom="24">
          <n-icon>
            <BookPulse24Regular />
          </n-icon>
        </n-float-button>
      </template>
      <div style="max-width: 360px">
        <pre style="white-space: pre-wrap; word-break: break-all; margin: 0; font-size: 12px">{{ buildLog || t('buildPage.floatPlaceholder') }}</pre>
        <div style="text-align: right; margin-top: 8px">
          <n-button size="tiny" @click="copyBuildLog">
            {{ t('buildPage.copyLog') }}
          </n-button>
        </div>
      </div>
    </n-popover>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { BuildAPK, GetIconPaths, SelectDirectory, SelectFile, PrepareFileInput, OpenFolder, GeneratePackageName } from '../../wailsjs/go/main/App'
import { useAppStore } from '@/stores/appStore'
import CheckmarkCircle24Regular from '@vicons/fluent/es/CheckmarkCircle24Regular'
import BookPulse24Regular from '@vicons/fluent/es/BookPulse24Regular'
import BuildButton from '@/components/BuildButton.vue'
import {
  NInput, NButton, NTag, NIcon, NTabPane, NTabs, NCard,
  NInputGroup, NCollapseItem, NCollapse,
  NFloatButton, NPopover,
  useMessage,
} from 'naive-ui'

const { t } = useI18n()
const appStore = useAppStore()
const message = useMessage()

// ── State ──
const inputTab = ref('folder')
const projectPath = ref('')
const filePath = ref('')
const appName = ref('')
const building = ref(false)
const buildLog = ref('')

// Package name segments
const pkgSegment1 = ref('com')
const pkgSegment2 = ref('')
const pkgSegment3 = ref('')
const lastSeg2 = ref('')

// Custom version number
const versionName = ref('')

// Icon
const fileInputRef = ref<HTMLInputElement>()
const customIconData = ref<string | null>(null)
const customIconFileName = ref('')

// ── Computed ──
const sourceStatus = computed(() => {
  if (projectPath.value) return 'Folder'
  if (filePath.value) {
    const lower = filePath.value.toLowerCase()
    if (lower.endsWith('.zip')) return 'Zip'
    if (lower.endsWith('.html') || lower.endsWith('.htm')) return 'Html'
  }
  return null
})

const iconLetter = computed(() => {
  if (customIconData.value) return ''
  return appName.value ? appName.value.charAt(0).toUpperCase() : '∴'
})

const canBuild = computed(() => {
  if (building.value) return false
  if (inputTab.value === 'folder') return !!projectPath.value && !!appName.value
  return !!filePath.value && !!appName.value
})

const previewPkgName = computed(() => {
  const seg1 = pkgSegment1.value || 'com'
  const seg2 = pkgSegment2.value || 'app'
  const seg3 = pkgSegment3.value || 'app'
  return `${seg1}.${seg2}.${seg3}`
})

// ── Tab switch ──
function onTabChange(name: string) {
  inputTab.value = name
}

// ── File / Folder dialogs ──
async function selectDir() {
  const dir = await SelectDirectory()
  if (dir) {
    projectPath.value = dir.replace(/^"(.*)"$/, '$1').trim()
    filePath.value = '' // selecting folder clears file
  }
}

function onPathEnter() {
  projectPath.value = projectPath.value.replace(/^"(.*)"$/, '$1').trim()
  filePath.value = '' // manual folder entry clears file
}

async function selectFile() {
  const file = await SelectFile()
  if (file) {
    filePath.value = file
    projectPath.value = '' // selecting file clears folder
  }
}

// ── Icon ──
function openIconPicker() {
  fileInputRef.value?.click()
}

function onIconFileSelected(e: Event) {
  const input = e.target as HTMLInputElement
  if (!input.files?.length) return
  const file = input.files[0]
  customIconFileName.value = file.name

  const reader = new FileReader()
  reader.onload = () => {
    customIconData.value = reader.result as string
  }
  reader.readAsDataURL(file)
  input.value = ''
}

function clearIcon() {
  customIconData.value = null
  customIconFileName.value = ''
}

// ── Package name validation ──
function onlyAllowPkg(value: string) {
  return !value || /^[a-z][a-z0-9_.]*$/.test(value)
}

function onlyAllowVersion(value: string) {
  if (!value) return true
  // Only digits and dots
  if (!/^[\d.]*$/.test(value)) return false
  // No consecutive dots
  if (value.includes('..')) return false
  // At most 2 dots
  return (value.match(/\./g) || []).length <= 2
}

function previewPkg() {
  // preview is updated reactively via computed property
}

// ── Build ──

// Restore company name from store on mount
const __restoreCompany = () => {
  if (appStore.rememberCompany && appStore.companyName) {
    pkgSegment2.value = appStore.companyName
  }
}
__restoreCompany()

// Auto-save company name when enabled
watch(pkgSegment2, (val) => {
  if (appStore.rememberCompany) {
    appStore.companyName = val
    appStore.saveConfig()
  }
})

async function buildAPK() {
  if (!canBuild.value || building.value) return

  building.value = true
  buildLog.value = ''

  try {
    let projectDir = ''
    if (inputTab.value === 'folder') {
      projectDir = projectPath.value
    } else {
      projectDir = await PrepareFileInput(filePath.value)
    }

    const seg1 = pkgSegment1.value || 'com'

    let customPkg: string
    if (pkgSegment2.value && pkgSegment3.value) {
      customPkg = `${seg1}.${pkgSegment2.value}.${pkgSegment3.value}`
    } else if (pkgSegment2.value && !pkgSegment3.value) {
      const randomPkg = await GeneratePackageName(appName.value)
      const seg = randomPkg.split('.').pop()!
      customPkg = `${seg1}.${pkgSegment2.value}.${seg}`
    } else if (!pkgSegment2.value && pkgSegment3.value) {
      const randomPkg = await GeneratePackageName(appName.value)
      const seg = randomPkg.split('.').pop()!
      customPkg = `${seg1}.${pkgSegment3.value}${seg}`
    } else {
      const randomPkg = await GeneratePackageName(appName.value)
      const seg = randomPkg.split('.').pop()!
      customPkg = `${seg1}.${seg}`
    }
    if (customPkg) {
      buildLog.value += `${t('buildPage.logCustomPkg')}${customPkg}\n`
    }

    // Generate icons as WebP
    const icons: Record<string, string> = {}
    const iconPaths = await GetIconPaths()

    if (customIconData.value) {
      for (const iconPath of iconPaths) {
        const size = extractIconSize(iconPath)
        if (size > 0) {
          const blob = await resizeImage(customIconData.value, size, iconPath.includes('foreground'))
          const b64 = await blobToBase64(blob)
          icons[iconPath] = b64
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
          const b64 = await blobToBase64(blob)
          icons[iconPath] = b64
        }
      }
    }

    buildLog.value += t('buildPage.logBuilding')
    const res = await BuildAPK(projectDir, appName.value, customPkg, versionName.value.replace(/\.+$/, ''), icons)
    buildLog.value += res.Log || ''
    message.success(t('buildPage.successTitle') + '\n' + res.APKPath, { duration: 4000, keepAliveOnHover: true })
    if (appStore.openAfterBuild) {
      const dir = res.APKPath.substring(0, res.APKPath.lastIndexOf('\\'))
      if (dir) OpenFolder(dir)
    }
  } catch (err: any) {
    const msg = String(err?.message || err || '')
    buildLog.value += msg + '\n'
    // Map common errors to localized messages
    const errHint = msg.includes('index.html')
      ? t('buildPage.errorNoIndex')
      : t('buildPage.errorGeneric')
    message.error(t('buildPage.failTitle') + '\n' + errHint, { duration: 5000, keepAliveOnHover: true })
  } finally {
    building.value = false
  }
}

// ── Icon generation helpers ──

function extractIconSize(path: string): number {
  const sizes: Record<string, number> = { mdpi: 48, hdpi: 72, xhdpi: 96, xxhdpi: 144, xxxhdpi: 192 }
  for (const [key, size] of Object.entries(sizes)) {
    if (path.includes(key)) {
      return path.includes('foreground') ? Math.round(size * (108 / 48)) : size
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
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = () => {
      const result = reader.result as string
      const comma = result.indexOf(',')
      resolve(comma >= 0 ? result.substring(comma + 1) : result)
    }
    reader.onerror = () => reject(new Error('Failed to read blob'))
    reader.readAsDataURL(blob)
  })
}

async function copyBuildLog() {
  try {
    await navigator.clipboard.writeText(buildLog.value)
  } catch {
    // fallback: silently ignore
  }
}
</script>
