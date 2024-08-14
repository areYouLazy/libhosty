// Package libhosty is a pure golang library to manipulate the hosts file
package libhosty

import (
	"net"
	"os"
	"regexp"
	"strings"
	"sync"
)

const (
	//Version exposes library version
	Version = "2.0.0-alfa.1"
)

const (
	// defines default path for windows os
	windowsFilePath = "C:\\Windows\\System32\\drivers\\etc\\"

	// defines default path for linux os
	unixFilePath = "/etc/"

	// defines default filename
	hostsFileName = "hosts"
)

// LineType define a safe type for line type enumeration
type LineType int

const (
	//LineTypeUnknown defines unknown lines
	LineTypeUnknown LineType = 0

	//LineTypeEmpty defines empty lines
	LineTypeEmpty LineType = 10

	//LineTypeComment defines comment lines (starts with #)
	LineTypeComment LineType = 20

	//LineTypeAddress defines address lines (actual hosts lines)
	LineTypeAddress LineType = 30
)

// HostsFileLine holds hosts file lines data
type HostsFileLine struct {
	//Number is the original line number
	Number int

	//LineType defines the line type
	Type LineType

	//Address is a net.IP representation of the address
	Address net.IP

	//Parts is a slice of the line splitted by '#'
	Parts []string

	//Hostnames is a slice of hostnames for the relative IP
	Hostnames []string

	//Raw is the raw representation of the line, as it is in the hosts file
	Raw string

	//Comment is the comment part of the line (if present in an ADDRESS line)
	Comment string

	//IsCommented to know if the current ADDRESS line is commented out (starts with '#')
	IsCommented bool

	//trimed is a trimed version (no spaces before and after) of the line
	trimed string
}

// HostsFile is a reference for the hosts file configuration and lines
type HostsFile struct {
	sync.Mutex

	//Path is the file path and name
	Path string

	//HostsFileLines slice of HostsFileLine objects
	HostsFileLines []HostsFileLine
}

// Init returns a new instance of a hostsfile.
func Init(path string) (*HostsFile, error) {
	var hostsFileLines []HostsFileLine
	var err error

	// use custom path if provided
	// if not, check for glibc env variable
	// if not, go with the default file path
	if path != "" {
		hostsFileLines, err = ParseHostsFile(path)
		if err != nil {
			return nil, err
		}
	} else {
		// check if we have the glibc HOSTALIASES env variable defined
		fpath := os.Getenv("HOSTALIASES")
		if fpath != "" {
			hostsFileLines, err = ParseHostsFile(fpath)
			if err != nil {
				return nil, err
			}
		} else {
			// fallback to default location
			hostsFileLines, err = ParseHostsFile(GetOSHostsFilePath())
			if err != nil {
				return nil, err
			}
		}
	}

	// allocate a new HostsFile object
	hf := &HostsFile{
		// use default configuration
		Path: path,

		// allocate a new slice of HostsFileLine objects
		HostsFileLines: hostsFileLines,
	}

	//return HostsFile
	return hf, nil
}

// GetHostsFileLines returns every address row
func (h *HostsFile) GetHostsFileLines() []*HostsFileLine {
	var hfl []*HostsFileLine

	for idx := range h.HostsFileLines {
		if h.HostsFileLines[idx].Type == LineTypeAddress {
			hfl = append(hfl, h.GetHostsFileLineByRow(idx))
		}
	}

	return hfl
}

// GetHostsFileLineByRow returns a ponter to the given HostsFileLine row
func (h *HostsFile) GetHostsFileLineByRow(row int) *HostsFileLine {
	return &h.HostsFileLines[row]
}

// GetHostsFileLineByIP returns the index of the line and a ponter to the given HostsFileLine line
func (h *HostsFile) GetHostsFileLineByIP(ip net.IP) (int, *HostsFileLine) {
	if ip == nil {
		return -1, nil
	}

	for idx := range h.HostsFileLines {
		if net.IP.Equal(ip, h.HostsFileLines[idx].Address) {
			return idx, &h.HostsFileLines[idx]
		}
	}

	return -1, nil
}

