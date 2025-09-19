package app

import (
	"context"
	"testing"

	"github.com/goregion/hexago/internal/app"
	unit "github.com/goregion/hexago/tests/unit"
)

func TestTickProcessor_NewTickProcessor(t *testing.T) {
	mockPublisher1 := unit.NewMockTickPublisher()
	mockPublisher2 := unit.NewMockTickPublisher()

	processor := app.NewTickProcessor(mockPublisher1, mockPublisher2)

	if processor == nil {
		t.Fatal("TickProcessor should not be nil")
	}
}

func TestTickProcessor_ConsumeLPTick_FirstTime(t *testing.T) {
	mockPublisher := unit.NewMockTickPublisher()
	processor := app.NewTickProcessor(mockPublisher)

	lpTick := unit.CreateTestLPTick("BTCUSDT", 100.0, 101.0)

	ctx := context.Background()
	err := processor.ConsumeLPTick(ctx, lpTick)

	if err != nil {
		t.Fatalf("ConsumeLPTick should not return error: %v", err)
	}

	if mockPublisher.CallCount != 1 {
		t.Errorf("Expected 1 call to PublishTick, got %d", mockPublisher.CallCount)
	}

	if len(mockPublisher.PublishedTicks) != 1 {
		t.Fatalf("Expected 1 published tick, got %d", len(mockPublisher.PublishedTicks))
	}

	publishedTick := mockPublisher.PublishedTicks[0]
	if publishedTick.Symbol != lpTick.Symbol {
		t.Errorf("Expected symbol %s, got %s", lpTick.Symbol, publishedTick.Symbol)
	}
	if publishedTick.BestBidPrice != lpTick.BestBidPrice {
		t.Errorf("Expected bid price %f, got %f", lpTick.BestBidPrice, publishedTick.BestBidPrice)
	}
	if publishedTick.BestAskPrice != lpTick.BestAskPrice {
		t.Errorf("Expected ask price %f, got %f", lpTick.BestAskPrice, publishedTick.BestAskPrice)
	}
	if publishedTick.TimestampMs == 0 {
		t.Error("Published tick should have non-zero timestamp")
	}
}

func TestTickProcessor_ConsumeLPTick_SamePrices_NoDuplicate(t *testing.T) {
	mockPublisher := unit.NewMockTickPublisher()
	processor := app.NewTickProcessor(mockPublisher)

	lpTick := unit.CreateTestLPTick("BTCUSDT", 100.0, 101.0)

	ctx := context.Background()

	// First consumption - should publish
	err := processor.ConsumeLPTick(ctx, lpTick)
	if err != nil {
		t.Fatalf("First ConsumeLPTick should not return error: %v", err)
	}

	// Second consumption with same prices - should NOT publish
	err = processor.ConsumeLPTick(ctx, lpTick)
	if err != nil {
		t.Fatalf("Second ConsumeLPTick should not return error: %v", err)
	}

	// Should only have published once
	if mockPublisher.CallCount != 1 {
		t.Errorf("Expected 1 call to PublishTick for duplicate prices, got %d", mockPublisher.CallCount)
	}

	if len(mockPublisher.PublishedTicks) != 1 {
		t.Errorf("Expected 1 published tick for duplicate prices, got %d", len(mockPublisher.PublishedTicks))
	}
}

func TestTickProcessor_ConsumeLPTick_DifferentPrices_ShouldPublish(t *testing.T) {
	mockPublisher := unit.NewMockTickPublisher()
	processor := app.NewTickProcessor(mockPublisher)

	lpTick1 := unit.CreateTestLPTick("BTCUSDT", 100.0, 101.0)
	lpTick2 := unit.CreateTestLPTick("BTCUSDT", 102.0, 103.0) // Different prices

	ctx := context.Background()

	// First consumption
	err := processor.ConsumeLPTick(ctx, lpTick1)
	if err != nil {
		t.Fatalf("First ConsumeLPTick should not return error: %v", err)
	}

	// Second consumption with different prices - should publish
	err = processor.ConsumeLPTick(ctx, lpTick2)
	if err != nil {
		t.Fatalf("Second ConsumeLPTick should not return error: %v", err)
	}

	// Should have published twice
	if mockPublisher.CallCount != 2 {
		t.Errorf("Expected 2 calls to PublishTick for different prices, got %d", mockPublisher.CallCount)
	}

	if len(mockPublisher.PublishedTicks) != 2 {
		t.Errorf("Expected 2 published ticks for different prices, got %d", len(mockPublisher.PublishedTicks))
	}

	// Verify both ticks were published correctly
	tick1 := mockPublisher.PublishedTicks[0]
	tick2 := mockPublisher.PublishedTicks[1]

	if tick1.BestBidPrice != 100.0 || tick1.BestAskPrice != 101.0 {
		t.Errorf("First tick prices incorrect: bid=%f, ask=%f", tick1.BestBidPrice, tick1.BestAskPrice)
	}
	if tick2.BestBidPrice != 102.0 || tick2.BestAskPrice != 103.0 {
		t.Errorf("Second tick prices incorrect: bid=%f, ask=%f", tick2.BestBidPrice, tick2.BestAskPrice)
	}
}

