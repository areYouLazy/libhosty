package libhosty

import (
	"runtime"
	"strings"
	"testing"
)

//TestParseHostsFile implicitly tests parser()
// assuming the test is executed on a computer,
// we should be able to load the default hosts file
func TestParseHostsFile(t *testing.T) {
	var path string

	// define file path
	switch runtime.GOOS {
	case "win":
		path = "C:\\Windows\\System32\\drivers\\etc\\hosts"
	default:
		path = "/etc/hosts"
	}

	// parse file
	hf, err := ParseHostsFile(path)

	// check for errors
	if err != nil {
		t.Fatalf("error parsing hosts file: %s", err)
	}

	// check for unknown lines
	// this can be a false positive if the line is actually invalid
	// but we left it to ensure correct lines are not recognized as unknown lines
	for k, v := range hf {
		if v.Type == LineTypeUnknown {
			t.Fatalf("unable to parse line at index %d: %s", k, v.Raw)
		}
	}

	// for every address line, ensure we have a valid ip address
	for k, v := range hf {
		if v.Type == LineTypeAddress {
			if v.Address == nil {
				t.Fatalf("address line without address at index %d: %v", k, v)
			}
		}
	}

	// for every comment line ensure it starts with #
	for k, v := range hf {
		if v.Type == LineTypeComment {
			if !strings.HasPrefix(v.Raw, "#") {
				t.Fatalf("comment line does not starts with # at index %d: %v", k, v)
			}
		}
	}
}
