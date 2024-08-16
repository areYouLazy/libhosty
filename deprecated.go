package libhosty

import "os"

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

// SaveHostsFile write hosts file to configured path.
// error is not nil if something goes wrong
func (h *HostsFile) SaveHostsFile() error {
	return h.SaveHostsFileAs(h.Path)
}

// SaveHostsFileAs write hosts file to the given path.
// error is not nil if something goes wrong
func (h *HostsFile) SaveHostsFileAs(path string) error {
	// render the file as a byte slice
	dataBytes := []byte(h.RenderHostsFile())

	// write file to disk
	err := os.WriteFile(path, dataBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}
