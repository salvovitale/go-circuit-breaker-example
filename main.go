package main

import (
	"errors"
	"log"
	"math/rand"
	"time"

	"github.com/sony/gobreaker"
)

var cb *gobreaker.CircuitBreaker

func main() {
	cb = gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "demo",
		MaxRequests: 3,
		Timeout:     4,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return failureRatio >= 0.6
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			log.Printf("Circuit Breaker: %s, changed from %v, to %v", name, from, to)
		},
	})
	for {
		cbRes, cbErr := cb.Execute(func() (interface{}, error) {
			res, isErr := isError()
			if isErr {
				return nil, errors.New("error")
			}
			return res, nil
		})
		if cbErr != nil {
			log.Printf("Circuit breaker error %v", cbErr)
		} else {
			log.Printf("Circuit breaker result %v", cbRes)
		}
		time.Sleep(1 * time.Second)
	}
}

// this function will generate a failure 50% of the time
func isError() (int, bool) {
	min := 0
	max := 10
	result := rand.Intn(max-min) + min
	return result, result%2 != 0
}
