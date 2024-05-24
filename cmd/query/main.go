package main

import (
	"fmt"
	"log"
	"logger/cmd/internal/docs"
	"logger/cmd/internal/env"
	"logger/cmd/internal/flags"
	"logger/cmd/internal/printconfig"
	"logger/cmd/internal/status"
	"logger/pkg/config"
	"logger/pkg/metrics"
	"logger/pkg/version"
	"logger/plugin/storage"
	"logger/ports"

	"os"

	"logger/cmd/query/app"
	"logger/cmd/query/app/querysvc"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main(){
	svc := flags.NewService(ports.QueryAdminHTTP)
	storageFactory, err := storage.NewFactory(storage.FactoryConfigFromEnvAndCLI(os.Args, os.Stderr))
	if err != nil {
		log.Fatalf("Cannot initialize storage factory: %v", err)
	}
	v := viper.New()
	command := &cobra.Command{
		Use:   "query",
		Short: "query service provides an API for accessing log data.",
		Long:  `query service provides an API for accessing log data.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := svc.Start(v); err != nil {
				return err
			}
			baseFactory := svc.MetricsFactory.Namespace(metrics.NSOptions{Name: "logger"})

			logger := svc.Logger // shortcut
			queryOpts, err := new(app.QueryOptions).InitFromViper(v, logger)
			if err != nil {
				logger.Fatal("Failed to configure query service", zap.Error(err))
			}
			storageFactory.InitFromViper(v,logger)
			if err := storageFactory.Initialize(baseFactory, logger); err != nil {
				logger.Fatal("Failed to init storage factory", zap.Error(err))
			}
			logReader, err := storageFactory.CreateLogReader()
			if err != nil {
				logger.Fatal("Failed to create span reader", zap.Error(err))
			}
			// queryServiceOptions := queryOpts.BuildQueryServiceOptions(storageFactory, logger)
			queryService := querysvc.NewQueryService(logReader)
			server, err := app.NewServer(svc.Logger, svc.HC(), queryService, queryOpts)
			if err != nil {
				logger.Fatal("Failed to create server", zap.Error(err))
			}

			if err := server.Start(); err != nil {
				logger.Fatal("Could not start servers", zap.Error(err))
			}
			
			svc.RunAndThen(func() {
				server.Close()
				if err := storageFactory.Close(); err != nil {
					logger.Error("Failed to close storage factory", zap.Error(err))
				}
				// if err = jt.Close(context.Background()); err != nil {
				// 	logger.Fatal("Error shutting down tracer provider", zap.Error(err))
				// }
			})
			return nil

		},
	}

	command.AddCommand(version.Command())
	command.AddCommand(env.Command())
	command.AddCommand(docs.Command(v))
	command.AddCommand(status.Command(v, ports.QueryAdminHTTP))
	command.AddCommand(printconfig.Command(v))

	config.AddFlags(
		v,
		command,
		svc.AddFlags,
		storageFactory.AddFlags,
		app.AddFlags,
		// metricsReaderFactory.AddFlags,
		// add tenancy flags here to avoid panic caused by double registration in all-in-one
		// tenancy.AddFlags,
	)

	if err := command.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}