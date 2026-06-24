import './assets/main.css'
import './assets/font.css'

import { createApp } from 'vue'
import { createI18n } from 'vue-i18n'
import { createPinia } from 'pinia'
import App from './App.vue'
import router, { setupRouterGuards } from '@/router'
import zhCN from '@/locales/zh-CN.ts'
import en from '@/locales/en.ts'

const pinia = createPinia()

const i18n = createI18n({
  legacy: false,
  locale: 'zh-CN',
  fallbackLocale: 'en',
  messages: {
    'zh-CN': zhCN,
    'en': en,
  },
})

// Setup router guards after pinia is created (guard needs pinia instance)
setupRouterGuards(pinia)

createApp(App)
  .use(pinia)
  .use(router)
  .use(i18n)
  .mount('#app')
