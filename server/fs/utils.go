package server

import (
	"os"
	"path/filepath"
	"runtime"
)

// selects folder based on OS and env variables of the user
func selectFolder() string {
	// user select
	if folder := os.Getenv("EM_FOLDER"); folder != "" {
		return folder
	}
	// defaults
	if runtime.GOOS == "windows" {
		// Use %AppData% on Windows
		return filepath.Join(os.Getenv("AppData"), "emmer")
	} else {
		// Use XDG_DATA_HOME on linux (if exists)
		xdgData := os.Getenv("XDG_DATA_HOME")
		if xdgData == "" {
			return filepath.Join(os.Getenv("HOME"), ".local", "share", "emmer")
		}
		return filepath.Join(xdgData, "emmer")
	}
}
