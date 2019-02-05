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
	"net/url"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"

	"github.com/alerting/alerts/pkg/alerts"
	"github.com/alerting/alerts/pkg/resources"

	"github.com/alerting/alerts/pkg/cap"
	raven "github.com/getsentry/raven-go"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes"
	"github.com/gorilla/mux"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func toStatusCode(err error) int {
	st := status.Convert(err)
	if st.Code() == codes.NotFound {
		return http.StatusNotFound
	} else {
		return http.StatusInternalServerError
	}
}

func getTimeQuery(query *url.Values, prefix string) (*alerts.TimeConditions, error) {
	timeQuery := alerts.TimeConditions{}
	if val := query.Get(prefix + "_gte"); val != "" {
		t, err := time.Parse(time.RFC3339, val)
		if err != nil {
			return nil, err
		}
		timeQuery.Gte, err = ptypes.TimestampProto(t)
		if err != nil {
			return nil, err
		}
	}
	if val := query.Get(prefix + "_gt"); val != "" {
		t, err := time.Parse(time.RFC3339, val)
		if err != nil {
			return nil, err
		}
		timeQuery.Gt, err = ptypes.TimestampProto(t)
		if err != nil {
			return nil, err
		}
	}
	if val := query.Get(prefix + "_lte"); val != "" {
		t, err := time.Parse(time.RFC3339, val)
		if err != nil {
			return nil, err
		}
		timeQuery.Lte, err = ptypes.TimestampProto(t)
		if err != nil {
			return nil, err
		}
	}
	if val := query.Get(prefix + "_lt"); val != "" {
		t, err := time.Parse(time.RFC3339, val)
		if err != nil {
			return nil, err
		}
		timeQuery.Lt, err = ptypes.TimestampProto(t)
		if err != nil {
			return nil, err
		}
	}

	if timeQuery.Gte != nil || timeQuery.Gt != nil || timeQuery.Lte != nil || timeQuery.Lt != nil {
		return &timeQuery, nil
	}

	return nil, nil
}

func getAlerts(w http.ResponseWriter, r *http.Request) {
	span := opentracing.SpanFromContext(r.Context())
	ctx := opentracing.ContextWithSpan(context.Background(), span)

	log.Info("Get alerts")

	req := &alerts.FindCriteria{}

	// Assemble the search request
	query := r.URL.Query()

	if val := query.Get("start"); val != "" {
		ival, err := strconv.ParseInt(val, 10, 32)
		if err != nil {
			log.WithError(err).Error("Failed to get start")
			raven.CaptureError(err, nil)
			ext.Error.Set(span, true)
			span.LogFields(otlog.Error(err))
			w.WriteHeader(400)
			w.Write([]byte("Invalid start value"))
			return
		}

		req.Start = int32(ival)
	}

	if val := query.Get("count"); val != "" {
		ival, err := strconv.ParseInt(val, 10, 32)
		if err != nil {
			log.WithError(err).Error("Failed to get count")
			raven.CaptureError(err, nil)
			ext.Error.Set(span, true)
			span.LogFields(otlog.Error(err))
			w.WriteHeader(400)
			w.Write([]byte("Invalid count value"))
			return
		}
		req.Count = int32(ival)
	}

	if val := query.Get("sort"); val != "" {
		req.Sort = strings.Split(val, ",")
	}

	if val := query.Get("superseded"); val != "" {
		b, err := strconv.ParseBool(val)
		if err != nil {
			log.WithError(err).Error("Failed to get superseded")
			raven.CaptureError(err, nil)
			ext.Error.Set(span, true)
			span.LogFields(otlog.Error(err))
			w.WriteHeader(400)
			w.Write([]byte("Invalid superseded value"))
			return
		}

		if b {
			req.Superseded = true
		} else {
			req.NotSuperseded = true
		}
	}

	if val := query.Get("status"); val != "" {
		if status, ok := cap.Alert_Status_value[strings.ToUpper(val)]; ok {
			req.Status = cap.Alert_Status(status)
		} else {
			w.WriteHeader(400)
			w.Write([]byte("Invalid status value"))
			return
		}
	}

	if val := query.Get("messageType"); val != "" {
		if msgType, ok := cap.Alert_MessageType_value[strings.ToUpper(val)]; ok {
			req.MessageType = cap.Alert_MessageType(msgType)
		} else {
			w.WriteHeader(400)
			w.Write([]byte("Invalid messageType value"))
			return
		}
	}

	if val := query.Get("scope"); val != "" {
		if scope, ok := cap.Alert_Scope_value[strings.ToUpper(val)]; ok {
			req.Scope = cap.Alert_Scope(scope)
		} else {
			w.WriteHeader(400)
			w.Write([]byte("Invalid scope value"))
			return
		}
	}

	if val := query.Get("language"); val != "" {
		req.Language = val
	}

	if val := query.Get("certainty"); val != "" {
		if certainty, ok := cap.Info_Certainty_value[strings.ToUpper(val)]; ok {
			req.Certainty = cap.Info_Certainty(certainty)
		} else {
			w.WriteHeader(400)
			w.Write([]byte("Invalid certainty value"))
			return
		}
	}

	if val := query.Get("severity"); val != "" {
		if severity, ok := cap.Info_Severity_value[strings.ToUpper(val)]; ok {
			req.Severity = cap.Info_Severity(severity)
		} else {
			w.WriteHeader(400)
			w.Write([]byte("Invalid severity value"))
			return
		}
	}

	if val := query.Get("urgency"); val != "" {
		if urgency, ok := cap.Info_Urgency_value[strings.ToUpper(val)]; ok {
			req.Urgency = cap.Info_Urgency(urgency)
		} else {
			w.WriteHeader(400)
			w.Write([]byte("Invalid urgency value"))
			return
		}
	}

	var err error
	req.Effective, err = getTimeQuery(&query, "effective")
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Invalid effective value"))
		return
	}

	req.Onset, err = getTimeQuery(&query, "onset")
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Invalid onset value"))
		return
	}

	req.Expires, err = getTimeQuery(&query, "expires")
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Invalid expires value"))
		return
	}

	if val := query.Get("headline"); val != "" {
		req.Headline = val
	}

	if val := query.Get("description"); val != "" {
		req.Description = val
	}

	if val := query.Get("instruction"); val != "" {
		req.Instruction = val
	}

	if val := query.Get("area_description"); val != "" {
		req.AreaDescription = val
	}

	if val := query.Get("point"); val != "" {
		components := strings.Split(val, ",")
		if len(components) != 2 {
			w.WriteHeader(400)
			w.Write([]byte("Invalid point value"))
			return
		}

		lat, err := strconv.ParseFloat(components[0], 64)
		if err != nil {
			log.WithError(err).Error("Failed to get point.lat")
			raven.CaptureError(err, nil)
			ext.Error.Set(span, true)
			span.LogFields(otlog.Error(err))
			w.WriteHeader(400)
			w.Write([]byte("Invalid point.lat value"))
			return
		}

		lon, err := strconv.ParseFloat(components[1], 64)
		if err != nil {
			log.WithError(err).Error("Failed to get point.lon")
			raven.CaptureError(err, nil)
			ext.Error.Set(span, true)
			span.LogFields(otlog.Error(err))
			w.WriteHeader(400)
			w.Write([]byte("Invalid point.lon value"))
			return
		}

		req.Point = &alerts.Coordinate{
			Lat: lat,
			Lon: lon,
		}
	}

	// Do the find
	alerts, err := alertsClient.Find(ctx, req)
	if err != nil {
		log.WithError(err).Error("Failed to get alerts")
		raven.CaptureError(err, nil)
		ext.Error.Set(span, true)
		span.LogFields(otlog.Error(err))
		w.WriteHeader(toStatusCode(err))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	enc := jsonpb.Marshaler{
		EmitDefaults: alerts.Total == 0,
	}
	if err := enc.Marshal(w, alerts); err != nil {
		log.WithError(err).Error("Failed to get marshal result")
		raven.CaptureError(err, nil)
		ext.Error.Set(span, true)
		span.LogFields(otlog.Error(err))
	}
}

