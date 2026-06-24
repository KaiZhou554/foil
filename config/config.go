package config

import "fmt"

// Config is the root configuration structure.
// All tag names use camelCase to match the frontend store keys directly —
// no conversion needed anywhere in the chain.
type Config struct {
	Version     string        `toml:"version"      json:"version"`
	FirstLaunch bool          `toml:"firstLaunch"  json:"firstLaunch"`
	Language    string        `toml:"language"     json:"language"`
	General     GeneralConfig `toml:"general"      json:"general"`
	Privacy     PrivacyConfig `toml:"privacy"      json:"privacy"`
}

type GeneralConfig struct {
	AutoStart      bool `toml:"autoStart"       json:"autoStart"`
	Notifications  bool `toml:"notifications"   json:"notifications"`
	MinimizeToTray bool `toml:"minimizeToTray"  json:"minimizeToTray"`
}

// ── Privacy / permissions ──

type PrivacyConfig struct {
	Options PrivacyFields `toml:"options" json:"options"`
	Display PrivacyFields `toml:"display" json:"display"`
	Meta    PrivacyMeta   `toml:"meta"    json:"meta"`
}

// PrivacyFields holds the boolean toggles shared by both option values and option visibility.
// Every key here matches a key in the frontend privacyConfigStore exactly.
type PrivacyFields struct {
	// Section display toggles
	SectionStorage  bool `toml:"sectionStorage"  json:"sectionStorage"`
	SectionDevice   bool `toml:"sectionDevice"   json:"sectionDevice"`
	SectionOnline   bool `toml:"sectionOnline"   json:"sectionOnline"`
	SectionBehavior bool `toml:"sectionBehavior" json:"sectionBehavior"`
	SectionPrivacy  bool `toml:"sectionPrivacy"  json:"sectionPrivacy"`

	// Storage
	StorageDefault bool `toml:"storageDefault" json:"storageDefault"`
	StandardFolder bool `toml:"standardFolder" json:"standardFolder"`
	StorageCustom  bool `toml:"storageCustom"  json:"storageCustom"`

	// Device & System
	AdminPrivilege       bool `toml:"adminPrivilege"        json:"adminPrivilege"`
	IdentityVerification bool `toml:"identityVerification"  json:"identityVerification"`
	ModifySystemSettings bool `toml:"modifySystemSettings"  json:"modifySystemSettings"`
	InputMonitor         bool `toml:"inputMonitor"          json:"inputMonitor"`
	Microphone           bool `toml:"microphone"            json:"microphone"`
	Camera               bool `toml:"camera"                json:"camera"`
	Location             bool `toml:"location"              json:"location"`
	Sensor               bool `toml:"sensor"                json:"sensor"`
	Clipboard            bool `toml:"clipboard"             json:"clipboard"`
	Screenshot           bool `toml:"screenshot"            json:"screenshot"`
	OtherDevice          bool `toml:"otherDevice"           json:"otherDevice"`
	AccessOtherAppData   bool `toml:"accessOtherAppData"    json:"accessOtherAppData"`
	ManageOtherApps      bool `toml:"manageOtherApps"       json:"manageOtherApps"`

	// Online features
	Network           bool `toml:"network"           json:"network"`
	ThirdPartyService bool `toml:"thirdPartyService" json:"thirdPartyService"`
	CloudAI           bool `toml:"cloudAI"           json:"cloudAI"`

	// App behavior
	Personalization bool `toml:"personalization" json:"personalization"`
	Notification    bool `toml:"notification"    json:"notification"`
	Update          bool `toml:"update"          json:"update"`
	AutoStart       bool `toml:"autoStart"       json:"autoStart"`
	ContextMenu     bool `toml:"contextMenu"     json:"contextMenu"`
	TrayIcon        bool `toml:"trayIcon"        json:"trayIcon"`
	Shortcut        bool `toml:"shortcut"        json:"shortcut"`
	OpenOtherApp    bool `toml:"openOtherApp"    json:"openOtherApp"`

	// Privacy & analytics
	PersonalInfo     bool `toml:"personalInfo"     json:"personalInfo"`
	PlaintextStorage bool `toml:"plaintextStorage" json:"plaintextStorage"`
	Analytics        bool `toml:"analytics"        json:"analytics"`
}

