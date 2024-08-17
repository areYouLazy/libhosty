package libhosty

import (
	"net"
	"strings"
	"testing"
)

var hf *HostsFile

func TestInit(t *testing.T) {
	var err error

	hf, err = Init()
	if err != nil {
		t.Fatal(err)
	}

	if hf == nil {
		t.Fatal("hostsFile is nil")
	}

	if len(hf.HostsFileLines) <= 0 {
		t.Fatalf("we should have at least 1 line")
	}
}

func TestInitFromCustomPath(t *testing.T) {
	//TODO(areYouLazy): Test missing
}

func TestInitFromString(t *testing.T) {
	//TODO(areYouLazy): Test missing
}

func TestGetHostsFileLines(t *testing.T) {
	//TODO(areYouLazy): Test missing
}

func TestGetHostsFileLineByRow(t *testing.T) {
	idx, _, _ := hf.AddHostsFileLine("9.9.9.9", "gethostsfilelinebyrow.libhosty.local", "")
	hfl := hf.GetHostsFileLineByRow(idx)
	if hfl.Number != idx {
		t.Fatalf("error: wants %d got %d", idx, hfl.Number)
	}

	if len(hfl.Hostnames) > 1 {
		t.Fatalf("error: wants %d got %d", 1, len(hfl.Hostnames))
	}

	if hfl.Hostnames[0] != "gethostsfilelinebyrow.libhosty.local" {
		t.Fatalf("error: wants %s got %s", "gethostsfilelinebyrow.libhosty.local", hfl.Hostnames[0])
	}
}

func TestGetHostsFileLinesByIP(t *testing.T) {
	_, _, err := hf.AddHostsFileLine("8.8.8.8", "gethostsfilelinebyip.libhosty.local", "")
	if err != nil {
		t.Fatal(err)
	}

	ip := net.ParseIP("8.8.8.8")

	for _, hfl := range hf.GetHostsFileLinesByIP(ip) {
		if !net.IP.Equal(ip, hfl.Address) {
			t.Fatalf("error: wants %q got %q", ip, hfl.Address)
		}
	}
}

func TestGetHostsFileLinesByAddress(t *testing.T) {
	hf.AddHostsFileLine("7.7.7.7", "gethostsfilelinebyaddress.libhosty.local", "")

	for _, hfl := range hf.GetHostsFileLinesByAddress("7.7.7.7") {
		if res := strings.Compare(hfl.Address.String(), "7.7.7.7"); res != 0 {
			t.Fatalf("error: wants %q got %q", "7.7.7.7", hfl.Address.String())
		}
	}
}

func TestGetHostsFileLinesByHostname(t *testing.T) {
	hf.AddHostsFileLine("6.6.6.6", "gethostsfilelinesbyhostname.libhosty.local", "")

	res := false
	for _, hfl := range hf.GetHostsFileLinesByHostname("gethostsfilelinesbyhostname.libhosty.local") {
		for _, v := range hfl.Hostnames {
			if v == "gethostsfilelinesbyhostname.libhosty.local" {
				res = true
			}
		}

		if res != true {
			t.Fatalf("error: missing localhost in hostnames: %s", hfl.Hostnames)
		}
	}
}

func TestGetHostsFileLinesByRegexp(t *testing.T) {
	//TODO(areYouLazy): Test missing
}

func TestRemoveHostsFileLineByRow(t *testing.T) {
	//TODO(areYouLazy): Test missing
}

func TestRemoveHostsFileLinesByIP(t *testing.T) {
	//TODO(areYouLazy): Test missing
}

func TestRemoveHostsFileLinesByAddress(t *testing.T) {
	//TODO(areYouLazy): Test missing
}

func TestRemoveHostsFileLinesByHostname(t *testing.T) {
	//TODO(areYouLazy): Test missing
}

func TestRemoveHostsFileLinesByRegexp(t *testing.T) {
	//TODO(areYouLazy): Test missing
}

func TestLookupByHostname(t *testing.T) {
	//TODO(areYouLazy): Test missing
}

func TestAddHostsFileLineRaw(t *testing.T) {
	idx, _, err := hf.AddHostsFileLineRaw("4.4.4.4", "addhostraw.libhosty.local", "my comment")
	if err != nil {
		t.Fatal(err)
	}

	hfl := hf.HostsFileLines[idx]
	if hfl.Type != LineTypeAddress {
		t.Fatalf("expecting address line, found %q", hfl.Type)
	}

	if res := strings.Compare(hfl.Address.String(), "4.4.4.4"); res != 0 {
		t.Fatalf("expecting %q, found %q", "4.4.4.4", hfl.Address.String())
	}

	r := false
	for _, v := range hfl.Hostnames {
		if res := strings.Compare(v, "addhostraw.libhosty.local"); res == 0 {
			r = true
		}
	}
	if !r {
		t.Fatalf("mssing %q in hostnames: got %q", "addhostraw.libhosty.local", hfl.Hostnames)
	}

	if res := strings.Compare(hfl.Comment, "my comment"); res != 0 {
		t.Fatalf("expecting %q, found %q", "my comment", hfl.Comment)
	}

	idx, _, _ = hf.AddHostsFileLineRaw("fa.ke.i.p", "addhostraw.libhosty.local", "")
	if idx != -1 {
		t.Fatalf("wants %d got %d", -1, idx)
	}
}

