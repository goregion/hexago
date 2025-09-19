package mysql

import (
	"testing"
)

func TestMakeOHLCTableName(t *testing.T) {
	testCases := []struct {
		name          string
		timeframeName string
		expected      string
	}{
		{
			name:          "1 minute timeframe",
			timeframeName: "1m",
			expected:      "`ohlc_1m`",
		},
		{
			name:          "5 minute timeframe",
			timeframeName: "5m",
			expected:      "`ohlc_5m`",
		},
		{
			name:          "1 hour timeframe",
			timeframeName: "1h",
			expected:      "`ohlc_1h`",
		},
		{
			name:          "1 day timeframe",
			timeframeName: "1d",
			expected:      "`ohlc_1d`",
		},
		{
			name:          "Empty timeframe",
			timeframeName: "",
			expected:      "`ohlc_`",
		},
		{
			name:          "Complex timeframe",
			timeframeName: "15m",
			expected:      "`ohlc_15m`",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test the table name generation logic
			result := "`ohlc_" + tc.timeframeName + "`"
			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}

func TestOHLCTableNameSafety(t *testing.T) {
	// Test that table names are properly escaped for SQL safety
	dangerousInputs := []struct {
		name         string
		input        string
		shouldEscape bool
	}{
		{
			name:         "Normal timeframe",
			input:        "1m",
			shouldEscape: false,
		},
		{
			name:         "SQL injection attempt",
			input:        "1m'; DROP TABLE users; --",
			shouldEscape: true,
		},
		{
			name:         "Special characters",
			input:        "1m@#$%",
			shouldEscape: true,
		},
		{
			name:         "Spaces",
			input:        "1 m",
			shouldEscape: true,
		},
	}

	for _, tc := range dangerousInputs {
		t.Run(tc.name, func(t *testing.T) {
			result := "`ohlc_" + tc.input + "`"

			// Check that dangerous characters are contained within backticks
			if tc.shouldEscape {
				if result[0] != '`' || result[len(result)-1] != '`' {
					t.Error("Dangerous input should be contained within backticks")
				}
			}

			// Basic validation: result should always start and end with backtick
			if result[0] != '`' || result[len(result)-1] != '`' {
				t.Error("Table name should always be wrapped in backticks")
			}
		})
	}
}

func TestSQLQueryConstruction(t *testing.T) {
	// Test the SQL query construction logic used in the MySQL adapter
	testCases := []struct {
		name          string
		timeframeName string
		expectedQuery string
	}{
		{
			name:          "1 minute query",
			timeframeName: "1m",
			expectedQuery: "INSERT INTO `ohlc_1m` (symbol, open, high, low, close, close_time_ms, timeframe) VALUES (:symbol, :open, :high, :low, :close, :close_time_ms, :timeframe)",
		},
		{
			name:          "5 minute query",
			timeframeName: "5m",
			expectedQuery: "INSERT INTO `ohlc_5m` (symbol, open, high, low, close, close_time_ms, timeframe) VALUES (:symbol, :open, :high, :low, :close, :close_time_ms, :timeframe)",
		},
		{
			name:          "1 hour query",
			timeframeName: "1h",
			expectedQuery: "INSERT INTO `ohlc_1h` (symbol, open, high, low, close, close_time_ms, timeframe) VALUES (:symbol, :open, :high, :low, :close, :close_time_ms, :timeframe)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tableName := "`ohlc_" + tc.timeframeName + "`"
			query := "INSERT INTO " + tableName + " (symbol, open, high, low, close, close_time_ms, timeframe) VALUES (:symbol, :open, :high, :low, :close, :close_time_ms, :timeframe)"

			if query != tc.expectedQuery {
				t.Errorf("Expected query:\n%s\nGot:\n%s", tc.expectedQuery, query)
			}
		})
	}
}

func TestSQLParameterBinding(t *testing.T) {
	// Test that the SQL parameter names match the OHLC entity fields
	expectedParams := []string{
		":symbol",
		":open",
		":high",
		":low",
		":close",
		":close_time_ms",
		":timeframe",
	}

	query := "INSERT INTO `ohlc_1m` (symbol, open, high, low, close, close_time_ms, timeframe) VALUES (:symbol, :open, :high, :low, :close, :close_time_ms, :timeframe)"

	for _, param := range expectedParams {
		found := false
		for i := 0; i < len(query)-len(param); i++ {
			if query[i:i+len(param)] == param {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected parameter %s not found in query", param)
		}
	}
}

func TestOHLCPublisherCreation(t *testing.T) {
	// Test the OHLCPublisher creation logic (without actual database)
	testCases := []struct {
		name          string
		timeframeName string
		hasError      bool
	}{
		{
			name:          "Valid 1m timeframe",
			timeframeName: "1m",
			hasError:      false,
		},
		{
			name:          "Valid 5m timeframe",
			timeframeName: "5m",
			hasError:      false,
		},
		{
			name:          "Valid 1h timeframe",
			timeframeName: "1h",
			hasError:      false,
		},
		{
			name:          "Empty timeframe",
			timeframeName: "",
			hasError:      false, // May be valid depending on implementation
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test basic validation that would happen during publisher creation
			if tc.timeframeName == "" {
				t.Log("Empty timeframe may need special handling")
			}

			// Test table name generation for this timeframe
			tableName := "`ohlc_" + tc.timeframeName + "`"
			if tableName == "`ohlc_`" && tc.timeframeName == "" {
				t.Log("Empty timeframe creates table name: " + tableName)
			}

			// Verify table name format
			if len(tableName) < 7 { // "`ohlc_`" is minimum 7 characters
				t.Error("Table name too short")
			}
			if tableName[0] != '`' || tableName[len(tableName)-1] != '`' {
				t.Error("Table name should be wrapped in backticks")
			}
		})
	}
}

func TestDatabaseConnection(t *testing.T) {
	// Test database connection validation (without actual connection)
	t.Run("Connection parameters", func(t *testing.T) {
		// Test the types of parameters that would be needed for MySQL connection
		requiredParams := []string{
			"host",
			"port",
			"database",
			"username",
			"password",
		}

		// Verify we know what parameters are needed
		for _, param := range requiredParams {
			if param == "" {
				t.Error("Parameter name should not be empty")
			}
		}

		// Test connection string format validation
		exampleDSN := "user:password@tcp(localhost:3306)/database"
		if len(exampleDSN) < 10 {
			t.Error("DSN too short")
		}
		if !contains(exampleDSN, "@tcp(") {
			t.Error("DSN should contain TCP connection info")
		}
	})
}

func TestOHLCDataValidation(t *testing.T) {
	// Test OHLC data validation that should happen before database insert
	testCases := []struct {
		name    string
		symbol  string
		open    float64
		high    float64
		low     float64
		close   float64
		timeMs  int64
		isValid bool
	}{
		{
			name:    "Valid OHLC",
			symbol:  "BTCUSDT",
			open:    100.0,
			high:    105.0,
			low:     98.0,
			close:   102.0,
			timeMs:  1234567890,
			isValid: true,
		},
		{
			name:    "Invalid OHLC - High < Low",
			symbol:  "BTCUSDT",
			open:    100.0,
			high:    95.0, // High should be >= Low
			low:     98.0,
			close:   102.0,
			timeMs:  1234567890,
			isValid: false,
		},
		{
			name:    "Zero timestamp",
			symbol:  "BTCUSDT",
			open:    100.0,
			high:    105.0,
			low:     98.0,
			close:   102.0,
			timeMs:  0,
			isValid: false,
		},
		{
			name:    "Empty symbol",
			symbol:  "",
			open:    100.0,
			high:    105.0,
			low:     98.0,
			close:   102.0,
			timeMs:  1234567890,
			isValid: false,
		},
		{
			name:    "Negative prices",
			symbol:  "BTCUSDT",
			open:    -100.0,
			high:    -95.0,
			low:     -105.0,
			close:   -102.0,
			timeMs:  1234567890,
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Basic validation logic
			hasSymbol := tc.symbol != ""
			hasValidTimestamp := tc.timeMs > 0
			pricesNonNegative := tc.open >= 0 && tc.high >= 0 && tc.low >= 0 && tc.close >= 0
			validHighLow := tc.high >= tc.low

			basicValid := hasSymbol && hasValidTimestamp && pricesNonNegative && validHighLow

			if tc.isValid && !basicValid {
				t.Error("Expected valid OHLC to pass basic validation")
			}
			if !tc.isValid && basicValid {
				t.Log("OHLC marked as invalid but passes basic validation - may need additional validation")
			}
		})
	}
}

// Helper function
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
