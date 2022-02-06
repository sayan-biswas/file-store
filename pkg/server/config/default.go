package config

import "github.com/spf13/viper"

func init() {
	viper.SetDefault("server.host", "")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.log", "info")
	viper.SetDefault("server.debug", false)
	viper.SetDefault("server.tls", false)
	viper.SetDefault("server.certificate", "certificate.crt")
	viper.SetDefault("server.privateKey", "private.key")

	viper.SetDefault("cors.allowOrigins", []string{"*"})
	viper.SetDefault("cors.allowMethods", []string{"*"})
	viper.SetDefault("cors.allowHeaders", []string{"*"})
	viper.SetDefault("cors.exposeHeaders", []string{"*"})
	viper.SetDefault("cors.allowCredentials", true)
	viper.SetDefault("cors.maxAge", 86400)

	viper.SetDefault("database.diskless", false)
	viper.SetDefault("database.encryption", false)
	viper.SetDefault("database.cacheSize", 100)
	viper.SetDefault("database.path", "database")
}
