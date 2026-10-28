//go:build !windows

package main

import "log"

func showError(s string) { log.Println(s) }
