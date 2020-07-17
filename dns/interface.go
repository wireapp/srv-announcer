package dns

import (
	"net"
)

// ISRVManager provides the interface to manage an SRV record
type ISRVManager interface {
	// Add record to rrset. If the same already exists, it's a no-op
	Add(*net.SRV) error

	// Remove from an rrset
	Remove(*net.SRV) error
}
