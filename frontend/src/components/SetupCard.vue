<template>
  <div class="flex items-center justify-center p-4"
    :style="{ '--theme-color': themeColor }">
    <div
      class="select-none w-full max-w-3xl h-[85vh] rounded-xl relative overflow-hidden">
      <!-- 置顶固定栏 -->
      <div
        class="absolute top-0 left-0 right-0 h-18 px-6 z-10 flex justify-between items-center
        pointer-events-none bg-linear-to-b from-gray-50 via-gray-50 to-transparent dark:from-gray-900 dark:via-gray-900 dark:to-transparent">
        <h1 class="text-xl font-medium text-gray-800 dark:text-gray-100">{{ t('setup.continueSetup') }}</h1>
        <button
          @click="$emit('finish')"
          class="setup-btn pointer-events-auto cursor-pointer px-6 py-2 text-white rounded-xl font-medium transition-all flex items-center gap-3 outline-none focus-visible:outline-none">
          <span>{{ t('setup.finish') }}</span>
          <span class="flex items-center gap-1 text-xs opacity-75">
            <kbd class="px-1.5 py-0.5 rounded bg-white/10 border border-white/10">Ctrl</kbd>
            <span>+</span>
            <kbd class="px-1.5 py-0.5 rounded bg-white/10 border border-white/10">Enter</kbd>
          </span>
        </button>
      </div>

      <!-- 内容区 -->
      <n-scrollbar>
        <div class="h-full pt-20 px-6 pb-6 overflow-y-auto">

          <!-- 引导文案 -->
          <section v-if="showIntro" class="mb-8">
            <h2 class="text-lg font-medium text-gray-800 dark:text-gray-100 mb-3">{{ t('setup.introTitle', { appName }) }}</h2>
            <div class="mt-3 flex flex-col gap-2 text-sm text-gray-500 dark:text-gray-400">
              <p>{{ t('setup.introDefault', { appName }) }}</p>
              <p>{{ t('setup.introChangeLater') }}</p>
            </div>
          </section>

          <slot />
          <div class="h-6"></div>
        </div>
      </n-scrollbar>
      <div
        class="absolute bottom-0 left-0 right-0 h-6 z-10
        pointer-events-none bg-linear-to-t from-gray-50 via-gray-50 to-transparent dark:from-gray-900 dark:via-gray-900 dark:to-transparent">
      </div>
    </div>
  </div>
</template>

<script setup>
import { useI18n } from 'vue-i18n'
import { NScrollbar } from 'naive-ui'

const { t } = useI18n()

defineProps({
  appName: {
    type: String,
    default: 'App',
  },
  themeColor: {
    type: String,
    default: 'oklch(0.716 0.191 2.477)',
  },
  showIntro: {
    type: Boolean,
    default: true,
  },
})

defineEmits(['finish'])
</script>

<style scoped>
.setup-btn {
  background-color: var(--theme-color);
}

@media (prefers-color-scheme: light) {
  .setup-btn:hover {
    background-color: color-mix(in oklch, var(--theme-color), black 15%);
  }
}

@media (prefers-color-scheme: dark) {
  .setup-btn:hover {
    background-color: color-mix(in oklch, var(--theme-color), white 15%);
  }
}

@media (prefers-color-scheme: light) {
  .setup-btn:focus-visible {
    box-shadow: 0 0 0 2px color-mix(in srgb, var(--theme-color), transparent 60%);
  }
}

@media (prefers-color-scheme: dark) {
  .setup-btn:focus-visible {
    box-shadow: 0 0 0 2px color-mix(in srgb, var(--theme-color), white 40%);
  }
}
</style>
