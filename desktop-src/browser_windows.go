//go:build windows

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func launchBrowser() error {
	candidates := []string{"msedge.exe"}
	for _, root := range []string{os.Getenv("ProgramFiles(x86)"), os.Getenv("ProgramFiles"), os.Getenv("LOCALAPPDATA")} {
		if root != "" {
			candidates = append(candidates, filepath.Join(root, "Microsoft", "Edge", "Application", "msedge.exe"))
		}
	}
	var browser string
	for _, p := range candidates {
		if q, err := exec.LookPath(p); err == nil {
			browser = q
			break
		}
	}
	if browser == "" {
		return fmt.Errorf("Для отдельного окна нужен Microsoft Edge")
	}
	profile := filepath.Join(baseDataDir(), "BrowserProfile")
	if err := os.MkdirAll(profile, 0700); err != nil {
		return fmt.Errorf("сохранения: %w", err)
	}
	c := exec.Command(browser, "--app="+origin, "--user-data-dir="+profile, "--no-first-run", "--no-default-browser-check")
	if err := c.Start(); err != nil {
		return err
	}
	go func() { _ = c.Wait() }()
	return nil
}
