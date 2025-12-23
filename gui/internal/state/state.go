package state

// InstallConfig holds the configuration for the installation
type InstallConfig struct {
	// Storage
	Disk               string
	ManualPartitioning bool
	TargetRoot         string // For manual
	TargetEFI          string // For manual (UEFI)
	FormatRoot         bool
	FormatEFI          bool

	// Encryption
	Encrypt      bool
	LuksPassword string

	// Filesystem
	Filesystem string // ext4, btrfs

	// Account
	Hostname     string
	FullName     string
	Username     string
	RootPassword string
	UserPassword string

	// Localization
	Timezone string
	Locale   string
	Keymap   string

	// Desktop
	Desktop       string // xfce, gnome, etc.
	Graphics      string // nvidia, etc. (bool in backend, but keeping flexible)
	InstallNvidia bool

	// Shell
	Shell string
}

func NewInstallConfig() *InstallConfig {
	return &InstallConfig{
		Hostname:   "archlinux",
		Username:   "user",
		Filesystem: "ext4",
		Desktop:    "xfce",
		Shell:      "bash",
		Timezone:   "UTC",
		Locale:     "en_US",
		Keymap:     "us",
		FormatRoot: true, // Default to format even in manual unless unchecked
	}
}
