package main

import (
	"context"
	"fmt"
	"time"

	"github.com/goregion/hexago/pkg/grexit"
)

func main() {
	fmt.Println("=== Grexit Timeout Demo ===")

	// Demo 1: Basic timeout context
	fmt.Println("\n1. Testing basic timeout context...")
	ctx := context.Background()
	timeoutCtx := grexit.WithGrexitTimeoutDuration(ctx, 2*time.Second)

	fmt.Println("Press Ctrl+C to test graceful shutdown with timeout...")
	fmt.Println("Timeout is set to 2 seconds after signal is received.")

	// Simulate some work
	go func() {
		for {
			select {
			case <-timeoutCtx.Done():
				fmt.Println("Background worker received shutdown signal")
				return
			case <-time.After(500 * time.Millisecond):
				fmt.Println("Background worker is running...")
			}
		}
	}()

	// Wait for shutdown signal
	<-timeoutCtx.Done()
	fmt.Println("Main received shutdown signal, starting cleanup...")

	// Simulate cleanup work
	for i := 1; i <= 5; i++ {
		select {
		case <-time.After(400 * time.Millisecond):
			fmt.Printf("Cleanup step %d/5 completed\n", i)
		case <-timeoutCtx.Done():
			fmt.Println("Timeout expired! Forcing shutdown...")
			return
		}
	}

	fmt.Println("Graceful shutdown completed successfully!")
}