func (h *HostsFile) GetHostsFileLinesByIP(ip net.IP) []*HostsFileLine {
	if ip == nil {
		return nil
	}

	hfl := make([]*HostsFileLine, 0)

	for idx := range h.HostsFileLines {
		if net.IP.Equal(ip, h.HostsFileLines[idx].Address) {
			hfl = append(hfl, &h.HostsFileLines[idx])
		}
	}

	return hfl
}

// GetHostsFileLineByAddress returns the index of the line and a ponter to the given HostsFileLine line
func (h *HostsFile) GetHostsFileLineByAddress(address string) (int, *HostsFileLine) {
	ip := net.ParseIP(address)
	return h.GetHostsFileLineByIP(ip)
}

func (h *HostsFile) GetHostsFileLinesByAddress(address string) []*HostsFileLine {
	ip := net.ParseIP(address)
	return h.GetHostsFileLinesByIP(ip)
}

// GetHostsFileLineByHostname returns the index of the line and a ponter to the given HostsFileLine line
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

func (h *HostsFile) GetHostsFileLinesByHostname(hostname string) []*HostsFileLine {
	hfl := make([]*HostsFileLine, 0)

	for idx := range h.HostsFileLines {
		for _, hn := range h.HostsFileLines[idx].Hostnames {
			if hn == hostname {
				hfl = append(hfl, &h.HostsFileLines[idx])
				continue
			}
		}
	}

	return hfl
}

func (h *HostsFile) GetHostsFileLinesByHostnameAsRegexp(hostname string) []*HostsFileLine {
	hfl := make([]*HostsFileLine, 0)

	reg := regexp.MustCompile(hostname)

	for idx := range h.HostsFileLines {
		for _, hn := range h.HostsFileLines[idx].Hostnames {
			if reg.MatchString(hn) {
				hfl = append(hfl, &h.HostsFileLines[idx])
				continue
			}
		}
	}

	return hfl
}

// RemoveHostsFileLineByRow remove row at given index from HostsFileLines
func (h *HostsFile) RemoveHostsFileLineByRow(row int) {
	// prevent out-of-index
	if row < len(h.HostsFileLines) {
		h.Lock()
		h.HostsFileLines = append(h.HostsFileLines[:row], h.HostsFileLines[row+1:]...)
		h.Unlock()
	}
}

func (h *HostsFile) RemoveHostsFileLineByIP(ip net.IP) {
	for idx := len(h.HostsFileLines) - 1; idx >= 0; idx-- {
		if net.IP.Equal(ip, h.HostsFileLines[idx].Address) {
			h.RemoveHostsFileLineByRow(idx)
			return
		}
	}
}

func (h *HostsFile) RemoveHostsFileLinesByIP(ip net.IP) {
	for idx := len(h.HostsFileLines) - 1; idx >= 0; idx-- {
		if net.IP.Equal(ip, h.HostsFileLines[idx].Address) {
			h.RemoveHostsFileLineByRow(idx)
		}
	}
}

func (h *HostsFile) RemoveHostsFileLineByAddress(address string) {
	ip := net.ParseIP(address)

	h.RemoveHostsFileLineByIP(ip)
}

func (h *HostsFile) RemoveHostsFileLinesByAddress(address string) {
	ip := net.ParseIP(address)

	h.RemoveHostsFileLinesByIP(ip)
}

func (h *HostsFile) RemoveHostsFileLineByHostname(hostname string) {
	for idx := len(h.HostsFileLines) - 1; idx >= 0; idx-- {
		if h.HostsFileLines[idx].Type == LineTypeAddress {
			for _, hn := range h.HostsFileLines[idx].Hostnames {
				if hn == hostname {
					h.RemoveHostsFileLineByRow(idx)
					return
				}
			}
		}
	}
}

