package paceconfig

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// Config is the config struct
type Config struct {
	DocumentDBUrl  string
	FirebaseConfig string
	RollbarToken   string
}

// GetConf gets a config file from local disk
func GetConf(path string) (*Config, error) {
	conf := &Config{}
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		return conf, err
	}

	err = viper.Unmarshal(conf)
	if err != nil {
		return conf, err
	}

	log.Println(conf.FirebaseConfig)
	conf.FirebaseConfig = fmt.Sprintf("%s%s", path, conf.FirebaseConfig)
	log.Println(conf.FirebaseConfig)

	return conf, err
}
