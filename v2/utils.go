package libhosty

import (
	"path/filepath"
	"runtime"
)

// GetOSHOstsFilePath returns the hostsfile absolute path based on runtime.GOOS result
func GetOSHostsFilePath() string {
	switch runtime.GOOS {
	case "windows":
		return filepath.Join(windowsFilePath, hostsFileName)
	default:
		return filepath.Join(unixFilePath, hostsFileName)
	}
}
