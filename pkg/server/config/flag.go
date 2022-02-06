package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	pflag.Bool("debug", false, "Debug Mode")
	pflag.String("log", "info", "Logger Mode - (debug, info, warn, error, fatal, panic)")
	pflag.String("config", "store.yaml", "Configuration File")
	pflag.Int("port", 8080, "Server Port")
	pflag.String("host", "", "Server Hostname")
	pflag.Bool("tls", false, "Enable TLS")
	pflag.Parse()

	viper.BindPFlag("server.log", pflag.Lookup("log"))
	viper.BindPFlag("server.debug", pflag.Lookup("debug"))
	viper.BindPFlag("server.tls", pflag.Lookup("tls"))
	viper.BindPFlag("server.host", pflag.Lookup("host"))
	viper.BindPFlag("server.port", pflag.Lookup("port"))

	viper.BindPFlag("database.location", pflag.Lookup("database"))
}
