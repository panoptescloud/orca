// To override the name of a config field in the yaml file,  you need to use the
// mapstructure tag instead of the yaml tag as you may expect. Viper starts by
// unmarshalling the yaml to map[string]interface. access_log is one example
package config

type ConfigLogging struct {
	Level  string
	Format string
}

type Config struct {
	Logging ConfigLogging
}

func NewDefault() *Config {
	return &Config{
		Logging: ConfigLogging{
			Level:  "info",
			Format: "json",
		},
	}
}
