package tcp

import (
	"context"
	"net"
	"time"

	"github.com/lthibault/jitterbug"
	log "github.com/sirupsen/logrus"

	"github.com/zinfra/srv-announcer/checker/healthchecks"
)

// ensure TCPHealthcheck implements IHealthcheck
var _ healthchecks.IHealthcheck = &Healthcheck{}

// Healthcheck checks whether it's able to connect to a given target
// via TCP
type Healthcheck struct {
	target         string
	connectTimeout time.Duration
	checkInterval  time.Duration
}

// NewHealthcheck creates a new Healthcheck
func NewHealthcheck(target string, connectTimeout time.Duration, checkInterval time.Duration) *Healthcheck {
	return &Healthcheck{
		target:         target,
		connectTimeout: connectTimeout,
		checkInterval:  checkInterval,
	}
}

// Run regularily tries to connect to the target via TCP,
// sends true if it was able to, false otherwise
func (hc *Healthcheck) Run(ctx context.Context, healthyChan chan<- bool) {
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
