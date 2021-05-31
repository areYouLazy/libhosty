package libhosty

import (
	"net"
	"runtime"
	"strings"
	"testing"
)

var hf *HostsFile
var hc *HostsConfig

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

func TestNewHostsConfig(t *testing.T) {
	var err error

	hc, err = NewHostsConfig("")
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

	hc, err := NewHostsConfig("/etc")
	if err != nil {
		t.Fatal(err)
	}
	if res := strings.Compare(hf.Config.FilePath, "/etc/hosts"); res != 0 {
		t.Fatalf("should have %q got %q", "/etc/hosts", hc.FilePath)
	}

	hc, err = NewHostsConfig("/etc/passwd")
	if err == nil && hc.FilePath != "/etc/hosts" {
		t.Fatalf("should have error or default hostsConfig")
	}

	hc, err = NewHostsConfig("/etc/hosts")
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
	hf, err = InitWithConf(nil)
	if err != nil {
		t.Fatal(err)
	}

	// test with conf
	hf, err = InitWithConf(hc)
	if err != nil {
		t.Fatal(err)
	}

	if res := strings.Compare(hf.Config.FilePath, "/etc/hosts"); res != 0 {
		t.Fatalf("error in InitWithConfig path: wants %q got %q", "/etc/hosts", hc.FilePath)
	}
}

func TestGetHostsFileLineByRow(t *testing.T) {
	hfl := hf.GetHostsFileLineByRow(0)
	if hfl.Number != 0 {
		t.Fatalf("error: wants 0 got %q", hfl.Number)
	}
}

func TestGetHostsFileLineByIP(t *testing.T) {
	ip := net.ParseIP("127.0.0.1")
	_, hfl := hf.GetHostsFileLineByIP(ip)

	if !net.IP.Equal(ip, hfl.Address) {
		t.Fatalf("error: wants %q got %q", ip, hfl.Address)
	}

	ip = net.ParseIP("321.321.321.321")
	idx, _ := hf.GetHostsFileLineByIP(ip)
	if idx != -1 {
		t.Fatalf("wants %d got %d", -1, idx)
	}
}

func TestGetHostsFileLineByAddress(t *testing.T) {
	_, hfl := hf.GetHostsFileLineByAddress("127.0.0.1")

	if res := strings.Compare(hfl.Address.String(), "127.0.0.1"); res != 0 {
		t.Fatalf("error: wants 127.0.0.1 got %q", hfl.Address.String())
	}
}

func TestGetHostsFileLineByHostname(t *testing.T) {
	_, hfl := hf.GetHostsFileLineByHostname("localhost")

	res := false
	for _, v := range hfl.Hostnames {
		if v == "localhost" {
			res = true
		}
	}

	if res != true {
		t.Fatalf("error: missing localhost in hostnames: %s", hfl.Hostnames)
	}
}

func TestAddComment(t *testing.T) {
	hf.AddComment("comment")

	hfl := hf.HostsFileLines[len(hf.HostsFileLines)-1]
	if hfl.Type != LineTypeComment {
		t.Fatalf("expecting comment line, found %q", hfl.Type)
	}

	if res := strings.Compare(hfl.Raw, "# comment"); res != 0 {
		t.Fatalf("wants %q got %q", "# comment", hfl.Raw)
	}
}

func TestAddEmpty(t *testing.T) {
	idx, _, err := hf.AddEmpty()
	if err != nil {
		t.Fatal(err)
	}

	hfl := hf.HostsFileLines[idx]
	if hfl.Type != LineTypeEmpty {
		t.Fatalf("expecting empty line, found %q", hfl.Type)
	}
}

func TestAddHost(t *testing.T) {
	idx, _, err := hf.AddHost("127.0.0.1", "my.host.name", "my comment")
	if err != nil {
		t.Fatal(err)
	}

	hfl := hf.HostsFileLines[idx]
	if hfl.Type != LineTypeAddress {
		t.Fatalf("expecting address line, found %q", hfl.Type)
	}

	if res := strings.Compare(hfl.Address.String(), "127.0.0.1"); res != 0 {
		t.Fatalf("expecting %q, found %q", "127.0.0.1", hfl.Address.String())
	}

	r := false
	for _, v := range hfl.Hostnames {
		if res := strings.Compare(v, "my.host.name"); res == 0 {
			r = true
		}
	}
	if !r {
		t.Fatalf("mssing %q in hostnames: got %q", "my.host.name", hfl.Hostnames)
	}

	if res := strings.Compare(hfl.Comment, "my comment"); res != 0 {
		t.Fatalf("expecting %q, found %q", "my comment", hfl.Comment)
	}
}

func TestAddHostRaw(t *testing.T) {
	idx, _, err := hf.AddHostRaw("127.0.0.1", "my.host.name", "my comment")
	if err != nil {
		t.Fatal(err)
	}

	hfl := hf.HostsFileLines[idx]
	if hfl.Type != LineTypeAddress {
		t.Fatalf("expecting address line, found %q", hfl.Type)
	}

	if res := strings.Compare(hfl.Address.String(), "127.0.0.1"); res != 0 {
		t.Fatalf("expecting %q, found %q", "127.0.0.1", hfl.Address.String())
	}

	r := false
	for _, v := range hfl.Hostnames {
		if res := strings.Compare(v, "my.host.name"); res == 0 {
			r = true
		}
	}
	if !r {
		t.Fatalf("mssing %q in hostnames: got %q", "my.host.name", hfl.Hostnames)
	}

	if res := strings.Compare(hfl.Comment, "my comment"); res != 0 {
		t.Fatalf("expecting %q, found %q", "my comment", hfl.Comment)
	}
}
