package main

import (
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/gomqtt/")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal: cannot read config file: %s\n"), err)
	}
}