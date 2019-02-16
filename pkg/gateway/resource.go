package gateway

import (
	"context"
	"net/http"

	"github.com/alerting/alerts/pkg/resources"

	raven "github.com/getsentry/raven-go"
	"github.com/gorilla/mux"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"
	log "github.com/sirupsen/logrus"
)

func GetResource(resourcesClient resources.ResourcesServiceClient) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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
}
