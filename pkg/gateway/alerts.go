package gateway

import (
	"context"
	"net/http"

	"github.com/alerting/alerts/pkg/alerts"
	raven "github.com/getsentry/raven-go"
	"github.com/golang/protobuf/jsonpb"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"
	log "github.com/sirupsen/logrus"
)

func GetAlerts(alertsClient alerts.AlertsServiceClient) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		span := opentracing.SpanFromContext(r.Context())
		ctx := opentracing.ContextWithSpan(context.Background(), span)

		log.Info("Get alerts")

		req, err := getAlertsQuery(ctx, r.URL.Query())

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
}
