package main

type ColorScheme struct {
	UsernameFg int
	UsernameBg int

	HostnameFg int
	HostnameBg int

	HomeFg      int
	HomeBg      int
	PathFg      int
	PathBg      int
	CwdFg       int
	SeparatorFg int

	ReadOnlyFg int
	ReadOnlyBg int

	RepoCleanFg int
	RepoCleanBg int
	RepoDirtyFg int
	RepoDirtyBg int

	PromptFg int
	PromptBg int
}

func InitColors() map[string]ColorScheme {
	colors := make(map[string]ColorScheme)
	colors["default"] = ColorScheme{
		UsernameFg:  250,
		UsernameBg:  240,
		HostnameFg:  250,
		HostnameBg:  238,
		HomeFg:      15,
		HomeBg:      31,
		PathFg:      250,
		PathBg:      237,
		CwdFg:       254,
		SeparatorFg: 244,
		ReadOnlyFg:  254,
		ReadOnlyBg:  124,
		RepoCleanFg: 0,
		RepoCleanBg: 148,
		RepoDirtyFg: 15,
		RepoDirtyBg: 161,
		PromptFg:    15,
		PromptBg:    236}
	return colors
}
