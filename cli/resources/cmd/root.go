// Copyright © 2018 Zachary Seguin <zachary@zacharyseguin.ca>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/alerting/alerts/resources"
	"github.com/alerting/alerts/resources/filesystem"
	"github.com/alerting/alerts/tracing"
	raven "github.com/getsentry/raven-go"
	homedir "github.com/mitchellh/go-homedir"
	opentracing "github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var debug bool

// Storage
var storageType string
var filesystemPath string

var storage resources.Storage

// OpenTracing
var closer io.Closer

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "alerts",
	Short: "Resources microservice",
	Long:  `Microservice responsible for retrieving and serving alert resources.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Initialize tracing.
		log.Debug("Initializing Jaeger + OpenTracing")

		var tracer opentracing.Tracer
		var err error
		tracer, closer, err = tracing.Init()
		if err != nil {
			log.WithError(err).Warn("Unable to initialize tracer")
			raven.CaptureErrorAndWait(err, nil)

			tracer = opentracing.NoopTracer{}
		}
		opentracing.SetGlobalTracer(tracer)

		// Initialize storage.
		log.Debug("Initialize storage")
		switch storageType {
		case "filesystem":
			log.WithField("path", filesystemPath).Info("Initializing filesystem storage")
			storage, err = filesystem.NewStorage(filesystemPath)

			if err != nil {
				log.WithError(err).Error("Unable to create storage")
				raven.CaptureErrorAndWait(err, nil)
			}
		default:
			log.WithField("storage-type", storageType).Panic("Unknown storage type")
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		log.Debug("Closing tracing")
		if closer != nil {
			closer.Close()
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	raven.CapturePanicAndWait(func() {
		if err := rootCmd.Execute(); err != nil {
			raven.CaptureErrorAndWait(err, nil)
			os.Exit(1)
		}
	}, nil)
}

func init() {
	cobra.OnInitialize(initConfig, initLog, initErrorReporting)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.alerts.yaml)")

	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Show debug messages")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// Storage flags
	rootCmd.PersistentFlags().StringVar(&storageType, "storage-type", "filesystem", "Storage type (options: filesystem, azure)")
	rootCmd.PersistentFlags().StringVar(&filesystemPath, "filesystem-path", "resources", "Path to store resources")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".alerts" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".alerts")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func initLog() {
	if debug {
		log.SetLevel(log.DebugLevel)
	}
}

func initErrorReporting() {
	log.WithField("sentry_dsn", os.Getenv("SENTRY_DSN")).Debug("Initializing error reporting")
	raven.SetDSN(os.Getenv("SENTRY_DSN"))
}
