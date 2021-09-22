package main

import (
	"flag"
	"fmt"
	"log"

	internalHTTP "github.com/iKOPKACtraxa/wb-bank/internal/server/http"
	memorystorage "github.com/iKOPKACtraxa/wb-bank/internal/storage/memory"
)

var configFilePath string

func init() {
	flag.StringVar(&configFilePath, "config", "../../configs/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()
	config := NewConfig()
	storage := memorystorage.New()
	server := internalHTTP.NewServer(storage, config.HTTPServer.HostPort)
	fmt.Println("WB-bank is running...")
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
