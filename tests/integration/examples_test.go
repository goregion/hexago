package integration_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/goregion/hexago/internal/adapter/grpc-api/gen"
	integration "github.com/goregion/hexago/tests/integration"
)

// Example_basicUsage demonstrates basic usage of the test framework
func Example_basicUsage() {
	// This example shows how to use the integration test framework
	// in your own tests or for manual testing

	// Create a test instance (normally you'd do this in a test function)
	// testServer := integration.NewTestServer(t)
	// testServer.Start(t)
	// defer testServer.Stop()

	fmt.Println("Integration test framework example")
	// Output: Integration test framework example
}

// Example_customTest shows how to write a custom integration test
func Example_customTest() {
	// This would be inside a proper test function
	var t *testing.T // In real usage, this comes from the test function

	// Setup
	testServer := integration.NewTestServer(t)
	testServer.Start(t)
	defer testServer.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create a subscription
	stream, err := testServer.Client().SubscribeToOHLCStream(ctx, &gen.SubscribeToOHLCStreamRequest{
		Symbol:        "BTCUSDT",
		TimeframeName: "1m",
	})
	if err != nil {
		fmt.Printf("Failed to subscribe: %v", err)
		return
	}
	defer stream.CloseSend()

	// Publish some data
	testOHLC := integration.CreateTestOHLC("BTCUSDT", time.Now().UnixMilli())
	go func() {
		time.Sleep(100 * time.Millisecond)
		testServer.Server().PublishOHLC(ctx, testOHLC)
	}()

	// Wait for the message
	receivedOHLC, err := integration.WaitForMessage(stream, 3*time.Second)
	if err != nil {
		fmt.Printf("Failed to receive message: %v", err)
		return
	}

	fmt.Printf("Received OHLC: Open=%.2f, Close=%.2f", receivedOHLC.Open, receivedOHLC.Close)
}

// Benchmark for testing performance
func BenchmarkOHLCPublishing(b *testing.B) {
	// Note: This benchmark requires converting the test framework
	// to work with *testing.B instead of *testing.T

	// For now, this is just a placeholder showing how benchmarks could be structured
	b.Skip("Benchmark requires adaptation of test framework for *testing.B")

	/* Example implementation:
	testServer := NewBenchmarkTestServer(b)
	testServer.Start(b)
	defer testServer.Stop()

	ctx := context.Background()
	stream, _ := testServer.Client().SubscribeToOHLCStream(ctx, &gen.SubscribeToOHLCStreamRequest{
		Symbol:        "BTCUSDT",
		TimeframeName: "1m",
	})
	defer stream.CloseSend()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			testOHLC := CreateTestOHLC("BTCUSDT", time.Now().UnixMilli())
			testServer.Server().PublishOHLC(ctx, testOHLC)
		}
	})
	*/
}

// TestRealWorldScenario demonstrates a more complex real-world scenario
func TestRealWorldScenario(t *testing.T) {
	testServer := integration.NewTestServer(t)
	testServer.Start(t)
	defer testServer.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	symbols := []string{"BTCUSDT", "ETHUSDT", "ADAUSDT"}

	// Simulate multiple clients subscribing to different symbols
	streams := make([]gen.OHLCService_SubscribeToOHLCStreamClient, len(symbols))
	for i, symbol := range symbols {
		stream, err := testServer.Client().SubscribeToOHLCStream(ctx, &gen.SubscribeToOHLCStreamRequest{
			Symbol:        symbol,
			TimeframeName: "1m",
		})
		if err != nil {
			t.Fatalf("Failed to subscribe to %s: %v", symbol, err)
		}
		streams[i] = stream
		defer stream.CloseSend()
	}

	// Simulate market data coming in for different symbols
	go func() {
		time.Sleep(100 * time.Millisecond)
		for i := 0; i < 5; i++ {
			for _, symbol := range symbols {
				testOHLC := integration.CreateTestOHLC(symbol, time.Now().UnixMilli()+int64(i*1000))
				err := testServer.Server().PublishOHLC(ctx, testOHLC)
				if err != nil {
					t.Errorf("Failed to publish OHLC for %s: %v", symbol, err)
				}
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	// Each client should receive messages for their symbol
	receivedCounts := make([]int, len(symbols))
	for i := range streams {
		go func(streamIndex int) {
			for {
				_, err := integration.WaitForMessage(streams[streamIndex], 1*time.Second)
				if err != nil {
					return // Timeout or error, stop receiving
				}
				receivedCounts[streamIndex]++
			}
		}(i)
	}

	// Wait for messages to be processed
	time.Sleep(2 * time.Second)

	// Verify that each client received some messages
	for i, count := range receivedCounts {
		t.Logf("Client %d (%s) received %d messages", i, symbols[i], count)
		if count == 0 {
			t.Errorf("Client %d (%s) should have received at least one message", i, symbols[i])
		}
	}
}
