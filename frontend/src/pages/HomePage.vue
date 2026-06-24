<template>
  <div class="h-full overflow-y-auto p-6">
    <!-- Header -->
    <div class="mb-8">
      <h1 class="text-2xl font-bold text-gray-800 dark:text-gray-100">
        Foil
      </h1>
      <p class="text-sm text-gray-500 dark:text-gray-400 mt-1">
        HTML → APK 一键打包
      </p>
    </div>

    <!-- Build Form -->
    <div class="max-w-2xl space-y-6">
      <!-- Project Directory -->
      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
          HTML 项目目录
        </label>
        <div class="flex gap-2">
          <input
            v-model="projectDir"
            type="text"
            readonly
            placeholder="选择包含 index.html 的文件夹..."
            class="flex-1 px-3 py-2 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-gray-800 dark:text-gray-200 text-sm cursor-pointer"
            @click="selectDir"
          />
          <button
            class="px-4 py-2 rounded-lg bg-blue-500 hover:bg-blue-600 text-white text-sm font-medium transition-colors"
            @click="selectDir"
          >
            选择
          </button>
        </div>
      </div>

      <!-- App Name -->
      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
          应用名称
        </label>
        <input
          v-model="appName"
          type="text"
          placeholder="例如：我的应用"
          class="w-full px-3 py-2 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-gray-800 dark:text-gray-200 text-sm"
          @input="onAppNameChange"
        />
      </div>

      <!-- Auto-generated info -->
      <div class="bg-gray-50 dark:bg-gray-800/50 rounded-lg p-4 space-y-2">
        <div class="flex justify-between text-sm">
          <span class="text-gray-500 dark:text-gray-400">包名</span>
          <span class="text-gray-700 dark:text-gray-300 font-mono text-xs">{{ autoPkgName }}</span>
        </div>
        <div class="flex justify-between text-sm">
          <span class="text-gray-500 dark:text-gray-400">版本</span>
          <span class="text-gray-700 dark:text-gray-300 font-mono text-xs">{{ autoVersionName }} ({{ autoVersionCode }})</span>
        </div>
      </div>

      <!-- Icon Preview -->
      <div v-if="iconPreviewUrl">
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
          自动生成图标预览
        </label>
        <img
          :src="iconPreviewUrl"
          alt="App Icon"
          class="w-16 h-16 rounded-xl shadow-sm border border-gray-200 dark:border-gray-700"
        />
      </div>

      <!-- Advanced Toggle -->
      <details class="text-sm">
        <summary class="cursor-pointer text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200">
          高级选项
        </summary>
        <div class="mt-3 space-y-3 pl-4 border-l-2 border-gray-200 dark:border-gray-700">
          <div>
            <label class="block text-xs text-gray-500 dark:text-gray-400 mb-1">自定义包名（留空自动生成）</label>
            <input
              v-model="customPkgName"
              type="text"
              placeholder="com.example.myapp"
              class="w-full px-3 py-1.5 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-gray-800 dark:text-gray-200 text-xs font-mono"
            />
          </div>
        </div>
      </details>

      <!-- Generate Button -->
      <button
        :disabled="!canBuild || building"
        class="w-full py-3 rounded-xl text-white font-semibold text-base transition-all duration-200"
        :class="canBuild && !building
          ? 'bg-gradient-to-r from-blue-500 to-purple-600 hover:from-blue-600 hover:to-purple-700 shadow-lg hover:shadow-xl active:scale-[0.98]'
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

      <!-- Progress / Log -->
      <div v-if="buildLog" class="bg-black/5 dark:bg-white/5 rounded-lg p-3">
        <pre class="text-xs font-mono text-gray-600 dark:text-gray-400 whitespace-pre-wrap max-h-48 overflow-y-auto">{{ buildLog }}</pre>
      </div>

      <!-- Result -->
      <div
        v-if="result"
        class="bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-lg p-4 text-sm"
      >
        <p class="text-green-700 dark:text-green-300 font-medium">✓ APK 生成成功</p>
        <p class="text-green-600 dark:text-green-400 mt-1 font-mono text-xs">{{ result.apkPath }}</p>
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
import { BuildAPK, GetIconPaths, GeneratePackageName } from '../../wailsjs/go/main/App'

const projectDir = ref('')
const appName = ref('')
const customPkgName = ref('')
const autoVersionName = ref('')
const autoVersionCode = ref('')
const autoPkgName = ref('')
const building = ref(false)
const buildLog = ref('')
const result = ref<{ apkPath: string; packageName: string; versionName: string; versionCode: number } | null>(null)
const error = ref('')
const iconPreviewUrl = ref('')

