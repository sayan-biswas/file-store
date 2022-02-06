package config

import (
	"errors"
	"path"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {

	// auto environment variables
	viper.AllowEmptyEnv(false)
	viper.SetEnvPrefix("STORE")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// config search paths
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("/etc/store/")
	viper.AddConfigPath("$HOME/store/")

}

func Load() error {

	// set config file path
	file := pflag.Lookup("config").Value.String()
	dir, file := path.Split(file)
	viper.AddConfigPath(dir)
	viper.SetConfigName(strings.TrimSuffix(file, path.Ext(file)))

	// read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return errors.New("config file not found, loading default config")
		} else {
			return errors.New("error reading config file, loading default config")
		}
	}
	return nil
}

func Get() (*Config, error) {

	// unmarshal config in to Config struct
	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}
	return config, nil
}
