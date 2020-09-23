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
	log.Infof("Looking up SRV resource record set for %s", s.recordName)
	resourceRecordSet, err := s.client.GetResourceRecordSetByName(s.hostedZoneID, s.recordName, "SRV")
	if err != nil {
		return err
	}

	if resourceRecordSet == nil {
		log.Infof("Resource Record set for %s didn't exist, will create", s.recordName)
		resourceRecordSet = &route53Client.ResourceRecordSet{
			TTL:  aws.Int64(int64(s.ttl)),
			Name: aws.String(s.recordName),
			Type: aws.String("SRV"),
		}
	}

	newResourceRecords := editResourceRecords(add, resourceRecordSet.ResourceRecords, srv)

	if !resourceRecordsDiffer(resourceRecordSet.ResourceRecords, newResourceRecords) {
		log.Infof("skipped update, no change needed.")
	} else {
		resourceRecordSet.ResourceRecords = newResourceRecords

		_, err := s.client.ChangeRecord(s.hostedZoneID, route53Client.ChangeActionUpsert, resourceRecordSet)
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
