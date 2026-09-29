package configuration

import (
	"os"
	"testing"
)

// TestScanRangeDefaultsToMaxRangeKilometers tests that MaxScanRangeKilometers defaults to MaxRangeKilometers
// when MAX_SCAN_RANGE_KILOMETERS is not set
func TestScanRangeDefaultsToMaxRangeKilometers(t *testing.T) {
	// Save current environment
	oldMaxRange := os.Getenv("MAX_RANGE_KILOMETERS")
	oldMaxScanRange := os.Getenv("MAX_SCAN_RANGE_KILOMETERS")
	defer func() {
		os.Setenv("MAX_RANGE_KILOMETERS", oldMaxRange)
		os.Setenv("MAX_SCAN_RANGE_KILOMETERS", oldMaxScanRange)
	}()

	// Set MAX_RANGE_KILOMETERS, but not MAX_SCAN_RANGE_KILOMETERS
	os.Setenv("MAX_RANGE_KILOMETERS", "50")
	os.Unsetenv("MAX_SCAN_RANGE_KILOMETERS")

	config, err := GetConfig()
	if err != nil {
		t.Fatalf("Failed to get config: %v", err)
	}

	if config.MaxRangeKilometers != 50 {
		t.Fatalf("expected MaxRangeKilometers to be 50, got %d", config.MaxRangeKilometers)
	}

	if config.MaxScanRangeKilometers != config.MaxRangeKilometers {
		t.Fatalf("expected MaxScanRangeKilometers to equal MaxRangeKilometers (%d), got %d",
			config.MaxRangeKilometers, config.MaxScanRangeKilometers)
	}
}

// TestScanRangeCanBeDifferentFromMaxRange tests that MaxScanRangeKilometers can be set to a different
// value than MaxRangeKilometers
func TestScanRangeCanBeDifferentFromMaxRange(t *testing.T) {
	// Save current environment
	oldMaxRange := os.Getenv("MAX_RANGE_KILOMETERS")
	oldMaxScanRange := os.Getenv("MAX_SCAN_RANGE_KILOMETERS")
	defer func() {
		os.Setenv("MAX_RANGE_KILOMETERS", oldMaxRange)
		os.Setenv("MAX_SCAN_RANGE_KILOMETERS", oldMaxScanRange)
	}()

	// Set both environment variables to different values
	os.Setenv("MAX_RANGE_KILOMETERS", "30")
	os.Setenv("MAX_SCAN_RANGE_KILOMETERS", "100")

	config, err := GetConfig()
	if err != nil {
		t.Fatalf("Failed to get config: %v", err)
	}

	if config.MaxRangeKilometers != 30 {
		t.Fatalf("expected MaxRangeKilometers to be 30, got %d", config.MaxRangeKilometers)
	}

	if config.MaxScanRangeKilometers != 100 {
		t.Fatalf("expected MaxScanRangeKilometers to be 100, got %d", config.MaxScanRangeKilometers)
	}
}

func TestRegistrationsDefaultsToAll(t *testing.T) {
	old, wasSet := os.LookupEnv("REGISTRATIONS")
	defer func() {
		if wasSet {
			os.Setenv("REGISTRATIONS", old)
		} else {
			os.Unsetenv("REGISTRATIONS")
		}
	}()
	os.Unsetenv("REGISTRATIONS")

	config, err := GetConfig()
	if err != nil {
		t.Fatalf("Failed to get config: %v", err)
	}

	if len(config.Registrations) != 1 || config.Registrations[0] != "ALL" {
		t.Fatalf("expected Registrations to be [ALL], got %v", config.Registrations)
	}
}

func TestRegistrationsAreParsed(t *testing.T) {
	old, wasSet := os.LookupEnv("REGISTRATIONS")
	defer func() {
		if wasSet {
			os.Setenv("REGISTRATIONS", old)
		} else {
			os.Unsetenv("REGISTRATIONS")
		}
	}()
	os.Setenv("REGISTRATIONS", " oo-abc, n12345 ,,fa-* ")

	config, err := GetConfig()
	if err != nil {
		t.Fatalf("Failed to get config: %v", err)
	}

	expected := []string{"OO-ABC", "N12345", "FA-*"}
	if len(config.Registrations) != len(expected) {
		t.Fatalf("expected %v, got %v", expected, config.Registrations)
	}
	for i := range expected {
		if config.Registrations[i] != expected[i] {
			t.Fatalf("expected %v, got %v", expected, config.Registrations)
		}
	}
}
