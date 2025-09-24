package integration

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/goregion/hexago/internal/adapter/grpc-api/gen"
	adapter_grpc_api "github.com/goregion/hexago/internal/adapter/grpc-api/impl"
	"github.com/goregion/hexago/internal/entity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TestServer holds the test gRPC server and client
type TestServer struct {
	server   *adapter_grpc_api.Server
	grpcConn *grpc.ClientConn
	client   gen.OHLCServiceClient
	addr     string
	cancel   context.CancelFunc
}

// NewTestServer creates a new test server instance
func NewTestServer(t *testing.T) *TestServer {
	// Find available port
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Failed to find available port: %v", err)
	}
	addr := listener.Addr().String()
	listener.Close()

	ts := &TestServer{
		server: adapter_grpc_api.NewServer(addr),
		addr:   addr,
	}

	return ts
}

// Start starts the test server
func (ts *TestServer) Start(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	ts.cancel = cancel

	// Start server in goroutine
	go func() {
		if err := ts.server.Launch(ctx); err != nil && ctx.Err() == nil {
			t.Errorf("Server failed: %v", err)
		}
	}()

	// Wait for server to start
	time.Sleep(100 * time.Millisecond)

	// Create client connection
	conn, err := grpc.NewClient(ts.addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	ts.grpcConn = conn
	ts.client = gen.NewOHLCServiceClient(conn)
}

// Stop stops the test server
func (ts *TestServer) Stop() {
	if ts.grpcConn != nil {
		ts.grpcConn.Close()
	}
	if ts.cancel != nil {
		ts.cancel()
	}
}

// Client returns the gRPC client
func (ts *TestServer) Client() gen.OHLCServiceClient {
	return ts.client
}

// Server returns the gRPC server
func (ts *TestServer) Server() *adapter_grpc_api.Server {
	return ts.server
}

// CreateTestOHLC creates a test OHLC entity
func CreateTestOHLC(symbol string, timeMs int64) *entity.OHLC {
	return &entity.OHLC{
		Symbol:      symbol,
		Open:        100.0,
		High:        105.0,
		Low:         98.0,
		Close:       102.0,
		TimestampMs: timeMs,
	}
}

// WaitForMessage waits for a message on the stream with timeout
func WaitForMessage(stream gen.OHLCService_SubscribeToOHLCStreamClient, timeout time.Duration) (*gen.OHLC, error) {
	ch := make(chan *gen.OHLC, 1)
	errCh := make(chan error, 1)

	go func() {
		msg, err := stream.Recv()
		if err != nil {
			errCh <- err
			return
		}
		ch <- msg
	}()

	select {
	case msg := <-ch:
		return msg, nil
	case err := <-errCh:
		return nil, err
	case <-time.After(timeout):
		return nil, fmt.Errorf("timeout waiting for message")
	}
}
