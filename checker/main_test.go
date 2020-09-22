package checker

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zinfra/srv-announcer/checker/healthchecks/mock"
	"github.com/zinfra/srv-announcer/dns/dummy"
)

// TestChecker tests Healthcheck.Run.
// It passes in a mock.Healthcheck and a dummySrvManager, then
// uses the mock.Healthcheck.HealthC channel to change health
// and inspects dummySrvManagers rrset afterwards
// it also tests the record is removed from the set once the context is cancelled
func TestChecker(t *testing.T) {
	// setup tooling
	mockHealthcheck := &mock.Healthcheck{
		HealthC: make(chan bool, 1),
	}

	dummySrv := &net.SRV{
		Priority: 10,
		Weight:   20,
		Port:     4242,
		Target:   "foobar.example.com.",
	}

	dummySrvManager := &dummy.SrvManager{}

	// setup context
	ctx, cancelCtx := context.WithCancel(context.Background())

	// kick off the checker
	go Run(ctx, mockHealthcheck, dummySrv, dummySrvManager)

	assert.Len(t, dummySrvManager.SrvRecordSet, 0, "Initially rrset should be empty")

	// let service become healthy
	mockHealthcheck.HealthC <- true

	assert.Eventually(t, func() bool {
		return len(dummySrvManager.SrvRecordSet) != 0
	}, 50*time.Millisecond, time.Millisecond, "rrset should eventually not be empty anymore")
	assert.Equal(t, *dummySrv, dummySrvManager.SrvRecordSet[0], "rrset should contain dummySrv")

	// let service become unhealthy
	mockHealthcheck.HealthC <- false

	assert.Eventually(t, func() bool {
		return len(dummySrvManager.SrvRecordSet) == 0
	}, 50*time.Millisecond, time.Millisecond, "rrset should eventually become empty again")

	// make service healthy again so we can test teardown
	mockHealthcheck.HealthC <- true

	assert.Eventually(t, func() bool {
		return len(dummySrvManager.SrvRecordSet) != 0
	}, 50*time.Millisecond, time.Millisecond, "rrset should eventually not be empty anymore")

	// cancel context
	cancelCtx()

	assert.Eventually(t, func() bool {
		return len(dummySrvManager.SrvRecordSet) == 0
	}, 50*time.Millisecond, time.Millisecond, "rrset should eventually become empty on teardown")

}
