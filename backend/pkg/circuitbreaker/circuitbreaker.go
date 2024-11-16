package circuitbreaker

import (
	"time"

	"github.com/sony/gobreaker"
)

func NewCircuitBreaker(name string) *gobreaker.CircuitBreaker {
	return gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        name,
		MaxRequests: 3,
		Interval:    10 * time.Second,
		Timeout:     60 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failureRatio >= 0.6
		},
	})
} 