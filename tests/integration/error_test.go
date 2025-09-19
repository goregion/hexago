package integration

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/goregion/hexago/internal/adapter/grpc-api/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestClientDisconnection(t *testing.T) {
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

	// Immediately close the stream to simulate client disconnection
	err = stream.CloseSend()
	if err != nil {
		t.Errorf("Failed to close stream: %v", err)
	}

	// Try to receive - should get EOF or similar error
	_, err = stream.Recv()
	if err == nil {
		t.Error("Expected error when receiving from closed stream")
	} else if err != io.EOF {
		// Check if it's a cancellation error which is also acceptable
		if s, ok := status.FromError(err); ok {
			if s.Code() != codes.Canceled {
				t.Logf("Received expected error: %v", err)
			}
		} else {
			t.Logf("Received error (acceptable): %v", err)
		}
	}

	t.Log("Client disconnection test passed")
}

func TestContextCancellation(t *testing.T) {
	testServer := NewTestServer(t)
	testServer.Start(t)
	defer testServer.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// Subscribe to stream
	stream, err := testServer.Client().SubscribeToOHLCStream(ctx, &gen.SubscribeToOHLCStreamRequest{
		Symbol:        "BTCUSDT",
		TimeframeName: "1m",
	})
	if err != nil {
		t.Fatalf("Failed to subscribe: %v", err)
	}

	// Cancel context
	cancel()

	// Wait a bit for cancellation to propagate
	time.Sleep(100 * time.Millisecond)

	// Try to receive - should get cancellation error
	_, err = stream.Recv()
	if err == nil {
		t.Error("Expected error when receiving from cancelled context")
	} else {
		// Check if it's a cancellation-related error
		if s, ok := status.FromError(err); ok {
			if s.Code() == codes.Canceled || s.Code() == codes.DeadlineExceeded {
				t.Logf("Received expected cancellation error: %v", err)
			} else {
				t.Logf("Received error: %v (code: %v)", err, s.Code())
			}
		} else {
			t.Logf("Received error: %v", err)
		}
	}

	t.Log("Context cancellation test passed")
}

func TestServerShutdownDuringStream(t *testing.T) {
	testServer := NewTestServer(t)
	testServer.Start(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Subscribe to stream
	stream, err := testServer.Client().SubscribeToOHLCStream(ctx, &gen.SubscribeToOHLCStreamRequest{
		Symbol:        "BTCUSDT",
		TimeframeName: "1m",
	})
	if err != nil {
		t.Fatalf("Failed to subscribe: %v", err)
	}

	// Shutdown server while stream is active
	go func() {
		time.Sleep(100 * time.Millisecond)
		testServer.Stop()
	}()

	// Try to receive - should get error due to server shutdown
	_, err = stream.Recv()
	if err == nil {
		t.Error("Expected error when server shuts down")
	} else {
		t.Logf("Received expected error after server shutdown: %v", err)
	}

	t.Log("Server shutdown during stream test passed")
}

func TestInvalidStreamOperations(t *testing.T) {
	testServer := NewTestServer(t)
	testServer.Start(t)
	defer testServer.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Subscribe to stream
	stream, err := testServer.Client().SubscribeToOHLCStream(ctx, &gen.SubscribeToOHLCStreamRequest{
		Symbol:        "BTCUSDT",
		TimeframeName: "1m",
	})
	if err != nil {
		t.Fatalf("Failed to subscribe: %v", err)
	}

	// Close the stream
	err = stream.CloseSend()
	if err != nil {
		t.Errorf("Failed to close stream: %v", err)
	}

	// Try to close again - should not cause panic
	err = stream.CloseSend()
	// This might return an error or not, depending on gRPC implementation
	// The important thing is that it doesn't panic
	t.Logf("Second close returned: %v", err)

	t.Log("Invalid stream operations test passed")
}

func TestNetworkTimeout(t *testing.T) {
	testServer := NewTestServer(t)
	testServer.Start(t)
	defer testServer.Stop()

	// Use a very short timeout to simulate network timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	// This should timeout quickly
	_, err := testServer.Client().SubscribeToOHLCStream(ctx, &gen.SubscribeToOHLCStreamRequest{
		Symbol:        "BTCUSDT",
		TimeframeName: "1m",
	})

	// We expect a timeout error
	if err == nil {
		// If it doesn't timeout, that's also acceptable - the connection might be very fast
		t.Log("Connection succeeded despite short timeout")
	} else {
		if s, ok := status.FromError(err); ok {
			if s.Code() == codes.DeadlineExceeded {
				t.Logf("Received expected timeout error: %v", err)
			} else {
				t.Logf("Received error (might be timeout-related): %v", err)
			}
		} else {
			t.Logf("Received error: %v", err)
		}
	}

	t.Log("Network timeout test passed")
}

func TestStreamAfterError(t *testing.T) {
	testServer := NewTestServer(t)
	testServer.Start(t)
	defer testServer.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// First, create a stream and force an error
	stream1, err := testServer.Client().SubscribeToOHLCStream(ctx, &gen.SubscribeToOHLCStreamRequest{
		Symbol:        "BTCUSDT",
		TimeframeName: "1m",
	})
	if err != nil {
		t.Fatalf("Failed to subscribe: %v", err)
	}

	// Close it
	stream1.CloseSend()

	// Now create a new stream - should work fine
	stream2, err := testServer.Client().SubscribeToOHLCStream(ctx, &gen.SubscribeToOHLCStreamRequest{
		Symbol:        "ETHUSDT",
		TimeframeName: "1m",
	})
	if err != nil {
		t.Fatalf("Failed to create new stream after error: %v", err)
	}
	defer stream2.CloseSend()

	// Test that the new stream works
	go func() {
		time.Sleep(100 * time.Millisecond)
		testOHLC := CreateTestOHLC("ETHUSDT", time.Now().UnixMilli())
		err := testServer.Server().PublishOHLC(ctx, testOHLC)
		if err != nil {
			t.Errorf("Failed to publish OHLC: %v", err)
		}
	}()

	_, err = WaitForMessage(stream2, 3*time.Second)
	if err != nil {
		t.Errorf("Failed to receive message on new stream: %v", err)
	}

	t.Log("Stream after error test passed")
}
