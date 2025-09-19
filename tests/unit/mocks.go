package unit

import (
	"context"
	"errors"

	"github.com/goregion/hexago/internal/entity"
)

// MockOHLCPublisher is a mock implementation of port.OHLCPublisher
type MockOHLCPublisher struct {
	PublishOHLCFunc func(ctx context.Context, ohlc *entity.OHLC) error
	PublishedOHLCs  []*entity.OHLC
	CallCount       int
	ShouldError     bool
	ErrorMessage    string
}

func NewMockOHLCPublisher() *MockOHLCPublisher {
	mock := &MockOHLCPublisher{
		PublishedOHLCs: make([]*entity.OHLC, 0),
	}
	mock.PublishOHLCFunc = mock.defaultPublishOHLC
	return mock
}

func (m *MockOHLCPublisher) PublishOHLC(ctx context.Context, ohlc *entity.OHLC) error {
	return m.PublishOHLCFunc(ctx, ohlc)
}

func (m *MockOHLCPublisher) defaultPublishOHLC(ctx context.Context, ohlc *entity.OHLC) error {
	m.CallCount++
	if m.ShouldError {
		return errors.New(m.ErrorMessage)
	}
	// Handle nil OHLC gracefully
	if ohlc == nil {
		return errors.New("nil OHLC provided")
	}
	// Create a copy to avoid reference issues
	ohlcCopy := *ohlc
	m.PublishedOHLCs = append(m.PublishedOHLCs, &ohlcCopy)
	return nil
}

func (m *MockOHLCPublisher) Reset() {
	m.PublishedOHLCs = make([]*entity.OHLC, 0)
	m.CallCount = 0
	m.ShouldError = false
	m.ErrorMessage = ""
}

// MockTickPublisher is a mock implementation of port.TickPublisher
type MockTickPublisher struct {
	PublishTickFunc func(ctx context.Context, tick *entity.Tick) error
	PublishedTicks  []*entity.Tick
	CallCount       int
	ShouldError     bool
	ErrorMessage    string
}

func NewMockTickPublisher() *MockTickPublisher {
	mock := &MockTickPublisher{
		PublishedTicks: make([]*entity.Tick, 0),
	}
	mock.PublishTickFunc = mock.defaultPublishTick
	return mock
}

func (m *MockTickPublisher) PublishTick(ctx context.Context, tick *entity.Tick) error {
	return m.PublishTickFunc(ctx, tick)
}

func (m *MockTickPublisher) defaultPublishTick(ctx context.Context, tick *entity.Tick) error {
	m.CallCount++
	if m.ShouldError {
		return errors.New(m.ErrorMessage)
	}
	// Create a copy to avoid reference issues
	tickCopy := *tick
	m.PublishedTicks = append(m.PublishedTicks, &tickCopy)
	return nil
}

func (m *MockTickPublisher) Reset() {
	m.PublishedTicks = make([]*entity.Tick, 0)
	m.CallCount = 0
	m.ShouldError = false
	m.ErrorMessage = ""
}

// MockLPTickConsumer is a mock implementation of port.LPTickConsumer
type MockLPTickConsumer struct {
	ConsumeLPTickFunc func(ctx context.Context, lpTick *entity.LPTick) error
	ConsumedLPTicks   []*entity.LPTick
	CallCount         int
	ShouldError       bool
	ErrorMessage      string
}

func NewMockLPTickConsumer() *MockLPTickConsumer {
	mock := &MockLPTickConsumer{
		ConsumedLPTicks: make([]*entity.LPTick, 0),
	}
	mock.ConsumeLPTickFunc = mock.defaultConsumeLPTick
	return mock
}

func (m *MockLPTickConsumer) ConsumeLPTick(ctx context.Context, lpTick *entity.LPTick) error {
	return m.ConsumeLPTickFunc(ctx, lpTick)
}

func (m *MockLPTickConsumer) defaultConsumeLPTick(ctx context.Context, lpTick *entity.LPTick) error {
	m.CallCount++
	if m.ShouldError {
		return errors.New(m.ErrorMessage)
	}
	// Create a copy to avoid reference issues
	lpTickCopy := *lpTick
	m.ConsumedLPTicks = append(m.ConsumedLPTicks, &lpTickCopy)
	return nil
}

func (m *MockLPTickConsumer) Reset() {
	m.ConsumedLPTicks = make([]*entity.LPTick, 0)
	m.CallCount = 0
	m.ShouldError = false
	m.ErrorMessage = ""
}

// Test data creation helpers

// CreateTestTick creates a test tick with specified values
func CreateTestTick(symbol string, bidPrice, askPrice float64, timestampMs int64) *entity.Tick {
	return &entity.Tick{
		Symbol:       symbol,
		BestBidPrice: bidPrice,
		BestAskPrice: askPrice,
		TimestampMs:  timestampMs,
	}
}

// CreateTestLPTick creates a test LP tick with specified values
func CreateTestLPTick(symbol string, bidPrice, askPrice float64) *entity.LPTick {
	return &entity.LPTick{
		Symbol:       symbol,
		BestBidPrice: bidPrice,
		BestAskPrice: askPrice,
	}
}

// CreateTestTickRange creates a test tick range with specified ticks
func CreateTestTickRange(symbol string, fromMs, toMs int64, ticks []*entity.Tick) *entity.TickRange {
	return &entity.TickRange{
		Symbol:    symbol,
		FromMs:    fromMs,
		ToMs:      toMs,
		TickSlice: ticks,
	}
}

// CreateTestOHLC creates a test OHLC with specified values
func CreateTestOHLC(symbol string, open, high, low, close float64, closeTimeMs int64) *entity.OHLC {
	return &entity.OHLC{
		Symbol:      symbol,
		Open:        open,
		High:        high,
		Low:         low,
		Close:       close,
		CloseTimeMs: closeTimeMs,
	}
}
