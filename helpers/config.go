package helpers

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Mysql struct {
		Host string `json:"host"`
		Port string `json:"port"`
		User string `json:"user"`
		Pass string `json:"pass"`
		Name string `json:"name"`
	} `json:"mysql"`
	JWT struct {
		Secret  string `json:"secret"`
		Expired int    `json:"expired"`
	} `json:"jwt"`
}

func GetConfig() Config {
	var conf Config

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("error: ", err)
	}

	if err := viper.Unmarshal(&conf); err != nil {
		panic(err)
	}
	return conf
}
