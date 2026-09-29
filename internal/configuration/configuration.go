package configuration

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/jftuga/geodist"
)

// Config is the type used for all user configuration.
// All parameters can be set using ENV variables.
// The comments below are structured as following:
// ENV_VARIABLE_NAME DEFAULT_VALUE
type Config struct {
	// Latitude and Longitude coordinates of the location you want to use.
	// LOCATION_LATITUDE 51.17348
	// LOCATION_LONGITUDE 5.45921
	Location geodist.Coord

	// Maximum range in kilometers from the location that you want aircraft to be spotted.
	// Note that this is an approximation due to roundings.
	// MAX_RANGE_KILOMETERS 30
	MaxRangeKilometers int

	// Maximum range in kilometers from the location that you want to scan for aircraft.
	// This is used to determine the area to query from the ADS-B API.
	// If not set, MAX_RANGE_KILOMETERS will be used.
	// MAX_SCAN_RANGE_KILOMETERS 30
	MaxScanRangeKilometers int

	// Maximum altitude in feet that you want to spot aircraft at.
	// Set to 0 to disable the filter.
	// MAX_ALTITUDE_FEET 0
	MaxAltitudeFeet int

	// A comma seperated list of types that you want to spot
	// If not set, 'ALL' will be used, which will disable the filter and show all aircraft within range.
	// Full list can be found at https://www.icao.int/publications/doc8643/pages/search.aspx in 'Type Designator' column.
	// AIRCRAFT_TYPES ALL
	// EXAMPLES
	// AIRCRAFT_TYPES F16,F35
	// To spot all military aircraft, you can use MILITARY.
	// AIRCRAFT_TYPES MILITARY
	AircraftTypes []string

	// A comma seperated list of registrations (tail numbers) that you want to spot.
	// If not set, 'ALL' will be used, which will disable the filter.
	// Matching is case-insensitive and ignores hyphens, so 'OO-ABC' and 'ooabc' are equivalent.
	// A trailing '*' can be used as a wildcard to match a prefix, for example 'FA-*' matches FA-101, FA-102, ...
	// This filter is combined with AIRCRAFT_TYPES, an aircraft must match both filters to be spotted.
	// REGISTRATIONS ALL
	// EXAMPLES
	// REGISTRATIONS OO-ABC,N12345
	// REGISTRATIONS FA-*
	Registrations []string

	// Webhook used to send notifications to Slack. If not set, no messages will be sent to Slack.
	// SLACK_WEBHOOK_URL ""
	SlackWebHookURL string

	// Webhook used to send notifications to Discord. If not set, no messages will be sent to Discord.
	// DISCORD_WEBHOOK_URL ""
	DiscordWebHookURL string

	// Discord notifications use an embed color based on the alitute of the aircraft.
	// DISCORD_COLOR_ALTITUDE "true"
	DiscordColorAltitude string

	// Custom text that is displayed above the aircraft in Discord notifications.
	// Discord markdown, emojis and mentions (for example <@&ROLE_ID>) are supported.
	// The placeholders {count}, {registrations} and {callsigns} are replaced by the values of the aircraft in the message.
	// If set to an empty string, only the aircraft embeds are sent. Discord limits this text to 2000 characters.
	// DISCORD_MESSAGE ":airplane: A jet has been spotted! :airplane:"
	// EXAMPLE
	// DISCORD_MESSAGE ":police_car: Police aircraft nearby: {registrations}"
	DiscordMessage string

	// Interval in seconds between fetching aircraft, minimum is 60 due to API rate limiting.
	// FETCH_INTERVAL 60
	FetchInterval int

	// Token to authenticate with the gotify server.
	// GOTIFY_TOKEN ""
	GotifyToken string

	// URL of the gotify server.
	// GOTIFY_URL ""
	GotifyURL string

	// Port where metrics will be exposed on
	// METRICS_PORT "7070"
	MetricsPort string

	// Port where API will be exposed on
	// API_PORT "8085"
	APIPort string

	// Enable or disable the web UI
	// WEB_UI_ENABLED "true"
	WebUIEnabled bool

	// Port where web UI will be exposed on
	// WEB_UI_PORT "8080"
	WebUIPort string

	// Topic to publish message to
	// NTFY_TOPIC ""
	NtfyTopic string

	// URL of the ntfy server.
	// NTFY_SERVER "https://ntfy.sh"
	NtfyServer string

	// Token for ntfy server authentication.
	// NTFY_TOKEN ""
	NtfyToken string
}

