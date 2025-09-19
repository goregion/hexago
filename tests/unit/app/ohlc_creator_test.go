package app

import (
	"context"
	"testing"

	"github.com/goregion/hexago/internal/app"
	"github.com/goregion/hexago/internal/entity"
	unit "github.com/goregion/hexago/tests/unit"
)

func TestOHLCCreator_NewOHLCCreator(t *testing.T) {
	mockPublisher1 := unit.NewMockOHLCPublisher()
	mockPublisher2 := unit.NewMockOHLCPublisher()

	creator := app.NewOHLCCreator(app.USE_BID_PRICE, mockPublisher1, mockPublisher2)

	if creator == nil {
		t.Fatal("OHLCCreator should not be nil")
	}
}

func TestOHLCCreator_ConsumeTickRange_WithBidPrice(t *testing.T) {
	mockPublisher := unit.NewMockOHLCPublisher()
	creator := app.NewOHLCCreator(app.USE_BID_PRICE, mockPublisher)

	// Create test ticks with bid prices
	ticks := []*entity.Tick{
		unit.CreateTestTick("BTCUSDT", 100.0, 101.0, 1000), // Open: 100.0
		unit.CreateTestTick("BTCUSDT", 105.0, 106.0, 2000), // High: 105.0
		unit.CreateTestTick("BTCUSDT", 95.0, 96.0, 3000),   // Low: 95.0
		unit.CreateTestTick("BTCUSDT", 102.0, 103.0, 4000), // Close: 102.0
	}

	tickRange := unit.CreateTestTickRange("BTCUSDT", 1000, 5000, ticks)

	ctx := context.Background()
	err := creator.ConsumeTickRange(ctx, tickRange)

	if err != nil {
		t.Fatalf("ConsumeTickRange should not return error: %v", err)
	}

	if mockPublisher.CallCount != 1 {
		t.Errorf("Expected 1 call to PublishOHLC, got %d", mockPublisher.CallCount)
	}

	if len(mockPublisher.PublishedOHLCs) != 1 {
		t.Fatalf("Expected 1 published OHLC, got %d", len(mockPublisher.PublishedOHLCs))
	}

	ohlc := mockPublisher.PublishedOHLCs[0]
	if ohlc.Symbol != "BTCUSDT" {
		t.Errorf("Expected symbol BTCUSDT, got %s", ohlc.Symbol)
	}
	if ohlc.Open != 100.0 {
		t.Errorf("Expected open 100.0, got %f", ohlc.Open)
	}
	if ohlc.High != 105.0 {
		t.Errorf("Expected high 105.0, got %f", ohlc.High)
	}
	if ohlc.Low != 95.0 {
		t.Errorf("Expected low 95.0, got %f", ohlc.Low)
	}
	if ohlc.Close != 102.0 {
		t.Errorf("Expected close 102.0, got %f", ohlc.Close)
	}
	if ohlc.TimestampMs != 5000 {
		t.Errorf("Expected close time 5000, got %d", ohlc.TimestampMs)
	}
}

func TestOHLCCreator_ConsumeTickRange_WithAskPrice(t *testing.T) {
	mockPublisher := unit.NewMockOHLCPublisher()
	creator := app.NewOHLCCreator(app.USE_ASK_PRICE, mockPublisher)

	// Create test ticks with ask prices
	ticks := []*entity.Tick{
		unit.CreateTestTick("ETHUSDT", 200.0, 201.0, 1000), // Open: 201.0 (ask)
		unit.CreateTestTick("ETHUSDT", 205.0, 210.0, 2000), // High: 210.0 (ask)
		unit.CreateTestTick("ETHUSDT", 190.0, 195.0, 3000), // Low: 195.0 (ask)
		unit.CreateTestTick("ETHUSDT", 198.0, 203.0, 4000), // Close: 203.0 (ask)
	}

	tickRange := unit.CreateTestTickRange("ETHUSDT", 1000, 5000, ticks)

	ctx := context.Background()
	err := creator.ConsumeTickRange(ctx, tickRange)

	if err != nil {
		t.Fatalf("ConsumeTickRange should not return error: %v", err)
	}

	ohlc := mockPublisher.PublishedOHLCs[0]
	if ohlc.Open != 201.0 {
		t.Errorf("Expected open 201.0 (ask), got %f", ohlc.Open)
	}
	if ohlc.High != 210.0 {
		t.Errorf("Expected high 210.0 (ask), got %f", ohlc.High)
	}
	if ohlc.Low != 195.0 {
		t.Errorf("Expected low 195.0 (ask), got %f", ohlc.Low)
	}
	if ohlc.Close != 203.0 {
		t.Errorf("Expected close 203.0 (ask), got %f", ohlc.Close)
	}
}

