package dns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSRV(t *testing.T) {
	rrText := "10 20 443 foobar.baz.bar."
	srv, err := ParseSRV(rrText)
	if assert.NoError(t, err) {
		assert.Equal(t, uint16(10), srv.Priority, "Priority should be parsed correctly")
		assert.Equal(t, uint16(20), srv.Weight, "Weight should be parsed correctly")
		assert.Equal(t, uint16(443), srv.Port, "Port should be parsed correctly")
		assert.Equal(t, "foobar.baz.bar.", srv.Target, "Target should be parsed correctly")
	}
}
