package dummy

import (
	"net"

	log "github.com/sirupsen/logrus"
	"github.com/zinfra/srv-announcer/dns"
)

// ensure SrvManager implements ISRVManager
var _ dns.ISRVManager = &SrvManager{}

// SrvManager implements the ISrvManager interface, but doesn't talk to any real service
// instead, it just logs and exposes its struct.
type SrvManager struct {
	SrvRecordSet []net.SRV
}

// Add adds the record to the record set if it doesn't already exist
func (s *SrvManager) Add(srv *net.SRV) error {
	log.Info("add called")
	// if that record is in the set already, we're done
	for _, aSrv := range s.SrvRecordSet {
		if aSrv.Port == srv.Port && aSrv.Priority == srv.Priority &&
			aSrv.Target == srv.Target && aSrv.Weight == aSrv.Weight {
			log.Debugf("Record %+v already exists, doing nothing", aSrv)
			return nil
		}
	}
	// else append it
	log.Infof("Adding record %+v to record set", srv)
	s.SrvRecordSet = append(s.SrvRecordSet, *srv)
	return nil
}

// Remove removes the record to the record set if it exists
func (s *SrvManager) Remove(srv *net.SRV) error {
	log.Info("remove called")
	newRecordSet := make([]net.SRV, 0)
	for _, aSrv := range s.SrvRecordSet {
		// if that record is in the set already, remove it, else copy it over
		if aSrv.Port == srv.Port && aSrv.Priority == srv.Priority &&
			aSrv.Target == srv.Target && aSrv.Weight == aSrv.Weight {
			log.Debugf("Removing record %+v from list", aSrv)
		} else {
			newRecordSet = append(newRecordSet, aSrv)
		}
	}
	s.SrvRecordSet = newRecordSet
	return nil
}
