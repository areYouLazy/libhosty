package libhosty

import (
	"net"
	"runtime"
	"strings"
	"testing"
)

var hf *HostsFile
var hc *HostsFileConfig

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

func TestNewHostsFileConfig(t *testing.T) {
	var err error

	hc, err = NewHostsFileConfig("")
	if err != nil {
		t.Fatal(err)
	}
	switch runtime.GOOS {
	case "windows":
		if res := strings.Compare(hc.FilePath, "C:\\Windows\\System32\\drivers\\etc\\hosts"); res != 0 {
			t.Fatalf("error in hostsConfig path: wants %q got %q", "C:\\Windows\\System32\\drivers\\etc\\hosts", hc.FilePath)
		}
	default:
		if res := strings.Compare(hc.FilePath, "/etc/hosts"); res != 0 {
			t.Fatalf("error in hostsConfig path: wants %q got %q", "/etc/hosts", hc.FilePath)
		}
	}

	hc, err := NewHostsFileConfig("/etc")
	if err != nil {
		t.Fatal(err)
	}
	if res := strings.Compare(hf.Config.FilePath, "/etc/hosts"); res != 0 {
		t.Fatalf("should have %q got %q", "/etc/hosts", hc.FilePath)
	}

	hc, err = NewHostsFileConfig("/etc/hosts")
	if err != nil {
		t.Fatal(err)
	}

	if res := strings.Compare(hc.FilePath, "/etc/hosts"); res != 0 {
		t.Fatalf("error in hostsConfig path: wants %q got %q", "/etc/hosts", hc.FilePath)
	}
}

func TestInitWithConf(t *testing.T) {
	var err error

	//test with conf = nil
	hf, err = InitWithConfig(nil)
	if err != nil {
		t.Fatal(err)
	}

	// test with conf
	hf, err = InitWithConfig(hc)
	if err != nil {
		t.Fatal(err)
	}

	if res := strings.Compare(hf.Config.FilePath, "/etc/hosts"); res != 0 {
		t.Fatalf("error in InitWithConfig path: wants %q got %q", "/etc/hosts", hc.FilePath)
	}
}

func TestGetHostsFileLineByRow(t *testing.T) {
	idx, _, _ := hf.AddHostsFileLine("9.9.9.9", "gethostsfilelinebyrow.libhosty.local", "")
	hfl := hf.GetHostsFileLineByRow(idx)
	if hfl.Number != 0 {
		t.Fatalf("error: wants %d got %d", idx, hfl.Number)
	}
}

func TestGetHostsFileLineByIP(t *testing.T) {
	_, _, err := hf.AddHostsFileLine("8.8.8.8", "gethostsfilelinebyip.libhosty.local", "")
	if err != nil {
		t.Fatal(err)
	}

	ip := net.ParseIP("8.8.8.8")
	_, hfl := hf.GetHostsFileLineByIP(ip)

	if !net.IP.Equal(ip, hfl.Address) {
		t.Fatalf("error: wants %q got %q", ip, hfl.Address)
	}

	ip = net.ParseIP("fa.ke.i.p")
	idx, _ := hf.GetHostsFileLineByIP(ip)
	if idx != -1 {
		t.Fatalf("wants %d got %d", -1, idx)
	}
}

func TestGetHostsFileLineByAddress(t *testing.T) {
	hf.AddHostsFileLine("7.7.7.7", "gethostsfilelinebyaddress.libhosty.local", "")
	_, hfl := hf.GetHostsFileLineByAddress("7.7.7.7")

	if res := strings.Compare(hfl.Address.String(), "7.7.7.7"); res != 0 {
		t.Fatalf("error: wants %q got %q", "7.7.7.7", hfl.Address.String())
	}
}

