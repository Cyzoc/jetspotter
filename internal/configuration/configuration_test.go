package configuration

import (
	"os"
	"testing"
)

// TestScanRangeDefaultsToMaxRangeKilometers tests that MaxScanRangeKilometers defaults to MaxRangeKilometers
// when MAX_SCAN_RANGE_KILOMETERS is not set
func TestScanRangeDefaultsToMaxRangeKilometers(t *testing.T) {
	// t.Setenv restores the original environment when the test finishes
	t.Setenv("MAX_RANGE_KILOMETERS", "")
	t.Setenv("MAX_SCAN_RANGE_KILOMETERS", "")

	// Set MAX_RANGE_KILOMETERS, but not MAX_SCAN_RANGE_KILOMETERS
	t.Setenv("MAX_RANGE_KILOMETERS", "50")
	t.Setenv("MAX_SCAN_RANGE_KILOMETERS", "")

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
	// t.Setenv restores the original environment when the test finishes
	t.Setenv("MAX_RANGE_KILOMETERS", "")
	t.Setenv("MAX_SCAN_RANGE_KILOMETERS", "")

	// Set both environment variables to different values
	t.Setenv("MAX_RANGE_KILOMETERS", "30")
	t.Setenv("MAX_SCAN_RANGE_KILOMETERS", "100")

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

// setValidBaseEnv makes sure the required numeric settings are valid, and restores them after the test.
func setValidBaseEnv(t *testing.T) {
	t.Helper()
	t.Setenv("MAX_RANGE_KILOMETERS", "30")
	t.Setenv("MAX_SCAN_RANGE_KILOMETERS", "")
}

func TestRegistrationsDefaultsToAll(t *testing.T) {
	setValidBaseEnv(t)
	t.Setenv("REGISTRATIONS", "")
	os.Unsetenv("REGISTRATIONS") // t.Setenv restores the original value afterwards

	config, err := GetConfig()
	if err != nil {
		t.Fatalf("Failed to get config: %v", err)
	}

	if len(config.Registrations) != 1 || config.Registrations[0] != "ALL" {
		t.Fatalf("expected Registrations to be [ALL], got %v", config.Registrations)
	}
}

func TestRegistrationsAreParsed(t *testing.T) {
	setValidBaseEnv(t)
	t.Setenv("REGISTRATIONS", " oo-abc, n12345 ,,fa-* ")

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

func TestDiscordMessageDefaultAndCustom(t *testing.T) {
	setValidBaseEnv(t)
	t.Setenv("DISCORD_MESSAGE", "")
	os.Unsetenv("DISCORD_MESSAGE") // t.Setenv restores the original value afterwards

	config, err := GetConfig()
	if err != nil {
		t.Fatalf("Failed to get config: %v", err)
	}
	if config.DiscordMessage != ":airplane: A jet has been spotted! :airplane:" {
		t.Fatalf("unexpected default: %q", config.DiscordMessage)
	}

	t.Setenv("DISCORD_MESSAGE", "Police nearby: {registrations}")
	config, err = GetConfig()
	if err != nil {
		t.Fatalf("Failed to get config: %v", err)
	}
	if config.DiscordMessage != "Police nearby: {registrations}" {
		t.Fatalf("unexpected message: %q", config.DiscordMessage)
	}
}
