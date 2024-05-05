package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
)

// const RABBITMQ_URL = "amqp://guest:guest@localhost:5672"

type Config struct {
	RABBITMQ_URL string `mapstructure:"RABBITMQ_URL"`
}

var C *Config

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
func loadConfig() {
	v := viper.New()
	v.AddConfigPath(".")
	v.SetConfigName(".env")
	v.SetConfigType("env")
	err := v.ReadInConfig()
	if err != nil {
		failOnError(err, "Failed to read config")
	}
	v.AutomaticEnv()

	err = v.Unmarshal(&C)
	if err != nil {
		failOnError(err, "Failed to read enivroment")
	}
}
func main() {
	fmt.Println("Starting Peril server...")
	loadConfig()
	conn, err := amqp.Dial(C.RABBITMQ_URL)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	log.Printf("connect to %s successfully", C.RABBITMQ_URL)
	// wait for ctrl + c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	log.Printf("receive signal %v, closing connection and shuting down", os.Interrupt)
}
