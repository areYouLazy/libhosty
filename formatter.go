package libhosty

import (
	"fmt"
	"strings"
)

// lineFormatter return a readable form for the given HostsFileLine object
func lineFormatter(hfl HostsFileLine) string {
	// returns raw for UNKNOWN linetype
	if hfl.Type == LineTypeUnknown {
		return hfl.Raw
	}

	// returns empty for EMPTY linetype
	if hfl.Type == LineTypeEmpty {
		return ""
	}

	// return a well formatted comment for COMEMNT linetype
	if hfl.Type == LineTypeComment {
		return fmt.Sprintf("# %s", hfl.Comment)
	}

	// address lines

	// check if it's a commented line
	if hfl.IsCommented {
		// check if there's a comment for that line
		if len(hfl.Comment) > 0 {
			return fmt.Sprintf("# %-16s %s #%s", hfl.Address, strings.Join(hfl.Hostnames, " "), hfl.Comment)
		}

		return fmt.Sprintf("# %-16s %s", hfl.Address, strings.Join(hfl.Hostnames, " "))
	}

	// return the actual hosts entry
	if len(hfl.Comment) > 0 {
		return fmt.Sprintf("%-16s %s #%s", hfl.Address, strings.Join(hfl.Hostnames, " "), hfl.Comment)
	}

	return fmt.Sprintf("%-16s %s", hfl.Address, strings.Join(hfl.Hostnames, " "))
}
