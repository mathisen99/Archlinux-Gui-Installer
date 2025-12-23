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

	config := generateConfig(disk, host, user, rootPass, userPass)

	checks := map[string]string{
		"DISK":           disk,
		"HOSTNAME":       host,
		"USERNAME":       user,
		"ROOT_PASSWORD":  rootPass,
		"USER_PASSWORD":  userPass,
		"NONINTERACTIVE": "yes",
	}

	for key, expected := range checks {
		expectedLine := key + "=" + expected
		if !strings.Contains(config, expectedLine) {
			t.Errorf("Config missing or incorrect for %s. Expected to find '%s' in:\n%s", key, expectedLine, config)
		}
	}
}
