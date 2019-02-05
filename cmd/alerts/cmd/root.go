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
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/alerting/alerts/internal/tracing"
	"github.com/alerting/alerts/pkg/alerts"
	alertselastic "github.com/alerting/alerts/pkg/alerts/elasticsearch"
	"github.com/alerting/alerts/pkg/resources"
	raven "github.com/getsentry/raven-go"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/olivere/elastic"
	opentracing "github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var cfgFile string
var debug bool
var storageType, elasticsearchURL, elasticsearchIndex string
var elasticsearchSniff, elasticsearchHealthCheck bool

var storage alerts.Storage

// Resources client
var resourcesService, resourcesCA string
var resourcesInsecure bool

var resourcesGrpcConn *grpc.ClientConn
var resourcesClient resources.ResourcesServiceClient

// OpenTracing
var closer io.Closer

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "alerts",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//Run: func(cmd *cobra.Command, args []string) {
	//
	//},

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
		case "elasticsearch":
			opts := []elastic.ClientOptionFunc{
				elastic.SetSniff(elasticsearchSniff),
				elastic.SetHealthcheck(elasticsearchHealthCheck),
			}

			var err error
			storage, err = alertselastic.NewStorage(elasticsearchURL, elasticsearchIndex, opts...)
			if err != nil {
				log.WithError(err).Error("Unable to create storage")
				raven.CaptureErrorAndWait(err, nil)
			}
		default:
			log.WithField("storage-type", storageType).Warn("Unknown storage type")
		}

		// Connect to resources service.
		var resourcesSecurity grpc.DialOption

		if resourcesInsecure {
			resourcesSecurity = grpc.WithInsecure()
		} else {
			var roots *x509.CertPool

			if resourcesCA == "" {
				roots, err = x509.SystemCertPool()
				if err != nil {
					log.WithError(err).Error("Unable to load system certificate store")
					raven.CaptureErrorAndWait(err, nil)
					return
				}
			} else {
				roots = x509.NewCertPool()
				certs, err := ioutil.ReadFile(resourcesCA)
				if err != nil {
					log.WithError(err).WithField("ca", resourcesCA).Error("Unable to load CA certificate")
					raven.CaptureErrorAndWait(err, nil)
					return
				}

				if !roots.AppendCertsFromPEM(certs) {
					log.WithField("ca", resourcesCA).Warn("No certificates loaded")
				}
			}

			resourcesSecurity = grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
				RootCAs:            roots,
				InsecureSkipVerify: false,
			}))
		}

		log.WithFields(log.Fields{
			"target":   resourcesService,
			"ca":       resourcesCA,
			"insecure": fmt.Sprintf("%t", resourcesInsecure),
		}).Info("Dialing resources service")

		resourcesGrpcConn, err = grpc.Dial(
			resourcesService,
			grpc.WithMaxMsgSize(1024*1024*1024),
			resourcesSecurity,
			grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
			grpc.WithStreamInterceptor(otgrpc.OpenTracingStreamClientInterceptor(opentracing.GlobalTracer())),
		)
		if err != nil {
			log.WithError(err).Error("Unable to connect to resources service")
			raven.CaptureErrorAndWait(err, map[string]string{
				"service":  resourcesService,
				"ca":       resourcesCA,
				"insecure": fmt.Sprintf("%t", resourcesInsecure),
			})
			return
		}
		log.Info("Connected to resources service")

		resourcesClient = resources.NewResourcesServiceClient(resourcesGrpcConn)
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		log.Debug("Closing tracing")
		if closer != nil {
			closer.Close()
		}

		log.Debug("Closing resources connection")
		if resourcesGrpcConn != nil {
			resourcesGrpcConn.Close()
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
	rootCmd.PersistentFlags().StringVar(&storageType, "storage-type", "elasticsearch", "Storage type (options: elasticsearch)")
	rootCmd.PersistentFlags().StringVar(&elasticsearchURL, "elasticsearch-url", "http://localhost:9200", "ElasticSearch URL")
	rootCmd.PersistentFlags().StringVar(&elasticsearchIndex, "elasticsearch-index", "alerts", "ElasticSearch index")
	rootCmd.PersistentFlags().BoolVar(&elasticsearchSniff, "elasticsearch-sniff", true, "ElasticSearch sniff")
	rootCmd.PersistentFlags().BoolVar(&elasticsearchHealthCheck, "elasticsearch-healthcheck", true, "ElasticSearch health check")

	// Resources flags
	rootCmd.PersistentFlags().StringVar(&resourcesService, "resources-service", "localhost:2401", "Resources service address")
	rootCmd.PersistentFlags().StringVar(&resourcesCA, "resources-ca", "", "CA certificate for resources service")
	rootCmd.PersistentFlags().BoolVar(&resourcesInsecure, "resources-insecure", false, "Use an insecure connection to resources")
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
