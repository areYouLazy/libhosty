package libhosty

import (
	"strings"
	"testing"
)

var (
	// for this test use a custom hosts file
	customHostsFile = `## Custom hosts file
		# To check if everything works
		1.1.1.1		my.cloudflare.domain
		8.8.8.8		my.google.dns
		6.6.6.6		an.evil.domain
		# 12.12.12.12	commented.evil.domain # with comments
		# 2.3.4.5		first.domain second.domain third.domain
		5.6.7.8		first.domain second.domain third.domain
		5.6.7.8		a.domain b.domain c.domain`
)

// TestParseHostsFile implicitly tests parser()
// assuming the test is executed on a computer,
// we should be able to load the default hosts file
func TestParseHostsFile(t *testing.T) {
	// check for invalid path
	path := "/my/custom/invalid/path/hosts"
	_, err := InitFromCustomPath(path)
	if err == nil {
		t.Fatalf("should fail with invalid path: %s", path)
	}

	// parse file
	hf, err := Init()
	// check for errors
	if err != nil {
		t.Fatalf("error parsing hosts file: %s", err)
	}

	// load custom file for tests
	hf.HostsFileLines, err = ParseHostsFileFromString(customHostsFile)
	if err != nil {
		t.Fatalf("error parsing custom hosts file: %s", err)
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

	// for every  line,
	for k, v := range hf.HostsFileLines {
		// ensure addresses have a valid ip
		if v.Type == LineTypeAddress {
			if v.Address == nil {
				t.Fatalf("address line without address at index %d: %v", k, v)
			}
		}

		// ensure comments starts with #
		if v.Type == LineTypeComment {
			if !strings.HasPrefix(v.Raw, "#") {
				t.Fatalf("comment line does not starts with # at index %d: %v", k, v)
			}
		}
	}
}

// TestParseHostsFileFromString test unknown line type
// actual parsing is already tested in TestParseHostsFile()
func TestParseHostsFileFromString(t *testing.T) {
	// parse custom hosts file from string
	hf, err := ParseHostsFileFromString(customHostsFile)
	if err != nil {
		t.Fatalf("error parsing customHostsFile: %s", err)
	}

	// for every  line,
	for k, v := range hf {
		// ensure addresses have a valid ip
		if v.Type == LineTypeAddress {
			if v.Address == nil {
				t.Fatalf("address line without address at index %d: %v", k, v)
			}
		}

		// ensure comments starts with #
		if v.Type == LineTypeComment {
			if !strings.HasPrefix(v.Raw, "#") {
				t.Fatalf("comment line does not starts with # at index %d: %v", k, v)
			}
		}
	}
}
