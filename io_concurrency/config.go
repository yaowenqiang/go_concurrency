package main

import (
	"fmt"
	"github.com/spf13/viper"
)

func main() {
	viper.SetDefault("userName", "jack")
	viper.SetConfigName("config")
	viper.AddConfigPath("conf")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()

	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	userName := viper.GetString("username")
	fmt.Printf("%s\n", userName)

}
