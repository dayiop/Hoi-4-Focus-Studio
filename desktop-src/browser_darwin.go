//go:build darwin

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
		return fmt.Errorf("не удалось создать папку сохранений: %w", err)
	}
	for _, name := range []string{"Google Chrome", "Microsoft Edge", "Chromium", "Brave Browser"} {
		if err := exec.Command("/usr/bin/open", "-Ra", name).Run(); err != nil {
			continue
		}
		cmd := exec.Command("/usr/bin/open", "-na", name, "--args", "--app="+origin, "--user-data-dir="+profile, "--no-first-run", "--no-default-browser-check")
		if err := cmd.Run(); err == nil {
			return nil
		}
	}
	// Safari / default browser: the editor opens in a normal tab rather than app mode.
	if err := exec.Command("/usr/bin/open", origin).Start(); err != nil {
		return fmt.Errorf("не удалось открыть браузер: %w", err)
	}
	return nil
}