func (h *HostsFile) RemoveHostsFileLinesByHostnameAsRegexp(hostname string) {
	reg := regexp.MustCompile(hostname)

	for idx := len(h.HostsFileLines) - 1; idx >= 0; idx-- {
		for _, hn := range h.HostsFileLines[idx].Hostnames {
			if reg.MatchString(hn) {
				h.RemoveHostsFileLineByRow(idx)
				continue
			}
		}
	}
}

func (h *HostsFile) RemoveHostsFileLinesByHostname(hostname string) {
	for idx := len(h.HostsFileLines) - 1; idx >= 0; idx-- {
		if h.HostsFileLines[idx].Type == LineTypeAddress {
			for _, hn := range h.HostsFileLines[idx].Hostnames {
				if hn == hostname {
					h.RemoveHostsFileLineByRow(idx)
					continue
				}
			}
		}
	}
}

// LookupByHostname check if the given fqdn exists.
// if yes, it returns the index of the address and the associated address.
// error is not nil if something goes wrong
func (h *HostsFile) LookupByHostname(hostname string) (int, net.IP, error) {
	for idx, hfl := range h.HostsFileLines {
		for _, hn := range hfl.Hostnames {
			if hn == hostname {
				return idx, h.HostsFileLines[idx].Address, nil
			}
		}
	}

	return -1, nil, ErrHostnameNotFound
}

// AddHostsFileLineRaw add the given ip/fqdn/comment pair
// this is different from AddHostFileLine because it does not take care of duplicates
// this just append the new entry to the hosts file
func (h *HostsFile) AddHostsFileLineRaw(ipRaw, fqdnRaw, comment string) (int, *HostsFileLine, error) {
	// hostname to lowercase
	hostname := strings.ToLower(fqdnRaw)
	// parse ip to net.IP
	ip := net.ParseIP(ipRaw)

	// if we have a valid IP
	if ip != nil {
		// create a new hosts line
		hfl := HostsFileLine{
			Type:        LineTypeAddress,
			Address:     ip,
			Hostnames:   []string{hostname},
			Comment:     comment,
			IsCommented: false,
		}

		// append to hosts
		h.HostsFileLines = append(h.HostsFileLines, hfl)

		// get index
		idx := len(h.HostsFileLines) - 1

		// return created entry
		return idx, &h.HostsFileLines[idx], nil
	}

	// return error
	return -1, nil, ErrCannotParseIPAddress(ipRaw)
}

