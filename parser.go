package libhosty

import (
	"io/ioutil"
	"net"
	"strings"
)

//ParseHostsFile parse a hosts file from the given location.
// error is not nil if something goes wrong
func ParseHostsFile(path string) ([]HostsFileLine, error) {
	byteData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return parser(byteData)
}

//ParseHostsFileAsString parse a hosts file from a given string.
// error is not nil if something goes wrong
func ParseHostsFileAsString(stringData string) ([]HostsFileLine, error) {
	bytesData := []byte(stringData)
	return parser(bytesData)
}

func parser(bytesData []byte) ([]HostsFileLine, error) {
	byteDataNormalized := strings.Replace(string(bytesData), "\r\n", "\n", -1)
	fileLines := strings.Split(byteDataNormalized, "\n")
	hostsFileLines := make([]HostsFileLine, len(fileLines))

	// trim leading an trailing whitespace
	for i, l := range fileLines {
		curLine := &hostsFileLines[i]
		curLine.LineNumber = i
		curLine.Raw = l

		// trim line
		curLine.Trimed = strings.TrimSpace(l)

		// check if it's an empty line
		if curLine.Trimed == "" {
			curLine.LineType = EMPTY
			continue
		}

		// check if line starts with a #
		if strings.HasPrefix(curLine.Trimed, "#") {
			// this can be a comment or a commented host line
			// so remove the 1st char (#), trim spaces
			// and try to parse the line as a host line
			noCommentLine := strings.TrimPrefix(curLine.Trimed, "#")
			tmpParts := strings.Fields(strings.TrimSpace(noCommentLine))
			address := net.ParseIP(tmpParts[0])

			// if address is nil this line is definetly a comment
			if address == nil {
				curLine.LineType = COMMENT
				continue
			}

			// otherwise it is a commented line so let's try to parse it as a normal line
			curLine.IsCommented = true
			curLine.Trimed = noCommentLine
		}

		// not a comment or empty line so try to parse it
		// check if it contains a comment
		curLineSplit := strings.SplitN(curLine.Trimed, "#", 2)
		if len(curLineSplit) > 1 {
			curLine.Comment = curLineSplit[1]
		}

		curLine.Trimed = curLineSplit[0]
		curLine.Parts = strings.Fields(curLine.Trimed)

		if len(curLine.Parts) > 1 {
			curLine.LineType = ADDRESS
			curLine.Address = net.ParseIP(curLine.Parts[0])
			// lower case all
			for _, p := range curLine.Parts[1:] {
				curLine.Hostnames = append(curLine.Hostnames, strings.ToLower(p))
			}

			continue
		}

		// if we can't figure out what this line is mark it as unknown
		curLine.LineType = UNKNOWN
	}

	return hostsFileLines, nil
}
