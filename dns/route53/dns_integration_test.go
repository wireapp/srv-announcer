// +build integration

package route53

import (
	"testing"

	"fmt"
	"net"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// TestIntegrationAddAndRemoveATarget requires valid AWS credentials being available
// DOCS: https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html#specifying-credentials
func TestIntegrationAddAndRemoveATarget(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	zoneName := os.Getenv("TEST_ZONE_NAME")
	if zoneName == "" {
		t.Fatal("Missing environment variable: TEST_ZONE_NAME. Mark as failed")
	}

	recordName := fmt.Sprintf("%s.%s", "_srv._tcp.test", zoneName)
	ttl := uint16(60)
	client := NewClient()
	dnsZone, err := client.GetZoneByName(zoneName)
	if err != nil {
		t.Fatal(err)
	}
	hostedZoneID := aws.StringValue(dnsZone.Id)

	srvManager := NewSRVManager(
		client,
		hostedZoneID,
		recordName,
		ttl,
	)

	record := &net.SRV{
		Priority: 10,
		Weight:   20,
		Port:     4242,
		Target:   "sub.domain.test.",
	}

	err = srvManager.Add(record)
	if assert.NoError(t, err) {
		resourceRecordSet, err := srvManager.client.GetResourceRecordSetByName(hostedZoneID, recordName, "SRV")
		assert.NoError(t, err)
		assert.NotNil(t, resourceRecordSet, "Record should exist")
	}

	err = srvManager.Remove(record)
	if assert.NoError(t, err) {
		resourceRecordSet, err := srvManager.client.GetResourceRecordSetByName(hostedZoneID, recordName, "SRV")
		assert.NoError(t, err)
		assert.Nil(t, resourceRecordSet, "Record should not exist")
	}
}
