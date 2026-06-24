<template>
  <div>
    <!-- ==================== 第一部分：存储 ==================== -->
    <section class="mb-8">
      <h2 v-if="store.optionDisplay.sectionStorage" class="text-lg font-medium text-gray-800 dark:text-gray-100 mb-3">{{ t('privacy.sections.storage') }}</h2>
      <div class="text-base">
        <div v-if="store.optionDisplay.storageDefault" class="hover:bg-gray-600/5 dark:hover:bg-gray-300/5 rounded-full py-2 px-4 transition-colors duration-300 flex justify-between items-center mb-2">
          <div class="flex items-center gap-2">
            <Icon class="text-gray-600 dark:text-gray-300"><Folder16Regular /></Icon>
            <span class="text-gray-600 dark:text-gray-300">{{ t('privacy.options.storageDefault') }}</span>
          </div>
          <n-switch disabled v-model:value="store.optionValues.storageDefault" />
        </div>

        <div v-if="store.optionDisplay.standardFolder" class="hover:bg-gray-600/5 dark:hover:bg-gray-300/5 rounded-full py-2 px-4 transition-colors duration-300 flex justify-between items-center mb-2">
          <div class="flex items-center gap-2">
            <Icon class="text-gray-600 dark:text-gray-300"><Image20Regular /></Icon>
            <span class="text-gray-600 dark:text-gray-300">{{ t('privacy.options.standardFolder') }}</span>
          </div>
          <n-switch v-model:value="store.optionValues.standardFolder" />
        </div>

        <div v-if="store.optionDisplay.storageCustom" class="hover:bg-gray-600/5 dark:hover:bg-gray-300/5 rounded-full py-2 px-4 transition-colors duration-300 flex justify-between items-center mb-2">
          <div class="flex items-center gap-2">
            <Icon class="text-gray-600 dark:text-gray-300"><Storage24Regular /></Icon>
            <span class="text-gray-600 dark:text-gray-300">{{ t('privacy.options.storageCustom') }}</span>
          </div>
          <n-switch v-model:value="store.optionValues.storageCustom" />
        </div>
      </div>
    </section>

    <!-- ==================== 第二部分：设备与系统 ==================== -->
    <section class="mb-8">
      <h2 v-if="store.optionDisplay.sectionDevice" class="text-lg font-medium text-gray-800 dark:text-gray-100 mb-3">{{ t('privacy.sections.device') }}</h2>
      <div class="text-base">
        <div v-if="store.optionDisplay.adminPrivilege" class="hover:bg-gray-600/5 dark:hover:bg-gray-300/5 rounded-full py-2 px-4 transition-colors duration-300 flex justify-between items-center mb-2">
          <div class="flex items-center gap-2">
            <Icon class="text-gray-600 dark:text-gray-300"><ShieldKeyhole20Regular /></Icon>
            <span class="text-gray-600 dark:text-gray-300">{{ t('privacy.options.adminPrivilege') }}</span>
          </div>
          <n-switch v-model:value="store.optionValues.adminPrivilege" />
        </div>

        <div v-if="store.optionDisplay.identityVerification" class="hover:bg-gray-600/5 dark:hover:bg-gray-300/5 rounded-full py-2 px-4 transition-colors duration-300 flex justify-between items-center mb-2">
          <div class="flex items-center gap-2">
            <Icon class="text-gray-600 dark:text-gray-300"><Fingerprint24Regular /></Icon>
            <span class="text-gray-600 dark:text-gray-300">{{ t('privacy.options.identityVerification') }}</span>
          </div>
          <n-switch v-model:value="store.optionValues.identityVerification" />
        </div>

        <div v-if="store.optionDisplay.modifySystemSettings" class="hover:bg-gray-600/5 dark:hover:bg-gray-300/5 rounded-full py-2 px-4 transition-colors duration-300 flex justify-between items-center mb-2">
          <div class="flex items-center gap-2">
            <Icon class="text-gray-600 dark:text-gray-300"><Settings16Regular /></Icon>
            <span class="text-gray-600 dark:text-gray-300">{{ t('privacy.options.modifySystemSettings') }}</span>
          </div>
          <n-switch v-model:value="store.optionValues.modifySystemSettings" />
        </div>

        <div v-if="store.optionDisplay.inputMonitor" class="hover:bg-gray-600/5 dark:hover:bg-gray-300/5 rounded-full py-2 px-4 transition-colors duration-300 flex justify-between items-center mb-2">
          <div class="flex items-center gap-2">
            <Icon class="text-gray-600 dark:text-gray-300"><Keyboard20Regular /></Icon>
            <span class="text-gray-600 dark:text-gray-300">{{ t('privacy.options.inputMonitor') }}</span>
          </div>
          <n-switch v-model:value="store.optionValues.inputMonitor" />
        </div>

        <div v-if="store.optionDisplay.microphone" class="hover:bg-gray-600/5 dark:hover:bg-gray-300/5 rounded-full py-2 px-4 transition-colors duration-300 flex justify-between items-center mb-2">
          <div class="flex items-center gap-2">
            <Icon class="text-gray-600 dark:text-gray-300"><Mic20Regular /></Icon>
            <span class="text-gray-600 dark:text-gray-300">{{ t('privacy.options.microphone') }}</span>
          </div>
          <n-switch v-model:value="store.optionValues.microphone" />
        </div>

        <div v-if="store.optionDisplay.camera" class="hover:bg-gray-600/5 dark:hover:bg-gray-300/5 rounded-full py-2 px-4 transition-colors duration-300 flex justify-between items-center mb-2">
          <div class="flex items-center gap-2">
            <Icon class="text-gray-600 dark:text-gray-300"><Camera16Regular /></Icon>
            <span class="text-gray-600 dark:text-gray-300">{{ t('privacy.options.camera') }}</span>
          </div>
          <n-switch v-model:value="store.optionValues.camera" />
        </div>

        <div v-if="store.optionDisplay.location" class="hover:bg-gray-600/5 dark:hover:bg-gray-300/5 rounded-full py-2 px-4 transition-colors duration-300 flex justify-between items-center mb-2">
          <div class="flex items-center gap-2">
            <Icon class="text-gray-600 dark:text-gray-300"><Location16Regular /></Icon>
            <span class="text-gray-600 dark:text-gray-300">{{ t('privacy.options.location') }}</span>
          </div>
          <n-switch v-model:value="store.optionValues.location" />
        </div>

        <div v-if="store.optionDisplay.sensor" class="hover:bg-gray-600/5 dark:hover:bg-gray-300/5 rounded-full py-2 px-4 transition-colors duration-300 flex justify-between items-center mb-2">
          <div class="flex items-center gap-2">
            <Icon class="text-gray-600 dark:text-gray-300"><Live20Regular /></Icon>
            <span class="text-gray-600 dark:text-gray-300">{{ t('privacy.options.sensor') }}</span>
          </div>
          <n-switch v-model:value="store.optionValues.sensor" />
        </div>

        <div v-if="store.optionDisplay.clipboard" class="hover:bg-gray-600/5 dark:hover:bg-gray-300/5 rounded-full py-2 px-4 transition-colors duration-300 flex justify-between items-center mb-2">
          <div class="flex items-center gap-2">
            <Icon class="text-gray-600 dark:text-gray-300"><Clipboard16Regular /></Icon>
            <span class="text-gray-600 dark:text-gray-300">{{ t('privacy.options.clipboard') }}</span>
          </div>
          <n-switch v-model:value="store.optionValues.clipboard" />
        </div>

        <div v-if="store.optionDisplay.screenshot" class="hover:bg-gray-600/5 dark:hover:bg-gray-300/5 rounded-full py-2 px-4 transition-colors duration-300 flex justify-between items-center mb-2">
          <div class="flex items-center gap-2">
            <Icon class="text-gray-600 dark:text-gray-300"><Screenshot20Regular /></Icon>
            <span class="text-gray-600 dark:text-gray-300">{{ t('privacy.options.screenshot') }}</span>
          </div>
          <n-switch v-model:value="store.optionValues.screenshot" />
        </div>

        <div v-if="store.optionDisplay.otherDevice" class="hover:bg-gray-600/5 dark:hover:bg-gray-300/5 rounded-full py-2 px-4 transition-colors duration-300 flex justify-between items-center mb-2">
          <div class="flex items-center gap-2">
            <Icon class="text-gray-600 dark:text-gray-300"><UsbStick20Regular /></Icon>
            <span class="text-gray-600 dark:text-gray-300">{{ t('privacy.options.otherDevice') }}</span>
          </div>
          <n-switch v-model:value="store.optionValues.otherDevice" />
        </div>

        <div v-if="store.optionDisplay.accessOtherAppData" class="hover:bg-gray-600/5 dark:hover:bg-gray-300/5 rounded-full py-2 px-4 transition-colors duration-300 flex justify-between items-center mb-2">
          <div class="flex items-center gap-2">
            <Icon class="text-gray-600 dark:text-gray-300"><Database20Regular /></Icon>
            <span class="text-gray-600 dark:text-gray-300">{{ t('privacy.options.accessOtherAppData') }}</span>
          </div>
          <n-switch v-model:value="store.optionValues.accessOtherAppData" />
        </div>

        <div v-if="store.optionDisplay.manageOtherApps" class="hover:bg-gray-600/5 dark:hover:bg-gray-300/5 rounded-full py-2 px-4 transition-colors duration-300 flex justify-between items-center mb-2">
          <div class="flex items-center gap-2">
            <Icon class="text-gray-600 dark:text-gray-300"><WindowApps20Regular /></Icon>
            <span class="text-gray-600 dark:text-gray-300">{{ t('privacy.options.manageOtherApps') }}</span>
          </div>
          <n-switch v-model:value="store.optionValues.manageOtherApps" />
        </div>
      </div>
    </section>

    <!-- ==================== 第三部分：在线功能 ==================== -->
    <section class="mb-8">
      <h2 v-if="store.optionDisplay.sectionOnline" class="text-lg font-medium text-gray-800 dark:text-gray-100 mb-3">{{ t('privacy.sections.online') }}</h2>
      <div class="text-base">
        <div v-if="store.optionDisplay.network" class="hover:bg-gray-600/5 dark:hover:bg-gray-300/5 rounded-full py-2 px-4 transition-colors duration-300 flex justify-between items-center mb-2">
          <div class="flex items-center gap-2">
            <Icon class="text-gray-600 dark:text-gray-300"><Globe16Regular /></Icon>
            <span class="text-gray-600 dark:text-gray-300">{{ t('privacy.options.network') }}</span>
          </div>
          <n-switch v-model:value="store.optionValues.network" />
        </div>

        <div v-if="store.optionDisplay.thirdPartyService" class="hover:bg-gray-600/5 dark:hover:bg-gray-300/5 rounded-full py-2 px-4 transition-colors duration-300 flex justify-between items-center mb-2">
          <div class="flex items-center gap-2">
            <Icon class="text-gray-600 dark:text-gray-300"><Molecule20Regular /></Icon>
            <span class="text-gray-600 dark:text-gray-300">{{ t('privacy.options.thirdPartyService') }}</span>
          </div>
          <n-switch v-model:value="store.optionValues.thirdPartyService" />
        </div>

        <div v-if="store.optionDisplay.cloudAI" class="hover:bg-gray-600/5 dark:hover:bg-gray-300/5 rounded-full py-2 px-4 transition-colors duration-300 flex justify-between items-center mb-2">
          <div class="flex items-center gap-2">
            <Icon class="text-gray-600 dark:text-gray-300"><Thinking20Regular /></Icon>
            <span class="text-gray-600 dark:text-gray-300">{{ t('privacy.options.cloudAI') }}</span>
          </div>
          <n-switch v-model:value="store.optionValues.cloudAI" />
        </div>
      </div>
    </section>

    <!-- ==================== 第四部分：应用行为 ==================== -->
    <section class="mb-8">
      <h2 v-if="store.optionDisplay.sectionBehavior" class="text-lg font-medium text-gray-800 dark:text-gray-100 mb-3">{{ t('privacy.sections.behavior') }}</h2>
      <div class="text-base">
        <div v-if="store.optionDisplay.notification" class="hover:bg-gray-600/5 dark:hover:bg-gray-300/5 rounded-full py-2 px-4 transition-colors duration-300 flex justify-between items-center mb-2">
          <div class="flex items-center gap-2">
            <Icon class="text-gray-600 dark:text-gray-300"><Alert16Regular /></Icon>
            <span class="text-gray-600 dark:text-gray-300">{{ t('privacy.options.notification') }}</span>
          </div>
          <n-switch v-model:value="store.optionValues.notification" />
        </div>

        <div v-if="store.optionDisplay.update" class="hover:bg-gray-600/5 dark:hover:bg-gray-300/5 rounded-full py-2 px-4 transition-colors duration-300 flex justify-between items-center mb-2">
          <div class="flex items-center gap-2">
            <Icon class="text-gray-600 dark:text-gray-300"><ArrowSyncCircle16Regular /></Icon>
            <span class="text-gray-600 dark:text-gray-300">{{ t('privacy.options.update') }}</span>
          </div>
          <n-switch v-model:value="store.optionValues.update" />
        </div>

        <div v-if="store.optionDisplay.autoStart" class="hover:bg-gray-600/5 dark:hover:bg-gray-300/5 rounded-full py-2 px-4 transition-colors duration-300 flex justify-between items-center mb-2">
          <div class="flex items-center gap-2">
            <Icon class="text-gray-600 dark:text-gray-300"><Gauge20Regular /></Icon>
            <span class="text-gray-600 dark:text-gray-300">{{ t('privacy.options.autoStart') }}</span>
          </div>
          <n-switch v-model:value="store.optionValues.autoStart" />
        </div>

        <div v-if="store.optionDisplay.contextMenu" class="hover:bg-gray-600/5 dark:hover:bg-gray-300/5 rounded-full py-2 px-4 transition-colors duration-300 flex justify-between items-center mb-2">
          <div class="flex items-center gap-2">
            <Icon class="text-gray-600 dark:text-gray-300"><CursorClick20Regular /></Icon>
            <span class="text-gray-600 dark:text-gray-300">{{ t('privacy.options.contextMenu') }}</span>
          </div>
          <n-switch v-model:value="store.optionValues.contextMenu" />
        </div>

        <div v-if="store.optionDisplay.trayIcon" class="hover:bg-gray-600/5 dark:hover:bg-gray-300/5 rounded-full py-2 px-4 transition-colors duration-300 flex justify-between items-center mb-2">
          <div class="flex items-center gap-2">
            <Icon class="text-gray-600 dark:text-gray-300"><TrayItemAdd20Regular /></Icon>
            <span class="text-gray-600 dark:text-gray-300">{{ t('privacy.options.trayIcon') }}</span>
          </div>
          <n-switch v-model:value="store.optionValues.trayIcon" />
        </div>

        <div v-if="store.optionDisplay.shortcut" class="hover:bg-gray-600/5 dark:hover:bg-gray-300/5 rounded-full py-2 px-4 transition-colors duration-300 flex justify-between items-center mb-2">
          <div class="flex items-center gap-2">
            <Icon class="text-gray-600 dark:text-gray-300"><Apps20Regular /></Icon>
            <span class="text-gray-600 dark:text-gray-300">{{ t('privacy.options.shortcut') }}</span>
          </div>
          <n-switch v-model:value="store.optionValues.shortcut" />
        </div>

        <div v-if="store.optionDisplay.openOtherApp" class="hover:bg-gray-600/5 dark:hover:bg-gray-300/5 rounded-full py-2 px-4 transition-colors duration-300 flex justify-between items-center mb-2">
          <div class="flex items-center gap-2">
            <Icon class="text-gray-600 dark:text-gray-300"><Open20Regular /></Icon>
            <span class="text-gray-600 dark:text-gray-300">{{ t('privacy.options.openOtherApp') }}</span>
          </div>
          <n-switch v-model:value="store.optionValues.openOtherApp" />
        </div>
      </div>
    </section>

    <!-- ==================== 第五部分：隐私 ==================== -->
    <section class="mb-8">
      <h2 v-if="store.optionDisplay.sectionPrivacy" class="text-lg font-medium text-gray-800 dark:text-gray-100 mb-3">{{ t('privacy.sections.privacy') }}</h2>
      <p class="text-gray-500 dark:text-gray-400 mb-4">{{ t('privacy.adjustHint') }}</p>
      <div class="text-base">
        <div v-if="store.optionDisplay.personalization" class="hover:bg-gray-600/5 dark:hover:bg-gray-300/5 rounded-full py-2 px-4 transition-colors duration-300 flex justify-between items-center mb-2">
          <div class="flex items-center gap-2">
            <Icon class="text-gray-600 dark:text-gray-300"><DataWaterfall20Regular /></Icon>
            <span class="text-gray-600 dark:text-gray-300">{{ t('privacy.options.personalization') }}</span>
          </div>
          <n-switch v-model:value="store.optionValues.personalization" />
        </div>

        <div v-if="store.optionDisplay.personalInfo" class="hover:bg-gray-600/5 dark:hover:bg-gray-300/5 rounded-full py-2 px-4 transition-colors duration-300 flex justify-between items-center mb-2">
          <div class="flex items-center gap-2">
            <Icon class="text-gray-600 dark:text-gray-300"><Person16Regular /></Icon>
            <span class="text-gray-600 dark:text-gray-300">{{ t('privacy.options.personalInfo') }}</span>
          </div>
          <n-switch v-model:value="store.optionValues.personalInfo" />
        </div>

        <div v-if="store.optionDisplay.plaintextStorage" class="hover:bg-gray-600/5 dark:hover:bg-gray-300/5 rounded-full py-2 px-4 transition-colors duration-300 flex justify-between items-center mb-2">
          <div class="flex items-center gap-2">
            <Icon class="text-gray-600 dark:text-gray-300"><LockOpen16Regular /></Icon>
            <span class="text-gray-600 dark:text-gray-300">{{ t('privacy.options.plaintextStorage') }}</span>
          </div>
          <n-switch v-model:value="store.optionValues.plaintextStorage" />
        </div>

        <div v-if="store.optionDisplay.analytics" class="hover:bg-gray-600/5 dark:hover:bg-gray-300/5 rounded-full py-2 px-4 transition-colors duration-300 flex justify-between items-center mb-2">
          <div class="flex items-center gap-2">
            <Icon class="text-gray-600 dark:text-gray-300"><Stethoscope20Regular /></Icon>
            <span class="text-gray-600 dark:text-gray-300">{{ t('privacy.options.analytics') }}</span>
          </div>
          <n-switch v-model:value="store.optionValues.analytics" />
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
import { useI18n } from 'vue-i18n'
import { NSwitch } from 'naive-ui'
import { Icon } from '@vicons/utils'
import Apps20Regular from '@vicons/fluent/es/Apps20Regular'
import Alert16Regular from '@vicons/fluent/es/Alert16Regular'
import Folder16Regular from '@vicons/fluent/es/Folder16Regular'
import Storage24Regular from '@vicons/fluent/es/Storage24Regular'
import Stethoscope20Regular from '@vicons/fluent/es/Stethoscope20Regular'
import ArrowSyncCircle16Regular from '@vicons/fluent/es/ArrowSyncCircle16Regular'
import Globe16Regular from '@vicons/fluent/es/Globe16Regular'
import Image20Regular from '@vicons/fluent/es/Image20Regular'
import Gauge20Regular from '@vicons/fluent/es/Gauge20Regular'
import Mic20Regular from '@vicons/fluent/es/Mic20Regular'
import Open20Regular from '@vicons/fluent/es/Open20Regular'
import ShieldKeyhole20Regular from '@vicons/fluent/es/ShieldKeyhole20Regular'
import UsbStick20Regular from '@vicons/fluent/es/UsbStick20Regular'
import Camera16Regular from '@vicons/fluent/es/Camera16Regular'
import Clipboard16Regular from '@vicons/fluent/es/Clipboard16Regular'
import Location16Regular from '@vicons/fluent/es/Location16Regular'
import Screenshot20Regular from '@vicons/fluent/es/Screenshot20Regular'
import TrayItemAdd20Regular from '@vicons/fluent/es/TrayItemAdd20Regular'
import Live20Regular from '@vicons/fluent/es/Live20Regular'
import Thinking20Regular from '@vicons/fluent/es/Thinking20Regular'
import Molecule20Regular from '@vicons/fluent/es/Molecule20Regular'
import DataWaterfall20Regular from '@vicons/fluent/es/DataWaterfall20Regular'
import Fingerprint24Regular from '@vicons/fluent/es/Fingerprint24Regular'
import Settings16Regular from '@vicons/fluent/es/Settings16Regular'
import Database20Regular from '@vicons/fluent/es/Database20Regular'
import Person16Regular from '@vicons/fluent/es/Person16Regular'
import LockOpen16Regular from '@vicons/fluent/es/LockOpen16Regular'
import Keyboard20Regular from '@vicons/fluent/es/Keyboard20Regular'
import WindowApps20Regular from '@vicons/fluent/es/WindowApps20Regular'
import CursorClick20Regular from '@vicons/fluent/es/CursorClick20Regular'
import { usePrivacyConfigStore } from '@/stores/privacyConfigStore'

const { t } = useI18n()

defineProps({
  appName: {
    type: String,
    default: 'App',
  },
})

const store = usePrivacyConfigStore()
</script>
