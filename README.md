# Let's config

> 基于 Wails 2.12 + Vue 3 的桌面配置管理工具。

## 技术栈

| 层 | 技术 |
|---|---|
| 桌面框架 | [Wails](https://wails.io/) 2.12 |
| 前端 | Vue 3 + TypeScript + Vite |
| 状态管理 | Pinia 3 |
| 路由 | Vue Router 5 (hash 模式) |
| UI 组件 | Naive UI |
| 样式 | Tailwind CSS 4 |
| 国际化 | vue-i18n |
| 配置存储 | [BurntSushi/toml](https://github.com/BurntSushi/toml) |

## 项目结构

```
lets_config/
├── main.go                  # Wails 入口
├── app.go                   # App 结构体 & 前端绑定方法
├── config/
│   ├── config.go            # 配置结构体 (TOML + JSON 标签)
│   └── manager.go           # 线程安全的配置读写管理器
├── frontend/
│   └── src/
│       ├── main.ts          # Vue 应用入口
│       ├── App.vue          # 根组件 (TitleBar、路由视图、窗口状态同步)
│       ├── router/          # 路由配置 & 首次启动守卫
│       ├── stores/          # Pinia stores
│       │   ├── appStore.ts          # 应用状态 (语言、首次启动、通用设置)
│       │   └── privacyConfigStore.ts # 隐私权限选项
│       ├── composables/
│       │   └── useWindowState.ts     # 窗口状态轮询同步
│       ├── pages/           # 页面组件
│       │   ├── WelcomePage.vue       # 欢迎页 (语言选择)
│       │   ├── SetupPage.vue         # 引导设置页
│       │   ├── MainLayout.vue        # 主布局 (侧边栏 + 内容区)
│       │   ├── HomePage.vue          # 主页
│       │   ├── SettingsPage.vue      # 设置页
│       │   └── PrivacyPage.vue       # 隐私权限详情页
│       ├── components/      # 可复用组件
│       ├── locales/         # i18n 语言文件 (zh-CN / en)
│       └── assets/          # 图标、字体、CSS
└── wails.json               # Wails 项目配置
```

## 配置系统

配置文件存储在 `%APPDATA%/unieditdept/lets_config/config.toml`，首次启动自动创建。

全链路统一使用 camelCase 命名：

```
前端 Store → Wails Models → JSON Bridge → Go Struct Tags → TOML 文件
```

全部配置字段定义在 [`config/config.go`](config/config.go) 中，扩展新字段只需在结构体中添加一行。

### 读写流程

```
加载：Go 读取 TOML → Wails GetConfig() → Store.loadConfig() → 所有页面
保存：UI 修改 → Store.saveConfig() → Wails SaveConfig() → Go 写入 TOML
```

设置页面**实时保存**（每次开关切换立即写入），隐私设置 400ms 防抖保存。

## 窗口管理

- **无边框窗口** + 亚克力背景（Windows）
- 自定义 `TitleBar` 支持最小化 / 最大化 / 关闭，可拖拽
- 窗口状态通过 [`useWindowState`](frontend/src/composables/useWindowState.ts) composable **轮询同步**，兼容：
  - 双击标题栏最大化/恢复
  - Win + ↑ / Win + ↓
  - Snap Layout 窗口贴靠
  - 拖动最大化窗口恢复

## 首次启动引导

```
首次启动 → /welcome (选择语言) → /setup (隐私权限配置) → /main/home
再次启动 → 直接进入 /main/home
```

`firstLaunch` 标记持久化在配置文件中，路由守卫根据该标记控制访问。

## 开发

```bash
# 安装前端依赖
cd frontend && npm install

# 开发模式 (热重载)
wails dev

# 构建
wails build
```

## 许可

MIT License — 详见 [LICENSE](LICENSE)。
