package checker

import (
	"context"
	"net"
	"time"

	"github.com/lthibault/jitterbug"
	log "github.com/sirupsen/logrus"
)

// ensure TCPHealthcheck implements IHealthcheck
var _ IHealthcheck = &TCPHealthcheck{}

// TCPHealthcheck checks whether it's able to connect to a given target
// via TCP
type TCPHealthcheck struct {
	target         string
	connectTimeout time.Duration
	checkInterval  time.Duration
}

// NewTCPHealthcheck creates a new TCPHealthcheck
func NewTCPHealthcheck(target string, connectTimeout time.Duration, checkInterval time.Duration) *TCPHealthcheck {
	return &TCPHealthcheck{
		target:         target,
		connectTimeout: connectTimeout,
		checkInterval:  checkInterval,
	}
}

// Run regularily tries to connect to the target via TCP,
// sends true if it was able to, false otherwise
func (hc *TCPHealthcheck) Run(ctx context.Context, healthyChan chan<- bool) {
	// jitter around a 10th of the configured interval
	t := jitterbug.New(hc.checkInterval, &jitterbug.Norm{Stdev: hc.checkInterval / 10})

	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			conn, err := net.DialTimeout("tcp", hc.target, hc.connectTimeout)
			if err != nil {
				log.Infof("%s is unreachable", hc.target)
				healthyChan <- false
				continue
			}
			defer conn.Close()
			log.Infof("%s is reachable", hc.target)
			healthyChan <- true
		}
	}
}