func TestOHLCCreator_ConsumeTickRange_EmptyTickRange(t *testing.T) {
	mockPublisher := unit.NewMockOHLCPublisher()
	creator := app.NewOHLCCreator(app.USE_BID_PRICE, mockPublisher)

	// Empty tick range
	tickRange := unit.CreateTestTickRange("BTCUSDT", 1000, 5000, []*entity.Tick{})

	ctx := context.Background()
	err := creator.ConsumeTickRange(ctx, tickRange)

	if err != nil {
		t.Fatalf("ConsumeTickRange should not return error for empty range: %v", err)
	}

	if mockPublisher.CallCount != 1 {
		t.Errorf("Expected 1 call to PublishOHLC even for empty range, got %d", mockPublisher.CallCount)
	}

	ohlc := mockPublisher.PublishedOHLCs[0]
	if ohlc.Open != 0 {
		t.Errorf("Expected open 0 for empty range, got %f", ohlc.Open)
	}
	if ohlc.High != 0 {
		t.Errorf("Expected high 0 for empty range, got %f", ohlc.High)
	}
	if ohlc.Low != 0 {
		t.Errorf("Expected low 0 for empty range, got %f", ohlc.Low)
	}
	if ohlc.Close != 0 {
		t.Errorf("Expected close 0 for empty range, got %f", ohlc.Close)
	}
}

func TestOHLCCreator_ConsumeTickRange_SingleTick(t *testing.T) {
	mockPublisher := unit.NewMockOHLCPublisher()
	creator := app.NewOHLCCreator(app.USE_BID_PRICE, mockPublisher)

	// Single tick
	ticks := []*entity.Tick{
		unit.CreateTestTick("ADAUSDT", 1.5, 1.6, 1000),
	}

	tickRange := unit.CreateTestTickRange("ADAUSDT", 1000, 2000, ticks)

	ctx := context.Background()
	err := creator.ConsumeTickRange(ctx, tickRange)

	if err != nil {
		t.Fatalf("ConsumeTickRange should not return error: %v", err)
	}

	ohlc := mockPublisher.PublishedOHLCs[0]
	// All OHLC values should be the same for single tick
	if ohlc.Open != 1.5 {
		t.Errorf("Expected open 1.5, got %f", ohlc.Open)
	}
	if ohlc.High != 1.5 {
		t.Errorf("Expected high 1.5, got %f", ohlc.High)
	}
	if ohlc.Low != 1.5 {
		t.Errorf("Expected low 1.5, got %f", ohlc.Low)
	}
	if ohlc.Close != 1.5 {
		t.Errorf("Expected close 1.5, got %f", ohlc.Close)
	}
}

