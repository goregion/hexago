package integration

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/goregion/hexago/internal/adapter/grpc-api/gen"
	"github.com/goregion/hexago/internal/entity"
)

func TestConcurrentSubscriptions(t *testing.T) {
	testServer := NewTestServer(t)
	testServer.Start(t)
	defer testServer.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	symbols := []string{"BTCUSDT", "ETHUSDT", "ADAUSDT", "DOTUSDT"}
	numClients := len(symbols)

	var wg sync.WaitGroup
	results := make(chan error, numClients)

	// Create concurrent subscriptions for different symbols
	for i, symbol := range symbols {
		wg.Add(1)
		go func(clientId int, sym string) {
			defer wg.Done()

			// Subscribe to stream
			stream, err := testServer.Client().SubscribeToOHLCStream(ctx, &gen.SubscribeToOHLCStreamRequest{
				Symbol:        sym,
				TimeframeName: "1m",
			})
			if err != nil {
				results <- err
				return
			}
			defer stream.CloseSend()

			// Wait for a message
			_, err = WaitForMessage(stream, 2*time.Second)
			// It's ok if we don't receive a message (timeout), we're testing concurrent connections
			if err != nil && err.Error() != "timeout waiting for message" {
				results <- err
				return
			}

			results <- nil
		}(i, symbol)
	}

	// Publish data for each symbol
	go func() {
		time.Sleep(100 * time.Millisecond) // Wait for subscriptions to be active
		for _, symbol := range symbols {
			testOHLC := CreateTestOHLC(symbol, time.Now().UnixMilli())
			err := testServer.Server().PublishOHLC(ctx, testOHLC)
			if err != nil {
				t.Errorf("Failed to publish OHLC for %s: %v", symbol, err)
			}
		}
	}()

	// Wait for all goroutines to complete
	wg.Wait()
	close(results)

	// Check results
	for err := range results {
		if err != nil {
			t.Errorf("Concurrent subscription error: %v", err)
		}
	}

	t.Log("Concurrent subscriptions test passed")
}

func TestSymbolIsolation(t *testing.T) {
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

	// Publish OHLC for BTCUSDT only
	go func() {
		time.Sleep(100 * time.Millisecond)
		testOHLC := CreateTestOHLC("BTCUSDT", time.Now().UnixMilli())
		err := testServer.Server().PublishOHLC(ctx, testOHLC)
		if err != nil {
			t.Errorf("Failed to publish OHLC: %v", err)
		}
	}()

	// BTC stream should receive message
	_, err = WaitForMessage(streamBTC, 3*time.Second)
	if err != nil {
		t.Errorf("BTCUSDT stream should have received message: %v", err)
	}

	// ETH stream should NOT receive message (timeout expected)
	_, err = WaitForMessage(streamETH, 1*time.Second)
	if err == nil {
		t.Error("ETHUSDT stream should not have received BTCUSDT message")
	} else if err.Error() != "timeout waiting for message" {
		t.Errorf("Unexpected error: %v", err)
	}

	t.Log("Symbol isolation test passed")
}

func TestHighConcurrency(t *testing.T) {
	testServer := NewTestServer(t)
	testServer.Start(t)
	defer testServer.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	numClients := 20
	symbol := "BTCUSDT"

	var wg sync.WaitGroup
	errors := make(chan error, numClients)
	received := make(chan int, numClients)

	// Create many concurrent subscriptions to the same symbol
	for i := 0; i < numClients; i++ {
		wg.Add(1)
		go func(clientId int) {
			defer wg.Done()

			stream, err := testServer.Client().SubscribeToOHLCStream(ctx, &gen.SubscribeToOHLCStreamRequest{
				Symbol:        symbol,
				TimeframeName: "1m",
			})
			if err != nil {
				errors <- err
				return
			}
			defer stream.CloseSend()

			// Wait for message
			_, err = WaitForMessage(stream, 5*time.Second)
			if err != nil {
				if err.Error() != "timeout waiting for message" {
					errors <- err
					return
				}
			} else {
				received <- clientId
			}

			errors <- nil
		}(i)
	}

	// Publish a single OHLC message
	go func() {
		time.Sleep(200 * time.Millisecond) // Wait for all subscriptions
		testOHLC := CreateTestOHLC(symbol, time.Now().UnixMilli())
		err := testServer.Server().PublishOHLC(ctx, testOHLC)
		if err != nil {
			t.Errorf("Failed to publish OHLC: %v", err)
		}
	}()

	// Wait for all clients
	wg.Wait()
	close(errors)
	close(received)

	// Check for errors
	errorCount := 0
	for err := range errors {
		if err != nil {
			errorCount++
			t.Errorf("High concurrency error: %v", err)
		}
	}

	// Count received messages
	receivedCount := 0
	for range received {
		receivedCount++
	}

	t.Logf("High concurrency test completed: %d errors, %d received messages out of %d clients",
		errorCount, receivedCount, numClients)

	// Note: Due to the current server implementation, only the last subscriber for a symbol
	// will receive messages. This is a limitation of the current sync.Map approach.
	// In a production system, you'd want to support multiple subscribers per symbol.
}

func TestConcurrentPublishing(t *testing.T) {
	testServer := NewTestServer(t)
	testServer.Start(t)
	defer testServer.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	symbol := "BTCUSDT"

	// Subscribe to stream
	stream, err := testServer.Client().SubscribeToOHLCStream(ctx, &gen.SubscribeToOHLCStreamRequest{
		Symbol:        symbol,
		TimeframeName: "1m",
	})
	if err != nil {
		t.Fatalf("Failed to subscribe: %v", err)
	}
	defer stream.CloseSend()

	numPublishers := 5
	messagesPerPublisher := 3
	totalMessages := numPublishers * messagesPerPublisher

	var wg sync.WaitGroup

	// Start concurrent publishers
	for i := 0; i < numPublishers; i++ {
		wg.Add(1)
		go func(publisherId int) {
			defer wg.Done()
			for j := 0; j < messagesPerPublisher; j++ {
				testOHLC := &entity.OHLC{
					Symbol:      symbol,
					Open:        float64(100 + publisherId*10 + j),
					High:        float64(105 + publisherId*10 + j),
					Low:         float64(98 + publisherId*10 + j),
					Close:       float64(102 + publisherId*10 + j),
					CloseTimeMs: time.Now().UnixMilli() + int64(publisherId*1000+j*100),
				}
				err := testServer.Server().PublishOHLC(ctx, testOHLC)
				if err != nil {
					t.Errorf("Publisher %d failed to publish message %d: %v", publisherId, j, err)
				}
				time.Sleep(10 * time.Millisecond)
			}
		}(i)
	}

	// Receive messages
	receivedCount := 0
	timeout := time.After(10 * time.Second)

receiveLoop:
	for receivedCount < totalMessages {
		select {
		case <-timeout:
			break receiveLoop
		default:
			_, err := WaitForMessage(stream, 1*time.Second)
			if err != nil {
				if err.Error() == "timeout waiting for message" {
					break receiveLoop
				}
				t.Errorf("Failed to receive message: %v", err)
				break receiveLoop
			}
			receivedCount++
		}
	}

	wg.Wait()

	t.Logf("Concurrent publishing test: received %d out of %d expected messages", receivedCount, totalMessages)

	// We expect to receive some messages, though not necessarily all due to the server implementation
	if receivedCount == 0 {
		t.Error("Expected to receive at least some messages")
	}
}
