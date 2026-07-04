import { createRouter, createWebHashHistory } from 'vue-router'
import { useAppStore } from '@/stores/appStore'
import type { Pinia } from 'pinia'

// Lazy-loaded pages for code splitting
const WelcomePage = () => import('@/pages/WelcomePage.vue')
const MainLayout = () => import('@/pages/MainLayout.vue')
const HomePage = () => import('@/pages/HomePage.vue')
const SettingsPage = () => import('@/pages/SettingsPage.vue')

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    {
      path: '/',
      redirect: '/welcome',
    },
    {
      path: '/welcome',
      name: 'welcome',
      component: WelcomePage,
    },
    {
      path: '/main',
      component: MainLayout,
      children: [
        {
          path: '',
          redirect: { name: 'home' },
        },
        {
          path: 'home',
          name: 'home',
          component: HomePage,
        },
        {
          path: 'settings',
          name: 'settings',
          component: SettingsPage,
        },
      ],
    },
  ],
})

/**
 * Setup the navigation guard.
 * Must be called after the Pinia instance is created.
 */
export function setupRouterGuards(pinia: Pinia) {
  router.beforeEach((to, _from) => {
    const appStore = useAppStore(pinia)

    if (!appStore.configLoaded) {
      return true
    }

    // Once onboarding is done, redirect welcome to home
    if (!appStore.isFirstLaunch && to.path === '/welcome') {
      return '/main/home'
    }

    return true
  })
}

export default router