func TestAddHostsFileLine(t *testing.T) {
	idx, _, err := hf.AddHostsFileLine("5.5.5.5", "1.libhosty.local", "my comment")
	if err != nil {
		t.Fatal(err)
	}

	hfl := hf.HostsFileLines[idx]
	if hfl.Type != LineTypeAddress {
		t.Fatalf("expecting address line, found %q", hfl.Type)
	}

	if res := strings.Compare(hfl.Address.String(), "5.5.5.5"); res != 0 {
		t.Fatalf("expecting %q, found %q", "5.5.5.5", hfl.Address.String())
	}

	r := false
	for _, v := range hfl.Hostnames {
		if res := strings.Compare(v, "1.libhosty.local"); res == 0 {
			r = true
		}
	}
	if !r {
		t.Fatalf("mssing %q in hostnames: got %q", "1.libhosty.local", hfl.Hostnames)
	}

	if res := strings.Compare(hfl.Comment, "my comment"); res != 0 {
		t.Fatalf("expecting %q, found %q", "my comment", hfl.Comment)
	}

	_, _, err = hf.AddHostsFileLine("5.5.5.5", "addhost.libhosty.local", "")
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = hf.AddHostsFileLine("5.5.5.5", "addhost2.libhosty.local", "")
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = hf.AddHostsFileLine("5.5.5.6", "addhost2.libhosty.local", "")
	if err != nil {
		t.Fatal(err)
	}

	// add more than 6 hostnames
	hf.AddHostsFileLine("5.5.5.5", "2.libhosty.local", "")
	hf.AddHostsFileLine("5.5.5.5", "3.libhosty.local", "")
	hf.AddHostsFileLine("5.5.5.5", "4.libhosty.local", "")
	hf.AddHostsFileLine("5.5.5.5", "5.libhosty.local", "")
	hf.AddHostsFileLine("5.5.5.5", "6.libhosty.local", "")
	hf.AddHostsFileLine("5.5.5.5", "7.libhosty.local", "")
}

func TestAddCommentFileLine(t *testing.T) {
	hf.AddCommentFileLine("comment")

	hfl := hf.HostsFileLines[len(hf.HostsFileLines)-1]
	if hfl.Type != LineTypeComment {
		t.Fatalf("expecting comment line, found %q", hfl.Type)
	}

	if res := strings.Compare(hfl.Raw, "# comment"); res != 0 {
		t.Fatalf("wants %q got %q", "# comment", hfl.Raw)
	}
}

func TestAddEmptyFileLine(t *testing.T) {
	idx, _, err := hf.AddEmptyFileLine()
	if err != nil {
		t.Fatal(err)
	}

	hfl := hf.HostsFileLines[idx]
	if hfl.Type != LineTypeEmpty {
		t.Fatalf("expecting empty line, found %q", hfl.Type)
	}
}

func TestCommentHostsFileLineByRow(t *testing.T) {
	idx, hfl, err := hf.AddHostsFileLine("3.3.3.3", "commentbyrow.host.name", "")
	if err != nil {
		t.Fatal(err)
	}

	if hfl.IsCommented {
		t.Fatal("new line is commented")
	}

	err = hf.CommentHostsFileLineByRow(idx)
	if err != nil {
		t.Fatal(err)
	}

	hfl2 := hf.GetHostsFileLineByRow(idx)
	if err != nil {
		t.Fatal(err)
	}

	if !hfl2.IsCommented {
		t.Fatal("line should be commented")
	}
}

func TestCommentHostsFileLinesByIP(t *testing.T) {
	//TODO(areYouLazy): Test missing
}

func TestCommentHostsFileLineByAddress(t *testing.T) {
	//TODO(areYouLazy): Test missing
}

func TestCommentHostsFileLineByHostname(t *testing.T) {
	//TODO(areYouLazy): Test missing
}

func TestCommentHostsFileLinesByRegexp(t *testing.T) {
	//TODO(areYouLazy): Test missing
}

func TestUncommentHostsFileLinesByRow(t *testing.T) {
	idx, _, err := hf.AddHostsFileLine("3.3.3.3", "commentbyrow.host.name", "")
	if err != nil {
		t.Fatal(err)
	}

	err = hf.UncommentHostsFileLineByRow(idx)
	if err != nil {
		t.Fatal(err)
	}

	hfl := hf.GetHostsFileLineByRow(idx)
	if hfl.IsCommented {
		t.Fatal("line should be uncommented")
	}
}

func TestUncommentHostsFileLineByIP(t *testing.T) {
	//TODO(areYouLazy): Test missing
}

func TestUncommentHostsFileLineByAddress(t *testing.T) {
	//TODO(areYouLazy): Test missing
}

func TestUncommentHostsFileLineByHostname(t *testing.T) {
	//TODO(areYouLazy): Test missing
}

func TestUncommentHostsFileLinesByRegexp(t *testing.T) {
	//TODO(areYouLazy): Test missing
}

func TestSaveHostsFile(t *testing.T) {
	//TODO(areYouLazy): Test missing
}

func TestSaveHostsFileAs(t *testing.T) {
	//TODO(areYouLazy): Test missing
}
