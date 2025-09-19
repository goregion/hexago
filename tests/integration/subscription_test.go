package integration

import (
	"context"
	"testing"
	"time"

	"github.com/goregion/hexago/internal/adapter/grpc-api/gen"
)

func TestOHLCStreamSubscription(t *testing.T) {
	testServer := NewTestServer(t)
	testServer.Start(t)
	defer testServer.Stop()

	tests := []struct {
		name          string
		symbol        string
		timeframeName string
		expectError   bool
	}{
		{
			name:          "Valid subscription BTCUSDT",
			symbol:        "BTCUSDT",
			timeframeName: "1m",
			expectError:   false,
		},
		{
			name:          "Valid subscription ETHUSDT",
			symbol:        "ETHUSDT",
			timeframeName: "5m",
			expectError:   false,
		},
		{
			name:          "Valid subscription with empty timeframe",
			symbol:        "ADAUSDT",
			timeframeName: "",
			expectError:   false, // Server doesn't validate timeframe currently
		},
		{
			name:          "Valid subscription with empty symbol",
			symbol:        "",
			timeframeName: "1h",
			expectError:   false, // Server doesn't validate symbol currently
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			stream, err := testServer.Client().SubscribeToOHLCStream(ctx, &gen.SubscribeToOHLCStreamRequest{
				Symbol:        tt.symbol,
				TimeframeName: tt.timeframeName,
			})

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Verify stream is created
			if stream == nil {
				t.Fatal("Stream is nil")
			}

			// Close the stream
			err = stream.CloseSend()
			if err != nil {
				t.Errorf("Failed to close stream: %v", err)
			}
		})
	}
}

func TestMultipleSubscriptions(t *testing.T) {
	testServer := NewTestServer(t)
	testServer.Start(t)
	defer testServer.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create multiple subscriptions
	symbols := []string{"BTCUSDT", "ETHUSDT", "ADAUSDT"}
	streams := make([]gen.OHLCService_SubscribeToOHLCStreamClient, len(symbols))

	// Subscribe to multiple symbols
	for i, symbol := range symbols {
		stream, err := testServer.Client().SubscribeToOHLCStream(ctx, &gen.SubscribeToOHLCStreamRequest{
			Symbol:        symbol,
			TimeframeName: "1m",
		})
		if err != nil {
			t.Fatalf("Failed to subscribe to %s: %v", symbol, err)
		}
		streams[i] = stream
	}

	// Verify all streams are created
	for i, stream := range streams {
		if stream == nil {
			t.Fatalf("Stream %d is nil", i)
		}
	}

	// Close all streams
	for i, stream := range streams {
		err := stream.CloseSend()
		if err != nil {
			t.Errorf("Failed to close stream %d: %v", i, err)
		}
	}

	t.Log("Multiple subscriptions test passed")
}

func TestSubscriptionWithCancelledContext(t *testing.T) {
	testServer := NewTestServer(t)
	testServer.Start(t)
	defer testServer.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	// Subscribe to stream
	stream, err := testServer.Client().SubscribeToOHLCStream(ctx, &gen.SubscribeToOHLCStreamRequest{
		Symbol:        "BTCUSDT",
		TimeframeName: "1m",
	})
	if err != nil {
		t.Fatalf("Failed to subscribe: %v", err)
	}

	// Cancel context immediately
	cancel()

	// Wait a bit for the cancellation to propagate
	time.Sleep(100 * time.Millisecond)

	// Stream should be cancelled, but we can still close it
	stream.CloseSend()

	t.Log("Subscription with cancelled context test passed")
}
