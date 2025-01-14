package config

import "github.com/spf13/viper"

type Conf struct {
	LimitPerIP           bool   `mapstructure:"LIMIT_PER_IP"`      // true or false
	LimitPerAPI_KEY      bool   `mapstructure:"LIMIT_PER_API_KEY"` // true or false
	MaxRequestsPerSecond int    `mapstructure:"MAX_REQUESTS_PER_SECOND"`
	BlockTime            int    `mapstructure:"BLOCK_TIME"` // in seconds
	HPPTCodeLimitReached int    `mapstructure:"HTTP_CODE_LIMIT_REACHED"`
	MessageLimitReached  string `mapstructure:"MESSAGE_LIMIT_REACHED"`
	RedisAddr            string `mapstructure:"REDIS_ADDR"`
}

func LoadConfig(path string) (*Conf, error) {
	var cfg *Conf

	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}

	return cfg, nil
}
