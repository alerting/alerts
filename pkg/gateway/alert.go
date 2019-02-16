package gateway

import (
	"context"
	"net/http"

	"github.com/alerting/alerts/pkg/alerts"
	"github.com/alerting/alerts/pkg/cap"
	raven "github.com/getsentry/raven-go"
	"github.com/golang/protobuf/jsonpb"
	"github.com/gorilla/mux"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"
	log "github.com/sirupsen/logrus"
)

func GetAlert(alertsClient alerts.AlertsServiceClient) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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
}
