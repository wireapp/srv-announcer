package route53

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	awsSession "github.com/aws/aws-sdk-go/aws/session"
	route53Client "github.com/aws/aws-sdk-go/service/route53"
)

// Client wraps the AWS Route53 service
type Client struct {
	Service *route53Client.Route53
}

// NewClient constructs a new Service, wrapping a aws route53 client under the hood.
func NewClient() *Client {
	return &Client{
		Service: route53Client.New(awsSession.Must(awsSession.NewSession())),
	}
}

// GetZoneByName looks up a zone by its DNSName.
// If the zone doesn't exist, the zone will be nil.
func (c *Client) GetZoneByName(name string) (*route53Client.HostedZone, error) {
	name = addDotSuffixIfNeeded(name)

	resp, err := c.Service.ListHostedZonesByName(&route53Client.ListHostedZonesByNameInput{
		DNSName:  aws.String(name),
		MaxItems: aws.String("1"),
	})

	if err != nil {
		return nil, err
	}

	// if we don't receive a record, it doesn't exist
	if len(resp.HostedZones) == 0 {
		return nil, nil
	}

	zone := resp.HostedZones[0]

	// if we received a record, we need to check if it's the one we asked for,
	// or just something lexicographically later than what we requested
	if *zone.Name != name {
		return nil, nil
	}

	return zone, nil
}

// GetResourceRecordSetByName returns a resource record zet of a given zone with the passed record name,
func (c *Client) GetResourceRecordSetByName(zoneID, recordName, recordType string) (*route53Client.ResourceRecordSet, error) {
	recordName = addDotSuffixIfNeeded(recordName)

	reqInput := &route53Client.ListResourceRecordSetsInput{
		HostedZoneId:    aws.String(zoneID),
		StartRecordName: aws.String(recordName),
		StartRecordType: aws.String(recordType),
		MaxItems:        aws.String("1"),
	}
	resp, err := c.Service.ListResourceRecordSets(reqInput)

	if err != nil {
		return nil, err
	}

	if len(resp.ResourceRecordSets) == 0 || aws.StringValue(resp.ResourceRecordSets[0].Name) != recordName {
		return nil, nil
	}

	return resp.ResourceRecordSets[0], nil
}

// ChangeRecord creates, updates or deletes a record depending on the given ChangeAction
func (c *Client) ChangeRecord(zoneID, action string, recordSet *route53Client.ResourceRecordSet) (*route53Client.ChangeInfo, error) {
	recordName := addDotSuffixIfNeeded(*recordSet.Name)

	recordSet.Name = aws.String(recordName)

	recordSetInput := &route53Client.ChangeResourceRecordSetsInput{
		HostedZoneId: aws.String(zoneID),
		ChangeBatch: &route53Client.ChangeBatch{
			Comment: aws.String(fmt.Sprintf("Updated automatically on %s", time.Now().Format(time.RFC3339))),
			Changes: []*route53Client.Change{{
				Action:            aws.String(action),
				ResourceRecordSet: recordSet,
			}},
		},
	}

	res, err := c.Service.ChangeResourceRecordSets(recordSetInput)
	if err != nil {
		return nil, err
	}
	return res.ChangeInfo, nil
}
