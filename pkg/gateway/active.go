package gateway

import (
	"context"
	"net/http"
	"time"

	"github.com/alerting/alerts/pkg/alerts"
	raven "github.com/getsentry/raven-go"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"
	log "github.com/sirupsen/logrus"
)

func GetActiveAlerts(alertsClient alerts.AlertsServiceClient) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		span := opentracing.SpanFromContext(r.Context())
		ctx := opentracing.ContextWithSpan(context.Background(), span)

		// Generate search criteria
		req, err := getAlertsQuery(ctx, r.URL.Query())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Println(req)

		// Overwrite values
		req.Superseded = false
		req.NotSuperseded = true

		now := time.Now()

		req.Effective = &alerts.TimeConditions{}
		req.Effective.Lte, err = ptypes.TimestampProto(now)
		if err != nil {
			log.WithError(err).Error("Failed to process current time")
			raven.CaptureError(err, nil)
			ext.Error.Set(span, true)
			span.LogFields(otlog.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		req.Expires = &alerts.TimeConditions{}
		req.Expires.Gte, err = ptypes.TimestampProto(now)
		if err != nil {
			log.WithError(err).Error("Failed to process current time")
			raven.CaptureError(err, nil)
			ext.Error.Set(span, true)
			span.LogFields(otlog.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Do the find
		log.Println(req)
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
}
