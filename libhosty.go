// Package libhosty is a pure golang library to manipulate the hosts file
package libhosty

import (
	"net"
	"regexp"
	"strings"
	"sync"

	"golang.org/x/exp/slices"
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

func (lt LineType) String() string {
	switch lt {
	case LineTypeEmpty:
		return "line-type-empty"
	case LineTypeComment:
		return "line-type-comment"
	case LineTypeAddress:
		return "line-type-address"
	default:
		return "line-type-unknown"
	}
}

// HostsFileLine holds hosts file lines data
type HostsFileLine struct {
	//LineType defines the line type
	Type LineType

	//Address is a net.IP representation of the address
	Address net.IP

	//Hostnames is a slice of hostnames for the relative IP
	Hostnames []string

	//Raw is the raw representation of the line, as it is in the hosts file
	Raw string

	//Comment is the comment part of the line (if present in an ADDRESS line)
	Comment string

	//IsCommented to know if the current ADDRESS line is commented out (starts with '#')
	IsCommented bool
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
func Init() (*HostsFile, error) {
	// get hosts file default path
	fpath := GetOSHostsFilePath()

	// parse hosts file from default path
	hostsFileLines, err := ParseHostsFile(fpath)
	if err != nil {
		return nil, err
	}

	// allocate a new HostsFile object
	hf := &HostsFile{
		// use default configuration
		Path: fpath,

		// allocate a new slice of HostsFileLine objects
		HostsFileLines: hostsFileLines,
	}

	//return HostsFile
	return hf, nil
}

func InitFromCustomPath(path string) (*HostsFile, error) {
	// parse hosts file from default path
	hostsFileLines, err := ParseHostsFile(path)
	if err != nil {
		return nil, err
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

func InitFromString(lines string) (*HostsFile, error) {
	// parse inline hosts file
	hostsFileLines, err := ParseHostsFileFromString(lines)
	if err != nil {
		return nil, err
	}

	// allocate a new HostsFile object
	hf := &HostsFile{
		// use default configuration
		Path: "",

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

// GetHostsFileLinesByIP returns every line that maches a given IP
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

// GetHostsFileLinesByAddress returns every line that maches a given IP as String
func (h *HostsFile) GetHostsFileLinesByAddress(address string) []*HostsFileLine {
	ip := net.ParseIP(address)
	return h.GetHostsFileLinesByIP(ip)
}

// GetHostsFileLinesByHostname returns every line that maches a given Hostname
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

// GetHostsFileLinesByRegexp returns every line that maches a given regexp
func (h *HostsFile) GetHostsFileLinesByRegexp(pattern string) []HostsFileLine {
	hfl := make([]HostsFileLine, 0)

	reg := regexp.MustCompile(pattern)

	for idx := range h.HostsFileLines {
		for _, hn := range h.HostsFileLines[idx].Hostnames {
			if reg.MatchString(hn) {
				hfl = append(hfl, h.HostsFileLines[idx])
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

// RemoveHostFileLinesByIP remove every line that matches a given IP
func (h *HostsFile) RemoveHostsFileLinesByIP(ip net.IP) {
	for idx := len(h.HostsFileLines) - 1; idx >= 0; idx-- {
		if net.IP.Equal(ip, h.HostsFileLines[idx].Address) {
			h.RemoveHostsFileLineByRow(idx)
		}
	}
}

// RemoveHostFileLinesByAddress remove every line that matches a given IP as String
func (h *HostsFile) RemoveHostsFileLinesByAddress(address string) {
	ip := net.ParseIP(address)

	h.RemoveHostsFileLinesByIP(ip)
}

// RemoveHostFileLinesByHostname remove every line that matches a given Hostname
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

// RemoveHostFileLinesByRegexp remove every line that matches a given regexp
func (h *HostsFile) RemoveHostsFileLinesByRegexp(pattern string) {
	reg := regexp.MustCompile(pattern)

	for idx := len(h.HostsFileLines) - 1; idx >= 0; idx-- {
		for _, hn := range h.HostsFileLines[idx].Hostnames {
			if reg.MatchString(hn) {
				h.RemoveHostsFileLineByRow(idx)
				continue
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

	// get index
	idx := len(h.HostsFileLines)

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

		// return created entry
		return idx, &hfl, nil
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
		// check if we alredy have the fqdn
		if idx, addr, err := h.LookupByHostname(hostname); err == nil {
			// if actual ip is the same as the given one, we are done
			if net.IP.Equal(addr, ip) {
				// handle comment
				if comment != "" {
					// just replace the current comment with the new one
					h.HostsFileLines[idx].Comment = comment
				}
				return idx, &h.HostsFileLines[idx], nil
			}

			// if address is different, we need to remove the hostname from the previous entry
			for hostIdx, hn := range h.HostsFileLines[idx].Hostnames {
				// if hostnames matches
				if hn == hostname {
					// remove hostname if there's at least one another hostname
					if len(h.HostsFileLines[idx].Hostnames) > 1 {
						h.Lock()
						h.HostsFileLines[idx].Hostnames = append(h.HostsFileLines[idx].Hostnames[:hostIdx], h.HostsFileLines[idx].Hostnames[hostIdx+1:]...)
						h.Unlock()
					} else {
						// remove the whole line
						h.RemoveHostsFileLineByRow(idx)
					}
				}
			}
		}

		// index saves the last matching index for the next for loop
		// in this way, if we find a matching line (same IP)
		// but needs to create a new line (hostname limit exceeded)
		// it would be nice to place the new line next to the existing one
		index := -1

		// we don't have the fqdn
		// if we already have the address, just add the hostname to that line
		for idx, hfl := range h.HostsFileLines {
			if net.IP.Equal(hfl.Address, ip) {
				// if we already have 6 hostnames, just continue
				// we'll either find another matching line
				// or we'll end up creating a new line
				if len(h.HostsFileLines[idx].Hostnames) >= 6 {
					// save index
					index = idx
					continue
				}

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

		var idx int
		// we found a matching line
		if index != -1 {
			// so set the index next to the matching line
			idx = index + 1
		} else {
			idx = len(h.HostsFileLines)
		}

		// at this point we need to create new host line
		hfl := HostsFileLine{
			Type:        LineTypeAddress,
			Address:     ip,
			Hostnames:   []string{hostname},
			Raw:         "",
			Comment:     comment,
			IsCommented: false,
		}

		// generate raw version of the line
		hfl.Raw = lineFormatter(hfl)

		// if we found a matching line (index != -1)
		// append the new line after the matched one
		if index != -1 {
			h.HostsFileLines = slices.Insert(h.HostsFileLines, idx, hfl)
		} else {
			// else append to hosts
			h.HostsFileLines = append(h.HostsFileLines, hfl)
		}

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

	idx := len(h.HostsFileLines)

	hfl := HostsFileLine{
		Type:        LineTypeComment,
		Address:     []byte{},
		Hostnames:   []string{},
		Raw:         "# " + comment,
		Comment:     comment,
		IsCommented: false,
	}

	hfl.Raw = lineFormatter(hfl)

	h.HostsFileLines = append(h.HostsFileLines, hfl)
	return idx, &h.HostsFileLines[idx], nil
}

// AddEmptyFileLine adds a new line of type empty.
// it returns the index of the edited (created) line and a pointer to the hostsfileline object.
// error is not nil if something goes wrong
func (h *HostsFile) AddEmptyFileLine() (int, *HostsFileLine, error) {
	h.Lock()
	defer h.Unlock()

	idx := len(h.HostsFileLines)

	hfl := HostsFileLine{
		Type:        LineTypeEmpty,
		Address:     []byte{},
		Hostnames:   []string{},
		Raw:         "",
		Comment:     "",
		IsCommented: false,
	}

	h.HostsFileLines = append(h.HostsFileLines, hfl)
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

// CommentHostsFileLinesByIP set IsCommented to true on every line that matches a given net.IP
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

// CommentHostsFileLinesByAddress set IsCommented to true on every line that matches a given net.IP
func (h *HostsFile) CommentHostsFileLinesByAddress(address string) {
	ip := net.ParseIP(address)

	h.CommentHostsFileLinesByIP(ip)
}

// CommentHostsFileLinesByHostname set IsCommented to true on every line that matches a given Hostname
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

// CommentHostsFileLinesByRegexp set IsCommented to true on every line that matches a given regexp
func (h *HostsFile) CommentHostsFileLinesByRegexp(pattern string) {
	h.Lock()
	defer h.Unlock()

	reg := regexp.MustCompile(pattern)

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

// UncommentHostsFileLinesByIP set IsCommented to false for every line that matches a given net.IP
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

// UncommentHostsFileLinesByAddress set IsCommented to false for every line that matches a given net.IP as String
func (h *HostsFile) UncommentHostsFileLinesByAddress(address string) {
	ip := net.ParseIP(address)
	h.UncommentHostsFileLinesByIP(ip)
}

// UncommentHostsFileLinesByHostname set IsCommented to false for every line that matches a given Hostname
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

// UncommentHostsFileLinesByRegexp set IsCommented to false for every line that matches a given regexp
func (h *HostsFile) UncommentHostsFileLinesByRegexp(pattern string) {
	h.Lock()
	defer h.Unlock()

	reg := regexp.MustCompile(pattern)

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
