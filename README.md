# Foil

English | [简体中文](README_zh.md)

A desktop app for building custom Android APK files from any HTML project. Built with Go + Wails2.

**No Java / Android SDK required** — everything needed is bundled: a minimal JRE, Apktool, template APK, Apksigner, and signing keys generator.

---

## Features

- **Build APK from HTML** — point to any local folder with an `index.html`, or upload a `.zip` / `.html` file
- **Custom metadata** — app name, package name, version — all configurable
- **App icon** — upload an image or auto-generate a letter icon with correct DPI scaling
- **APK signing**
  - Auto‑generated self‑signed certificate (Go‑based, private key encrypted with Windows DPAPI)
  - Or bring your own keystore (signed via bundled Android Apksigner)
- **Certificate management** — remember certificate path/passwords across sessions (encrypted with Windows DPAPI)
- **Floating build log** — hover to inspect build output without leaving the page
- **Bilingual UI** — Chinese and English interfaces
- **Single‑file distribution** — all assets embedded in the binary, extracted on first run

---

## Quick Start

### Development

```bash
# Install Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Clone and enter the project
git clone https://github.com/KaiZhou554/foil.git
cd foil

# Install frontend dependencies
cd frontend && npm install && cd ..

# Run in dev mode
wails dev
```

### Production Build

```bash
wails build
```

The output binary and installer will be in the `build/bin` directory.

---

## Tech Stack

| Layer | Technology |
|---|---|
| Desktop framework | [Wails v2](https://wails.io) (Go + WebView2) |
| Frontend | Vue 3 + TypeScript + Vite + Naive UI + Pinia |
| APK tooling | Apktool + Android Apksigner (both bundled) |
| Certificate encryption | Windows DPAPI (`CryptProtectData` / `CryptUnprotectData`) |
| Java runtime | Bundled minimal JRE (via jlink, ~43 MB) |
| Desktop path (Windows) | Read from `HKCU\…\User Shell Folders\Desktop` |

---

## Project Structure

```
foil/
├── app.go                  # Wails-bound API methods
├── main.go                 # Entry point, asset embedding
├── certstorage.go          # DPAPI-encrypted certificate storage
├── keytool.go              # List keystore aliases via bundled keytool
├── desktop_windows.go      # Registry-based desktop path
├── config/                 # TOML config management
├── internal/
│   ├── builder/            # APK build pipeline
│   │   ├── builder.go      # Orchestrator (Go sign & Apksigner)
│   │   ├── apktool.go      # Apktool integration
│   │   ├── naming.go       # Package / version generators
│   │   ├── icons.go        # Icon definitions
│   │   ├── genkey.go       # Auto‑generated key pair + DPAPI encryption
│   │   └── zipalign.go     # 4‑byte alignment
│   ├── dpapi/              # Windows DPAPI bindings
│   └── apksigner/          # Go APK v1 + v2 signing library (for auto certs)
├── frontend/
│   ├── src/
│   │   ├── pages/          # HomePage, AdvancedPage, SettingsPage, WelcomePage
│   │   ├── components/     # UI components (Sidebar, BuildButton, SetupCard…)
│   │   ├── stores/         # Pinia stores
│   │   └── locales/        # i18n (zh-CN, en)
│   └── wailsjs/            # Auto-generated Wails bindings
└── assets/                 # Bundled resources
    ├── foil-example.apk    # Template APK
    ├── apksigner.jar       # Android Apksigner
    ├── apktool.jar         # Apktool
    └── jre-minimal/        # Minimal JRE
```

---

## Configuration

Settings are persisted to `%APPDATA%\unieditdept\foil\config.toml` (Windows).

| Key | Default | Description |
|---|---|---|
| `language` | `zh-CN` | Display language (`zh-CN` / `en`) |
| `outputDir` | Desktop | APK output directory |
| `showFloatButton` | `false` | Show floating build log button |
| `openAfterBuild` | `true` | Open output folder in Explorer after build |
| `useCustomCert` | `false` | Use custom keystore instead of auto‑generated |
| `rememberLevel` | `off` | Certificate remember level (`off` / `path` / `full`) |
| `rememberCompany` | `false` | Remember company (package second segment) |
| `companyName` | `""` | Stored company name |

---

## License

MIT
