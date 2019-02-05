package resources

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	cap "github.com/alerting/alerts/pkg/cap"
	raven "github.com/getsentry/raven-go"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"
	log "github.com/sirupsen/logrus"
	context "golang.org/x/net/context"
)

// Server is GRPC server for the resource microservice.
type Server struct {
	Storage Storage
}

// NewServer creates a new Server.
func NewServer(storage Storage) (*Server, error) {
	return &Server{
		Storage: storage,
	}, nil
}

func decode(b64 []byte) (io.Reader, error) {
	var data []byte
	_, err := base64.StdEncoding.Decode(data, b64)
	return bytes.NewReader(data), err
}

func fetch(ctx context.Context, uri *url.URL) (io.ReadCloser, error) {
	// Start a new span.
	span, _ := opentracing.StartSpanFromContext(ctx, "fetch")
	defer span.Finish()

	log.WithField("url", uri).Info("Fetching")
	span.SetTag("url", uri)

	// Fetch the resource.
	res, err := http.Get(uri.String())
	if err != nil {
		return nil, err
	}

	log.WithField("status", res.StatusCode).Info("Got response")
	span.SetTag("status", res.StatusCode)

	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		return res.Body, nil
	}

	return nil, fmt.Errorf("Unexpected response code %d", res.StatusCode)
}

// Add adds a new resource.
func (s *Server) Add(ctx context.Context, resource *cap.Resource) (*cap.Resource, error) {
	// Start a new span.
	span, sctx := opentracing.StartSpanFromContext(ctx, "Server::Add")
	defer span.Finish()

	log.WithField("uri", resource.Uri).Info("Add resource")

	log.WithFields(log.Fields{
		"uri":      resource.Uri,
		"mime":     resource.MimeType,
		"checksum": resource.Checksum(),
	}).Debug("Processing resource")

	// If we have a URL resource, just return it the way it is.
	if resource.MimeType == "application/x-url" {
		return resource, nil
	}

	// Identify the destination filename.
	uri, err := url.Parse(resource.Uri)
	if err != nil {
		log.WithError(err).Error("Failed to parse URI")
		raven.CaptureError(err, nil)
		ext.Error.Set(span, true)
		span.LogFields(otlog.Error(err))
		return nil, err
	}

	var filename string

	if resource.Checksum() == "" {
		filename = strings.ToLower(filepath.Base(uri.Path))
	} else {
		filename = strings.ToLower(fmt.Sprintf("%s%s", resource.Checksum(), filepath.Ext(uri.Path)))
	}

	log.WithField("filename", filename).Info("Generated filename")

	// Determine if we have to fetch the file.
	has, err := s.Storage.Has(sctx, filename)
	if err != nil {
		log.WithError(err).Error("Failed to determine if file already exists")
		raven.CaptureError(err, nil)
		ext.Error.Set(span, true)
		span.LogFields(otlog.Error(err))
		return nil, err
	}

	if !has {
		// If we have a DerefUri, decode + write that to file.
		if len(resource.DerefUri) > 0 {
			data := bytes.NewReader(resource.DerefUri)
			err = s.Storage.Add(sctx, filename, resource.MimeType, data)
			if err != nil {
				log.WithError(err).Error("Failed to write DerefUri resource")
				raven.CaptureError(err, nil)
				ext.Error.Set(span, true)
				span.LogFields(otlog.Error(err))
				return nil, err
			}
		} else if uri.Hostname() != "" {
			// If we have a URL, fetch it.
			data, err := fetch(sctx, uri)
			if err != nil {
				log.WithError(err).Error("Failed to fetch resource")
				raven.CaptureError(err, nil)
				ext.Error.Set(span, true)
				span.LogFields(otlog.Error(err))
				return nil, err
			}
			defer data.Close()

			err = s.Storage.Add(sctx, filename, resource.MimeType, data)
			if err != nil {
				log.WithError(err).Error("Failed to write resource")
				raven.CaptureError(err, nil)
				ext.Error.Set(span, true)
				span.LogFields(otlog.Error(err))
				return nil, err
			}
		} else {
			// We don't have enough information to generate the resource.
			log.Warn("Unable to save resource... Not enough information provided to fetch.")
			return nil, errors.New("Not enough information to save resource")
		}
	} else {
		log.Debug("Resource already exists")
	}

	// Return the updated resource object
	return &cap.Resource{
		Uri:      filename,
		MimeType: resource.MimeType,
		Digest:   resource.Checksum(),
		Size:     resource.Size,
	}, nil
}

func (s *Server) Get(ctx context.Context, req *GetRequest) (*GetResponse, error) {
	// Start a new span.
	span, sctx := opentracing.StartSpanFromContext(ctx, "Server::Get")
	defer span.Finish()

	log.WithField("uri", req.Filename).Info("Get resource")

	data, mimeType, err := s.Storage.Get(sctx, req.Filename)
	if err != nil {
		log.WithError(err).Error("Failed to get resource")
		raven.CaptureError(err, nil)
		ext.Error.Set(span, true)
		span.LogFields(otlog.Error(err))
		return nil, err
	}

	return &GetResponse{
		Data:     data,
		MimeType: mimeType,
	}, nil
}
