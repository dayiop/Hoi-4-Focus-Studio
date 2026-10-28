//go:build linux

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func launchBrowser() error {
	profile := filepath.Join(baseDataDir(), "BrowserProfile")
	if err := os.MkdirAll(profile, 0700); err != nil {
		return fmt.Errorf("не удалось создать профиль браузера: %w", err)
	}
	for _, b := range []string{"google-chrome-stable", "google-chrome", "chromium", "chromium-browser", "microsoft-edge-stable", "microsoft-edge", "brave-browser", "vivaldi"} {
		path, err := exec.LookPath(b)
		if err != nil {
			continue
		}
		cmd := exec.Command(path, "--app="+origin, "--user-data-dir="+profile, "--no-first-run", "--no-default-browser-check")
		if err = cmd.Start(); err == nil {
			go func() { _ = cmd.Wait() }()
			return nil
		}
	}
	for _, b := range []string{"xdg-open", "gio"} {
		path, err := exec.LookPath(b)
		if err != nil {
			continue
		}
		args := []string{origin}
		if b == "gio" {
			args = append([]string{"open"}, args...)
		}
		cmd := exec.Command(path, args...)
		if err = cmd.Start(); err == nil {
			go func() { _ = cmd.Wait() }()
			return nil
		}
	}
	return fmt.Errorf("не удалось открыть редактор. Установите Chrome, Chromium или браузер с поддержкой xdg-open")
}
