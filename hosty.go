package libhosty

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"runtime"
	"strings"
	"sync"
)

const (
	//UNKNOWN defines unknown line type
	UNKNOWN = 0
	//EMPTY defines empty line type
	EMPTY = 10
	//COMMENT defines comment line type
	COMMENT = 20
	//ADDRESS defines address line type
	ADDRESS = 30

	// defines default path for windows os
	windowsFilePath = "C:\\Windows\\System32\\drivers\\etc\\"
	// defines default path for linux os
	unixFilePath = "/etc/"
	// defines default filename
	hostsFileName = "hosts"
)

//HostsConfig defines parameters to find hosts file.
// FilePath is the absolute path of the hosts file (filename included)
type HostsConfig struct {
	FilePath string
}

//HostsFileLine holds hosts file lines data
type HostsFileLine struct {
	//LineNumber is the original line number
	LineNumber int

	//LineType is one of the types: UNKNOWN, EMPTY, COMMENT, ADDRESS
	LineType int

	//Address is a net.IP representation of the address
	Address net.IP

	//Parts is a slice of the line splitted by '#'
	Parts []string

	//Hostnames is a slice of hostnames for the relative IP
	Hostnames []string

	//Raw is the raw representation of the line, as it is in the hosts file
	Raw string

	//Trimed is a trimed version (no spaces before and after) of the line
	Trimed string

	//Comment is the comment part of the line (if present in an ADDRESS line)
	Comment string

	//IsCommented to know if the current ADDRESS line is commented out (starts with '#')
	IsCommented bool
}

//HostsFile is a reference for the hosts file configuration and lines
type HostsFile struct {
	sync.Mutex

	//Config reference to a HostsConfig object
	Config *HostsConfig

	//HostsFileLines slice of HostsFileLine objects
	HostsFileLines []HostsFileLine
}

//Init returns a new instance of a hostsfile.
// Init gets the default Hosts configuratin and allocate
// an empty slice of HostsFileLine objects to store the parsed hosts file
func Init() *HostsFile {
	var err error

	// allocate a new HostsFile object
	hf := &HostsFile{
		// use default configuration
		Config: initHostsConfig(),

		// allocate a new slice of HostsFileLine objects
		HostsFileLines: make([]HostsFileLine, 0),
	}

	// parse the hosts file and load file lines
	hf.HostsFileLines, err = ParseHostsFile(hf.Config.FilePath)
	if err != nil {
		panic(err)
	}

	//return HostsFile
	return hf
}

// initHostsConfig loads hosts file based on environment.
// initHostsConfig initializa the default file path based
// on the OS since the location file cannot be changed
func initHostsConfig() *HostsConfig {
	var hc *HostsConfig
	if runtime.GOOS == "windows" {
		hc = &HostsConfig{
			FilePath: windowsFilePath + hostsFileName,
		}
	} else if runtime.GOOS == "linux" {
		hc = &HostsConfig{
			FilePath: unixFilePath + hostsFileName,
		}
	} else if runtime.GOOS == "darwin" {
		hc = &HostsConfig{
			FilePath: unixFilePath + hostsFileName,
		}
	} else {
		fmt.Printf("Unrecognized os: %s", runtime.GOOS)
		os.Exit(1)
	}

	return hc
}

//GetHostsFileLineByRow returns a ponter to the given HostsFileLine row
func (h *HostsFile) GetHostsFileLineByRow(row int) *HostsFileLine {
	return &h.HostsFileLines[row]
}

//GetHostsFileLineByAddress returns the index of the line and a ponter to the given HostsFileLine line
func (h *HostsFile) GetHostsFileLineByAddress(ip net.IP) (int, *HostsFileLine) {
	for idx := range h.HostsFileLines {
		if ip.String() == h.HostsFileLines[idx].Address.String() {
			return idx, &h.HostsFileLines[idx]
		}
	}

	return -1, nil
}

//GetHostsFileLineByAddressAsString returns the index of the line and a ponter to the given HostsFileLine line
func (h *HostsFile) GetHostsFileLineByAddressAsString(address string) (int, *HostsFileLine) {
	ip := net.ParseIP(address)
	return h.GetHostsFileLineByAddress(ip)
}

//GetHostsFileLineByHostname returns the index of the line and a ponter to the given HostsFileLine line
func (h *HostsFile) GetHostsFileLineByHostname(hostname string) (int, *HostsFileLine) {
	for idx := range h.HostsFileLines {
		for _, hn := range h.HostsFileLines[idx].Hostnames {
			if hn == hostname {
				return idx, &h.HostsFileLines[idx]
			}
		}
	}

	return -1, nil
}

