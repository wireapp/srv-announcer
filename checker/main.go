package checker

import (
	"net"
	"time"

	"github.com/lthibault/jitterbug"
	log "github.com/sirupsen/logrus"
	config "github.com/zinfra/srv-announcer/config"
	dns "github.com/zinfra/srv-announcer/dns"
)

func tcpReachable(target string, checkTimeout time.Duration) bool {
	timeout := checkTimeout
	conn, err := net.DialTimeout("tcp", target, timeout)

	if err != nil {
		log.Infof("%s is unreachable", target)
		return false
	}
	defer conn.Close()
	log.Infof("%s is reachable", target)
	return true
}

// Run runs the healthchecks and updates the SRV record
func Run(config *config.Config, srvManager dns.ISRVManager) error {
	// jitter around a 10th of the configured interval
	jitter := &jitterbug.Norm{Stdev: (config.CheckInterval / 10)}
	t := jitterbug.New(
		config.CheckInterval,
		jitter,
	)

	// initially wait the jitter time
	time.Sleep(t.Jitter.Jitter(config.CheckInterval / 10))

	for range t.C {
		var err error;
		if tcpReachable(config.CheckTarget, config.CheckTimeout) {
			err = srvManager.Add(&config.SRVRecord)
		} else {
			err = srvManager.Remove(&config.SRVRecord)
		}
		if err != nil {
			// only log the error here, don't exit the check loop.
			// It might be a networking blip - we usually want to
			// keep doing health checks.
			log.Errorf("%s", err.Error())
		}
	}
	return nil
}
