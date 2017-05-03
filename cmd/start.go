// mystack/mystack-logger
// https://github.com/topfreegames/mystack/mystack-controller
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/topfreegames/mystack-logger/api"
	"github.com/topfreegames/mystack-logger/logger"
	"github.com/topfreegames/mystack-logger/storage"
)

var bind string
var port int
var unsecure bool

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "starts mystack-logger",
	Long:  `starts mystack-logger`,
	Run: func(cmd *cobra.Command, args []string) {
		log := logrus.New()
		if debug {
			log.Level = logrus.DebugLevel
		} else {
			log.Level = logrus.InfoLevel
		}
		log.Info("starting mystack logger consumer and storage adapter...")

		storageAdapter, err := storage.NewRedisStorageAdapter(config)
		if err != nil {
			log.Panic(err)
		}
		storageAdapter.Start()
		defer storageAdapter.Stop()

		collector := logger.NewLogCollector(storageAdapter, config)
		err = collector.Start()
		if err != nil {
			log.Panic(err)
		}
		defer collector.Stop()
		log.Info("log collector running")

		app, err := api.NewApp(bind, port, config, log, storageAdapter, collector, unsecure)

		if err != nil {
			log.Panic(err)
		}

		go func() {
			closer, err := app.ListenAndServe()
			if closer != nil {
				defer closer.Close()
			}
			if err != nil {
				log.Panic(err)
			}
		}()

		log.Infof("api listening @ %s:%d", bind, port)

		stoppedCh := collector.Stopped()
		select {
		case stopErr := <-stoppedCh:
			if err != nil {
				log.Fatal("Log collector has stopped: ", stopErr)
			} else {
				log.Fatal("Log collector has stopped with no error")
			}
		}
	},
}

func init() {
	startCmd.Flags().BoolVar(&unsecure, "unsecure", false, "unsecure api (for development, will default user to testuser")
	startCmd.Flags().StringVarP(&bind, "bind", "b", "0.0.0.0", "the address that mystack logger api will bind to")
	startCmd.Flags().IntVarP(&port, "port", "p", 5000, "the address that mystack logger must listen")
	RootCmd.AddCommand(startCmd)
}
