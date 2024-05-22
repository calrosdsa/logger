package main

import (
	"fmt"
	// "io"
	"log"
	"logger/cmd/collector/app"
	"logger/cmd/collector/app/flags"
	cmdFlags "logger/cmd/internal/flags"

	"logger/cmd/internal/docs"
	"logger/cmd/internal/env"
	"logger/cmd/internal/printconfig"
	"logger/cmd/internal/status"

	"logger/pkg/config"

	"logger/internal/metrics/expvar"
	"logger/internal/metrics/fork"
	"logger/pkg/metrics"
	"logger/pkg/tenancy"
	"logger/pkg/version"
	"logger/plugin/storage"
	"logger/ports"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const serviceName = "collector"

func main(){
	svc := cmdFlags.NewService(ports.CollectorAdminHTTP)
	storageFactory, err := storage.NewFactory(storage.FactoryConfigFromEnvAndCLI(os.Args, os.Stderr))
	if err != nil {
		log.Fatalf("Cannot initialize storage factory: %v", err)
	}
	v := viper.New()

	command := &cobra.Command{
		Use:   "logger-collector",
		Short: "Logger collector receives and processes traces from Logger agents and clients",
		Long:  `Logger collector receives traces from Logger agents and runs them through a processing pipeline.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := svc.Start(v); err != nil {
				return err
			}
			logger := svc.Logger // shortcut
			baseFactory := svc.MetricsFactory.Namespace(metrics.NSOptions{Name: "logger"})
			metricsFactory := fork.New("internal",
				expvar.NewFactory(10), // backend for internal opts
				baseFactory.Namespace(metrics.NSOptions{Name: "collector"}))
			version.NewInfoMetrics(metricsFactory)

			storageFactory.InitFromViper(v, logger)
			if err := storageFactory.Initialize(baseFactory, logger); err != nil {
				logger.Fatal("Failed to init storage factory", zap.Error(err))
			}
			logWriter, err := storageFactory.CreateLogWriter()
			if err != nil {
				logger.Fatal("Failed to create span writer", zap.Error(err))
			}

			// ssFactory, err := storageFactory.CreateSamplingStoreFactory()
			// if err != nil {
			// 	logger.Fatal("Failed to create sampling store factory", zap.Error(err))
			// }

			// strategyStoreFactory.InitFromViper(v, logger)
			// if err := strategyStoreFactory.Initialize(metricsFactory, ssFactory, logger); err != nil {
			// 	logger.Fatal("Failed to init sampling strategy store factory", zap.Error(err))
			// }
			// strategyStore, aggregator, err := strategyStoreFactory.CreateStrategyStore()
			// if err != nil {
			// 	logger.Fatal("Failed to create sampling strategy store", zap.Error(err))
			// }
			collectorOpts, err := new(flags.CollectorOptions).InitFromViper(v, logger)
			if err != nil {
				logger.Fatal("Failed to initialize collector", zap.Error(err))
			}
			tm := tenancy.NewManager(&collectorOpts.GRPC.Tenancy)

			collector := app.New(&app.CollectorParams{
				ServiceName:    serviceName,
				Logger:         logger,
				MetricsFactory: metricsFactory,
				LogWriter:     logWriter,
				// StrategyStore:  strategyStore,
				// Aggregator:     aggregator,
				HealthCheck:    svc.HC(),
				TenancyMgr:     tm,
			})
			// Start all Collector services
			if err := collector.Start(collectorOpts); err != nil {
				logger.Fatal("Failed to start collector", zap.Error(err))
			}
			// Wait for shutdown
			// svc.RunAndThen(func() {
			// 	if err := collector.Close(); err != nil {
			// 		logger.Error("failed to cleanly close the collector", zap.Error(err))
			// 	}
			// 	if closer, ok := logWriter.(io.Closer); ok {
			// 		err := closer.Close()
			// 		if err != nil {
			// 			logger.Error("failed to close span writer", zap.Error(err))
			// 		}
			// 	}
			// 	if err := storageFactory.Close(); err != nil {
			// 		logger.Error("Failed to close storage factory", zap.Error(err))
			// 	}
			// })
			return nil
		},
	}

	command.AddCommand(version.Command())
	command.AddCommand(env.Command())
	command.AddCommand(docs.Command(v))
	command.AddCommand(status.Command(v, ports.CollectorAdminHTTP))
	command.AddCommand(printconfig.Command(v))

	config.AddFlags(
		v,
		command,
		svc.AddFlags,
		flags.AddFlags,
		// storageFactory.AddPipelineFlags,
		// strategyStoreFactory.AddFlags,
	)

	if err := command.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

}