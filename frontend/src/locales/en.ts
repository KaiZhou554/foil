export default {
  app: {
    versionPrefix: 'Version ',
    defaultAppName: 'Foil',
  },

  welcome: {
    title: 'Welcome to Foil',
    slogan: 'Generate APKs in a few steps.',
    selectLanguage: 'Language',
    getStarted: 'Continue',
  },

  sidebar: {
    home: 'Home',
    advanced: 'Advanced',
    settings: 'Settings',
  },

  settings: {
    title: 'Settings',
    language: 'Language',
    displayLanguage: 'Display Language',

    general: 'General',

    backToSettings: 'Settings',

    saveLocation: 'Save Location',
    locationDesktop: 'Desktop',
    locationCustom: 'Custom',
    locationPlaceholder: 'Select folder',
    btnBrowse: 'Browse…',
    showFloatButton: 'Show Build Log',
    openAfterBuild: 'Open Folder After Build',
    rememberCompany: 'Remember Company Name',
    about: 'About',
    aboutDesc: 'A desktop app for building custom Android APK files from any HTML project.',
  },

  titlebar: {
    minimize: 'Minimize',
    maximize: 'Maximize',
    close: 'Close',
  },

  buildPage: {
    header: 'Build APK',

    sourceCard: 'Project',

    statusNotSelected: 'Not selected',
    statusFolder: 'Folder',
    statusZip: 'ZIP',
    statusHtml: 'HTML',

    tabFolder: 'Folder',
    tabFile: 'ZIP or HTML',

    placeholderFolder: 'Select folder',
    placeholderFile: 'Select ZIP or HTML',

    btnBrowse: 'Browse…',
    btnSelectFile: 'Select…',

    appNameCard: 'Name',
    placeholderAppName: 'Enter app name',

    iconCard: 'Icon',
    iconHint: 'Auto-generated if not set',
    btnSelectImage: 'Select…',
    btnClear: 'Remove',

    advanced: 'Advanced',
    pkgPreview: 'Package Name',
    versionLabel: 'Version',
    versionPlaceholder: 'e.g. 1.0.0',

    btnBuild: 'Build APK',
    btnBuilding: 'Building…',

    successTitle: 'APK Built',
    failTitle: 'Build Failed',
    errorNoIndex: 'No index.html found in the archive',
    errorGeneric: 'Please check your input and try again',

    logCustomPkg: 'Using custom package: ',
    logBuilding: 'Building APK...\n',
    floatPlaceholder: 'No build yet',
  },

  advancedPage: {
    header: 'Advanced Build',
    certificate: 'Certificate',
    certAuto: 'Use auto-generated certificate',
    certCustom: 'Use my own certificate',
    certPathPlaceholder: 'Select certificate file…',
    certPasswordPlaceholder: 'Certificate password (optional)',
    rememberCert: 'Remember',
    remember: 'Remember',
    rememberOff: 'Off',
    rememberPath: 'Path only',
    rememberFull: 'Path & Password',
    rememberTooltip: 'Off: clear certificate info on close; Path only: remember only the file path, passwords not saved; Path & Password: remember path, passwords and alias, auto-fill next time.<br/><br/>Passwords are encrypted with Windows DPAPI, decryptable only by the current user.',
    aliasPlaceholder: 'Select or enter Alias…',
    keyPasswordPlaceholder: 'Key password (optional)',
    keySameAsStore: 'Key password same as Keystore',
    btnBrowse: 'Browse…',
    errorAliasRequired: 'Please enter a certificate alias',
  },
}
