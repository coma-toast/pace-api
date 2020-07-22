package paceconfig

import "github.com/spf13/viper"

// Config is the config struct
type Config struct {
	DocumentDBUrl  string
	FirebaseConfig string
	RollbarToken   string
}

// GetConf gets a config file from local disk
func GetConf() (*Config, error) {
	conf := &Config{}
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()

	if err != nil {
		return conf, err
	}

	err = viper.Unmarshal(conf)

	if err != nil {
		return conf, err
	}

	return conf, err
}
