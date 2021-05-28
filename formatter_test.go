package libhosty

import (
	"net"
	"strings"
	"testing"
)

func TestLineFormatter(t *testing.T) {
	// define custom hostsFileLine with IsCommented
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
	w := "# 1.1.1.1         \tmy.host.name\t#This is a host"

	// check
	if !strings.EqualFold(l, w) {
		t.Fatalf(`IsCommented=true and Comment: wants '%q' got '%q'`, w, l)
	}

	// test without IsCommented
	hfl.IsCommented = false
	l = lineFormatter(hfl)
	w = "1.1.1.1         \tmy.host.name\t#This is a host"
	if !strings.EqualFold(l, w) {
		t.Fatalf(`IsCommented=false and Comment: wants '%q' got '%q'`, w, l)
	}

	// test with IsCommented but without comment in line
	hfl.IsCommented = true
	hfl.Comment = ""
	l = lineFormatter(hfl)
	w = "# 1.1.1.1         \tmy.host.name"
	if !strings.EqualFold(l, w) {
		t.Fatalf(`IsCommented=true no Comment: wants '%q' got '%q'`, w, l)
	}

	// check without IsCommented
	hfl.IsCommented = false
	hfl.Comment = ""
	l = lineFormatter(hfl)
	w = "1.1.1.1         \tmy.host.name"
	if !strings.EqualFold(l, w) {
		t.Fatalf(`IsCommented=false no Comment: wants '%q' got '%q'`, w, l)
	}

	// define a comment line
	hfl = HostsFileLine{
		Number:      0,
		Type:        20,
		Address:     []byte{},
		Parts:       []string{},
		Hostnames:   []string{},
		Raw:         "# Comment Line",
		Comment:     "",
		IsCommented: false,
		trimed:      "",
	}
	w = "# Comment line"
	l = lineFormatter(hfl)
	if !strings.EqualFold(l, w) {
		t.Fatalf(`Comment: wants '%q' got '%q'`, w, l)
	}
}