const canBuild = computed(() => projectDir.value && appName.value)

function onAppNameChange() {
  if (appName.value) {
    // Generate package name for preview
    GeneratePackageName(appName.value).then((pkg: string) => {
      autoPkgName.value = pkg
    })
    generateIconPreview(appName.value)
  }
}

function generateIconPreview(name: string) {
  if (!name) return
  const canvas = document.createElement('canvas')
  canvas.width = 192
  canvas.height = 192
  const ctx = canvas.getContext('2d')!
  const letter = name.charAt(0).toUpperCase()

  // Background
  const colors = ['#4F46E5', '#7C3AED', '#EC4899', '#F59E0B', '#10B981', '#3B82F6', '#EF4444']
  const colorIndex = name.charCodeAt(0) % colors.length
  ctx.fillStyle = colors[colorIndex]
  ctx.beginPath()
  ctx.arc(96, 96, 88, 0, Math.PI * 2)
  ctx.fill()

  // Letter
  ctx.fillStyle = '#FFFFFF'
  ctx.font = 'bold 96px system-ui, -apple-system, sans-serif'
  ctx.textAlign = 'center'
  ctx.textBaseline = 'middle'
  ctx.fillText(letter, 96, 98)

  // Set preview
  iconPreviewUrl.value = canvas.toDataURL('image/png')
}

async function selectDir() {
  // In Wails, we can use the runtime dialog or a file input
  // For now, use a prompt (Wails can expose a dialog via Go)
  const dir = prompt('请输入 HTML 项目文件夹路径：')
  if (dir) {
    projectDir.value = dir
  }
}

async function buildAPK() {
  if (!canBuild.value || building.value) return

  building.value = true
  buildLog.value = ''
  result.value = null
  error.value = ''

  try {
    // Generate all icon sizes as WebP via canvas
    const icons: Record<string, number[]> = {}
    const iconPaths = await GetIconPaths()
    const letter = appName.value.charAt(0).toUpperCase()
    const colorIndex = appName.value.charCodeAt(0) % 7
    const colors = ['#4F46E5', '#7C3AED', '#EC4899', '#F59E0B', '#10B981', '#3B82F6', '#EF4444']

    for (const iconPath of iconPaths) {
      // Extract size from path (e.g., "mipmap-hdpi-v4/ic_launcher.webp" → 72)
      const size = extractIconSize(iconPath)
      if (size > 0) {
        const blob = await generateWebPIcon(letter, colors[colorIndex], size)
        const buf = await blob.arrayBuffer()
        icons[iconPath] = Array.from(new Uint8Array(buf))
      }
    }

    buildLog.value += '正在构建 APK...\n'
    const res = await BuildAPK(projectDir.value, appName.value, icons)
    result.value = { apkPath: res.APKPath, packageName: res.PackageName, versionName: res.VersionName, versionCode: res.VersionCode }
    buildLog.value += res.Log || ''
  } catch (err: any) {
    error.value = String(err?.message || err || '未知错误')
  } finally {
    building.value = false
  }
}

function extractIconSize(path: string): number {
  const sizes: Record<string, number> = {
    'mdpi': 48,
    'hdpi': 72,
    'xhdpi': 96,
    'xxhdpi': 144,
    'xxxhdpi': 192,
  }
  for (const [key, size] of Object.entries(sizes)) {
    if (path.includes(key)) {
      // Foreground icons are 1.5x or 2.25x larger depending on density
      if (path.includes('foreground')) {
        return Math.round(size * (108 / 48))
      }
      return size
    }
  }
  return 0
}

function generateWebPIcon(letter: string, bgColor: string, size: number): Promise<Blob> {
  return new Promise((resolve, reject) => {
    const canvas = document.createElement('canvas')
    canvas.width = size
    canvas.height = size
    const ctx = canvas.getContext('2d')!
    const half = size / 2

    // Draw circle background
    ctx.fillStyle = bgColor
    ctx.beginPath()
    ctx.arc(half, half, half - size * 0.04, 0, Math.PI * 2)
    ctx.fill()

    // Draw letter
    ctx.fillStyle = '#FFFFFF'
    const fontSize = Math.round(size * 0.5)
    ctx.font = `bold ${fontSize}px system-ui, -apple-system, sans-serif`
    ctx.textAlign = 'center'
    ctx.textBaseline = 'middle'
    ctx.fillText(letter, half, half + 1)

    // Encode as WebP
    canvas.toBlob((blob) => {
      if (blob) resolve(blob)
      else reject(new Error('Failed to encode WebP'))
    }, 'image/webp', 0.9)
  })
}
</script>