// AddHostsFileLine add the given ip/fqdn/comment pair, cleanup is done for previous entry.
// it returns the index of the edited (created) line and a pointer to the hostsfileline object.
// error is not nil if something goes wrong
func (h *HostsFile) AddHostsFileLine(ipRaw, fqdnRaw, comment string) (int, *HostsFileLine, error) {
	// hostname to lowercase
	hostname := strings.ToLower(fqdnRaw)
	// parse ip to net.IP
	ip := net.ParseIP(ipRaw)

	// if we have a valid IP
	if ip != nil {
		//check if we alredy have the fqdn
		if idx, addr, err := h.LookupByHostname(hostname); err == nil {
			//if actual ip is the same as the given one, we are done
			if net.IP.Equal(addr, ip) {
				// handle comment
				if comment != "" {
					// just replace the current comment with the new one
					h.HostsFileLines[idx].Comment = comment
				}
				return idx, &h.HostsFileLines[idx], nil
			}

			//if address is different, we need to remove the hostname from the previous entry
			for hostIdx, hn := range h.HostsFileLines[idx].Hostnames {
				if hn == hostname {
					if len(h.HostsFileLines[idx].Hostnames) > 1 {
						h.Lock()
						h.HostsFileLines[idx].Hostnames = append(h.HostsFileLines[idx].Hostnames[:hostIdx], h.HostsFileLines[idx].Hostnames[hostIdx+1:]...)
						h.Unlock()
					}

					//remove the line if there are no more hostnames (other than the actual one)
					if len(h.HostsFileLines[idx].Hostnames) <= 1 {
						h.RemoveHostsFileLineByRow(idx)
					}
				}
			}
		}

		//if we alredy have the address, just add the hostname to that line
		for idx, hfl := range h.HostsFileLines {
			if net.IP.Equal(hfl.Address, ip) {
				h.Lock()
				h.HostsFileLines[idx].Hostnames = append(h.HostsFileLines[idx].Hostnames, hostname)
				h.Unlock()

				// handle comment
				if comment != "" {
					// just replace the current comment with the new one
					h.HostsFileLines[idx].Comment = comment
				}

				// return edited entry
				return idx, &h.HostsFileLines[idx], nil
			}
		}

		// at this point we need to create new host line
		hfl := HostsFileLine{
			Type:        LineTypeAddress,
			Address:     ip,
			Hostnames:   []string{hostname},
			Comment:     comment,
			IsCommented: false,
		}

		// generate raw version of the line
		hfl.Raw = lineFormatter(hfl)

		// append to hosts
		h.HostsFileLines = append(h.HostsFileLines, hfl)

		// get index
		idx := len(h.HostsFileLines) - 1

		// return created entry
		return idx, &h.HostsFileLines[idx], nil
	}

	// return error
	return -1, nil, ErrCannotParseIPAddress(ipRaw)
}

// AddCommentFileLine adds a new line of type comment with the given comment.
// it returns the index of the edited (created) line and a pointer to the hostsfileline object.
// error is not nil if something goes wrong
func (h *HostsFile) AddCommentFileLine(comment string) (int, *HostsFileLine, error) {
	h.Lock()
	defer h.Unlock()

	hfl := HostsFileLine{
		Type:    LineTypeComment,
		Raw:     "# " + comment,
		Comment: comment,
	}

	hfl.Raw = lineFormatter(hfl)

	h.HostsFileLines = append(h.HostsFileLines, hfl)
	idx := len(h.HostsFileLines) - 1
	return idx, &h.HostsFileLines[idx], nil
}

// AddEmptyFileLine adds a new line of type empty.
// it returns the index of the edited (created) line and a pointer to the hostsfileline object.
// error is not nil if something goes wrong
func (h *HostsFile) AddEmptyFileLine() (int, *HostsFileLine, error) {
	h.Lock()
	defer h.Unlock()

	hfl := HostsFileLine{
		Type: LineTypeEmpty,
		Raw:  "",
	}

	h.HostsFileLines = append(h.HostsFileLines, hfl)
	idx := len(h.HostsFileLines) - 1
	return idx, &h.HostsFileLines[idx], nil
}

// CommentHostsFileLineByRow set the IsCommented bit for the given row to true
func (h *HostsFile) CommentHostsFileLineByRow(row int) error {
	h.Lock()
	defer h.Unlock()

	if len(h.HostsFileLines) > row {
		if h.HostsFileLines[row].Type == LineTypeAddress {
			if !h.HostsFileLines[row].IsCommented {
				h.HostsFileLines[row].IsCommented = true

				h.HostsFileLines[row].Raw = h.RenderHostsFileLine(row)
				return nil
			}

			return ErrAlredyCommentedLine
		}

		return ErrNotAnAddressLine
	}

	return ErrUnknown
}

// CommentHostsFileLineByIP set the IsCommented bit for the given address to true
func (h *HostsFile) CommentHostsFileLineByIP(ip net.IP) error {
	h.Lock()
	defer h.Unlock()

	for idx := range h.HostsFileLines {
		if net.IP.Equal(ip, h.HostsFileLines[idx].Address) {
			if !h.HostsFileLines[idx].IsCommented {
				h.HostsFileLines[idx].IsCommented = true

				h.HostsFileLines[idx].Raw = h.RenderHostsFileLine(idx)
				return nil
			}

			return ErrAlredyCommentedLine
		}
	}

	return ErrAddressNotFound
}

