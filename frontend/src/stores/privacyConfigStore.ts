// src/stores/privacyConfigStore.ts
import { defineStore } from 'pinia'

// 首次启动隐私权限配置
export const usePrivacyConfigStore = defineStore('privacyConfig', {
  state: () => ({
    // ============== 1. 控制【选项显示/隐藏】 ==============
    optionDisplay: {
      // 分组标题
      sectionStorage: true,       // 「存储」标题
      sectionDevice: true,        // 「设备与系统」标题
      sectionOnline: true,        // 「在线功能」标题
      sectionBehavior: true,      // 「应用行为」标题
      sectionPrivacy: true,       // 「隐私」标题

      // 存储位置
      storageDefault: true,       // 使用推荐位置
      standardFolder: true,       // 访问标准文件夹
      storageCustom: true,        // 访问所有文件

      // 系统权限与硬件
      adminPrivilege: true,       // 请求管理员权限
      identityVerification: true, // 用户身份验证
      modifySystemSettings: true, // 修改系统设置
      inputMonitor: true,         // 监听输入事件
      microphone: true,           // 麦克风
      camera: true,               // 摄像头
      location: true,             // 位置信息
      sensor: true,               // 传感器
      clipboard: true,            // 访问剪贴板
      screenshot: true,           // 截图与录屏
      otherDevice: true,          // 访问其他设备
      accessOtherAppData: true,   // 访问其他应用的数据
      manageOtherApps: true,      // 管理其他应用

      // 网络与 AI 服务
      network: true,              // 连接网络服务
      thirdPartyService: true,    // 与第三方在线服务交互
      cloudAI: true,              // 使用云端 AI 处理数据

      // 应用体验
      personalization: true,      // 个性化
      notification: true,         // 发送通知
      update: true,               // 自动检查更新
      autoStart: true,            // 开机时自动启动
      contextMenu: true,          // 启用上下文菜单集成
      trayIcon: true,             // 显示托盘图标
      shortcut: true,             // 桌面快捷方式
      openOtherApp: true,         // 允许打开其他应用

      // 隐私与改进
      personalInfo: true,         // 访问个人信息
      plaintextStorage: true,     // 明文存储数据
      analytics: true,            // 匿名发送使用情况统计
    },

    // ============== 2. 存储【选项开关的值】 ==============
    optionValues: {
      // 存储位置
      storageDefault: true,
      standardFolder: false,
      storageCustom: false,

      // 系统权限与硬件
      adminPrivilege: false,
      identityVerification: false,
      modifySystemSettings: false,
      inputMonitor: false,
      microphone: false,
      camera: false,
      location: false,
      sensor: false,
      clipboard: false,
      screenshot: false,
      otherDevice: false,
      accessOtherAppData: false,
      manageOtherApps: false,

      // 网络与 AI 服务
      network: false,
      thirdPartyService: false,
      cloudAI: false,

      // 应用体验
      personalization: true,
      notification: true,
      update: false,
      autoStart: false,
      contextMenu: true,
      trayIcon: true,
      shortcut: true,
      openOtherApp: false,

      // 隐私与改进
      personalInfo: false,
      plaintextStorage: false,
      analytics: false,
    },

    // ============== 3. 存储【元数据】 ==============
    meta: {
      displayVersion: 6,
      updateTime: '2025-01-01',
      author: 'system',
      desc: '首次启动隐私权限配置'
    }
  }),

  getters: {
    getDisplayVersion: (state) => state.meta.displayVersion,
    getAllOptionValues: (state) => state.optionValues
  },

  actions: {
    updateOptionValue(key: string, value: boolean) {
      if (key in this.optionValues) {
        this.optionValues[key as keyof typeof this.optionValues] = value
      }
    },

    updateOptionDisplay(key: string, isShow: boolean) {
      if (key in this.optionDisplay) {
        this.optionDisplay[key as keyof typeof this.optionDisplay] = isShow
      }
    },

    updateMeta(metaData: { displayVersion?: number; updateTime?: string; author?: string; desc?: string }) {
      this.meta = { ...this.meta, ...metaData }
    },

    resetAll() {
      this.$reset()
    }
  }
})