type PrivacyMeta struct {
	DisplayVersion int    `toml:"displayVersion" json:"displayVersion"`
	UpdateTime     string `toml:"updateTime"     json:"updateTime"`
	Author         string `toml:"author"         json:"author"`
	Desc           string `toml:"desc"           json:"desc"`
}

// ── Defaults ──

// Validate checks that the loaded configuration has sensible values.
// Returns nil if the config is valid, or an error describing the first problem found.
func (c *Config) Validate() error {
	if c.Version == "" {
		return fmt.Errorf("version must not be empty")
	}
	if c.Language == "" {
		return fmt.Errorf("language must not be empty")
	}
	if c.Privacy.Meta.DisplayVersion < 0 {
		return fmt.Errorf("privacy.meta.displayVersion must not be negative")
	}
	return nil
}

func DefaultConfig() *Config {
	return &Config{
		Version:     "0.1.0",
		FirstLaunch: true,
		Language:    "zh-CN",
		General: GeneralConfig{
			AutoStart:      false,
			Notifications:  true,
			MinimizeToTray: true,
		},
		Privacy: DefaultPrivacyConfig(),
	}
}

func DefaultPrivacyConfig() PrivacyConfig {
	return PrivacyConfig{
		Options: defaultPrivacyOptions(),
		Display: defaultPrivacyDisplay(),
		Meta: PrivacyMeta{
			DisplayVersion: 6,
			UpdateTime:     "2025-01-01",
			Author:         "system",
			Desc:           "首次启动隐私权限配置",
		},
	}
}

func defaultPrivacyOptions() PrivacyFields {
	return PrivacyFields{
		StorageDefault: true,
		StandardFolder: false,
		StorageCustom:  false,

		AdminPrivilege:       false,
		IdentityVerification: false,
		ModifySystemSettings: false,
		InputMonitor:         false,
		Microphone:           false,
		Camera:               false,
		Location:             false,
		Sensor:               false,
		Clipboard:            false,
		Screenshot:           false,
		OtherDevice:          false,
		AccessOtherAppData:   false,
		ManageOtherApps:      false,

		Network:           false,
		ThirdPartyService: false,
		CloudAI:           false,

		Personalization: true,
		Notification:    true,
		Update:          false,
		AutoStart:       false,
		ContextMenu:     true,
		TrayIcon:        true,
		Shortcut:        true,
		OpenOtherApp:    false,

		PersonalInfo:     false,
		PlaintextStorage: false,
		Analytics:        false,
	}
}

func defaultPrivacyDisplay() PrivacyFields {
	return PrivacyFields{
		SectionStorage:  true,
		SectionDevice:   true,
		SectionOnline:   true,
		SectionBehavior: true,
		SectionPrivacy:  true,

		StorageDefault:       true,
		StandardFolder:       true,
		StorageCustom:        true,
		AdminPrivilege:       true,
		IdentityVerification: true,
		ModifySystemSettings: true,
		InputMonitor:         true,
		Microphone:           true,
		Camera:               true,
		Location:             true,
		Sensor:               true,
		Clipboard:            true,
		Screenshot:           true,
		OtherDevice:          true,
		AccessOtherAppData:   true,
		ManageOtherApps:      true,
		Network:              true,
		ThirdPartyService:    true,
		CloudAI:              true,
		Personalization:      true,
		Notification:         true,
		Update:               true,
		AutoStart:            true,
		ContextMenu:          true,
		TrayIcon:             true,
		Shortcut:             true,
		OpenOtherApp:         true,
		PersonalInfo:         true,
		PlaintextStorage:     true,
		Analytics:            true,
	}
}
