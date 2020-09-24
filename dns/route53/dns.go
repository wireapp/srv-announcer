package route53

import (
	"net"

	"github.com/aws/aws-sdk-go/aws"
	route53Client "github.com/aws/aws-sdk-go/service/route53"
	log "github.com/sirupsen/logrus"
	dns "github.com/zinfra/srv-announcer/dns"
)

// SRVManager manages an SRV record inside Route53
type SRVManager struct {
	client       *Client
	hostedZoneID string
	recordName   string
	ttl          uint16
}

// ensure SRVManager implements dns.ISRVManager
var _ dns.ISRVManager = &SRVManager{}

// NewSRVManager initializes an SRV Manager by its zone ID and record name
func NewSRVManager(client *Client, hostedZoneID string, recordName string, ttl uint16) *SRVManager {
	return &SRVManager{
		client:       client,
		hostedZoneID: hostedZoneID,
		recordName:   recordName,
		ttl:          ttl,
	}
}

// edit provides both add and removal capabilities
func (s *SRVManager) edit(add bool, srv *net.SRV) error {
	log.Debugf("Looking up SRV resource record set for %s", s.recordName)
	currentResourceRecordSet, err := s.client.GetResourceRecordSetByName(s.hostedZoneID, s.recordName, "SRV")
	if err != nil {
		return err
	}

	var currentResourceRecords []*route53Client.ResourceRecord

	if currentResourceRecordSet == nil {
		log.Debugf("Resource Record set for %s doesn't exist", s.recordName)
		currentResourceRecords = []*route53Client.ResourceRecord{}
	} else {
		log.Debugf("Resource Record set for %s already exists", s.recordName)
		currentResourceRecords = currentResourceRecordSet.ResourceRecords
	}

	newResourceRecords := editResourceRecords(add, currentResourceRecords, srv)

	if !resourceRecordsDiffer(currentResourceRecords, newResourceRecords) {
		log.Debugf("Skipped update, no change needed.")
	} else {
		resourceRecordSet := &route53Client.ResourceRecordSet{
			TTL:  aws.Int64(int64(s.ttl)),
			Name: aws.String(s.recordName),
			Type: aws.String("SRV"),
		}

		var action string

		if len(newResourceRecords) > 0 {
			action = route53Client.ChangeActionUpsert
			resourceRecordSet.ResourceRecords = newResourceRecords
		} else {
			action = route53Client.ChangeActionDelete
			// NOTE: in order to delete the right record not only Name has to match
			//       but also the Value of that record
			resourceRecordSet.ResourceRecords = currentResourceRecords
		}

		log.Infof("%s record %+v", action, srv)
		log.Tracef("Full record set: %+v", resourceRecordSet.ResourceRecords)
		_, err := s.client.ChangeRecord(s.hostedZoneID, action, resourceRecordSet)
		if err != nil {
			return err
		}
	}

	return nil
}

// Add adds a SRV record into the record set
func (s *SRVManager) Add(srv *net.SRV) error {
	return s.edit(true, srv)
}

// Remove removes a SRV record from the record set
func (s *SRVManager) Remove(srv *net.SRV) error {
	return s.edit(false, srv)
}