//RenderHostsFile render and returns the hosts file with the lineFormatter() routine
func (h *HostsFile) RenderHostsFile() string {
	// allocate a buffer for file lines
	var sliceBuffer []string

	// iterate HostsFileLines and popolate the buffer with formatted lines
	for _, l := range h.HostsFileLines {
		sliceBuffer = append(sliceBuffer, lineFormatter(l))
	}

	// strings.Join() prevent the last line from being a new blank line
	// as opposite to a for loop with fmt.Printf(buffer + '\n')
	return strings.Join(sliceBuffer, "\n")
}

//RenderHostsFileLine render and returns the given hosts line with the lineFormatter() routine
func (h *HostsFile) RenderHostsFileLine(row int) string {
	// iterate to find the row to render
	for k, l := range h.HostsFileLines {
		if k == row {
			return lineFormatter(l)
		}
	}

	return ""
}

//SaveHostsFile write hosts file to configured path.
// error is not nil if something goes wrong
func (h *HostsFile) SaveHostsFile() error {
	return h.SaveHostsFileAs(h.Config.FilePath)
}

//SaveHostsFileAs write hosts file to the given path.
// error is not nil if something goes wrong
func (h *HostsFile) SaveHostsFileAs(path string) error {
	// render the file as a byte slice
	dataBytes := []byte(h.RenderHostsFile())

	// write file to disk
	err := ioutil.WriteFile(path, dataBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

//RemoveRow remove row at given index from HostsFileLines
func (h *HostsFile) RemoveRow(row int) {
	h.Lock()
	defer h.Unlock()

	// prevent out-of-index
	if row < len(h.HostsFileLines) {
		h.HostsFileLines = append(h.HostsFileLines[:row], h.HostsFileLines[row+1:]...)
	}
}

//LookupByHostname check if the given fqdn exists.
// if yes, it returns the index of the address and the associated address.
// error is not nil if something goes wrong
func (h *HostsFile) LookupByHostname(hostname string) (int, net.IP, error) {
	for i, v := range h.HostsFileLines {
		for _, k := range v.Hostnames {
			if k == hostname {
				return i, h.HostsFileLines[i].Address, nil
			}
		}
	}

	return -1, nil, errors.New("Hostname not found")
}

//AddHost add the given ip/fqdn/comment pair, cleanup is done for previous entry.
// it returns the index of the edited (created) line and a pointer to the hostsfileline object.
// error is not nil if something goes wrong
func (h *HostsFile) AddHost(ipRaw, fqdnRaw, comment string) (int, *HostsFileLine, error) {
	h.Lock()
	defer h.Unlock()

	// hostname to lowercase
	hostname := strings.ToLower(fqdnRaw)
	// parse ip to net.IP
	ip := net.ParseIP(ipRaw)

	// if we have a valid IP
	if ip != nil {
		//check if we alredy have the fqdn
		if idx, addr, err := h.LookupByHostname(hostname); err == nil {
			//if actual ip is the same as the given one, we are done
			if addr.String() == ip.String() {
				// handle comment
				if comment != "" {
					// just replace the current comment with the new one
					h.HostsFileLines[idx].Comment = comment
				}
				return idx, &h.HostsFileLines[idx], nil
			}

			//if address is different, we need to remove the hostname from the previous entry
			for hostIdx, fqdn := range h.HostsFileLines[idx].Hostnames {
				if fqdn == hostname {
					h.HostsFileLines[idx].Hostnames = append(h.HostsFileLines[idx].Hostnames[:hostIdx], h.HostsFileLines[idx].Hostnames[hostIdx+1:]...)

					//also remove the line if there are no more hostnames
					if len(h.HostsFileLines[idx].Hostnames) < 1 {
						h.RemoveRow(idx)
					}
				}
			}
		}

		//if we alredy have the address, just add the hostname to that line
		for k, v := range h.HostsFileLines {
			if v.Address.String() == ip.String() {
				h.HostsFileLines[k].Hostnames = append(h.HostsFileLines[k].Hostnames, hostname)
				// handle comment
				if comment != "" {
					// just replace the current comment with the new one
					h.HostsFileLines[k].Comment = comment
				}
				return k, &h.HostsFileLines[k], nil
			}
		}

		// at this point we need to create new host line
		hfl := HostsFileLine{
			LineType:    ADDRESS,
			Address:     ip,
			Hostnames:   []string{hostname},
			Comment:     comment,
			IsCommented: false,
		}

		hfl.Raw = lineFormatter(hfl)

		h.HostsFileLines = append(h.HostsFileLines, hfl)
		idx := len(h.HostsFileLines) - 1
		return idx, &h.HostsFileLines[idx], nil
	}

	return -1, nil, fmt.Errorf("Cannot parse IP address %s", ipRaw)
}

//AddComment adds a new line of type comment with the given comment.
// it returns the index of the edited (created) line and a pointer to the hostsfileline object.
// error is not nil if something goes wrong
func (h *HostsFile) AddComment(comment string) (int, *HostsFileLine, error) {
	h.Lock()
	defer h.Unlock()

	hfl := HostsFileLine{
		LineType: COMMENT,
		Raw:      "# " + comment,
	}

	hfl.Raw = lineFormatter(hfl)

	h.HostsFileLines = append(h.HostsFileLines, hfl)
	idx := len(h.HostsFileLines) - 1
	return idx, &h.HostsFileLines[idx], nil
}

//AddEmpty adds a new line of type empty.
// it returns the index of the edited (created) line and a pointer to the hostsfileline object.
// error is not nil if something goes wrong
func (h *HostsFile) AddEmpty() (int, *HostsFileLine, error) {
	h.Lock()
	defer h.Unlock()

	hfl := HostsFileLine{
		LineType: EMPTY,
	}

	hfl.Raw = ""

	h.HostsFileLines = append(h.HostsFileLines, hfl)
	idx := len(h.HostsFileLines) - 1
	return idx, &h.HostsFileLines[idx], nil
}

//CommentByRow set the IsCommented bit for the given row to true
func (h *HostsFile) CommentByRow(row int) {
	h.Lock()
	defer h.Unlock()

	if row <= len(h.HostsFileLines) {
		if h.HostsFileLines[row].LineType == ADDRESS {
			h.HostsFileLines[row].IsCommented = true
			return
		}

		return
	}

	return
}

//CommentByAddress set the IsCommented bit for the given address to true
func (h *HostsFile) CommentByAddress(ip net.IP) {
	h.Lock()
	defer h.Unlock()

	for idx, hfl := range h.HostsFileLines {
		if ip.String() == hfl.Address.String() {
			h.HostsFileLines[idx].IsCommented = true
		}
	}
}

//CommentByAddressAsString set the IsCommented bit for the given address as string to false
func (h *HostsFile) CommentByAddressAsString(address string) {
	ip := net.ParseIP(address)

	h.CommentByAddress(ip)
}

//CommentByHostname set the IsCommented bit for the given hostname to true
func (h *HostsFile) CommentByHostname(hostname string) {
	h.Lock()
	defer h.Unlock()

	for idx := range h.HostsFileLines {
		for _, hn := range h.HostsFileLines[idx].Hostnames {
			if hn == hostname {
				h.HostsFileLines[idx].IsCommented = true
			}
		}
	}
}

//UncommentByRow set the IsCommented bit for the given row to false
func (h *HostsFile) UncommentByRow(row int) {
	h.Lock()
	defer h.Unlock()

	if row <= len(h.HostsFileLines) {
		if h.HostsFileLines[row].LineType == ADDRESS {
			h.HostsFileLines[row].IsCommented = false
			return
		}

		return
	}

	return
}

//UncommentByAddress set the IsCommented bit for the given address to false
func (h *HostsFile) UncommentByAddress(ip net.IP) {
	h.Lock()
	defer h.Unlock()

	for idx, hfl := range h.HostsFileLines {
		if ip.String() == hfl.Address.String() {
			h.HostsFileLines[idx].IsCommented = false
		}
	}
}

//UncommentByAddressAsString set the IsCommented bit for the given address as string to false
func (h *HostsFile) UncommentByAddressAsString(address string) {
	ip := net.ParseIP(address)

	h.UncommentByAddress(ip)
}

//UncommentByHostname set the IsCommented bit for the given hostname to false
func (h *HostsFile) UncommentByHostname(hostname string) {
	h.Lock()
	defer h.Unlock()

	for idx := range h.HostsFileLines {
		for _, hn := range h.HostsFileLines[idx].Hostnames {
			if hn == hostname {
				h.HostsFileLines[idx].IsCommented = false
			}
		}
	}
}

//RestoreDefaultWindowsHostsFile loads the default windows hosts file
func (h *HostsFile) RestoreDefaultWindowsHostsFile() {
	hfl, _ := ParseHostsFileAsString(windowsHostsTemplate)
	h.HostsFileLines = hfl
}

//RestoreDefaultLinuxHostsFile loads the default linux hosts file
func (h *HostsFile) RestoreDefaultLinuxHostsFile() {
	hfl, _ := ParseHostsFileAsString(linuxHostsTemplate)
	h.HostsFileLines = hfl
}

//RestoreDefaultDarwinHostsFile loads the default darwin hosts file
func (h *HostsFile) RestoreDefaultDarwinHostsFile() {
	hfl, _ := ParseHostsFileAsString(darwinHostsTemplate)
	h.HostsFileLines = hfl
}
