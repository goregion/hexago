package app

import (
	"context"
	"testing"

	"github.com/goregion/hexago/internal/entity"
	service_ohlc "github.com/goregion/hexago/internal/service/ohlc"
	unit "github.com/goregion/hexago/tests/unit"
)

func TestOHLCPublisher_NewOHLCPublisher(t *testing.T) {
	mockPublisher := unit.NewMockOHLCPublisher()

	publisher := service_ohlc.NewOHLCPublisher(mockPublisher)

	if publisher == nil {
		t.Fatal("OHLCPublisher should not be nil")
	}
}

func TestOHLCPublisher_ConsumeOHLC_SinglePublisher(t *testing.T) {
	mockPublisher := unit.NewMockOHLCPublisher()
	publisher := service_ohlc.NewOHLCPublisher(mockPublisher)

	ohlc := unit.CreateTestOHLC("BTCUSDT", 100.0, 105.0, 98.0, 102.0, 1234567890)

	ctx := context.Background()
	err := publisher.ConsumeOHLC(ctx, ohlc)

	if err != nil {
		t.Fatalf("ConsumeOHLC should not return error: %v", err)
	}

	if mockPublisher.CallCount != 1 {
		t.Errorf("Expected 1 call to PublishOHLC, got %d", mockPublisher.CallCount)
	}

	if len(mockPublisher.PublishedOHLCs) != 1 {
		t.Fatalf("Expected 1 published OHLC, got %d", len(mockPublisher.PublishedOHLCs))
	}

	publishedOHLC := mockPublisher.PublishedOHLCs[0]
	if publishedOHLC.Symbol != ohlc.Symbol {
		t.Errorf("Expected symbol %s, got %s", ohlc.Symbol, publishedOHLC.Symbol)
	}
	if publishedOHLC.Open != ohlc.Open {
		t.Errorf("Expected open %f, got %f", ohlc.Open, publishedOHLC.Open)
	}
	if publishedOHLC.High != ohlc.High {
		t.Errorf("Expected high %f, got %f", ohlc.High, publishedOHLC.High)
	}
	if publishedOHLC.Low != ohlc.Low {
		t.Errorf("Expected low %f, got %f", ohlc.Low, publishedOHLC.Low)
	}
	if publishedOHLC.Close != ohlc.Close {
		t.Errorf("Expected close %f, got %f", ohlc.Close, publishedOHLC.Close)
	}
	if publishedOHLC.TimestampMs != ohlc.TimestampMs {
		t.Errorf("Expected close time %d, got %d", ohlc.TimestampMs, publishedOHLC.TimestampMs)
	}
}

func TestOHLCPublisher_ConsumeOHLC_MultiplePublishers(t *testing.T) {
	mockPublisher := unit.NewMockOHLCPublisher()

	publisher := service_ohlc.NewOHLCPublisher(mockPublisher)

	ohlc := unit.CreateTestOHLC("ETHUSDT", 200.0, 220.0, 190.0, 210.0, 1234567891)

	ctx := context.Background()
	err := publisher.ConsumeOHLC(ctx, ohlc)

	if err != nil {
		t.Fatalf("ConsumeOHLC should not return error: %v", err)
	}

	// Publisher should be called
	if mockPublisher.CallCount != 1 {
		t.Errorf("Expected 1 call to publisher, got %d", mockPublisher.CallCount)
	}

	// Publisher should receive the correct data
	if len(mockPublisher.PublishedOHLCs) != 1 {
		t.Fatalf("Publisher should have 1 published OHLC, got %d", len(mockPublisher.PublishedOHLCs))
	}

	publishedOHLC := mockPublisher.PublishedOHLCs[0]
	if publishedOHLC.Symbol != ohlc.Symbol {
		t.Errorf("Expected symbol %s, got %s", ohlc.Symbol, publishedOHLC.Symbol)
	}
}

func TestOHLCPublisher_ConsumeOHLC_NoPublishers(t *testing.T) {
	// This test is no longer valid since publisher requires at least one publisher
	t.Skip("Skipping test as NewOHLCPublisher now requires at least one publisher")
}

func TestOHLCPublisher_ConsumeOHLC_PublisherError(t *testing.T) {
	mockPublisher := unit.NewMockOHLCPublisher()

	// Make publisher return error
	mockPublisher.ShouldError = true
	mockPublisher.ErrorMessage = "publisher error"

	publisher := service_ohlc.NewOHLCPublisher(mockPublisher)

	ohlc := unit.CreateTestOHLC("BTCUSDT", 100.0, 105.0, 98.0, 102.0, 1234567893)

	ctx := context.Background()
	err := publisher.ConsumeOHLC(ctx, ohlc)

	if err == nil {
		t.Fatal("ConsumeOHLC should return error when publisher fails")
	}

	if mockPublisher.CallCount != 1 {
		t.Errorf("Expected 1 call to publisher, got %d", mockPublisher.CallCount)
	}
}

func TestOHLCPublisher_ConsumeOHLC_MultipleOHLCs(t *testing.T) {
	mockPublisher := unit.NewMockOHLCPublisher()
	publisher := service_ohlc.NewOHLCPublisher(mockPublisher)

	ohlcs := []*entity.OHLC{
		unit.CreateTestOHLC("BTCUSDT", 100.0, 105.0, 98.0, 102.0, 1000),
		unit.CreateTestOHLC("ETHUSDT", 200.0, 220.0, 190.0, 210.0, 2000),
		unit.CreateTestOHLC("ADAUSDT", 1.0, 1.2, 0.9, 1.1, 3000),
	}

	ctx := context.Background()
	for i, ohlc := range ohlcs {
		err := publisher.ConsumeOHLC(ctx, ohlc)
		if err != nil {
			t.Fatalf("ConsumeOHLC %d should not return error: %v", i, err)
		}
	}

	if mockPublisher.CallCount != 3 {
		t.Errorf("Expected 3 calls to PublishOHLC, got %d", mockPublisher.CallCount)
	}

	if len(mockPublisher.PublishedOHLCs) != 3 {
		t.Fatalf("Expected 3 published OHLCs, got %d", len(mockPublisher.PublishedOHLCs))
	}

	// Verify all OHLCs were published correctly
	for i, expectedOHLC := range ohlcs {
		publishedOHLC := mockPublisher.PublishedOHLCs[i]
		if publishedOHLC.Symbol != expectedOHLC.Symbol {
			t.Errorf("OHLC %d: Expected symbol %s, got %s", i, expectedOHLC.Symbol, publishedOHLC.Symbol)
		}
		if publishedOHLC.TimestampMs != expectedOHLC.TimestampMs {
			t.Errorf("OHLC %d: Expected close time %d, got %d", i, expectedOHLC.TimestampMs, publishedOHLC.TimestampMs)
		}
	}
}

func TestOHLCPublisher_ConsumeOHLC_NilOHLC(t *testing.T) {
	mockPublisher := unit.NewMockOHLCPublisher()
	publisher := service_ohlc.NewOHLCPublisher(mockPublisher)

	ctx := context.Background()
	err := publisher.ConsumeOHLC(ctx, nil)

	// Should return an error when nil OHLC is provided
	if err == nil {
		t.Error("Expected error when consuming nil OHLC, got nil")
	}

	if mockPublisher.CallCount != 1 {
		t.Errorf("Expected 1 call to PublishOHLC even with nil, got %d", mockPublisher.CallCount)
	}
}
