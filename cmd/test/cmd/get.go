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
	"fmt"
	"strings"
	"time"

	"github.com/alerting/alerts/pkg/cap"

	raven "github.com/getsentry/raven-go"
	opentracing "github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Start a span
		span := opentracing.StartSpan("get")
		defer span.Finish()
		ctx := opentracing.ContextWithSpan(context.Background(), span)

		alertID := "670bb380c3efd0548753999e87634e18c50bfd6b"

		// HAS.
		log.WithField("id", "aaaa").Info("Check if alert exists (shouldn't)")
		exists, err := alertsClient.Has(ctx, &cap.Reference{
			Id: "aaaa",
		})
		if err != nil {
			log.WithError(err).Error("Unable to get alert")
			raven.CaptureErrorAndWait(err, nil)
			return
		}

		log.WithField("exists", exists.Result).Info("Got response")

		log.WithField("id", alertID).Info("Check if alert exists")
		exists, err = alertsClient.Has(ctx, &cap.Reference{
			Id: alertID,
		})
		if err != nil {
			log.WithError(err).Error("Unable to get alert")
			raven.CaptureErrorAndWait(err, nil)
			return
		}

		log.WithField("exists", exists.Result).Info("Got response")

		// GET.
		log.WithField("id", alertID).Info("Loading alert")
		alert, err := alertsClient.Get(ctx, &cap.Reference{
			Id: alertID,
		})
		if err != nil {
			log.WithError(err).Error("Unable to get alert")
			raven.CaptureErrorAndWait(err, nil)
			return
		}

		fmt.Println(alert.Identifier)
		fmt.Println(alert.Sender)
		fmt.Println(time.Unix(alert.Sent.GetSeconds(), int64(alert.Sent.GetNanos())))
		fmt.Printf("Infos: (%d)\n", len(alert.Infos))
		for _, info := range alert.Infos {
			fmt.Printf("  %s\n", strings.ToUpper(info.Headline))
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
