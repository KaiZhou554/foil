package config

import "fmt"

// Config is the root configuration structure.
type Config struct {
	Version            string `toml:"version"            json:"version"`
	FirstLaunch        bool   `toml:"firstLaunch"        json:"firstLaunch"`
	Language           string `toml:"language"           json:"language"`
	OutputDir          string `toml:"outputDir"          json:"outputDir"`
	ShowFloatButton    bool   `toml:"showFloatButton"    json:"showFloatButton"`
	OpenAfterBuild     bool   `toml:"openAfterBuild"     json:"openAfterBuild"`
	UseCustomCert      bool   `toml:"useCustomCert"      json:"useCustomCert"`
	RememberLevel      string `toml:"rememberLevel"      json:"rememberLevel"`
	KeyPassSameAsStore bool   `toml:"keyPassSameAsStore" json:"keyPassSameAsStore"`
	CompanyName        string `toml:"companyName"        json:"companyName"`
	RememberCompany    bool   `toml:"rememberCompany"    json:"rememberCompany"`
}

// Validate checks that the loaded configuration has sensible values.
func (c *Config) Validate() error {
	if c.Version == "" {
		return fmt.Errorf("version must not be empty")
	}
	if c.Language == "" {
		return fmt.Errorf("language must not be empty")
	}
	return nil
}

func DefaultConfig() *Config {
	return &Config{
		Version:            "0.1.0",
		FirstLaunch:        true,
		Language:           "zh-CN",
		OutputDir:          "",
		ShowFloatButton:    false,
		OpenAfterBuild:     true,
		UseCustomCert:      false,
		RememberLevel:      "off",
		KeyPassSameAsStore: true,
		CompanyName:        "",
		RememberCompany:    false,
	}
}
