export default {
  app: {
    versionPrefix: '版本',
    defaultAppName: 'Foil',
  },

  welcome: {
    title: '欢迎使用 Foil',
    slogan: '几步即可生成 APK。',
    selectLanguage: '语言',
    getStarted: '继续',
  },

  sidebar: {
    home: '主页',
    advanced: '高级',
    settings: '设置',
  },

  settings: {
    title: '设置',
    language: '语言',
    displayLanguage: '显示语言',

    general: '通用',

    backToSettings: '设置',

    saveLocation: '保存位置',
    locationDesktop: '桌面',
    locationCustom: '自定义',
    locationPlaceholder: '选择文件夹',
    btnBrowse: '浏览…',
    showFloatButton: '显示日志',
    openAfterBuild: '构建后打开文件夹',
    about: '关于',
    aboutDesc: '从任意 HTML 项目构建 Android APK 的桌面应用。',
  },

  titlebar: {
    minimize: '最小化',
    maximize: '最大化',
    close: '关闭',
  },

  buildPage: {
    header: '生成 APK',

    sourceCard: '项目',

    statusNotSelected: '未选择',
    statusFolder: '文件夹',
    statusZip: 'ZIP',
    statusHtml: 'HTML',

    tabFolder: '文件夹',
    tabFile: 'ZIP 或 HTML',

    placeholderFolder: '选择文件夹',
    placeholderFile: '选择 ZIP 或 HTML',

    btnBrowse: '浏览…',
    btnSelectFile: '选择…',

    appNameCard: '名称',
    placeholderAppName: '输入应用名称',

    iconCard: '图标',
    iconHint: '未选择时将自动生成',
    btnSelectImage: '选择…',
    btnClear: '移除',

    advanced: '高级',
    pkgPreview: '包名',
    versionLabel: '版本号',
    versionPlaceholder: '例如 1.0.0',

    btnBuild: '生成 APK',
    btnBuilding: '正在生成…',

    successTitle: 'APK 已生成',
    failTitle: '无法生成',
    errorNoIndex: '压缩包内未找到 index.html',
    errorGeneric: '请检查输入后重试',

    logCustomPkg: '使用自定义包名：',
    logBuilding: '正在生成 APK...\n',
    floatPlaceholder: '暂无构建记录',
  },

  advancedPage: {
    header: '高级构建',
    certificate: '证书',
    certAuto: '使用自动生成的证书',
    certCustom: '使用自己的证书',
    certPathPlaceholder: '选择证书文件…',
    certPasswordPlaceholder: '证书密码（可选）',
    rememberCert: '记住此设置',
    btnBrowse: '浏览…',
  },
}
