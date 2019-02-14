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
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/alerting/alerts/pkg/alerts"
	raven "github.com/getsentry/raven-go"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	opentracing "github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve the Alerts microservice",
	Long:  `Provides the Alerts microservice.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize Alerts service.
		service, err := alerts.NewServer(storage, resourcesClient)
		if err != nil {
			log.WithError(err).Error("Unable to create server")
			raven.CaptureErrorAndWait(err, nil)
			return
		}

		// Setup the gRPC server.
		grpcServer := grpc.NewServer(
			grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer())),
			grpc.StreamInterceptor(otgrpc.OpenTracingStreamServerInterceptor(opentracing.GlobalTracer())),
			grpc.MaxMsgSize(1024*1024*1024),
		)
		alerts.RegisterAlertsServiceServer(grpcServer, service)

		// Setup gRPC security.
		var listener net.Listener

		cert := cmd.Flag("cert").Value.String()
		key := cmd.Flag("key").Value.String()
		listen := cmd.Flag("listen").Value.String()

		if cert != "" && key != "" {
			log.WithFields(log.Fields{
				"cert": cert,
				"key":  key,
			}).Info("Configuring TLS")

			keypair, err := tls.LoadX509KeyPair(cert, key)
			if err != nil {
				log.WithError(err).Error("Unable to load TLS keypair")
				raven.CaptureErrorAndWait(err, nil)
				return
			}

			listener, err = tls.Listen("tcp", listen, &tls.Config{
				Certificates: []tls.Certificate{keypair},
			})
			if err != nil {
				log.WithError(err).Error("Unable to listen")
				raven.CaptureErrorAndWait(err, nil)
				return
			}
		} else {
			listener, err = net.Listen("tcp", listen)
			if err != nil {
				log.WithError(err).Error("Unable to listen")
				raven.CaptureErrorAndWait(err, nil)
				return
			}
		}

		// Listen for signal to stop serving.
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			sig := <-sigs
			log.WithField("signal", sig).Info("Stopping server")
			grpcServer.Stop()
		}()

		// Initialize the HTTP server.
		log.Infof("Listening on %s", listen)
		err = grpcServer.Serve(listener)

		if err != nil {
			log.WithError(err).Error("Error serving gRPC")
			raven.CaptureErrorAndWait(err, nil)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().String("listen", ":2400", "Address to listen on.")
	serveCmd.Flags().String("cert", "", "TLS certificate")
	serveCmd.Flags().String("key", "", "TLS key")
}
