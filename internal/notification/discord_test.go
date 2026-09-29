package notification

import (
	"strings"
	"testing"

	"jetspotter/internal/jetspotter"
)

func TestColorByAltitude(t *testing.T) {

	expected := darkBlue
	actual := getColorByAltitude(25000)

	if expected != actual {
		t.Fatalf("expected '%v' to be the same as '%v'", expected, actual)
	}
}

func TestFormatDiscordContent(t *testing.T) {
	aircraft := []jetspotter.Aircraft{
		{Callsign: "POL30", Registration: "VH-PVO"},
		{Callsign: "POL35", Registration: "VH-PVE"},
	}

	tests := []struct {
		name     string
		template string
		expected string
	}{
		{"plain text", "Police nearby!", "Police nearby!"},
		{"empty disables the text", "", ""},
		{"registrations", "Spotted: {registrations}", "Spotted: VH-PVO, VH-PVE"},
		{"callsigns and count", "{count} aircraft ({callsigns})", "2 aircraft (POL30, POL35)"},
		{"role mention is kept", "<@&123456> {registrations}", "<@&123456> VH-PVO, VH-PVE"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := formatDiscordContent(tt.template, aircraft)
			if actual != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, actual)
			}
		})
	}
}

func TestFormatDiscordContentIsTruncated(t *testing.T) {
	actual := formatDiscordContent(strings.Repeat("a", 3000), nil)
	if len([]rune(actual)) != discordMaxContentLength {
		t.Fatalf("expected length %d, got %d", discordMaxContentLength, len([]rune(actual)))
	}
}
