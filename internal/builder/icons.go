package builder

// IconDef describes one Android icon file to generate.
type IconDef struct {
	Path string // relative path inside res/, e.g. "mipmap-hdpi-v4/ic_launcher.webp"
	Size int    // size in pixels (square)
}

// AllIcons lists the 15 icon files that need to be generated/replaced.
// These correspond to the standard Android launcher icon set.
var AllIcons = []IconDef{
	// Launcher icons
	{"mipmap-mdpi-v4/ic_launcher.webp", 48},
	{"mipmap-hdpi-v4/ic_launcher.webp", 72},
	{"mipmap-xhdpi-v4/ic_launcher.webp", 96},
	{"mipmap-xxhdpi-v4/ic_launcher.webp", 144},
	{"mipmap-xxxhdpi-v4/ic_launcher.webp", 192},

	// Round launcher icons
	{"mipmap-mdpi-v4/ic_launcher_round.webp", 48},
	{"mipmap-hdpi-v4/ic_launcher_round.webp", 72},
	{"mipmap-xhdpi-v4/ic_launcher_round.webp", 96},
	{"mipmap-xxhdpi-v4/ic_launcher_round.webp", 144},
	{"mipmap-xxxhdpi-v4/ic_launcher_round.webp", 192},

	// Adaptive icon foregrounds
	{"mipmap-mdpi-v4/ic_launcher_foreground.webp", 108},
	{"mipmap-hdpi-v4/ic_launcher_foreground.webp", 162},
	{"mipmap-xhdpi-v4/ic_launcher_foreground.webp", 216},
	{"mipmap-xxhdpi-v4/ic_launcher_foreground.webp", 324},
	{"mipmap-xxxhdpi-v4/ic_launcher_foreground.webp", 432},
}

// IconPaths returns just the path list (for frontend reference).
func IconPaths() []string {
	paths := make([]string, len(AllIcons))
	for i, def := range AllIcons {
		paths[i] = def.Path
	}
	return paths
}
