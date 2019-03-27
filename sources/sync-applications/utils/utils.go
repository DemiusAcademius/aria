package utils

import (
	"runtime"
	"path/filepath"
	"os"
)

// ExecutableDir get path to teleport executable
func ExecutableDir() string {
	executable, _ := os.Executable()
	return filepath.Dir(executable)
}

// UserHomeDir get path to User Home dir
func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}
