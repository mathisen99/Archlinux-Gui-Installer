package pages

import (
	"archgui/gui/internal/state"
	"strings"
	"testing"
)

func TestGenerateConfigEnv(t *testing.T) {
	// Setup a sample config
	config := state.NewInstallConfig()
	config.Disk = "/dev/sda"
	config.Hostname = "myarch"
	config.FullName = "Alice Smith"
	config.Username = "alice"
	config.RootPassword = "secretroot"
	config.UserPassword = "secretuser"
	config.Timezone = "UTC"
	config.Locale = "en_GB"
	config.Keymap = "uk"
	config.Filesystem = "btrfs"
	config.Encrypt = true
	config.LuksPassword = "cryptpass"
	config.Desktop = "kde"
	config.Shell = "zsh"
	config.InstallNvidia = true
	config.ManualPartitioning = false

	envStr := generateConfigEnv(config)

	checks := map[string]string{
		"DISK":           "/dev/sda",
		"HOSTNAME":       "myarch",
		"FULL_NAME":      "Alice Smith",
		"USERNAME":       "alice",
		"ROOT_PASSWORD":  "secretroot",
		"USER_PASSWORD":  "secretuser",
		"TIMEZONE":       "UTC",
		"LOCALE":         "en_GB",
		"KEYMAP":         "uk",
		"FS_TYPE":        "btrfs",
		"USE_LUKS":       "yes",
		"LUKS_PASSWORD":  "cryptpass",
		"DESKTOP_ENV":    "kde",
		"SHELL_CHOICE":   "zsh",
		"HAS_NVIDIA":     "yes",
		"NONINTERACTIVE": "yes",
	}

	for key, expected := range checks {
		expectedLine := key + "=" + expected
		if !strings.Contains(envStr, expectedLine) {
			t.Errorf("Config missing or incorrect for %s. Expected to find '%s' in:\n%s", key, expectedLine, envStr)
		}
	}
}
