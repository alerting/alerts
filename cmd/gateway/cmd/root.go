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
	"github.com/alerting/alerts/pkg/resources"
	raven "github.com/getsentry/raven-go"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	homedir "github.com/mitchellh/go-homedir"
	opentracing "github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var cfgFile string
var debug bool

// Alerts service
var alertsService, alertsCA string
var alertsInsecure bool

var resourcesService, resourcesCA string
var resourcesInsecure bool

var alertsGrpcConn *grpc.ClientConn
var alertsClient alerts.AlertsServiceClient

var resourcesGrpcConn *grpc.ClientConn
var resourcesClient resources.ResourcesServiceClient

// Tracing
var closer io.Closer

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gateway",
	Short: "Gateway microservice",
	Long:  `Provides a RESTful interface to alert and resources services.`,
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

		// Connect to alerts service.
		var alertsSecurity grpc.DialOption

		if alertsInsecure {
			alertsSecurity = grpc.WithInsecure()
		} else {
			var roots *x509.CertPool

			if alertsCA == "" {
				roots, err = x509.SystemCertPool()
				if err != nil {
					log.WithError(err).Error("Unable to load system certificate store")
					raven.CaptureErrorAndWait(err, nil)
					return
				}
			} else {
				roots = x509.NewCertPool()
				certs, err := ioutil.ReadFile(alertsCA)
				if err != nil {
					log.WithError(err).WithField("ca", alertsCA).Error("Unable to load CA certificate")
					raven.CaptureErrorAndWait(err, nil)
					return
				}

				if !roots.AppendCertsFromPEM(certs) {
					log.WithField("ca", alertsCA).Warn("No certificates loaded")
				}
			}

			alertsSecurity = grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
				RootCAs:            roots,
				InsecureSkipVerify: false,
			}))
		}

		log.WithFields(log.Fields{
			"target":   alertsService,
			"ca":       alertsCA,
			"insecure": fmt.Sprintf("%t", alertsInsecure),
		}).Info("Dialing alerts service")

		alertsGrpcConn, err = grpc.Dial(
			alertsService,
			grpc.WithMaxMsgSize(1024*1024*1024),
			alertsSecurity,
			grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
			grpc.WithStreamInterceptor(otgrpc.OpenTracingStreamClientInterceptor(opentracing.GlobalTracer())),
		)
		if err != nil {
			log.WithError(err).Error("Unable to connect to alerts service")
			raven.CaptureErrorAndWait(err, map[string]string{
				"service":  alertsService,
				"ca":       alertsCA,
				"insecure": fmt.Sprintf("%t", alertsInsecure),
			})
			return
		}
		log.Info("Connected to alerts service")

		alertsClient = alerts.NewAlertsServiceClient(alertsGrpcConn)

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

		if alertsGrpcConn != nil {
			alertsGrpcConn.Close()
		}

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
	cobra.OnInitialize(initConfig, initErrorReporting)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gateway.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Debug")

	rootCmd.PersistentFlags().StringVar(&alertsService, "alerts-service", "localhost:2400", "Alerts service address")
	rootCmd.PersistentFlags().StringVar(&alertsCA, "alerts-ca", "", "CA certificate for alerts service")
	rootCmd.PersistentFlags().BoolVar(&alertsInsecure, "alerts-insecure", false, "Use an insecure connection to alerts")

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

		// Search config in home directory with name ".gateway" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".gateway")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func initErrorReporting() {
	log.WithField("sentry_dsn", os.Getenv("SENTRY_DSN")).Debug("Initializing error reporting")
	raven.SetDSN(os.Getenv("SENTRY_DSN"))
}
