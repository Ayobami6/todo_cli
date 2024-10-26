package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Ayobami6/todo_cli/db"
	"github.com/Ayobami6/todo_cli/utils"
	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/viper"
)

func main() {
	mongoUrl := utils.GetEnv("MONGO_URL", "mongodb://localhost:27017")
	ctx, cancel := context.WithDeadline(context.TODO(), time.Now().Add(30*time.Second))
	defer cancel()
	viper.SetDefault("passcode", "1234")

	viper.SetConfigName("config") // Config file name (without extension)
	viper.SetConfigType("json")   // Config file type (could be json, yaml, etc.)
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	client, err := db.ConnectDb(ctx, mongoUrl)
	// client, err := db.ConnectDb(ctx, "mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}
	// fmt.Println(client)
	defer client.Disconnect(ctx)

}
