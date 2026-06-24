<template>
  <div class="h-full">
    <div
      class="w-full h-full ring-1 ring-black/5 dark:ring-white/10 overflow-hidden flex flex-col transition-colors"
    >
      <div class="h-full relative overflow-hidden">
        <!-- Backdrop overlay — closes sidebar when clicking outside -->
        <div
          v-if="appStore.sidebarOpen"
          class="absolute inset-0 z-[5]"
          @click="appStore.toggleSidebar()"
        />
        <Sidebar :sidebar-open="appStore.sidebarOpen" />

        <!-- Page content area with Naive UI scrollbar -->
        <div
          class="h-full ml-12 bg-gray-50 dark:bg-gray-900 rounded-tl-md"
        >
          <n-scrollbar>
            <router-view v-slot="{ Component, route }">
              <transition :name="nestedTransition" mode="out-in">
                <component :is="Component" :key="route.fullPath" />
              </transition>
            </router-view>
          </n-scrollbar>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { NScrollbar } from 'naive-ui'
import Sidebar from '@/components/Sidebar.vue'
import { useAppStore } from '@/stores/appStore'

const appStore = useAppStore()
const router = useRouter()

// ── Nested page transition (home ↔ settings ↔ settings-privacy) ──
const nestedTransition = ref('fade')

/** Route-name → Route-name → CSS transition name */
const nestedTransitionMap: Record<string, Record<string, string>> = {
  home:             { settings: 'slide-up' },
  settings:         { home: 'slide-up', 'settings-privacy': 'slide-left' },
  'settings-privacy': { settings: 'slide-right' },
}

watch(
  () => router.currentRoute.value.name,
  (to, from) => {
    const fromName = (from as string | undefined) ?? null
    const toName = (to as string | undefined) ?? null
    if (!fromName || !toName) {
      nestedTransition.value = 'fade'
      return
    }
    nestedTransition.value = nestedTransitionMap[fromName]?.[toName] ?? 'fade'
  },
)
</script>
