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
		description   string
	}{
		{
			name:          "Valid subscription BTCUSDT",
			symbol:        "BTCUSDT",
			timeframeName: "1m",
			expectError:   false,
			description:   "Should successfully create subscription for BTCUSDT with 1m timeframe",
		},
		{
			name:          "Valid subscription ETHUSDT",
			symbol:        "ETHUSDT",
			timeframeName: "5m",
			expectError:   false,
			description:   "Should successfully create subscription for ETHUSDT with 5m timeframe",
		},
		{
			name:          "Valid subscription with empty timeframe",
			symbol:        "ADAUSDT",
			timeframeName: "",
			expectError:   false,
			description:   "Should handle empty timeframe gracefully", // Server doesn't validate timeframe currently
		},
		{
			name:          "Valid subscription with empty symbol",
			symbol:        "",
			timeframeName: "1h",
			expectError:   false,
			description:   "Should handle empty symbol gracefully", // Server doesn't validate symbol currently
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
					t.Errorf("Expected error but got none for test: %s", tt.description)
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error for test %s: %v", tt.description, err)
			}

			// Verify stream is created
			if stream == nil {
				t.Fatal("Stream is nil")
			}

			// Test that stream context is properly set
			select {
			case <-stream.Context().Done():
				t.Error("Stream context should not be done immediately after creation")
			default:
				// Expected behavior
			}

			// Close the stream properly
			err = stream.CloseSend()
			if err != nil {
				t.Errorf("Failed to close stream: %v", err)
			}

			t.Logf("✓ %s: %s", tt.name, tt.description)
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
		t.Logf("✓ Successfully subscribed to %s", symbol)
	}

	// Verify all streams are created and have valid contexts
	for i, stream := range streams {
		if stream == nil {
			t.Fatalf("Stream %d (%s) is nil", i, symbols[i])
		}

		// Verify stream context is active
		select {
		case <-stream.Context().Done():
			t.Errorf("Stream %d (%s) context is already done", i, symbols[i])
		default:
			// Expected behavior - context should be active
		}
	}

	// Test that streams are independent - closing one doesn't affect others
	if len(streams) > 1 {
		err := streams[0].CloseSend()
		if err != nil {
			t.Errorf("Failed to close first stream: %v", err)
		}

		// Verify other streams are still active
		for i := 1; i < len(streams); i++ {
			select {
			case <-streams[i].Context().Done():
				t.Errorf("Stream %d should still be active after closing stream 0", i)
			default:
				// Expected behavior
			}
		}
	}

	// Close remaining streams
	for i := 1; i < len(streams); i++ {
		err := streams[i].CloseSend()
		if err != nil {
			t.Errorf("Failed to close stream %d: %v", i, err)
		}
	}

	t.Log("✓ Multiple subscriptions test passed - all streams created and closed properly")
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

	// Verify stream is initially active
	select {
	case <-stream.Context().Done():
		t.Error("Stream context should not be done immediately after creation")
	default:
		t.Log("✓ Stream context is active after creation")
	}

	// Cancel context
	cancel()

	// Wait for the cancellation to propagate
	maxWait := 2 * time.Second
	start := time.Now()
	for {
		select {
		case <-stream.Context().Done():
			t.Logf("✓ Stream context cancelled after %v", time.Since(start))
			goto contextCancelled
		default:
			if time.Since(start) > maxWait {
				t.Errorf("Stream context not cancelled after %v", maxWait)
				goto contextCancelled
			}
			time.Sleep(10 * time.Millisecond)
		}
	}

contextCancelled:
	// Attempt to receive from cancelled stream should return error
	_, err = stream.Recv()
	if err == nil {
		t.Error("Expected error when receiving from cancelled stream")
	} else {
		t.Logf("✓ Received expected error from cancelled stream: %v", err)
	}

	// Stream should still be closable gracefully
	err = stream.CloseSend()
	if err != nil {
		t.Logf("Note: CloseSend returned error (expected for cancelled stream): %v", err)
	}

	t.Log("✓ Subscription with cancelled context test passed")
}

