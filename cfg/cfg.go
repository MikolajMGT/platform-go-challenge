package cfg

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

type Configuration map[string]interface{}

func ServiceDefaults() Configuration {
	return Configuration{

		// Service
		"service.name": "assets",

		// WebServer
		"webserver.port": "8080",

		// Database
		"cassandra.cluster.ip":       "cassandra",
		"cassandra.cluster.keyspace": "assets_service",

		// Auth
		"auth.secret": "Ao8Qg52wYPhIzND",
	}
}

func Configure(configuration Configuration) (err error) {

	viper.SetConfigType("yaml")
	viper.SetConfigName("config")

	viper.AddConfigPath("/")
	viper.AddConfigPath(fmt.Sprintf("/etc/%s/", configuration["service.name"].(string)))
	viper.AddConfigPath(".")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	for k, v := range configuration {
		viper.SetDefault(k, v)
	}

	viper.SetEnvPrefix(configuration["service.name"].(string))
	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}
