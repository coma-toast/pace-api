package main

import (
	"github.com/spf13/viper"
)

// config is the configuration struct
type config struct {
	DocumentDBUrl  string
	FirebaseConfig string
}

// new config instance
var (
	conf *config
)

func getConf() *config {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()

	if err != nil {
		panic(err)
	}

	conf := &config{}
	err = viper.Unmarshal(conf)

	if err != nil {
		panic(err)
	}

	return conf
}
