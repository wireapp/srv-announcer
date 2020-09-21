package healthchecks

import "context"

// IHealthcheck describes the basic healthchecker interface
type IHealthcheck interface {
	// Run will run the healthchecker-specific check
	// in healthchecker-specific intervals.
	// They will regularily push their health status
	// to the healthyChan.
	Run(ctx context.Context, healthyChan chan<- bool)
}
