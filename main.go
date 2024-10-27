package main

import (
	"fmt"
	"os"

	"github.com/Ayobami6/todo_cli/cmd"
	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/viper"
)

func main() {

	viper.SetDefault("passcode", "1234")

	viper.SetConfigName("config") // Config file name (without extension)
	viper.SetConfigType("yaml")   // Config file type (could be json, yaml, etc.)
	viper.AddConfigPath(".")
	// check if config file exists, if not create
	_, err := os.Stat("config.yaml")
	// if file does not exist, create it
	if os.IsNotExist(err) {
		_, err := os.Create("config.yaml")
		if err != nil {
			panic(err)
		}
	}
	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	if err := cmd.RootCmd.Execute(); err != nil {
		panic(err)
	}

}
