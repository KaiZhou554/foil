package builder

// IconDef describes one Android icon file to generate.
type IconDef struct {
	Path string // relative path inside res/, e.g. "mipmap-hdpi/ic_launcher.webp"
	Size int    // size in pixels (square)
}

// AllIcons lists the 15 icon files that need to be generated/replaced.
// These correspond to the standard Android launcher icon set.
// Note: aapt2 strips legacy qualifiers like "-v4", so paths must NOT
// include "-v4" to match the rebuilt APK's resource directory layout.
var AllIcons = []IconDef{
	// Launcher icons
	{"mipmap-mdpi/ic_launcher.webp", 48},
	{"mipmap-hdpi/ic_launcher.webp", 72},
	{"mipmap-xhdpi/ic_launcher.webp", 96},
	{"mipmap-xxhdpi/ic_launcher.webp", 144},
	{"mipmap-xxxhdpi/ic_launcher.webp", 192},

	// Round launcher icons
	{"mipmap-mdpi/ic_launcher_round.webp", 48},
	{"mipmap-hdpi/ic_launcher_round.webp", 72},
	{"mipmap-xhdpi/ic_launcher_round.webp", 96},
	{"mipmap-xxhdpi/ic_launcher_round.webp", 144},
	{"mipmap-xxxhdpi/ic_launcher_round.webp", 192},

	// Adaptive icon foregrounds
	{"mipmap-mdpi/ic_launcher_foreground.webp", 108},
	{"mipmap-hdpi/ic_launcher_foreground.webp", 162},
	{"mipmap-xhdpi/ic_launcher_foreground.webp", 216},
	{"mipmap-xxhdpi/ic_launcher_foreground.webp", 324},
	{"mipmap-xxxhdpi/ic_launcher_foreground.webp", 432},
}

// IconPaths returns just the path list (for frontend reference).
func IconPaths() []string {
	paths := make([]string, len(AllIcons))
	for i, def := range AllIcons {
		paths[i] = def.Path
	}
	return paths
}
