package main

import (
	"log"
	"logger/plugin/storage"
	"os"
)

func main(){
	_, err := storage.NewFactory(storage.FactoryConfigFromEnvAndCLI(os.Args, os.Stderr))
	if err != nil {
		log.Fatalf("Cannot initialize storage factory: %v", err)
	}
}