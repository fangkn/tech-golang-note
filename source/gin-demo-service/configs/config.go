package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

type serverConfig struct {
	RunMode      string
	Port         int
	ReadTimeout  int
	WriteTimeout int
}

type appConfig struct {
	AreaType int
}

var (
	Server serverConfig
	App    appConfig
)

func Load() error {
	viper.SetConfigName("app")
	viper.AddConfigPath("./conf")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("viper ReadInConfig err: ", err)
		return err
	}
	app := struct {
		Server *serverConfig
		App    *appConfig
	}{
		Server: &Server,
		App:    &App,
	}
	if err := viper.Unmarshal(&app); err != nil {
		fmt.Println("viper decode err: ", err)
	}
	return nil
}
