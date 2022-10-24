package aggregatorv2

// ServerConfig represents the configuration of the aggregator server.
type ServerConfig struct {
	Host string `mapstructure:"Host"`
	Port int    `mapstructure:"Port"`
}