// Environment variable names
const (
	SlackWebhookURL        = "SLACK_WEBHOOK_URL"
	DiscordWebhookURL      = "DISCORD_WEBHOOK_URL"
	DiscordColorAltitude   = "DISCORD_COLOR_ALTITUDE"
	DiscordMessage         = "DISCORD_MESSAGE"
	LocationLatitude       = "LOCATION_LATITUDE"
	LocationLongitude      = "LOCATION_LONGITUDE"
	MaxRangeKilometers     = "MAX_RANGE_KILOMETERS"
	MaxScanRangeKilometers = "MAX_SCAN_RANGE_KILOMETERS"
	MaxAltitudeFeet        = "MAX_ALTITUDE_FEET"
	AircraftTypes          = "AIRCRAFT_TYPES"
	Registrations          = "REGISTRATIONS"
	FetchInterval          = "FETCH_INTERVAL"
	GotifyURL              = "GOTIFY_URL"
	NtfyTopic              = "NTFY_TOPIC"
	NtfyServer             = "NTFY_SERVER"
	NtfyToken              = "NTFY_TOKEN"
	GotifyToken            = "GOTIFY_TOKEN"
	MetricsPort            = "METRICS_PORT"
	APIPort                = "API_PORT"
	WebUIEnabled           = "WEB_UI_ENABLED"
	WebUIPort              = "WEB_UI_PORT"
)

// getEnvVariable looks up a specified environment variable, if not set the specified default is used
func getEnvVariable(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

// GetConfig attempts to read the configuration via environment variables and uses a default if the environment variable is not set
func GetConfig() (config Config, err error) {
	defaultFetchInterval := 60

	config.GotifyToken = getEnvVariable(GotifyToken, "")
	config.GotifyURL = getEnvVariable(GotifyURL, "")
	config.NtfyTopic = getEnvVariable(NtfyTopic, "")
	config.NtfyServer = getEnvVariable(NtfyServer, "https://ntfy.sh")
	config.NtfyToken = getEnvVariable(NtfyToken, "")
	config.SlackWebHookURL = getEnvVariable(SlackWebhookURL, "")
	config.DiscordWebHookURL = getEnvVariable(DiscordWebhookURL, "")
	config.DiscordColorAltitude = getEnvVariable(DiscordColorAltitude, "true")
	config.DiscordMessage = getEnvVariable(DiscordMessage, ":airplane: A jet has been spotted! :airplane:")
	config.MetricsPort = getEnvVariable(MetricsPort, "7070")
	config.APIPort = getEnvVariable(APIPort, "8085")
	config.WebUIPort = getEnvVariable(WebUIPort, "8080")

	// Convert WebUIEnabled string to bool
	webUIEnabledStr := getEnvVariable(WebUIEnabled, "true")
	config.WebUIEnabled, err = strconv.ParseBool(webUIEnabledStr)
	if err != nil {
		log.Printf("Invalid value for WEB_UI_ENABLED: %s, using default: true", webUIEnabledStr)
		config.WebUIEnabled = true
	}

	config.FetchInterval, err = strconv.Atoi(getEnvVariable(FetchInterval, strconv.Itoa(defaultFetchInterval)))
	if err != nil || config.FetchInterval < 5 {
		log.Printf("Fetch interval of %ds detected. You might hit rate limits, consider using the default of %ds instead.", config.FetchInterval, defaultFetchInterval)
	}

	config.Location.Lat, err = strconv.ParseFloat(getEnvVariable(LocationLatitude, "51.17348"), 64)
	if err != nil {
		return Config{}, err
	}

	config.Location.Lon, err = strconv.ParseFloat(getEnvVariable(LocationLongitude, "5.45921"), 64)
	if err != nil {
		return Config{}, err
	}

	config.MaxRangeKilometers, err = strconv.Atoi(getEnvVariable(MaxRangeKilometers, "30"))
	if err != nil {
		return Config{}, err
	}

	// Get MAX_SCAN_RANGE_KILOMETERS - If not set, use MAX_RANGE_KILOMETERS for backward compatibility
	maxScanRangeStr := getEnvVariable(MaxScanRangeKilometers, "")
	if maxScanRangeStr == "" {
		config.MaxScanRangeKilometers = config.MaxRangeKilometers
	} else {
		config.MaxScanRangeKilometers, err = strconv.Atoi(maxScanRangeStr)
		if err != nil {
			return Config{}, err
		}
	}

	config.MaxAltitudeFeet, err = strconv.Atoi(getEnvVariable(MaxAltitudeFeet, "0"))
	if err != nil {
		return Config{}, err
	}

	config.AircraftTypes = strings.Split(strings.ToUpper(strings.ReplaceAll(getEnvVariable(AircraftTypes, "ALL"), " ", "")), ",")
	config.Registrations = parseList(getEnvVariable(Registrations, "ALL"), "ALL")
	return config, nil
}

// parseList converts a comma separated string into an upper-cased list without whitespace or empty entries.
// If the result is empty, a list containing only the fallback is returned.
func parseList(value, fallback string) []string {
	var list []string
	for _, item := range strings.Split(strings.ToUpper(value), ",") {
		item = strings.ReplaceAll(strings.TrimSpace(item), " ", "")
		if item != "" {
			list = append(list, item)
		}
	}

	if len(list) == 0 {
		return []string{fallback}
	}

	return list
}
