package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DB           *Database `mapstructure:"db"`
	ModelPath    string    `mapstructure:"model_path"`
	ModelReplace string    `mapstructure:"model_replace"`
}

type Database struct {
	WriteDB *DB  `mapstructure:"write_db"`
	ShowSQL bool `mapstructure:"showSQL"`
}

type DB struct {
	Host     string `mapstructure:"host"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DB       string `mapstructure:"db"`
	Port     string `mapstructure:"port"`
}

// Init config init
func Init(config string) *Config {
	viper.SetConfigFile(config)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	cfg := &Config{}
	err = viper.Unmarshal(cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}