func TestTickProcessor_ConsumeLPTick_OnlyBidChanges(t *testing.T) {
	mockPublisher := unit.NewMockTickPublisher()
	processor := app.NewTickProcessor(mockPublisher)

	lpTick1 := unit.CreateTestLPTick("ETHUSDT", 200.0, 201.0)
	lpTick2 := unit.CreateTestLPTick("ETHUSDT", 199.0, 201.0) // Only bid changes

	ctx := context.Background()

	err := processor.ConsumeLPTick(ctx, lpTick1)
	if err != nil {
		t.Fatalf("First ConsumeLPTick should not return error: %v", err)
	}

	err = processor.ConsumeLPTick(ctx, lpTick2)
	if err != nil {
		t.Fatalf("Second ConsumeLPTick should not return error: %v", err)
	}

	// Should publish both times since bid price changed
	if mockPublisher.CallCount != 2 {
		t.Errorf("Expected 2 calls to PublishTick when bid changes, got %d", mockPublisher.CallCount)
	}
}

func TestTickProcessor_ConsumeLPTick_OnlyAskChanges(t *testing.T) {
	mockPublisher := unit.NewMockTickPublisher()
	processor := app.NewTickProcessor(mockPublisher)

	lpTick1 := unit.CreateTestLPTick("ADAUSDT", 1.0, 1.1)
	lpTick2 := unit.CreateTestLPTick("ADAUSDT", 1.0, 1.2) // Only ask changes

	ctx := context.Background()

	err := processor.ConsumeLPTick(ctx, lpTick1)
	if err != nil {
		t.Fatalf("First ConsumeLPTick should not return error: %v", err)
	}

	err = processor.ConsumeLPTick(ctx, lpTick2)
	if err != nil {
		t.Fatalf("Second ConsumeLPTick should not return error: %v", err)
	}

	// Should publish both times since ask price changed
	if mockPublisher.CallCount != 2 {
		t.Errorf("Expected 2 calls to PublishTick when ask changes, got %d", mockPublisher.CallCount)
	}
}

func TestTickProcessor_ConsumeLPTick_DifferentSymbols(t *testing.T) {
	mockPublisher := unit.NewMockTickPublisher()
	processor := app.NewTickProcessor(mockPublisher)

	lpTick1 := unit.CreateTestLPTick("BTCUSDT", 100.0, 101.0)
	lpTick2 := unit.CreateTestLPTick("ETHUSDT", 100.0, 101.0) // Same prices, different symbol

	ctx := context.Background()

	err := processor.ConsumeLPTick(ctx, lpTick1)
	if err != nil {
		t.Fatalf("First ConsumeLPTick should not return error: %v", err)
	}

	err = processor.ConsumeLPTick(ctx, lpTick2)
	if err != nil {
		t.Fatalf("Second ConsumeLPTick should not return error: %v", err)
	}

	// Should publish both times since symbols are different
	if mockPublisher.CallCount != 2 {
		t.Errorf("Expected 2 calls to PublishTick for different symbols, got %d", mockPublisher.CallCount)
	}

	// Verify symbols are correct
	tick1 := mockPublisher.PublishedTicks[0]
	tick2 := mockPublisher.PublishedTicks[1]

	if tick1.Symbol != "BTCUSDT" {
		t.Errorf("Expected first tick symbol BTCUSDT, got %s", tick1.Symbol)
	}
	if tick2.Symbol != "ETHUSDT" {
		t.Errorf("Expected second tick symbol ETHUSDT, got %s", tick2.Symbol)
	}
}

func TestTickProcessor_ConsumeLPTick_MultiplePublishers(t *testing.T) {
	mockPublisher1 := unit.NewMockTickPublisher()
	mockPublisher2 := unit.NewMockTickPublisher()
	processor := app.NewTickProcessor(mockPublisher1, mockPublisher2)

	lpTick := unit.CreateTestLPTick("BTCUSDT", 100.0, 101.0)

	ctx := context.Background()
	err := processor.ConsumeLPTick(ctx, lpTick)

	if err != nil {
		t.Fatalf("ConsumeLPTick should not return error: %v", err)
	}

	// Both publishers should be called
	if mockPublisher1.CallCount != 1 {
		t.Errorf("Expected 1 call to first publisher, got %d", mockPublisher1.CallCount)
	}
	if mockPublisher2.CallCount != 1 {
		t.Errorf("Expected 1 call to second publisher, got %d", mockPublisher2.CallCount)
	}

	// Both should have published the same tick data
	if len(mockPublisher1.PublishedTicks) != 1 || len(mockPublisher2.PublishedTicks) != 1 {
		t.Fatal("Both publishers should have published tick data")
	}

	tick1 := mockPublisher1.PublishedTicks[0]
	tick2 := mockPublisher2.PublishedTicks[0]

	if tick1.Symbol != tick2.Symbol || tick1.BestBidPrice != tick2.BestBidPrice ||
		tick1.BestAskPrice != tick2.BestAskPrice {
		t.Error("Both publishers should publish identical tick data")
	}
}