func TestGetHostsFileLineByHostname(t *testing.T) {
	hf.AddHostsFileLine("6.6.6.6", "gethostsfilelinebyhostname.libhosty.local", "")
	_, hfl := hf.GetHostsFileLineByHostname("gethostsfilelinebyhostname.libhosty.local")

	res := false
	for _, v := range hfl.Hostnames {
		if v == "gethostsfilelinebyhostname.libhosty.local" {
			res = true
		}
	}

	if res != true {
		t.Fatalf("error: missing localhost in hostnames: %s", hfl.Hostnames)
	}

	idx, _ := hf.GetHostsFileLineByHostname("")
	if idx != -1 {
		t.Fatalf("wants %d got %d", -1, idx)
	}
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

func TestAddHostsFileLine(t *testing.T) {
	idx, _, err := hf.AddHostsFileLine("5.5.5.5", "addhost.libhosty.local", "my comment")
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
		if res := strings.Compare(v, "addhost.libhosty.local"); res == 0 {
			r = true
		}
	}
	if !r {
		t.Fatalf("mssing %q in hostnames: got %q", "addhost.libhosty.local", hfl.Hostnames)
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

func TestUncommentHostsFileLineByRow(t *testing.T) {
	idx, _ := hf.GetHostsFileLineByAddress("3.3.3.3")

	err := hf.UncommentHostsFileLineByRow(idx)
	if err != nil {
		t.Fatal(err)
	}

	hfl := hf.GetHostsFileLineByRow(idx)
	if hfl.IsCommented {
		t.Fatal("line should be uncommented")
	}
}

func TestCommentHostsFileLineByIP(t *testing.T) {
	_, _, err := hf.AddHostsFileLine("2.2.2.2", "commentbyip.libhosty.local", "")
	if err != nil {
		t.Fatal(err)
	}

	ip := net.ParseIP("2.2.2.2")
	err = hf.CommentHostsFileLineByIP(ip)
	if err != nil {
		t.Fatal(err)
	}

	_, hfl2 := hf.GetHostsFileLineByIP(ip)
	if err != nil {
		t.Fatal(err)
	}

	if !hfl2.IsCommented {
		t.Fatal("line should be commented")
	}
}

func TestUncommentHostsFileLineByIP(t *testing.T) {
	ip := net.ParseIP("2.2.2.2")
	err := hf.UncommentHostsFileLineByIP(ip)
	if err != nil {
		t.Fatal(err)
	}

	_, hfl := hf.GetHostsFileLineByIP(ip)
	if hfl.IsCommented {
		t.Fatal("line should be uncommented")
	}
}

func TestCommentHostsFileLineByAddress(t *testing.T) {
	_, _, err := hf.AddHostsFileLine("2.2.2.3", "commentbyaddress.libhosty.local", "")
	if err != nil {
		t.Fatal(err)
	}

	err = hf.CommentHostsFileLineByAddress("2.2.2.3")
	if err != nil {
		t.Fatal(err)
	}

	_, hfl2 := hf.GetHostsFileLineByAddress("2.2.2.3")
	if err != nil {
		t.Fatal(err)
	}

	if !hfl2.IsCommented {
		t.Fatal("line should be commented")
	}
}

func TestUncommentHostsFileLineByAddress(t *testing.T) {
	err := hf.UncommentHostsFileLineByAddress("2.2.2.3")
	if err != nil {
		t.Fatal(err)
	}

	_, hfl := hf.GetHostsFileLineByAddress("2.2.2.3")
	if hfl.IsCommented {
		t.Fatal("line should be uncommented")
	}
}

func TestCommentHostsFileLineByHostname(t *testing.T) {
	_, _, err := hf.AddHostsFileLine("2.2.2.4", "commentbyhostname.libhosty.local", "")
	if err != nil {
		t.Fatal(err)
	}

	err = hf.CommentHostsFileLineByHostname("commentbyhostname.libhosty.local")
	if err != nil {
		t.Fatal(err)
	}

	_, hfl2 := hf.GetHostsFileLineByHostname("commentbyhostname.libhosty.local")
	if err != nil {
		t.Fatal(err)
	}

	if !hfl2.IsCommented {
		t.Fatal("line should be commented")
	}
}

func TestUncommentHostsFileLineByHostname(t *testing.T) {
	err := hf.UncommentHostsFileLineByHostname("commentbyhostname.libhosty.local")
	if err != nil {
		t.Fatal(err)
	}

	_, hfl := hf.GetHostsFileLineByHostname("commentbyhostname.libhosty.local")
	if hfl.IsCommented {
		t.Fatal("line should be uncommented")
	}
}
