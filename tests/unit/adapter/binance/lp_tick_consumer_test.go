package binance

import (
	"strconv"
	"testing"

	adapter_binance "github.com/goregion/hexago/internal/adapter/binance"
	unit "github.com/goregion/hexago/tests/unit"
)

func TestLPTickConsumer_NewLPTickConsumer(t *testing.T) {
	symbols := []string{"BTCUSDT", "ETHUSDT"}
	mockConsumer := unit.NewMockLPTickConsumer()
	errorHandler := func(err error) {}

	consumer := adapter_binance.NewLPTickConsumer(symbols, mockConsumer, errorHandler)

	if consumer == nil {
		t.Fatal("LPTickConsumer should not be nil")
	}
}

func TestLPTickConsumer_Configuration(t *testing.T) {
	testCases := []struct {
		name     string
		symbols  []string
		hasError bool
	}{
		{
			name:     "Single symbol",
			symbols:  []string{"BTCUSDT"},
			hasError: false,
		},
		{
			name:     "Multiple symbols",
			symbols:  []string{"BTCUSDT", "ETHUSDT", "ADAUSDT"},
			hasError: false,
		},
		{
			name:     "Empty symbols",
			symbols:  []string{},
			hasError: false,
		},
		{
			name:     "Nil symbols",
			symbols:  nil,
			hasError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockConsumer := unit.NewMockLPTickConsumer()
			errorHandler := func(err error) {}

			consumer := adapter_binance.NewLPTickConsumer(tc.symbols, mockConsumer, errorHandler)

			if consumer == nil {
				t.Fatal("LPTickConsumer should not be nil")
			}
		})
	}
}

// Test the price conversion logic that the adapter uses
func TestPriceConversion(t *testing.T) {
	testCases := []struct {
		name        string
		bidPrice    string
		askPrice    string
		expectedBid float64
		expectedAsk float64
		shouldError bool
	}{
		{
			name:        "Valid prices",
			bidPrice:    "50000.12",
			askPrice:    "50001.34",
			expectedBid: 50000.12,
			expectedAsk: 50001.34,
			shouldError: false,
		},
		{
			name:        "Zero prices",
			bidPrice:    "0.0",
			askPrice:    "0.0",
			expectedBid: 0.0,
			expectedAsk: 0.0,
			shouldError: false,
		},
		{
			name:        "Large numbers",
			bidPrice:    "999999.999999",
			askPrice:    "1000000.000001",
			expectedBid: 999999.999999,
			expectedAsk: 1000000.000001,
			shouldError: false,
		},
		{
			name:        "Small numbers",
			bidPrice:    "0.000001",
			askPrice:    "0.000002",
			expectedBid: 0.000001,
			expectedAsk: 0.000002,
			shouldError: false,
		},
		{
			name:        "Invalid bid price",
			bidPrice:    "invalid",
			askPrice:    "50001.34",
			shouldError: true,
		},
		{
			name:        "Invalid ask price",
			bidPrice:    "50000.12",
			askPrice:    "invalid",
			shouldError: true,
		},
		{
			name:        "Empty bid price",
			bidPrice:    "",
			askPrice:    "50001.34",
			shouldError: true,
		},
		{
			name:        "Empty ask price",
			bidPrice:    "50000.12",
			askPrice:    "",
			shouldError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test the price conversion logic that the adapter uses internally
			bidPrice, bidErr := strconv.ParseFloat(tc.bidPrice, 64)
			askPrice, askErr := strconv.ParseFloat(tc.askPrice, 64)

			if tc.shouldError {
				if bidErr == nil && askErr == nil {
					t.Error("Expected error but parsing succeeded")
				}
			} else {
				if bidErr != nil {
					t.Errorf("Unexpected error parsing bid price: %v", bidErr)
				}
				if askErr != nil {
					t.Errorf("Unexpected error parsing ask price: %v", askErr)
				}

				if bidErr == nil && bidPrice != tc.expectedBid {
					t.Errorf("Expected bid price %f, got %f", tc.expectedBid, bidPrice)
				}
				if askErr == nil && askPrice != tc.expectedAsk {
					t.Errorf("Expected ask price %f, got %f", tc.expectedAsk, askPrice)
				}
			}
		})
	}
}

func TestLPTickConsumer_ErrorHandlerFunction(t *testing.T) {
	errorCalled := false
	var receivedError error

	errorHandler := func(err error) {
		errorCalled = true
		receivedError = err
	}

	mockConsumer := unit.NewMockLPTickConsumer()
	consumer := adapter_binance.NewLPTickConsumer([]string{"BTCUSDT"}, mockConsumer, errorHandler)

	if consumer == nil {
		t.Fatal("Consumer should not be nil")
	}

	// Verify error handler is not called during initialization
	if errorCalled {
		t.Error("Error handler should not be called during initialization")
	}

	if receivedError != nil {
		t.Error("No error should be received during initialization")
	}
}

func TestLPTickConsumer_NilConsumer(t *testing.T) {
	errorHandler := func(err error) {}

	// Test with nil consumer - should not panic
	consumer := adapter_binance.NewLPTickConsumer([]string{"BTCUSDT"}, nil, errorHandler)

	if consumer == nil {
		t.Fatal("LPTickConsumer should not be nil even with nil consumer")
	}
}

func TestLPTickConsumer_NilErrorHandler(t *testing.T) {
	mockConsumer := unit.NewMockLPTickConsumer()

	// Test with nil error handler - should not panic
	consumer := adapter_binance.NewLPTickConsumer([]string{"BTCUSDT"}, mockConsumer, nil)

	if consumer == nil {
		t.Fatal("LPTickConsumer should not be nil even with nil error handler")
	}
}
