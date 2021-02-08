package config

import (
    "fmt"
    "github.com/spf13/viper"
    "log"
)

type Config struct {
    GinMode       string `mapstructure:"GIN_MODE"`
    ServerAddress string `mapstructure:"SERVER_ADDRESS"`

    MysqlUser string `mapstructure:"MYSQL_USER"`
    MysqlPass string `mapstructure:"MYSQL_PASS"`
    MysqlDb   string `mapstructure:"MYSQL_DB"`
    MysqlPort string `mapstructure:"MYSQL_PORT"`
    MysqlHost string `mapstructure:"MYSQL_HOST"`
    MysqlMaxConnection int `mapstructure:"MYSQL_MAX_CONNECTION"`

    JwtSecret string `mapstructure:"JWT_SECRET"`
    JwtIssuer string `mapstructure:"JWT_ISSUER"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config *Config, err error) {
    viper.SetConfigName("app")
    viper.SetConfigType("env")

    viper.AddConfigPath(path)
    viper.AddConfigPath(path + "/config")

    viper.AutomaticEnv()

    err = viper.ReadInConfig()
    if err != nil {
       log.Printf("Fatal error config file: %s \n", err)
    }

    err = viper.Unmarshal(&config)

    fmt.Printf("%#v", config)

    return
}

