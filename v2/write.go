package libhosty

import (
	"os"
)

// WriteHostsFile write hosts file to configured path.
// error is not nil if something goes wrong
func (h *HostsFile) WriteHostsFile() error {
	if h.Path != "" {
		return h.WriteHostsFileTo(h.Path)
	} else {
		return ErrPathNotConfigured
	}
}

// WriteHostsFileTo write hosts file to the given path.
// error is not nil if something goes wrong
func (h *HostsFile) WriteHostsFileTo(path string) error {
	// render the file as a byte slice
	dataBytes := []byte(h.RenderHostsFile())

	// write file to disk
	err := os.WriteFile(path, dataBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}
