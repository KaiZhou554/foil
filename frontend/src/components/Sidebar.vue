<template>
  <aside
    ref="sidebarRef"
    :class="[
      'absolute top-0 left-0 h-full z-10',
      'transition-all duration-300 ease-in-out overflow-hidden',
      sidebarOpen
        ? 'w-48 bg-white dark:bg-neutral-800 rounded-r-lg shadow-sm border-r border-neutral-200/50 dark:border-neutral-700/50'
        : 'w-12',
    ]"
  >
    <div class="relative h-full pt-2 px-2 flex flex-col items-start gap-1">
      <!-- Active indicator -->
      <span
        v-if="indicatorHeight !== '0px'"
        :style="{
          top: indicatorTop,
          height: indicatorHeight,
          transform: `scaleY(${indicatorScale})`,
        }"
        class="sidebar-indicator bg-taffy-400 dark:bg-taffy-200 transition-colors"
      />

      <button
        v-for="item in navItems"
        :key="item.route"
        :ref="el => setButtonRef(el as HTMLElement | null, item.route)"
        @click="navigate(item.route)"
        :class="[
          'w-full flex items-center rounded-sm text-sm transition-colors',
          isActive(item.route)
            ? 'bg-neutral-200/60 dark:bg-neutral-700/60 text-neutral-900 dark:text-white'
            : 'text-neutral-600 dark:text-neutral-400 hover:bg-neutral-100 dark:hover:bg-neutral-800',
        ]"
      >
        <span class="text-lg shrink-0">
          <Icon class="mt-2 ml-1.5" size="20">
            <component :is="item.icon" />
          </Icon>
        </span>
        <span
          :class="[
            'text-sm ml-2 whitespace-nowrap overflow-hidden',
            sidebarOpen ? 'block' : 'hidden',
          ]"
        >{{ item.label }}</span>
      </button>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { ref, watch, nextTick, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { Icon } from '@vicons/utils'
import Home20Regular from '@vicons/fluent/es/Home20Regular'
import Settings20Regular from '@vicons/fluent/es/Settings20Regular'

const { t } = useI18n()
const router = useRouter()
const route = useRoute()

const props = defineProps<{
  sidebarOpen: boolean
}>()

interface NavItem {
  route: string
  label: string
  icon: typeof Home20Regular
}

const navItems = computed<NavItem[]>(() => [
  {
    route: '/main/home',
    label: t('sidebar.home'),
    icon: Home20Regular,
  },
  {
    route: '/main/settings',
    label: t('sidebar.settings'),
    icon: Settings20Regular,
  },
])

function isActive(routePath: string): boolean {
  return route.path === routePath || route.path.startsWith(routePath + '/')
}

function navigate(routePath: string) {
  router.push(routePath)
}

// ---- Indicator animation (preserved from original) ----
const sidebarRef = ref<HTMLElement | null>(null)
const buttonRefs = ref<Record<string, HTMLElement | null>>({})
const indicatorTop = ref<string>('0px')
const indicatorHeight = ref<string>('0px')
const indicatorScale = ref<number>(1)

const setButtonRef = (el: HTMLElement | null, routeName: string) => {
  buttonRefs.value[routeName] = el
}

const updateIndicator = async () => {
  await nextTick()

  const sidebar = sidebarRef.value
  const activeRoute = navItems.value.find(item => isActive(item.route))
  if (!activeRoute) return
  const targetButton = buttonRefs.value[activeRoute.route]

  if (!sidebar || !targetButton) return

  const sidebarRect = sidebar.getBoundingClientRect()
  const targetRect = targetButton.getBoundingClientRect()

  const buttonCenter = targetRect.top - sidebarRect.top + targetRect.height / 2

  const desiredHeight = 20
  const indicatorTopPos = buttonCenter - desiredHeight / 2

  indicatorTop.value = `${indicatorTopPos}px`
  indicatorHeight.value = `${desiredHeight}px`

  const previousTop = parseFloat(indicatorTop.value) || buttonCenter
  const distance = Math.abs(buttonCenter - previousTop)
  indicatorScale.value = 1 + Math.min(distance / 100, 0.6)

  setTimeout(() => {
    indicatorScale.value = 1
  }, 240)
}

watch(() => route.path, updateIndicator)
onMounted(() => {
  updateIndicator()
})
</script>

<style scoped>
.sidebar-indicator {
  position: absolute;
  left: 0.45rem;
  width: 0.15rem;
  border-radius: 9999px;
  transition: top 0.3s ease, transform 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
}
</style>
