package config

// Configuration holds all configuration for modem-scraper.
type Configuration struct {
	Modem      Modem
	Polling    Polling
	MQTT       MQTT
	InfluxDB   InfluxDB
	BoltDB     BoltDB
	Prometheus Prometheus
}

// Modem holds modem configuration
type Modem struct {
	Url      string
	Username string
	Password string
}

// Polling holds polling configuration
type Polling struct {
	Schedule string
}

// MQTT holds MQTT connection configuration.
type MQTT struct {
	Enabled  bool
	Hostname string
	Port     string
	Username string
	Password string
	Topic    string
	ClientID string
}

// InfluxDB holds InfluxDB connection configuration.
type InfluxDB struct {
	Enabled       bool
	Url           string
	Database      string
	Username      string
	Password      string
	SkipVerifySsl bool
}

// BoltDB holds BoltDB configuration.
type BoltDB struct {
	Enabled bool
	Path    string
}

type Prometheus struct {
	Enabled bool
}
