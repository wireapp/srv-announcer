package config

import (
	"net"
	"time"
)

// Config represents the configuration of the service
type Config struct {
	DryRun        bool
	ZoneName      string
	SRVRecordName string
	TTL           uint16
	SRVRecord     *net.SRV
	CheckTarget   string
	CheckInterval time.Duration
	CheckTimeout  time.Duration
}