func (h *HostsFile) CommentHostsFileLinesByIP(ip net.IP) {
	h.Lock()
	defer h.Unlock()

	for idx := range h.HostsFileLines {
		if net.IP.Equal(ip, h.HostsFileLines[idx].Address) {
			if !h.HostsFileLines[idx].IsCommented {
				h.HostsFileLines[idx].IsCommented = true

				h.HostsFileLines[idx].Raw = h.RenderHostsFileLine(idx)
			}
		}
	}
}

// CommentHostsFileLineByAddress set the IsCommented bit for the given address as string to true
func (h *HostsFile) CommentHostsFileLineByAddress(address string) error {
	ip := net.ParseIP(address)

	return h.CommentHostsFileLineByIP(ip)
}

func (h *HostsFile) CommentHostsFileLinesByAddress(address string) {
	ip := net.ParseIP(address)
	h.CommentHostsFileLinesByIP(ip)
}

// CommentHostsFileLineByHostname set the IsCommented bit for the given hostname to true
func (h *HostsFile) CommentHostsFileLineByHostname(hostname string) error {
	h.Lock()
	defer h.Unlock()

	for idx := range h.HostsFileLines {
		for _, hn := range h.HostsFileLines[idx].Hostnames {
			if hn == hostname {
				if !h.HostsFileLines[idx].IsCommented {
					h.HostsFileLines[idx].IsCommented = true

					h.HostsFileLines[idx].Raw = h.RenderHostsFileLine(idx)
					return nil
				}

				return ErrAlredyCommentedLine
			}
		}
	}

	return ErrHostnameNotFound
}

func (h *HostsFile) CommentHostsFileLinesByHostname(hostname string) {
	h.Lock()
	defer h.Unlock()

	for idx := range h.HostsFileLines {
		for _, hn := range h.HostsFileLines[idx].Hostnames {
			if hn == hostname {
				if !h.HostsFileLines[idx].IsCommented {
					h.HostsFileLines[idx].IsCommented = true

					h.HostsFileLines[idx].Raw = h.RenderHostsFileLine(idx)
				}
			}
		}
	}
}

func (h *HostsFile) CommentHostsFileLinesByHostnameAsRegexp(hostname string) {
	h.Lock()
	defer h.Unlock()

	reg := regexp.MustCompile(hostname)

	for idx := range h.HostsFileLines {
		for _, hn := range h.HostsFileLines[idx].Hostnames {
			if reg.MatchString(hn) {
				if !h.HostsFileLines[idx].IsCommented {
					h.HostsFileLines[idx].IsCommented = true

					h.HostsFileLines[idx].Raw = h.RenderHostsFileLine(idx)
					continue
				}
			}
		}
	}
}

// UncommentHostsFileLineByRow set the IsCommented bit for the given row to false
func (h *HostsFile) UncommentHostsFileLineByRow(row int) error {
	h.Lock()
	defer h.Unlock()

	if len(h.HostsFileLines) > row {
		if h.HostsFileLines[row].Type == LineTypeAddress {
			if h.HostsFileLines[row].IsCommented {
				h.HostsFileLines[row].IsCommented = false

				h.HostsFileLines[row].Raw = h.RenderHostsFileLine(row)
				return nil
			}

			return ErrAlredyUncommentedLine
		}

		return ErrNotAnAddressLine
	}

	return ErrUnknown
}

// UncommentHostsFileLineByIP set the IsCommented bit for the given address to false
func (h *HostsFile) UncommentHostsFileLineByIP(ip net.IP) error {
	h.Lock()
	defer h.Unlock()

	for idx, hfl := range h.HostsFileLines {
		if net.IP.Equal(ip, hfl.Address) {
			if h.HostsFileLines[idx].IsCommented {
				h.HostsFileLines[idx].IsCommented = false

				h.HostsFileLines[idx].Raw = h.RenderHostsFileLine(idx)
				return nil
			}

			return ErrAlredyUncommentedLine
		}
	}

	return ErrNotAnAddressLine
}

