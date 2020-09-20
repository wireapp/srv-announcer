package checker

import (
	"context"

	log "github.com/sirupsen/logrus"
	config "github.com/zinfra/srv-announcer/config"
	dns "github.com/zinfra/srv-announcer/dns"
)

// IHealthcheck describes the basic healthchecker interface
type IHealthcheck interface {
	// Run will run the healthchecker-specific check
	// in healthchecker-specific intervals.
	// They will regularily push their health status
	// to the healthyChan.
	Run(ctx context.Context, healthyChan chan<- bool)
}

// Run runs the healthchecks and updates the SRV record
func Run(ctx context.Context, config *config.Config, srvManager dns.ISRVManager) error {
	var err error

	// start a TCP Health checker
	tcpHc := NewTCPHealthcheck(config.CheckTarget, config.CheckTimeout, config.CheckInterval)
	tcpHcC := make(chan bool, 1)
	go tcpHc.Run(ctx, tcpHcC)

	for {
		select {
		case <-ctx.Done():
			srvManager.Remove(&config.SRVRecord)
			return nil
		case isReachable := <-tcpHcC:
			log.Infof("got data on healthyC: %t", isReachable)
			if isReachable {
				err = srvManager.Add(&config.SRVRecord)
			} else {
				err = srvManager.Remove(&config.SRVRecord)
			}
			if err != nil {
				// only log the error here, don't exit the check loop.
				// It might be a networking blip - we usually want to
				// keep doing health checks.
				log.Errorf("%s", err.Error())
					err = srvManager.Add(srvRecord)
				} else {
					err = srvManager.Remove(srvRecord)
				}
				if err != nil {
					// only log the error here, don't exit the check loop.
					// It might be a networking blip - we usually want to
					// keep doing health checks.
					log.Errorf("%s", err.Error())
				}
			}
		}
	}
	return nil
}
