package libhosty

import (
	"fmt"
	"regexp"
	"testing"
)

var testHostsFile = `1.1.1.1 my.host.name # This is a host`

func TestLineFormatter(t *testing.T) {
	// parse fake hosts file
	hfl, err := ParseHostsFileAsString(testHostsFile)
	if err != nil {
		t.Fatalf("unable to parse testHostsFile: %s", err)
	}
	fmt.Println(hfl)

	// invoke lineFormatter on 1st hosts file line
	l := lineFormatter(hfl[0])

	// define what we expect
	w := regexp.MustCompile("1.1.1.1         \tmy.host.name\t#This is a host")

	// check
	if !w.MatchString(l) {
		t.Fatalf(`wants '%q' got '%q'`, w, l)
	}
}
