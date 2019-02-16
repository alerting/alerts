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
	"net/http"

	"github.com/alerting/alerts/pkg/gateway"

	raven "github.com/getsentry/raven-go"
	"github.com/gorilla/mux"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	opentracing "github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run the HTTP -> Alerts gateway",
	Run: func(cmd *cobra.Command, args []string) {
		router := mux.NewRouter()

		// Alerts
		router.HandleFunc("/alerts", gateway.GetAlerts(alertsClient)).Methods("GET")
		router.HandleFunc("/alerts/active", gateway.GetActiveAlerts(alertsClient)).Methods("GET")
		router.HandleFunc("/alerts/{id}", gateway.GetAlert(alertsClient)).Methods("GET")

		// Resources
		router.HandleFunc("/resources/{filename}", gateway.GetResource(resourcesClient)).Methods("GET")

		// Start listening
		cert := cmd.Flag("cert").Value.String()
		key := cmd.Flag("key").Value.String()
		listen := cmd.Flag("listen").Value.String()

		handler := nethttp.Middleware(opentracing.GlobalTracer(), router)
		var err error

		log.Infof("Listening on %s", listen)
		if cert != "" && key != "" {
			err = http.ListenAndServeTLS(listen, cert, key, handler)
		} else {
			err = http.ListenAndServe(listen, handler)
		}

		if err != nil {
			log.WithError(err).Error("Error serving")
			raven.CaptureErrorAndWait(err, nil)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().String("listen", ":8080", "Address to listen on.")
	serveCmd.Flags().String("cert", "", "TLS certificate")
	serveCmd.Flags().String("key", "", "TLS key")
}
