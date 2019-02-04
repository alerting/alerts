package alerts

import (
	"github.com/alerting/alerts/pkg/resources"

	cap "github.com/alerting/alerts/pkg/cap"
	protobuf "github.com/alerting/alerts/pkg/protobuf"
	raven "github.com/getsentry/raven-go"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"
	log "github.com/sirupsen/logrus"
	context "golang.org/x/net/context"
)

// Server is GRPC server for the alerts microservice.
type Server struct {
	Storage   Storage
	Resources resources.ResourcesServiceClient
}

// NewServer creates a new Server.
func NewServer(storage Storage, resources resources.ResourcesServiceClient) (*Server, error) {
	return &Server{
		Storage:   storage,
		Resources: resources,
	}, nil
}

// Add adds a new alert.
func (s *Server) Add(ctx context.Context, alert *cap.Alert) (*cap.Alert, error) {
	// Start a new span.
	span, sctx := opentracing.StartSpanFromContext(ctx, "Server::Add")
	defer span.Finish()

	log.WithField("alert", alert.ID()).Info("Add alert")

	// Run tasks against the incoming alert.
	err := cleanupAlert(sctx, alert)
	if err != nil {
		log.WithError(err).Error("Failed to cleanup alert")
		raven.CaptureError(err, nil)
		ext.Error.Set(span, true)
		span.LogFields(otlog.Error(err))
		return nil, err
	}

	// Process alert resources.
	log.Debug("Processing resources")
	for _, info := range alert.Infos {
		for i, resource := range info.Resources {
			log.WithField("resource", resource.Uri).Debug("Adding resource")
			newResource, err := s.Resources.Add(sctx, resource)
			if err != nil {
				log.WithError(err).Error("Failed to add resource")
				raven.CaptureError(err, nil)
				ext.Error.Set(span, true)
				span.LogFields(otlog.Error(err))
				return nil, err
			}

			info.Resources[i] = newResource
		}
	}

	// Figure out if the alert has already been superseded.
	alert.Superseded, err = s.Storage.IsSuperseded(sctx, alert.Reference())
	if err != nil {
		log.WithError(err).Error("Failed to check if alert has been superseded")
		raven.CaptureError(err, nil)
		ext.Error.Set(span, true)
		span.LogFields(otlog.Error(err))
		return nil, err
	}

	// If the alert is an UPDATE or a CANCEL, the supersede the old versions of the alert.
	if alert.MessageType == cap.Alert_UPDATE || alert.MessageType == cap.Alert_CANCEL {
		log.Debug("Superseding referenced alerts")

		for _, reference := range alert.References {
			log.WithField("id", reference.ID()).Debug("Superseding")
			err := s.Storage.Supersede(sctx, reference)
			if err != nil {
				log.WithError(err).Error("Failed to supersede alert")
				raven.CaptureError(err, nil)
				ext.Error.Set(span, true)
				span.LogFields(otlog.Error(err))
				//return nil, err
			}
		}
	}

	// Add the alert to storage.
	err = s.Storage.Add(sctx, alert)
	return alert, err
}

// Get returns the alert for the given reference.
func (s *Server) Get(ctx context.Context, reference *cap.Reference) (*cap.Alert, error) {
	// Start a new span.
	span, sctx := opentracing.StartSpanFromContext(ctx, "Server::Get")
	defer span.Finish()

	// Load the alert from storage.
	log.WithField("reference", reference).Info("Get alert")
	return s.Storage.Get(sctx, reference)
}

// Has returns whether or not an alert exists for the given reference.
func (s *Server) Has(ctx context.Context, reference *cap.Reference) (*protobuf.BooleanResult, error) {
	// Start a new span.
	span, sctx := opentracing.StartSpanFromContext(ctx, "Server::Has")
	defer span.Finish()

	// Check if the alert exists from storage.
	log.WithField("reference", reference).Info("Has alert")
	exists, err := s.Storage.Has(sctx, reference)

	return &protobuf.BooleanResult{
		Result: exists,
	}, err
}

// Find returns alerts matching the criteria. NOTE: Results are per Info block.
func (s *Server) Find(ctx context.Context, criteria *FindCriteria) (*FindResult, error) {
	// Start a new span.
	span, sctx := opentracing.StartSpanFromContext(ctx, "Server::Find")
	defer span.Finish()

	log.WithField("criteria", criteria).Info("Finding")
	return s.Storage.Find(sctx, criteria)
}