func TestTickProcessor_ConsumeLPTick_PublisherError(t *testing.T) {
	mockPublisher1 := unit.NewMockTickPublisher()
	mockPublisher2 := unit.NewMockTickPublisher()

	// Make first publisher return error
	mockPublisher1.ShouldError = true
	mockPublisher1.ErrorMessage = "publisher error"

	processor := app.NewTickProcessor(mockPublisher1, mockPublisher2)

	lpTick := unit.CreateTestLPTick("BTCUSDT", 100.0, 101.0)

	ctx := context.Background()
	err := processor.ConsumeLPTick(ctx, lpTick)

	if err == nil {
		t.Fatal("ConsumeLPTick should return error when publisher fails")
	}

	if mockPublisher1.CallCount != 1 {
		t.Errorf("Expected 1 call to first publisher, got %d", mockPublisher1.CallCount)
	}

	// Second publisher should not be called due to error in first
	if mockPublisher2.CallCount != 0 {
		t.Errorf("Expected 0 calls to second publisher after first fails, got %d", mockPublisher2.CallCount)
	}
}

func TestTickProcessor_ConsumeLPTick_NoPublishers(t *testing.T) {
	processor := app.NewTickProcessor() // No publishers

	lpTick := unit.CreateTestLPTick("BTCUSDT", 100.0, 101.0)

	ctx := context.Background()
	err := processor.ConsumeLPTick(ctx, lpTick)

	// Should not return error even with no publishers
	if err != nil {
		t.Fatalf("ConsumeLPTick should not return error with no publishers: %v", err)
	}
}

func TestTickProcessor_ConsumeLPTick_ZeroPrices(t *testing.T) {
	mockPublisher := unit.NewMockTickPublisher()
	processor := app.NewTickProcessor(mockPublisher)

	lpTick := unit.CreateTestLPTick("TESTUSDT", 0.0, 0.0)

	ctx := context.Background()
	err := processor.ConsumeLPTick(ctx, lpTick)

	if err != nil {
		t.Fatalf("ConsumeLPTick should handle zero prices: %v", err)
	}

	if mockPublisher.CallCount != 1 {
		t.Errorf("Expected 1 call to PublishTick with zero prices, got %d", mockPublisher.CallCount)
	}

	publishedTick := mockPublisher.PublishedTicks[0]
	if publishedTick.BestBidPrice != 0.0 || publishedTick.BestAskPrice != 0.0 {
		t.Errorf("Expected zero prices in published tick, got bid=%f, ask=%f",
			publishedTick.BestBidPrice, publishedTick.BestAskPrice)
	}
}

func TestTickProcessor_ConsumeLPTick_ConcurrentSymbols(t *testing.T) {
	mockPublisher := unit.NewMockTickPublisher()
	processor := app.NewTickProcessor(mockPublisher)

	symbols := []string{"BTCUSDT", "ETHUSDT", "ADAUSDT", "DOTUSDT"}

	ctx := context.Background()

	// Process ticks for different symbols
	for i, symbol := range symbols {
		lpTick := unit.CreateTestLPTick(symbol, float64(100+i), float64(101+i))
		err := processor.ConsumeLPTick(ctx, lpTick)
		if err != nil {
			t.Fatalf("ConsumeLPTick for %s should not return error: %v", symbol, err)
		}
	}

	// Should publish for each symbol
	if mockPublisher.CallCount != len(symbols) {
		t.Errorf("Expected %d calls to PublishTick for %d symbols, got %d",
			len(symbols), len(symbols), mockPublisher.CallCount)
	}

	if len(mockPublisher.PublishedTicks) != len(symbols) {
		t.Errorf("Expected %d published ticks for %d symbols, got %d",
			len(symbols), len(symbols), len(mockPublisher.PublishedTicks))
	}

	// Verify each symbol was published correctly
	for i, symbol := range symbols {
		tick := mockPublisher.PublishedTicks[i]
		if tick.Symbol != symbol {
			t.Errorf("Tick %d: expected symbol %s, got %s", i, symbol, tick.Symbol)
		}
		expectedBid := float64(100 + i)
		expectedAsk := float64(101 + i)
		if tick.BestBidPrice != expectedBid || tick.BestAskPrice != expectedAsk {
			t.Errorf("Tick %d: expected prices bid=%f, ask=%f, got bid=%f, ask=%f",
				i, expectedBid, expectedAsk, tick.BestBidPrice, tick.BestAskPrice)
		}
	}
}