func (h *HostsFile) UncommentHostsFileLinesByIP(ip net.IP) {
	h.Lock()
	defer h.Unlock()

	for idx := range h.HostsFileLines {
		if net.IP.Equal(ip, h.HostsFileLines[idx].Address) {
			if h.HostsFileLines[idx].IsCommented {
				h.HostsFileLines[idx].IsCommented = false

				h.HostsFileLines[idx].Raw = h.RenderHostsFileLine(idx)
			}
		}
	}
}

// UncommentHostsFileLineByAddress set the IsCommented bit for the given address as string to false
func (h *HostsFile) UncommentHostsFileLineByAddress(address string) error {
	ip := net.ParseIP(address)

	return h.UncommentHostsFileLineByIP(ip)
}

func (h *HostsFile) UncommentHostsFileLinesByAddress(address string) {
	ip := net.ParseIP(address)
	h.UncommentHostsFileLinesByIP(ip)
}

// UncommentHostsFileLineByHostname set the IsCommented bit for the given hostname to false
func (h *HostsFile) UncommentHostsFileLineByHostname(hostname string) error {
	h.Lock()
	defer h.Unlock()

	for idx := range h.HostsFileLines {
		for _, hn := range h.HostsFileLines[idx].Hostnames {
			if hn == hostname {
				if h.HostsFileLines[idx].IsCommented {
					h.HostsFileLines[idx].IsCommented = false

					h.HostsFileLines[idx].Raw = h.RenderHostsFileLine(idx)
					return nil
				}

				return ErrAlredyUncommentedLine
			}
		}
	}

	return ErrHostnameNotFound
}

func (h *HostsFile) UncommentHostsFileLinesByHostname(hostname string) {
	h.Lock()
	defer h.Unlock()

	for idx := range h.HostsFileLines {
		for _, hn := range h.HostsFileLines[idx].Hostnames {
			if hn == hostname {
				if h.HostsFileLines[idx].IsCommented {
					h.HostsFileLines[idx].IsCommented = false

					h.HostsFileLines[idx].Raw = h.RenderHostsFileLine(idx)
				}
			}
		}
	}
}

func (h *HostsFile) UncommentHostsFileLinesByHostnameAsRegexp(hostname string) {
	h.Lock()
	defer h.Unlock()

	reg := regexp.MustCompile(hostname)

	for idx := range h.HostsFileLines {
		for _, hn := range h.HostsFileLines[idx].Hostnames {
			if reg.MatchString(hn) {
				if h.HostsFileLines[idx].IsCommented {
					h.HostsFileLines[idx].IsCommented = false

					h.HostsFileLines[idx].Raw = h.RenderHostsFileLine(idx)
					continue
				}
			}
		}
	}
}

// RenderHostsFile render and returns the hosts file with the lineFormatter() routine
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

// RenderHostsFileLine render and returns the given hosts line with the lineFormatter() routine
func (h *HostsFile) RenderHostsFileLine(row int) string {
	// iterate to find the row to render
	if len(h.HostsFileLines) > row {
		return lineFormatter(h.HostsFileLines[row])
	}

	return ""
}

// SaveHostsFile write hosts file to configured path.
// error is not nil if something goes wrong
func (h *HostsFile) SaveHostsFile() error {
	return h.SaveHostsFileAs(h.Path)
}

// SaveHostsFileAs write hosts file to the given path.
// error is not nil if something goes wrong
func (h *HostsFile) SaveHostsFileAs(path string) error {
	// render the file as a byte slice
	dataBytes := []byte(h.RenderHostsFile())

	// write file to disk
	err := os.WriteFile(path, dataBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}
