package chassis

type Config interface {
	GetPostgresConfig() string
	GetListenAddress() string
}

// PostgresConfig is a struct that holds the configuration for the Postgres database
type PostgresConfig struct {
	ConnectionString      string `mapstructure:"postgres_connection_string"`        // Connection string
	MaxOpenConnections    int    `mapstructure:"postgres_max_open_connections"`     // Maximum number of open connections
	MaxIdleConnections    int    `mapstructure:"postgres_max_idle_connections"`     // Maximum number of idle connections
	ConnectionMaxLifetime int    `mapstructure:"postgres_connection_max_lifetime"`  // Maximum connection lifetime (seconds)
	ConnectionMaxIdleTime int    `mapstructure:"postgres_connection_max_idle_time"` // Maximum connection idle time (seconds)
}

func ParsePostgresConfig(config string) PostgresConfig {
	return PostgresConfig{
		ConnectionString:      config,
		MaxOpenConnections:    10,
		MaxIdleConnections:    5,
		ConnectionMaxLifetime: 5 * 60 * 60,
		ConnectionMaxIdleTime: 5 * 60 * 60,
	}
}

type MicroserviceConfig struct {
	PostgresURI   string
	ListenAddress string
}

func (c *MicroserviceConfig) GetPostgresConfig() string {
	return c.PostgresURI
}

func (c *MicroserviceConfig) GetListenAddress() string {
	return c.ListenAddress
}

func DefaultMicroserviceConfig() MicroserviceConfig {
	return MicroserviceConfig{
		PostgresURI:   "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable",
		ListenAddress: "0.0.0.0:8080",
	}
}
