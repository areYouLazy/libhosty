package libhosty

// RestoreDefaultWindowsHostsFile loads the default windows hosts file
func (h *HostsFile) RestoreDefaultWindowsHostsFile() {
	hfl, _ := ParseHostsFileFromString(windowsHostsTemplate)
	h.HostsFileLines = hfl
}

// RestoreDefaultLinuxHostsFile loads the default linux hosts file
func (h *HostsFile) RestoreDefaultLinuxHostsFile() {
	hfl, _ := ParseHostsFileFromString(linuxHostsTemplate)
	h.HostsFileLines = hfl
}

// RestoreDefaultDarwinHostsFile loads the default darwin hosts file
func (h *HostsFile) RestoreDefaultDarwinHostsFile() {
	hfl, _ := ParseHostsFileFromString(darwinHostsTemplate)
	h.HostsFileLines = hfl
}

// AddDockerDesktopTemplate adds the dockerDesktopTemplate to the actual hostsFile
func (h *HostsFile) AddDockerDesktopTemplate() {
	hfl, _ := ParseHostsFileFromString(dockerDesktopTemplate)
	h.HostsFileLines = append(h.HostsFileLines, hfl...)
}
