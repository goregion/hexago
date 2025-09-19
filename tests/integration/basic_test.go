package integration

import (
	"context"
	"testing"
	"time"

	"github.com/goregion/hexago/internal/adapter/grpc-api/gen"
)

func TestServerStartupAndShutdown(t *testing.T) {
	// Create and start test server
	testServer := NewTestServer(t)
	testServer.Start(t)
	defer testServer.Stop()

	// Test that we can connect to the server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Try to create a stream connection to verify server is running
	stream, err := testServer.Client().SubscribeToOHLCStream(ctx, &gen.SubscribeToOHLCStreamRequest{
		Symbol:        "BTCUSDT",
		TimeframeName: "1m",
	})
	if err != nil {
		t.Fatalf("Failed to create stream: %v", err)
	}

	// Close the stream
	stream.CloseSend()

	t.Log("Server startup and connection test passed")
}

func TestBasicSubscription(t *testing.T) {
	testServer := NewTestServer(t)
	testServer.Start(t)
	defer testServer.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Subscribe to OHLC stream
	stream, err := testServer.Client().SubscribeToOHLCStream(ctx, &gen.SubscribeToOHLCStreamRequest{
		Symbol:        "BTCUSDT",
		TimeframeName: "1m",
	})
	if err != nil {
		t.Fatalf("Failed to subscribe to OHLC stream: %v", err)
	}

	// Verify stream is created (it should block waiting for messages)
	// We'll test this by trying to close the stream without errors
	err = stream.CloseSend()
	if err != nil {
		t.Fatalf("Failed to close stream: %v", err)
	}

	t.Log("Basic subscription test passed")
}
