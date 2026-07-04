<template>
  <div class="h-full overflow-y-auto p-6">
    <!-- Header -->
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-gray-800 dark:text-gray-100">
        极速打包
      </h1>
    </div>

    <div class="max-w-2xl space-y-5">
      <!-- HTML 项目来源 (Tabs) -->
      <n-card title="HTML 项目来源" size="small">
        <n-tabs type="line" animated default-value="folder" @update:value="onTabChange">
          <!-- Tab 1: 文件夹路径 -->
          <n-tab-pane name="folder" tab="文件夹">
            <div class="space-y-3">
              <div class="flex items-center gap-2">
                <n-input
                  v-model:value="projectPath"
                  placeholder="粘贴文件夹路径，或点击右侧按钮选择..."
                  clearable
                  @keydown.enter="onPathEnter"
                />
                <n-button @click="selectDir" type="primary" ghost>
                  浏览
                </n-button>
              </div>
              <!-- Status tag -->
              <div class="flex items-center gap-2">
                <n-tag v-if="!projectPath" type="default" size="small">未选择</n-tag>
                <n-tag v-else-if="isFolderSelected" type="success" size="small" round>
                  <template #icon><n-icon :component="CheckmarkCircle24Regular" /></template>
                  已选择 — 文件夹
                </n-tag>
              </div>
            </div>
          </n-tab-pane>

          <!-- Tab 2: ZIP / HTML 文件 -->
          <n-tab-pane name="file" tab="ZIP 或 HTML 文件">
            <div class="space-y-3">
              <div class="flex items-center gap-2">
                <n-input
                  v-model:value="filePath"
                  placeholder="选择 .zip 或 .html 文件..."
                  clearable
                />
                <n-button @click="selectFile" type="primary" ghost>
                  选择文件
                </n-button>
              </div>
              <div class="flex items-center gap-2">
                <n-tag v-if="!filePath" type="default" size="small">未选择</n-tag>
                <n-tag v-else-if="isZipSelected" type="success" size="small" round>
                  <template #icon><n-icon :component="CheckmarkCircle24Regular" /></template>
                  已选择 — ZIP 包
                </n-tag>
                <n-tag v-else-if="isHtmlSelected" type="success" size="small" round>
                  <template #icon><n-icon :component="CheckmarkCircle24Regular" /></template>
                  已选择 — HTML 文件
                </n-tag>
              </div>
            </div>
          </n-tab-pane>
        </n-tabs>
      </n-card>

      <!-- 应用名称 -->
      <n-card title="应用名称" size="small">
        <n-input
          v-model:value="appName"
          placeholder="输入显示在手机上的应用名称"
          maxlength="30"
          :input-props="{ onInput: (e: Event) => { const t = (e.target as HTMLInputElement)?.value; if (t) previewPkg(); } }"
        />
      </n-card>

      <!-- 图标上传 -->
      <n-card title="应用图标" size="small">
        <div class="flex items-center gap-4">
          <!-- Preview: uploaded image or auto-generated letter -->
          <img
            v-if="customIconData"
            :src="customIconData"
            class="w-16 h-16 rounded-xl object-cover shadow-sm border border-gray-200 dark:border-gray-700 shrink-0"
          />
          <div
            v-else
            class="w-16 h-16 rounded-xl bg-linear-to-br from-blue-400 to-purple-500 flex items-center justify-center text-white text-3xl font-bold shadow-sm shrink-0"
          >
            {{ iconLetter }}
          </div>
          <div class="flex-1">
            <div class="text-xs text-gray-500 dark:text-gray-400 mb-2">
              上传一张图片自动处理为图标，不传则自动生成
            </div>
            <input
              ref="fileInputRef"
              type="file"
              accept="image/*"
              class="hidden"
              @change="onIconFileSelected"
            />
            <n-button size="small" @click="openIconPicker">选择图片</n-button>
            <n-button v-if="customIconData" size="small" tertiary @click="clearIcon" class="ml-2">
              清除
            </n-button>
          </div>
        </div>
      </n-card>

      <!-- 高级选项 -->
      <n-collapse>
        <n-collapse-item title="高级选项 — 包名 / 版本号">
          <div class="space-y-3">
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
                placeholder="自动填充"
                :allow-input="onlyAllowPkg"
                :style="{ width: '33%' }"
                maxlength="20"
              />
            </n-input-group>
            <div class="text-xs text-gray-400 font-mono">
              包名预览：{{ previewPkgName }}
            </div>
            <div class="border-t border-gray-200 dark:border-gray-700 pt-3">
              <label class="text-xs text-gray-500 dark:text-gray-400 block mb-1">自定义版本号（纯数字，留空自动生成）</label>
              <n-input
                v-model:value="versionName"
                placeholder="例如 7 或 20240701"
                :allow-input="onlyAllowDigits"
                maxlength="20"
              />
            </div>
          </div>
        </n-collapse-item>
      </n-collapse>

      <!-- 生成按钮 -->
      <button
        :disabled="!canBuild || building"
        class="w-full py-3 rounded-xl text-white font-semibold text-base transition-all duration-200"
        :class="canBuild && !building
          ? 'bg-linear-to-r from-blue-500 to-purple-600 hover:from-blue-600 hover:to-purple-700 shadow-lg hover:shadow-xl active:scale-[0.98]'
          : 'bg-gray-400 dark:bg-gray-600 cursor-not-allowed'"
        @click="buildAPK"
      >
        <span v-if="building" class="flex items-center justify-center gap-2">
          <svg class="animate-spin h-5 w-5" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" fill="none"/>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"/>
          </svg>
          正在生成 APK...
        </span>
        <span v-else>生成 APK</span>
      </button>

      <!-- 日志 -->
      <div v-if="buildLog" class="bg-black/5 dark:bg-white/5 rounded-lg p-3">
        <pre class="text-xs font-mono text-gray-600 dark:text-gray-400 whitespace-pre-wrap max-h-48 overflow-y-auto">{{ buildLog }}</pre>
      </div>

      <!-- 结果 -->
      <div v-if="result" class="bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-lg p-4 text-sm">
        <p class="text-green-700 dark:text-green-300 font-medium">✓ APK 生成成功</p>
        <p class="text-green-600 dark:text-green-400 mt-1 font-mono text-xs">{{ result.APKPath }}</p>
      </div>

      <div v-if="error" class="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4 text-sm">
        <p class="text-red-700 dark:text-red-300 font-medium">✗ 生成失败</p>
        <p class="text-red-600 dark:text-red-400 mt-1">{{ error }}</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { BuildAPK, GetIconPaths, SelectDirectory, SelectFile } from '../../wailsjs/go/main/App'
