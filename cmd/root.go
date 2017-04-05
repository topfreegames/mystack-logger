// mystack/mystack-cli api
// https://github.com/topfreegames/mystack/mystack-controller
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var debug bool
var config *viper.Viper

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "mystack-logger",
	Short: "mystack logger aggregator",
	Long:  `a logger aggregator for mystack deployed apps`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "./config/default.yaml", "config file")
	RootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "turn on debug logs mode")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	config = viper.New()
	if cfgFile != "" { // enable ability to specify config file via flag
		config.SetConfigFile(cfgFile)
	}
	config.SetConfigType("yaml")
	config.SetEnvPrefix("MYSTACK_LOGGER")
	config.AddConfigPath(".")
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	config.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	config.AutomaticEnv()

	// If a config file is found, read it in.
	if err := config.ReadInConfig(); err != nil {
		fmt.Printf("Config file %s failed to load: %s.\n", cfgFile, err.Error())
		panic("Failed to load config file")
	}
}
