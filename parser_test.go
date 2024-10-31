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

	// check for invalid path
	path = "/my/custom/invalid/path/hosts"
	_, err := ParseHostsFile(path)
	if err == nil {
		t.Fatalf("should fail with invalid path: %s", path)
	}

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
	// also we assume to have no unknown line in a standard system
	// but we left it to ensure correct lines are not recognized as unknown lines
	//
	// edit: this now breaks the github automatic build due to https://github.com/actions/virtual-environments/issues/3353
	// so we must comment this out
	//
	// for k, v := range hf {
	// 	if v.Type == LineTypeUnknown {
	// 		t.Fatalf("unable to parse line at index %d: %s", k, v.Raw)
	// 	}
	// }

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

//TestParseHostsFileAsString test unknown line type
// actual parsing is already tested in TestParseHostsFile()
func TestParseHostsFileAsString(t *testing.T) {
	// define a custom hosts file with an address line and an invalid line"
	var fakeHostsFile = `1.1.1.1 my.hosts.file # With Comment
129dj120isdj12i0 1092jd 210dk`

	hf, err := ParseHostsFileAsString(fakeHostsFile)
	if err != nil {
		t.Fatalf("error parsing fakeHostsFile: %s", err)
	}

	// we expect the 1st line to be of type address
	if hf[0].Type != LineTypeAddress {
		t.Fatalf("line should be of type address: %v", hf[0])
	}
	// and the 2nd line to be of type unknown
	if hf[1].Type != LineTypeUnknown {
		t.Fatalf("line should be of type address: %v", hf[1])
	}
}