import { CheckmarkCircle24Regular } from '@vicons/fluent'
import {
  NInput, NButton, NTag, NIcon, NTabPane, NTabs, NCard,
  NInputGroup, NCollapseItem, NCollapse,
} from 'naive-ui'

// ── State ──
const inputTab = ref('folder')
const projectPath = ref('')
const filePath = ref('')
const appName = ref('')
const building = ref(false)
const buildLog = ref('')
const result = ref<{ APKPath: string; PackageName: string; VersionName: string; VersionCode: number } | null>(null)
const error = ref('')

// Package name segments
const pkgSegment1 = ref('com')
const pkgSegment2 = ref('')
const pkgSegment3 = ref('')
const lastSeg2 = ref('')

// Custom version number (digits only)
const versionName = ref('')

// Icon
const fileInputRef = ref<HTMLInputElement>()
const customIconData = ref<string | null>(null)
const customIconFileName = ref('')

// ── Computed ──
const isFolderSelected = computed(() => !!projectPath.value)
const isZipSelected = computed(() => !!filePath.value && filePath.value.toLowerCase().endsWith('.zip'))
const isHtmlSelected = computed(() => !!filePath.value && (filePath.value.toLowerCase().endsWith('.html') || filePath.value.toLowerCase().endsWith('.htm')))

const iconLetter = computed(() => {
  if (customIconData.value) return ''
  return appName.value ? appName.value.charAt(0).toUpperCase() : '?'
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
  }
}

function onPathEnter() {
  projectPath.value = projectPath.value.replace(/^"(.*)"$/, '$1').trim()
}

async function selectFile() {
  const file = await SelectFile()
  if (file) {
    filePath.value = file
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
  return !value || /^[a-z0-9_.-]*$/.test(value)
}

function onlyAllowDigits(value: string) {
  return !value || /^\d*$/.test(value)
}

function previewPkg() {
  // preview is updated reactively via computed property
}

// ── Build ──
async function buildAPK() {
  if (!canBuild.value || building.value) return

  building.value = true
  buildLog.value = ''
  result.value = null
  error.value = ''

  try {
    // Determine project directory
    let projectDir = ''
    if (inputTab.value === 'folder') {
      projectDir = projectPath.value
    } else {
      // TODO: extract zip or use html file
      buildLog.value += '文件模式暂未实现，请使用文件夹模式\n'
      building.value = false
      return
    }

    // Generate package name if not custom
    const customPkg = pkgSegment3.value ? `${pkgSegment1.value}.${pkgSegment2.value}.${pkgSegment3.value}` : ''
    if (customPkg) {
      buildLog.value += `使用自定义包名: ${customPkg}\n`
    }

    // Generate icons as WebP
    const icons: Record<string, string> = {}
    const iconPaths = await GetIconPaths()

    if (customIconData.value) {
      // Uploaded image → resize to all icon sizes
      for (const iconPath of iconPaths) {
        const size = extractIconSize(iconPath)
        if (size > 0) {
          const blob = await resizeImage(customIconData.value, size, iconPath.includes('foreground'))
          const b64 = await blobToBase64(blob)
          icons[iconPath] = b64
        }
      }
    } else {
      // Auto-generated text icon
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

    buildLog.value += '正在构建 APK...\n'
    const res = await BuildAPK(projectDir, appName.value, versionName.value, icons)
    result.value = { APKPath: res.APKPath, PackageName: res.PackageName, VersionName: res.VersionName, VersionCode: res.VersionCode }
    buildLog.value += res.Log || ''
  } catch (err: any) {
    error.value = String(err?.message || err || '未知错误')
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

    ctx.clearRect(0, 0, size, size) // start fully transparent

    if (!isForeground) {
      // Launcher icon: rounded rectangle with color
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
      // Foreground: draw the bgColor as rounded rect too, so the icon
      // has a visible colored shape even before adaptive icon compositing.
      // Android's adaptive icon masks this with the device shape.
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

    // Draw center letter (white on both foreground and launcher)
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

      // Clip to rounded rectangle (same for foreground and launcher)
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

      // Crop image to square center and draw
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
      // Strip the data:...;base64, prefix
      const comma = result.indexOf(',')
      resolve(comma >= 0 ? result.substring(comma + 1) : result)
    }
    reader.onerror = () => reject(new Error('Failed to read blob'))
    reader.readAsDataURL(blob)
  })
}
</script>
