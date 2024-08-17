package libhosty

import "strings"

// RenderHostsFile render and returns the hosts file with the lineFormatter() routine
func (h *HostsFile) RenderHostsFile() string {
	// allocate a buffer for file lines
	var sliceBuffer []string

	// iterate HostsFileLines and popolate the buffer with formatted lines
	for _, l := range h.HostsFileLines {
		sliceBuffer = append(sliceBuffer, lineFormatter(l))
	}

	// strings.Join() prevent the last line from being a new blank line
	// as opposite to a for loop with fmt.Printf(buffer + '\n')
	return strings.Join(sliceBuffer, "\n")
}

// RenderHostsFileLine render and returns the given hosts line with the lineFormatter() routine
func (h *HostsFile) RenderHostsFileLine(row int) string {
	// iterate to find the row to render
	if len(h.HostsFileLines) > row {
		return lineFormatter(h.HostsFileLines[row])
	}

	return ""
}
