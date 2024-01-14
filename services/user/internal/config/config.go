package config

import "github.com/spf13/viper"

type ConfigInfo struct {
	PGHost     string
	PGPort     uint16
	PGUser     string
	PGDBName   string
	PGPassword string
}

func NewConfig(name, fileType, path string) (*ConfigInfo, error) {
	viper.SetConfigName(name)
	viper.SetConfigType(fileType)
	viper.AddConfigPath(path)
	viper.ReadInConfig()

	var config ConfigInfo

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
