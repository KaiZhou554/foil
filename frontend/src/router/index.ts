import { createRouter, createWebHashHistory } from 'vue-router'
import { useAppStore } from '@/stores/appStore'
import type { Pinia } from 'pinia'

// Lazy-loaded pages for code splitting
const WelcomePage = () => import('@/pages/WelcomePage.vue')
const SetupPage = () => import('@/pages/SetupPage.vue')
const MainLayout = () => import('@/pages/MainLayout.vue')
const HomePage = () => import('@/pages/HomePage.vue')
const SettingsPage = () => import('@/pages/SettingsPage.vue')
const PrivacyPage = () => import('@/pages/PrivacyPage.vue')

const router = createRouter({
  // Hash history is standard for Wails/Electron apps (no real server for pushState)
  history: createWebHashHistory(),
  routes: [
    {
      path: '/',
      redirect: () => {
        // Redirect based on first-launch state
        // We access the store via the pinia instance passed at creation time
        // For the initial redirect, default to /welcome (the route guard will adjust)
        return '/welcome'
      },
    },
    {
      path: '/welcome',
      name: 'welcome',
      component: WelcomePage,
    },
    {
      path: '/setup',
      name: 'setup',
      component: SetupPage,
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
        {
          path: 'settings/privacy',
          name: 'settings-privacy',
          component: PrivacyPage,
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

    // Allow all routes while config is still loading (App.vue will redirect after load)
    if (!appStore.configLoaded) {
      return true
    }

    // If still in first-launch onboarding, restrict to onboarding routes
    if (appStore.isFirstLaunch) {
      if (to.path === '/welcome' || to.path === '/setup') {
        return true
      }
      // Redirect any other route to /welcome
      return '/welcome'
    }

    // If onboarding is complete, redirect onboarding routes to main
    if (to.path === '/welcome' || to.path === '/setup') {
      return '/main/home'
    }

    return true
  })
}

export default router