func getAlert(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	span := opentracing.SpanFromContext(r.Context())
	ctx := opentracing.ContextWithSpan(context.Background(), span)

	log.WithField("id", vars["id"]).Info("Get alert")
	alert, err := alertsClient.Get(ctx, &cap.Reference{
		Id: vars["id"],
	})
	if err != nil {
		log.WithError(err).Error("Failed to get alert")
		raven.CaptureError(err, nil)
		ext.Error.Set(span, true)
		span.LogFields(otlog.Error(err))

		// Identify appropriate error code
		w.WriteHeader(toStatusCode(err))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	enc := jsonpb.Marshaler{}
	if err := enc.Marshal(w, alert); err != nil {
		log.WithError(err).Error("Failed to marshal result")
		raven.CaptureError(err, nil)
		ext.Error.Set(span, true)
		span.LogFields(otlog.Error(err))
	}
}

func getResource(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	span := opentracing.SpanFromContext(r.Context())
	ctx := opentracing.ContextWithSpan(context.Background(), span)

	log.WithField("filename", vars["filename"]).Info("Get resource")
	resource, err := resourcesClient.Get(ctx, &resources.GetRequest{
		Filename: vars["filename"],
	})
	if err != nil {
		log.WithError(err).Error("Failed to get resource")
		raven.CaptureError(err, nil)
		ext.Error.Set(span, true)
		span.LogFields(otlog.Error(err))
		w.WriteHeader(toStatusCode(err))
		return
	}

	w.Header().Set("Content-Type", resource.MimeType)
	w.Write(resource.Data)
}

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		mux := mux.NewRouter()

		// Alerts
		mux.HandleFunc("/alerts", getAlerts).Methods("GET")
		mux.HandleFunc("/alerts/{id}", getAlert).Methods("GET")

		// Resources
		mux.HandleFunc("/resources/{filename}", getResource).Methods("GET")

		// Start listening
		cert := cmd.Flag("cert").Value.String()
		key := cmd.Flag("key").Value.String()
		listen := cmd.Flag("listen").Value.String()

		handler := nethttp.Middleware(opentracing.GlobalTracer(), mux)
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
