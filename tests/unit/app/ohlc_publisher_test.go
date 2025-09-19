package app

import (
	"context"
	"testing"

	"github.com/goregion/hexago/internal/app"
	"github.com/goregion/hexago/internal/entity"
	unit "github.com/goregion/hexago/tests/unit"
)

func TestOHLCPublisher_NewOHLCPublisher(t *testing.T) {
	mockPublisher1 := unit.NewMockOHLCPublisher()
	mockPublisher2 := unit.NewMockOHLCPublisher()

	publisher := app.NewOHLCPublisher(mockPublisher1, mockPublisher2)

	if publisher == nil {
		t.Fatal("OHLCPublisher should not be nil")
	}
}

func TestOHLCPublisher_ConsumeOHLC_SinglePublisher(t *testing.T) {
	mockPublisher := unit.NewMockOHLCPublisher()
	publisher := app.NewOHLCPublisher(mockPublisher)

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
	if publishedOHLC.CloseTimeMs != ohlc.CloseTimeMs {
		t.Errorf("Expected close time %d, got %d", ohlc.CloseTimeMs, publishedOHLC.CloseTimeMs)
	}
}

func TestOHLCPublisher_ConsumeOHLC_MultiplePublishers(t *testing.T) {
	mockPublisher1 := unit.NewMockOHLCPublisher()
	mockPublisher2 := unit.NewMockOHLCPublisher()
	mockPublisher3 := unit.NewMockOHLCPublisher()

	publisher := app.NewOHLCPublisher(mockPublisher1, mockPublisher2, mockPublisher3)

	ohlc := unit.CreateTestOHLC("ETHUSDT", 200.0, 220.0, 190.0, 210.0, 1234567891)

	ctx := context.Background()
	err := publisher.ConsumeOHLC(ctx, ohlc)

	if err != nil {
		t.Fatalf("ConsumeOHLC should not return error: %v", err)
	}

	// All publishers should be called
	if mockPublisher1.CallCount != 1 {
		t.Errorf("Expected 1 call to first publisher, got %d", mockPublisher1.CallCount)
	}
	if mockPublisher2.CallCount != 1 {
		t.Errorf("Expected 1 call to second publisher, got %d", mockPublisher2.CallCount)
	}
	if mockPublisher3.CallCount != 1 {
		t.Errorf("Expected 1 call to third publisher, got %d", mockPublisher3.CallCount)
	}

	// All publishers should receive the same data
	publishers := []*unit.MockOHLCPublisher{mockPublisher1, mockPublisher2, mockPublisher3}
	for i, mock := range publishers {
		if len(mock.PublishedOHLCs) != 1 {
			t.Fatalf("Publisher %d should have 1 published OHLC, got %d", i+1, len(mock.PublishedOHLCs))
		}

		publishedOHLC := mock.PublishedOHLCs[0]
		if publishedOHLC.Symbol != ohlc.Symbol {
			t.Errorf("Publisher %d: Expected symbol %s, got %s", i+1, ohlc.Symbol, publishedOHLC.Symbol)
		}
		if publishedOHLC.Open != ohlc.Open {
			t.Errorf("Publisher %d: Expected open %f, got %f", i+1, ohlc.Open, publishedOHLC.Open)
		}
		if publishedOHLC.Close != ohlc.Close {
			t.Errorf("Publisher %d: Expected close %f, got %f", i+1, ohlc.Close, publishedOHLC.Close)
		}
	}
}

func TestOHLCPublisher_ConsumeOHLC_NoPublishers(t *testing.T) {
	publisher := app.NewOHLCPublisher() // No publishers

	ohlc := unit.CreateTestOHLC("ADAUSDT", 1.0, 1.2, 0.9, 1.1, 1234567892)

	ctx := context.Background()
	err := publisher.ConsumeOHLC(ctx, ohlc)

	// Should not return error even with no publishers
	if err != nil {
		t.Fatalf("ConsumeOHLC should not return error with no publishers: %v", err)
	}
}

func TestOHLCPublisher_ConsumeOHLC_FirstPublisherError(t *testing.T) {
	mockPublisher1 := unit.NewMockOHLCPublisher()
	mockPublisher2 := unit.NewMockOHLCPublisher()

	// Make first publisher return error
	mockPublisher1.ShouldError = true
	mockPublisher1.ErrorMessage = "first publisher error"

	publisher := app.NewOHLCPublisher(mockPublisher1, mockPublisher2)

	ohlc := unit.CreateTestOHLC("BTCUSDT", 100.0, 105.0, 98.0, 102.0, 1234567893)

	ctx := context.Background()
	err := publisher.ConsumeOHLC(ctx, ohlc)

	if err == nil {
		t.Fatal("ConsumeOHLC should return error when first publisher fails")
	}

	if mockPublisher1.CallCount != 1 {
		t.Errorf("Expected 1 call to first publisher, got %d", mockPublisher1.CallCount)
	}

	// Second publisher should not be called due to error in first
	if mockPublisher2.CallCount != 0 {
		t.Errorf("Expected 0 calls to second publisher after first fails, got %d", mockPublisher2.CallCount)
	}
}

func TestOHLCPublisher_ConsumeOHLC_SecondPublisherError(t *testing.T) {
	mockPublisher1 := unit.NewMockOHLCPublisher()
	mockPublisher2 := unit.NewMockOHLCPublisher()
	mockPublisher3 := unit.NewMockOHLCPublisher()

	// Make second publisher return error
	mockPublisher2.ShouldError = true
	mockPublisher2.ErrorMessage = "second publisher error"

	publisher := app.NewOHLCPublisher(mockPublisher1, mockPublisher2, mockPublisher3)

	ohlc := unit.CreateTestOHLC("ETHUSDT", 200.0, 220.0, 190.0, 210.0, 1234567894)

	ctx := context.Background()
	err := publisher.ConsumeOHLC(ctx, ohlc)

	if err == nil {
		t.Fatal("ConsumeOHLC should return error when second publisher fails")
	}

	// First publisher should be called successfully
	if mockPublisher1.CallCount != 1 {
		t.Errorf("Expected 1 call to first publisher, got %d", mockPublisher1.CallCount)
	}
	if len(mockPublisher1.PublishedOHLCs) != 1 {
		t.Errorf("First publisher should have published OHLC")
	}

	// Second publisher should be called and fail
	if mockPublisher2.CallCount != 1 {
		t.Errorf("Expected 1 call to second publisher, got %d", mockPublisher2.CallCount)
	}
	if len(mockPublisher2.PublishedOHLCs) != 0 {
		t.Errorf("Second publisher should not have published OHLC due to error")
	}

	// Third publisher should not be called due to error in second
	if mockPublisher3.CallCount != 0 {
		t.Errorf("Expected 0 calls to third publisher after second fails, got %d", mockPublisher3.CallCount)
	}
}

func TestOHLCPublisher_ConsumeOHLC_MultipleOHLCs(t *testing.T) {
	mockPublisher := unit.NewMockOHLCPublisher()
	publisher := app.NewOHLCPublisher(mockPublisher)

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
		if publishedOHLC.CloseTimeMs != expectedOHLC.CloseTimeMs {
			t.Errorf("OHLC %d: Expected close time %d, got %d", i, expectedOHLC.CloseTimeMs, publishedOHLC.CloseTimeMs)
		}
	}
}

func TestOHLCPublisher_ConsumeOHLC_NilOHLC(t *testing.T) {
	mockPublisher := unit.NewMockOHLCPublisher()
	publisher := app.NewOHLCPublisher(mockPublisher)

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