func TestOHLCCreator_ConsumeTickRange_MultiplePublishers(t *testing.T) {
	mockPublisher1 := unit.NewMockOHLCPublisher()
	mockPublisher2 := unit.NewMockOHLCPublisher()
	creator := app.NewOHLCCreator(app.USE_BID_PRICE, mockPublisher1, mockPublisher2)

	ticks := []*entity.Tick{
		unit.CreateTestTick("BTCUSDT", 100.0, 101.0, 1000),
	}

	tickRange := unit.CreateTestTickRange("BTCUSDT", 1000, 2000, ticks)

	ctx := context.Background()
	err := creator.ConsumeTickRange(ctx, tickRange)

	if err != nil {
		t.Fatalf("ConsumeTickRange should not return error: %v", err)
	}

	// Both publishers should be called
	if mockPublisher1.CallCount != 1 {
		t.Errorf("Expected 1 call to first publisher, got %d", mockPublisher1.CallCount)
	}
	if mockPublisher2.CallCount != 1 {
		t.Errorf("Expected 1 call to second publisher, got %d", mockPublisher2.CallCount)
	}

	// Both should have the same OHLC data
	if len(mockPublisher1.PublishedOHLCs) != 1 || len(mockPublisher2.PublishedOHLCs) != 1 {
		t.Fatal("Both publishers should have received OHLC data")
	}

	ohlc1 := mockPublisher1.PublishedOHLCs[0]
	ohlc2 := mockPublisher2.PublishedOHLCs[0]

	if ohlc1.Open != ohlc2.Open || ohlc1.High != ohlc2.High ||
		ohlc1.Low != ohlc2.Low || ohlc1.Close != ohlc2.Close {
		t.Error("Both publishers should receive identical OHLC data")
	}
}

func TestOHLCCreator_ConsumeTickRange_PublisherError(t *testing.T) {
	mockPublisher1 := unit.NewMockOHLCPublisher()
	mockPublisher2 := unit.NewMockOHLCPublisher()

	// Make first publisher return error
	mockPublisher1.ShouldError = true
	mockPublisher1.ErrorMessage = "publisher error"

	creator := app.NewOHLCCreator(app.USE_BID_PRICE, mockPublisher1, mockPublisher2)

	ticks := []*entity.Tick{
		unit.CreateTestTick("BTCUSDT", 100.0, 101.0, 1000),
	}

	tickRange := unit.CreateTestTickRange("BTCUSDT", 1000, 2000, ticks)

	ctx := context.Background()
	err := creator.ConsumeTickRange(ctx, tickRange)

	if err == nil {
		t.Fatal("ConsumeTickRange should return error when publisher fails")
	}

	if mockPublisher1.CallCount != 1 {
		t.Errorf("Expected 1 call to first publisher, got %d", mockPublisher1.CallCount)
	}

	// Second publisher should not be called due to error in first
	if mockPublisher2.CallCount != 0 {
		t.Errorf("Expected 0 calls to second publisher after first fails, got %d", mockPublisher2.CallCount)
	}
}

func TestOHLCCreator_ConsumeTickRange_ZeroPrices(t *testing.T) {
	mockPublisher := unit.NewMockOHLCPublisher()
	creator := app.NewOHLCCreator(app.USE_BID_PRICE, mockPublisher)

	// Ticks with zero and negative prices
	ticks := []*entity.Tick{
		unit.CreateTestTick("TESTUSDT", 0.0, 0.1, 1000), // Bid=0, Ask=0.1
		unit.CreateTestTick("TESTUSDT", 5.0, 5.1, 2000), // Bid=5, Ask=5.1
	}

	tickRange := unit.CreateTestTickRange("TESTUSDT", 1000, 3000, ticks)

	ctx := context.Background()
	err := creator.ConsumeTickRange(ctx, tickRange)

	if err != nil {
		t.Fatalf("ConsumeTickRange should handle zero prices: %v", err)
	}

	ohlc := mockPublisher.PublishedOHLCs[0]
	// Open should be from first non-zero tick (bid=5.0)
	if ohlc.Open != 5.0 {
		t.Errorf("Expected open 5.0, got %f", ohlc.Open)
	}
	if ohlc.High != 5.0 {
		t.Errorf("Expected high 5.0, got %f", ohlc.High)
	}
	// Low should be from first non-zero tick (bid=5.0)
	if ohlc.Low != 5.0 {
		t.Errorf("Expected low 5.0, got %f", ohlc.Low)
	}
	if ohlc.Close != 5.0 {
		t.Errorf("Expected close 5.0, got %f", ohlc.Close)
	}
}
