package main

import (
	"log"
	"logger/plugin/storage"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	"logger/pkg/version"
	"logger/pkg/config"
	"github.com/spf13/viper"
)

func main(){
	err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }
	log.Println(os.Getenv("LOG_STORAGE_TYPE"))
	storageFactory, err := storage.NewFactory(storage.FactoryConfigFromEnvAndCLI(os.Args, os.Stderr))
	if err != nil {
		log.Fatalf("Cannot initialize storage factory: %v", err)
	}
	v := viper.New()
	command := &cobra.Command{
		Use: "logger",
		Short: "Logger distribution with agent, collector and query in one process.",
		Long: `Logger distribution with agent, collector and query. Use with caution this version
by default uses only in-memory database.`,
        RunE: func(cmd *cobra.Command, args []string) error {
			log.Println("RUN LOGGER")
			return nil
		},
	}
	command.AddCommand(version.Command())
	config.AddFlags(v,command,storageFactory.AddFlags)
	if err := command.Execute(); err != nil {
		log.Fatal(err)
	}
}