English | [简体中文](README_zh.md)

# Foil

A desktop app for building custom Android APK files from any HTML project. Built with Go + Wails2.

**No Java / Android SDK required** — everything needed is bundled: a minimal JRE, Apktool, template APK, and signing keys generator.

---

## Features

- **Build APK from HTML** — point to any local folder with an `index.html`, or upload a `.zip` / `.html` file
- **Custom metadata** — app name, package name, version — all configurable
- **App icon** — upload an image or auto-generate a letter icon
- **APK signing** — v1 + v2 signing scheme, no external tools needed
- **Single‑file distribution** — all assets embedded in the binary, extracted on first run
- **i18n** — Chinese and English interfaces

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
| APK tooling | Apktool 3.0.2 (bundled) + custom Go signing |
| Java runtime | Bundled minimal JRE (via jlink, ~43 MB) |
| Desktop path (Windows) | Read from `HKCU\…\User Shell Folders\Desktop` |

---

## Project Structure

```
foil/
├── app.go                  # Wails-bound API methods
├── main.go                 # Entry point, asset embedding
├── desktop_windows.go      # Registry-based desktop path
├── config/                 # TOML config management
├── internal/
│   ├── builder/            # APK build pipeline
│   │   ├── builder.go      # Orchestrator
│   │   ├── apktool.go      # Apktool integration
│   │   ├── naming.go       # Package / version generators
│   │   ├── icons.go        # Icon definitions
│   │   └── zipalign.go     # 4-byte alignment
│   ├── axml/               # Binary AXML parser
│   └── apksigner/          # APK v1 + v2 signing
├── frontend/
│   ├── src/
│   │   ├── pages/          # HomePage, SettingsPage, WelcomePage
│   │   ├── components/     # UI components
│   │   ├── stores/         # Pinia stores
│   │   └── locales/        # i18n (zh-CN, en)
│   └── wailsjs/            # Auto-generated Wails bindings
└── assets/                 # Bundled resources
    ├── foil-example.apk    # Template APK
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
| `showFloatButton` | `false` | Show build log float button |
| `openAfterBuild` | `true` | Open output folder in Explorer after build |

---

## License

MIT
