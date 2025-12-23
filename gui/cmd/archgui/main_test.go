package main

import (
	"strings"
	"testing"
)

func TestGenerateConfig(t *testing.T) {
	disk := "/dev/sda"
	host := "myarch"
	user := "alice"
	rootPass := "secretroot"
	userPass := "secretuser"
	fs := "btrfs"
	luksPass := "cryptpass"
	desktop := "kde"
	shell := "zsh"
	encrypt := true
	nvidia := true

	config := generateConfig(disk, host, user, rootPass, userPass, fs, luksPass, desktop, shell, encrypt, nvidia)

	checks := map[string]string{
		"DISK":           disk,
		"HOSTNAME":       host,
		"USERNAME":       user,
		"ROOT_PASSWORD":  rootPass,
		"USER_PASSWORD":  userPass,
		"FS_TYPE":        fs,
		"USE_LUKS":       "yes",
		"LUKS_PASSWORD":  luksPass,
		"DESKTOP_ENV":    desktop,
		"SHELL_CHOICE":   shell,
		"HAS_NVIDIA":     "yes",
		"NONINTERACTIVE": "yes",
	}

	for key, expected := range checks {
		expectedLine := key + "=" + expected
		if !strings.Contains(config, expectedLine) {
			t.Errorf("Config missing or incorrect for %s. Expected to find '%s' in:\n%s", key, expectedLine, config)
		}
	}
}
