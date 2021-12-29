//go:build windows
// +build windows

package engine

import "os"

var CONFIG_PATH = os.Getenv("LOCALAPPDATA")