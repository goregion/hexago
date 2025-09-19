package integration

import (
	"context"
	"testing"
	"time"

	"github.com/goregion/hexago/internal/adapter/grpc-api/gen"
	"github.com/goregion/hexago/internal/entity"
)

func TestOHLCPublishingAndReceiving(t *testing.T) {
	testServer := NewTestServer(t)
	testServer.Start(t)
	defer testServer.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	symbol := "BTCUSDT"

	// Subscribe to OHLC stream
	stream, err := testServer.Client().SubscribeToOHLCStream(ctx, &gen.SubscribeToOHLCStreamRequest{
		Symbol:        symbol,
		TimeframeName: "1m",
	})
	if err != nil {
		t.Fatalf("Failed to subscribe to OHLC stream: %v", err)
	}
	defer stream.CloseSend()

	// Create test OHLC data
	testOHLC := CreateTestOHLC(symbol, time.Now().UnixMilli())

	// Publish OHLC data to the server
	go func() {
		// Wait a bit to ensure subscription is active
		time.Sleep(100 * time.Millisecond)
		err := testServer.Server().PublishOHLC(ctx, testOHLC)
		if err != nil {
			t.Errorf("Failed to publish OHLC: %v", err)
		}
	}()

	// Wait for the message
	receivedOHLC, err := WaitForMessage(stream, 5*time.Second)
	if err != nil {
		t.Fatalf("Failed to receive OHLC message: %v", err)
	}

	// Verify received data
	if receivedOHLC.Open != testOHLC.Open {
		t.Errorf("Expected Open %f, got %f", testOHLC.Open, receivedOHLC.Open)
	}
	if receivedOHLC.High != testOHLC.High {
		t.Errorf("Expected High %f, got %f", testOHLC.High, receivedOHLC.High)
	}
	if receivedOHLC.Low != testOHLC.Low {
		t.Errorf("Expected Low %f, got %f", testOHLC.Low, receivedOHLC.Low)
	}
	if receivedOHLC.Close != testOHLC.Close {
		t.Errorf("Expected Close %f, got %f", testOHLC.Close, receivedOHLC.Close)
	}
	if receivedOHLC.TimestampMs != testOHLC.TimestampMs {
		t.Errorf("Expected TimestampMs %d, got %d", testOHLC.TimestampMs, receivedOHLC.TimestampMs)
	}

	t.Log("OHLC publishing and receiving test passed")
}

func TestMultipleOHLCPublishing(t *testing.T) {
	testServer := NewTestServer(t)
	testServer.Start(t)
	defer testServer.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	symbol := "ETHUSDT"

	// Subscribe to OHLC stream
	stream, err := testServer.Client().SubscribeToOHLCStream(ctx, &gen.SubscribeToOHLCStreamRequest{
		Symbol:        symbol,
		TimeframeName: "1m",
	})
	if err != nil {
		t.Fatalf("Failed to subscribe to OHLC stream: %v", err)
	}
	defer stream.CloseSend()

	numMessages := 3
	testOHLCs := make([]*entity.OHLC, numMessages)

	// Create test OHLC data
	baseTime := time.Now().UnixMilli()
	for i := 0; i < numMessages; i++ {
		testOHLCs[i] = &entity.OHLC{
			Symbol:      symbol,
			Open:        float64(100 + i),
			High:        float64(105 + i),
			Low:         float64(98 + i),
			Close:       float64(102 + i),
			TimestampMs: baseTime + int64(i*1000),
		}
	}

	// Publish OHLC data
	go func() {
		time.Sleep(100 * time.Millisecond)
		for i, ohlc := range testOHLCs {
			err := testServer.Server().PublishOHLC(ctx, ohlc)
			if err != nil {
				t.Errorf("Failed to publish OHLC %d: %v", i, err)
			}
			time.Sleep(50 * time.Millisecond) // Small delay between messages
		}
	}()

	// Receive and verify messages
	for i := 0; i < numMessages; i++ {
		receivedOHLC, err := WaitForMessage(stream, 3*time.Second)
		if err != nil {
			t.Fatalf("Failed to receive OHLC message %d: %v", i, err)
		}

		expectedOHLC := testOHLCs[i]
		if receivedOHLC.Open != expectedOHLC.Open {
			t.Errorf("Message %d: Expected Open %f, got %f", i, expectedOHLC.Open, receivedOHLC.Open)
		}
		if receivedOHLC.Close != expectedOHLC.Close {
			t.Errorf("Message %d: Expected Close %f, got %f", i, expectedOHLC.Close, receivedOHLC.Close)
		}
		if receivedOHLC.TimestampMs != expectedOHLC.TimestampMs {
			t.Errorf("Message %d: Expected TimestampMs %d, got %d", i, expectedOHLC.TimestampMs, receivedOHLC.TimestampMs)
		}
	}

	t.Log("Multiple OHLC publishing test passed")
}

func TestPublishingToNonExistentSubscription(t *testing.T) {
	testServer := NewTestServer(t)
	testServer.Start(t)
	defer testServer.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create test OHLC data for a symbol with no subscribers
	testOHLC := CreateTestOHLC("NONEXISTENT", time.Now().UnixMilli())

	// This should not cause any error, just no messages sent
	err := testServer.Server().PublishOHLC(ctx, testOHLC)
	if err != nil {
		t.Errorf("Unexpected error when publishing to non-existent subscription: %v", err)
	}

	t.Log("Publishing to non-existent subscription test passed")
}
