package libhosty

import (
	"net"
	"regexp"
	"testing"
)

func TestLineFormatter(t *testing.T) {
	// define custom hostsFileLine
	hfl := HostsFileLine{
		Number:      0,
		Type:        30,
		Address:     net.ParseIP("1.1.1.1"),
		Parts:       []string{""},
		Hostnames:   []string{"my.host.name"},
		Raw:         "1.1.1.1 my.host.name # This is a host",
		Comment:     "This is a host",
		IsCommented: true,
		trimed:      "1.1.1.1 my.host.name",
	}

	// invoke lineFormatter hosts file line
	l := lineFormatter(hfl)

	// define what we expect
	w := regexp.MustCompile("1.1.1.1         \tmy.host.name\t#This is a host")

	// check
	if !w.MatchString(l) {
		t.Fatalf(`wants '%q' got '%q'`, w, l)
	}
}
