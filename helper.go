package libhosty

import "runtime"

// RestoreTemplate tries to restore the default hostsfile based on runtime.GOOS result
// returns true if restore goes well
func (h *HostsFile) RestoreTemplate() bool {
	var hfl []HostsFileLine
	// var err error

	switch runtime.GOOS {
	case "windows":
		hfl, _ = ParseHostsFileFromString(windowsHostsTemplate)
	case "docker":
		hfl, _ = ParseHostsFileFromString(dockerDesktopTemplate)
	case "linux|unix":
		hfl, _ = ParseHostsFileFromString(linuxHostsTemplate)
	case "darwin":
		hfl, _ = ParseHostsFileFromString(darwinHostsTemplate)
	default:
		// hfl = nil
		hfl, _ = ParseHostsFileFromString(linuxHostsTemplate)
	}

	if hfl != nil {
		h.HostsFileLines = hfl
		return true
	}

	return false
}

// RestoreNamedTemplate restored the named template as the current hostsfile
// returns true if restore goes well
func (h *HostsFile) RestoreNamedTemplate(template string) bool {
	var hfl []HostsFileLine
	// var err error

	switch template {
	case "windows":
		hfl, _ = ParseHostsFileFromString(windowsHostsTemplate)
	case "docker":
		hfl, _ = ParseHostsFileFromString(dockerDesktopTemplate)
	case "linux|unix":
		hfl, _ = ParseHostsFileFromString(linuxHostsTemplate)
	case "darwin":
		hfl, _ = ParseHostsFileFromString(darwinHostsTemplate)
	default:
		// hfl = nil
		hfl, _ = ParseHostsFileFromString(linuxHostsTemplate)
	}

	if hfl != nil {
		h.HostsFileLines = hfl
		return true
	}

	return false
}

// AppendNamedTemplate appends the named template to the current hostsfile
// returns true if restore goes well
func (h *HostsFile) AppendNamedTemplate(template string) bool {
	var hfl []HostsFileLine
	// var err error

	switch template {
	case "windows":
		hfl, _ = ParseHostsFileFromString(windowsHostsTemplate)
	case "docker":
		hfl, _ = ParseHostsFileFromString(dockerDesktopTemplate)
	case "linux|unix":
		hfl, _ = ParseHostsFileFromString(linuxHostsTemplate)
	case "darwin":
		hfl, _ = ParseHostsFileFromString(darwinHostsTemplate)
	default:
		// hfl = nil
		hfl, _ = ParseHostsFileFromString(linuxHostsTemplate)
	}

	if hfl != nil {
		h.HostsFileLines = append(h.HostsFileLines, hfl...)
		return true
	}

	return false
}
