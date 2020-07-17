package route53

import (
	"testing"

	"net"

	"github.com/aws/aws-sdk-go/aws"
	route53Client "github.com/aws/aws-sdk-go/service/route53"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestEditResourceRecords(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	dummySrv1 := &net.SRV{
		Priority: 10,
		Weight:   20,
		Port:     4242,
		Target:   "foobar.example.com",
	}
	dummySrv2 := &net.SRV{
		Priority: 11,
		Weight:   21,
		Port:     4243,
		Target:   "foobaz.example.com",
	}

	rr := []*route53Client.ResourceRecord{}
	newRr := editResourceRecords(true, rr, dummySrv1)

	assert.Len(t, newRr, 1, "New RRs should contain one element")
	assert.Equal(t, newRr, []*route53Client.ResourceRecord{{
		Value: aws.String("10 20 4242 foobar.example.com"),
	}})
	rr = newRr

	// trying to add the same record should not modify from the previous result
	newRr = editResourceRecords(true, rr, dummySrv1)
	assert.Equal(t, rr, newRr, "Adding the same already existing record should be a no-op")
	rr = newRr // just for consistency

	// let's add a record to the set that we don't understand, and ensure it stays untouched.
	rr = append(rr, &route53Client.ResourceRecord{
		Value: aws.String("invalid record"),
	})
	newRr = editResourceRecords(true, rr, dummySrv1)
	assert.Equal(t, rr, newRr, "Adding the same already existing record should be a no-op")
	rr = newRr // just for consistency

	// adding the second (proper) record should add it to the bottom
	newRr = editResourceRecords(true, rr, dummySrv2)
	assert.Len(t, newRr, 3, "new RRs should contain 3 elements")
	assert.Equal(t, *rr[0], *newRr[0], "The existing record shouldn't get touched")
	assert.Equal(t, *rr[1], *newRr[1], "The existing record shouldn't get touched")
	assert.Equal(t, &route53Client.ResourceRecord{
		Value: aws.String("11 21 4243 foobaz.example.com"),
	}, newRr[2])
	rr = newRr

	// add srv1 (again), should be a no-op. More importantly, it shouldn't reorder the list.
	newRr = editResourceRecords(true, rr, dummySrv1)
	assert.Equal(t, rr, newRr, "Adding the same already existing record should be a no-op")
	rr = newRr // just for consistency
	log.Debugf("fooooo - %+v", rr)

	// remove srv2, which was on the third position
	newRr = editResourceRecords(false, rr, dummySrv2)
	assert.Len(t, newRr, 2, "we should now only contain 2 elements")
	assert.Equal(t, *rr[0], *newRr[0])
	assert.Equal(t, *rr[1], *newRr[1])
	rr = newRr

	// remove srv1, after that we should only be left with the invalid record
	newRr = editResourceRecords(false, rr, dummySrv1)
	assert.Len(t, newRr, 1, "we should now only contain 1 elements")
	assert.Equal(t, route53Client.ResourceRecord{
		Value: aws.String("invalid record"),
	}, *newRr[0])
}
