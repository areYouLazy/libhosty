package libhosty

import (
	"net"
	"os"
	"strings"
)

// ParseHostsFile parse a hosts file from the given location.
// error is not nil if something goes wrong
func ParseHostsFile(path string) ([]HostsFileLine, error) {
	byteData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return parser(byteData)
}

// ParseHostsFileFromString parse a hosts file from a given string.
// error is not nil if something goes wrong
func ParseHostsFileFromString(stringData string) ([]HostsFileLine, error) {
	bytesData := []byte(stringData)
	return parser(bytesData)
}

// parser, the line parser
func parser(bytesData []byte) ([]HostsFileLine, error) {
	// normalize input
	byteDataNormalized := strings.Replace(string(bytesData), "\r\n", "\n", -1)

	// split by line
	fileLines := strings.Split(byteDataNormalized, "\n")

	// init hostsFileLines buffer
	hostsFileLines := make([]HostsFileLine, 0)

	// iterate file lines
	for _, line := range fileLines {
		// instantiate a new HostsFileLine
		curLine := HostsFileLine{
			Type:        0,
			Address:     []byte{},
			Hostnames:   []string{},
			Raw:         "",
			Comment:     "",
			IsCommented: false,
		}

		// trim line (remove spaces after and before)
		rawLine := strings.TrimSpace(line)

		// save a raw version of the line, after only TrimSpace sanitization
		curLine.Raw = rawLine

		// check if it's an empty line
		if rawLine == "" {
			curLine.Type = LineTypeEmpty
			hostsFileLines = append(hostsFileLines, curLine)
			continue
		}

		// check if line starts with a #
		if strings.HasPrefix(rawLine, "#") {
			// this can be a comment or a commented host line
			// ensure to remove every # char at the beginning of the line
			// keep track of how many # has been removed
			hashCounter := 0
			for !strings.HasPrefix(rawLine, "#") {
				rawLine = strings.TrimPrefix(rawLine, "#")
				// also trim spaces to avoid "# #" situations
				rawLine = strings.TrimSpace(rawLine)
				// increment hash counter
				hashCounter++
			}

			// this can be a hashes, comment or commented hosts line
			// no parts == hashes
			// 1st part != net.IP == comment
			// 1st part == net.IP == hosts line
			rawLineParts := strings.Fields(rawLine)

			// nothing except hashes, comment line
			if len(rawLineParts) == 0 {
				curLine.Type = LineTypeComment
				hostsFileLines = append(hostsFileLines, curLine)
				continue
			}

			// try to parse 1st field as an ip address
			// if address is nil this line is a comment
			if address := net.ParseIP(rawLineParts[0]); address == nil {
				comment := rawLine

				// mark line as comment, normalize and save comment
				curLine.Type = LineTypeComment

				// since there can be more than one hash, remove those in excess
				for !strings.HasPrefix(rawLine, "#") {
					comment = strings.TrimPrefix(comment, "#")
					// also trim spaces to avoid "# #" situations
					rawLine = strings.TrimSpace(rawLine)
				}

				// save comment
				curLine.Comment = comment
				hostsFileLines = append(hostsFileLines, curLine)
				continue
			}

			// if address is not nil, this line is a commented address line
			// so let's try to parse it as a normal line
			curLine.IsCommented = true
		}

		// So this line is not a comment or empty line, try to parse it

		// check if it contains a comment
		// len == 1 == no comment
		// len > 1 == comment
		rawLineSplit := strings.SplitN(rawLine, "#", 2)

		// if we have a comment, trim spaces and save it
		if len(rawLineSplit) > 1 {
			curLine.Comment = strings.TrimSpace(rawLineSplit[1])
		}

		// split the effective line by spaces
		addressAndHostnames := strings.Fields(rawLineSplit[0])

		// we should have at least 2 fields, the address at [0]
		// and the 1st hostname at [1], other hostnames at [2:...]
		if len(addressAndHostnames) > 1 {
			// sanitize address
			rawAddress := strings.TrimSpace(addressAndHostnames[0])

			// parse address to ensure we have a valid address line
			if address := net.ParseIP(rawAddress); address != nil {
				// set linetype as address and save it
				curLine.Type = LineTypeAddress
				curLine.Address = address

				// parse and lower case all hostnames
				for _, hostname := range addressAndHostnames[1:] {
					// sanitize hostname
					rawHostname := strings.TrimSpace(hostname)

					// if there are 6 or less hostnames, just proceede
					if len(rawHostname) <= 6 {
						// add hostname to hostnames slice
						curLine.Hostnames = append(curLine.Hostnames, strings.ToLower(rawHostname))
						hostsFileLines = append(hostsFileLines, curLine)
					} else {
						// if more than 6 we need to generate 2 lines
						// save first 6 hostnames to current line
						// save other to new line

						// get a new line (copy of currentLine, except for Hostnames and Raw)
						newLine := HostsFileLine{
							Type:        curLine.Type,
							Address:     curLine.Address,
							Hostnames:   make([]string, 0),
							Raw:         "",
							Comment:     curLine.Comment,
							IsCommented: curLine.IsCommented,
						}

						// save first 6 hostnames
						curLine.Hostnames = append(curLine.Hostnames, strings.ToLower(rawHostname[:6]))
						// append current line
						hostsFileLines = append(hostsFileLines, curLine)

						// save other hostnames
						newLine.Hostnames = append(newLine.Hostnames, rawHostname[6:])
						// update Raw field
						newLine.Raw = lineFormatter(newLine)
						// append new line
						hostsFileLines = append(hostsFileLines, newLine)
					}
				}

				// we got a line, go on to the next one
				continue
			}
		}

		// if we can't figure out what this line is mark it as unknown
		curLine.Type = LineTypeUnknown
	}

	// normalize slice
	hostsFileLines = hostsFileLines[:]

	return hostsFileLines, nil
}
