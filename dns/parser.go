package dns

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

//ParseSRV parses a SRV record value into a *net.SRV
func ParseSRV(rrText string) (*net.SRV, error) {
	fields := strings.Fields(rrText)
	if len(fields) != 4 {
		return nil, fmt.Errorf("number of fields != 4 in %s", rrText)
	}

	// parse fields
	var err error

	priority, err := strconv.ParseUint(fields[0], 10, 16)
	weight, err := strconv.ParseUint(fields[1], 10, 16)
	port, err := strconv.ParseUint(fields[2], 10, 16)
	target := fields[3]

	if err != nil {
		return nil, err
	}

	return &net.SRV{
		Priority: uint16(priority),
		Weight:   uint16(weight),
		Port:     uint16(port),
		Target:   target,
	}, nil
}
