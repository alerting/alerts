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
	"context"
	"net/http"

	"github.com/alerting/alerts/pkg/cap"
	raven "github.com/getsentry/raven-go"
	"github.com/golang/protobuf/jsonpb"
	opentracing "github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// resourceCmd represents the resource command
var resourceCmd = &cobra.Command{
	Use:   "resource",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Load an alert from alerts.zacharyseguin.ca
		res, err := http.Get("https://alerts.zacharyseguin.ca/api/alerts/670bb380c3efd0548753999e87634e18c50bfd6b")
		if err != nil {
			log.WithError(err).Error("Unable to fetch alert from alerts.zacharyseguin.ca")
			raven.CaptureErrorAndWait(err, nil)
			return
		}
		defer res.Body.Close()

		var alert cap.Alert
		dec := jsonpb.Unmarshaler{}
		err = dec.Unmarshal(res.Body, &alert)
		if err != nil {
			log.WithError(err).Error("Unable to unmarshal alert")
			raven.CaptureErrorAndWait(err, nil)
			return
		}

		span := opentracing.StartSpan("resources")
		defer span.Finish()
		ctx := opentracing.ContextWithSpan(context.Background(), span)

		for _, info := range alert.Infos {
			for _, resource := range info.Resources {
				log.WithField("uri", resource.Uri).Info("Adding resource")

				newResource, err := resourcesClient.Add(ctx, resource)
				if err != nil {
					log.WithError(err).Error("Unable to unmarshal alert")
					raven.CaptureErrorAndWait(err, nil)
					return
				}

				log.Println(newResource)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(resourceCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// resourceCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// resourceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
