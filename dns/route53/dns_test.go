package route53

import (
	"testing"
	"time"

	"fmt"
	"net"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type mockR53 struct {
	recordSet *route53.ResourceRecordSet
	zoneID    string
}

func (r53 *mockR53) ChangeResourceRecordSets(input *route53.ChangeResourceRecordSetsInput) (*route53.ChangeResourceRecordSetsOutput, error) {
	// NOTE: if the ChangeBatch would contain more than just one Change, they would be ignored
	change := input.ChangeBatch.Changes[0]

	switch *change.Action {
	case route53.ChangeActionCreate, route53.ChangeActionUpsert:
		r53.recordSet = change.ResourceRecordSet

	case route53.ChangeActionDelete:
		r53.recordSet = nil

	default:
		panic(fmt.Sprintf("Invalid action: %s", *change.Action))
	}

	return &route53.ChangeResourceRecordSetsOutput{
		ChangeInfo: &route53.ChangeInfo{
			Id:          aws.String("MOCKED-AWS-API-ID"),
			Status:      aws.String("INSYNC"),
			SubmittedAt: aws.Time(time.Now()),
		},
	}, nil
}

func (r53 *mockR53) ListResourceRecordSets(input *route53.ListResourceRecordSetsInput) (*route53.ListResourceRecordSetsOutput, error) {
	recordName := *input.StartRecordName

	resourceRecordSets := []*route53.ResourceRecordSet{}
	if r53.recordSet != nil && *r53.recordSet.Name == recordName {
		resourceRecordSets = append(resourceRecordSets, r53.recordSet)
	}

	return &route53.ListResourceRecordSetsOutput{
		ResourceRecordSets: resourceRecordSets,
		IsTruncated:        aws.Bool(false),
		MaxItems:           input.MaxItems,
	}, nil
}

// NOTE: always returns a list of one item matching the queried DNSName
func (r53 *mockR53) ListHostedZonesByName(input *route53.ListHostedZonesByNameInput) (*route53.ListHostedZonesByNameOutput, error) {
	return &route53.ListHostedZonesByNameOutput{
		HostedZones: []*route53.HostedZone{
			&route53.HostedZone{
				Id:              aws.String(r53.zoneID),
				CallerReference: aws.String("MOCKED-CAL-REF"),
				Name:            input.DNSName,
			},
		},
		IsTruncated: aws.Bool(false),
		MaxItems:    input.MaxItems,
	}, nil
}

func TestAddAndRemoveATarget(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	mockedZoneID := "20N31D"
	recordName := fmt.Sprintf("%s.%s", "_srv._tcp.test", "example.com")
	ttl := uint16(60)
	client := &Client{
		Service: &mockR53{
			recordSet: nil,
			zoneID:    mockedZoneID,
		},
	}

	srvManager := NewSRVManager(
		client,
		mockedZoneID,
		recordName,
		ttl,
	)

	record := &net.SRV{
		Priority: 10,
		Weight:   20,
		Port:     4242,
		Target:   "sub.domain.test.",
	}

	err := srvManager.Add(record)
	if assert.NoError(t, err) {
		resourceRecordSet, err := srvManager.client.GetResourceRecordSetByName(mockedZoneID, recordName, "SRV")
		assert.NoError(t, err)
		assert.NotNil(t, resourceRecordSet, "Record should exist")
	}

	err = srvManager.Remove(record)
	if assert.NoError(t, err) {
		resourceRecordSet, err := srvManager.client.GetResourceRecordSetByName(mockedZoneID, recordName, "SRV")
		assert.NoError(t, err)
		assert.Nil(t, resourceRecordSet, "Record should not exist")
	}
}
