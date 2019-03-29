package core

import (
	"fmt"
	"runtime"
	"path/filepath"
	"os"

	"github.com/fatih/color"
)

// ExecutableDir get path to publish-project executable
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

// WorkingDir get path to current working dir
func WorkingDir() string {
	dir, _ := os.Getwd()
	return dir
}

// PrintBlue print label bw and text in blue
func PrintBlue(label, text string) {
	print(label)
	color.Blue(text)
}

// PrintMagenta print label bw and text in blue
func PrintMagenta(label, text string) {
	print(label)
	color.Magenta(text)
}

// PrintErrorAndPanic prints error text in Red
func PrintErrorAndPanic(errorText error) {
	color.Red(fmt.Sprintf("%s\n", errorText))
	os.Exit(-1)
}

// FileExists return True if file exists
func FileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
        return false
	}
	return true
}