// TestSubscriptionDataReceiving tests actual data receiving from subscriptions
func TestSubscriptionDataReceiving(t *testing.T) {
	testServer := NewTestServer(t)
	testServer.Start(t)
	defer testServer.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Subscribe to stream
	stream, err := testServer.Client().SubscribeToOHLCStream(ctx, &gen.SubscribeToOHLCStreamRequest{
		Symbol:        "BTCUSDT",
		TimeframeName: "1m",
	})
	if err != nil {
		t.Fatalf("Failed to subscribe: %v", err)
	}
	defer stream.CloseSend()

	// Create test OHLC data
	testOHLC := CreateTestOHLC("BTCUSDT", time.Now().UnixMilli())

	// Publish OHLC data in a goroutine
	go func() {
		time.Sleep(100 * time.Millisecond) // Small delay to ensure subscription is ready
		err := testServer.Server().PublishOHLC(ctx, testOHLC)
		if err != nil {
			t.Errorf("Failed to publish OHLC: %v", err)
		}
	}()

	// Wait for message
	message, err := WaitForMessage(stream, 5*time.Second)
	if err != nil {
		t.Fatalf("Failed to receive message: %v", err)
	}

	// Verify received data
	if message == nil {
		t.Fatal("Received message is nil")
	}

	// Check OHLC values
	if message.Open != testOHLC.Open {
		t.Errorf("Expected Open %f, got %f", testOHLC.Open, message.Open)
	}
	if message.High != testOHLC.High {
		t.Errorf("Expected High %f, got %f", testOHLC.High, message.High)
	}
	if message.Low != testOHLC.Low {
		t.Errorf("Expected Low %f, got %f", testOHLC.Low, message.Low)
	}
	if message.Close != testOHLC.Close {
		t.Errorf("Expected Close %f, got %f", testOHLC.Close, message.Close)
	}
	if message.TimestampMs != testOHLC.TimestampMs {
		t.Errorf("Expected TimestampMs %d, got %d", testOHLC.TimestampMs, message.TimestampMs)
	}

	t.Log("✓ Subscription data receiving test passed - OHLC data received correctly")
}

// TestSubscriptionSymbolFiltering tests that subscriptions only receive data for their symbol
func TestSubscriptionSymbolFiltering(t *testing.T) {
	testServer := NewTestServer(t)
	testServer.Start(t)
	defer testServer.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Subscribe to BTCUSDT
	streamBTC, err := testServer.Client().SubscribeToOHLCStream(ctx, &gen.SubscribeToOHLCStreamRequest{
		Symbol:        "BTCUSDT",
		TimeframeName: "1m",
	})
	if err != nil {
		t.Fatalf("Failed to subscribe to BTCUSDT: %v", err)
	}
	defer streamBTC.CloseSend()

	// Subscribe to ETHUSDT
	streamETH, err := testServer.Client().SubscribeToOHLCStream(ctx, &gen.SubscribeToOHLCStreamRequest{
		Symbol:        "ETHUSDT",
		TimeframeName: "1m",
	})
	if err != nil {
		t.Fatalf("Failed to subscribe to ETHUSDT: %v", err)
	}
	defer streamETH.CloseSend()

	// Publish BTCUSDT data
	btcOHLC := CreateTestOHLC("BTCUSDT", time.Now().UnixMilli())
	go func() {
		time.Sleep(100 * time.Millisecond)
		testServer.Server().PublishOHLC(ctx, btcOHLC)
	}()

	// BTCUSDT stream should receive the message
	btcMessage, err := WaitForMessage(streamBTC, 2*time.Second)
	if err != nil {
		t.Fatalf("BTCUSDT stream should have received message: %v", err)
	}
	if btcMessage == nil {
		t.Fatal("BTCUSDT message is nil")
	}

	// ETHUSDT stream should NOT receive the message (timeout expected)
	ethMessage, err := WaitForMessage(streamETH, 1*time.Second)
	if err == nil {
		t.Error("ETHUSDT stream should not have received BTCUSDT message")
	}
	if ethMessage != nil {
		t.Error("ETHUSDT stream received unexpected message")
	}

	t.Log("✓ Subscription symbol filtering test passed - streams only receive their symbol's data")
}

// TestStreamRecoveryAfterError tests stream behavior after network errors
func TestStreamRecoveryAfterError(t *testing.T) {
	testServer := NewTestServer(t)
	testServer.Start(t)
	defer testServer.Stop()

	ctx1, cancel1 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel1()

	// Create subscription
	stream, err := testServer.Client().SubscribeToOHLCStream(ctx1, &gen.SubscribeToOHLCStreamRequest{
		Symbol:        "BTCUSDT",
		TimeframeName: "1m",
	})
	if err != nil {
		t.Fatalf("Failed to create subscription: %v", err)
	}

	// Force an error by closing the stream
	err = stream.CloseSend()
	if err != nil {
		t.Logf("CloseSend returned error (may be expected): %v", err)
	}

	// Try to receive - should get an error
	_, err = stream.Recv()
	if err == nil {
		t.Error("Expected error when receiving from closed stream")
	} else {
		t.Logf("✓ Received expected error from closed stream: %v", err)
	}

	// Create new context and subscription (simulating recovery)
	ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel2()

	newStream, err := testServer.Client().SubscribeToOHLCStream(ctx2, &gen.SubscribeToOHLCStreamRequest{
		Symbol:        "BTCUSDT",
		TimeframeName: "1m",
	})
	if err != nil {
		t.Fatalf("Failed to create new subscription after error: %v", err)
	}
	defer newStream.CloseSend()

	// Verify new stream works by checking it's not nil and context is active
	if newStream == nil {
		t.Fatal("New stream is nil")
	}

	select {
	case <-newStream.Context().Done():
		t.Error("New stream context should not be done immediately")
	default:
		t.Log("✓ New stream context is active")
	}

	t.Log("✓ Stream recovery after error test passed - new subscription works after error")
}
