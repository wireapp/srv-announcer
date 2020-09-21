package checker

import (
	"context"
	"net"

	log "github.com/sirupsen/logrus"
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

// Run runs a healthcheck and updates the SRV record whenever its status changes
func Run(ctx context.Context, healthcheck IHealthcheck, srvRecord *net.SRV, srvManager dns.ISRVManager) error {
	var err error
	var healthyC chan bool

	// initialize the healthcheck
	healthyC = make(chan bool, 1)
	go healthcheck.Run(ctx, healthyC)

	for {
		select {
		case <-ctx.Done():
			srvManager.Remove(srvRecord)
			return nil
		case isReachable := <-healthyC:
			log.Infof("got data on healthyC: %t", isReachable)
			if isReachable {
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
