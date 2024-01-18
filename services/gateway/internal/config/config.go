package config

import "github.com/spf13/viper"

type ConfigInfo struct {
	JWTSignKey string
}

func NewConfig(name, fileType, path string) (*ConfigInfo, error) {
	viper.SetConfigName(name)
	viper.SetConfigType(fileType)
	viper.AddConfigPath(path)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config ConfigInfo

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
