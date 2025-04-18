//go:build linux

package config

import (
	"os/user"
	"path/filepath"
)

func GetConfigPath() (string, error) {
	if user, err := user.Current(); err != nil {
		return "", err
	} else {
		return filepath.Join(user.HomeDir, ".config/thumbnailinago/config.yaml"), nil
	}
}

func GetInkscapePath() string {
	return "inkscape"
}
