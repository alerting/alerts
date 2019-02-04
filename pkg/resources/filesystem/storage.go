package filesystem

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/alerting/alerts/pkg/resources"
	raven "github.com/getsentry/raven-go"
	opentracing "github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
)

// Storage defines an filesystem resources storage.
type Storage struct {
	path string
}

// NewStorage creates the storage.
func NewStorage(path string) (resources.Storage, error) {
	// Initialize the path, if it does not exist
	ok, err := exists(path)
	if err != nil {
		log.WithError(err).Error("Failed to check resources directory")
		raven.CaptureErrorAndWait(err, nil)
		return nil, err
	}

	if !ok {
		log.WithField("path", path).Info("Creating resources directory")
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			log.WithError(err).Error("Failed to create resources directory")
			raven.CaptureErrorAndWait(err, nil)
			return nil, err
		}
	}

	// Return storage.
	return &Storage{
		path: path,
	}, nil
}

func (s *Storage) getPath(filename string) string {
	return filepath.Join(s.path, filename)
}

func exists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// Add saves the resource into the resources system.
func (s *Storage) Add(ctx context.Context, filename string, mimetype string, reader io.Reader) error {
	// Write the file.
	f, err := os.Create(s.getPath(filename))
	if err != nil {
		log.WithError(err).Error("Failed to create output file")
		raven.CaptureErrorAndWait(err, nil)
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, reader)
	return err
}

// Has checks if a resource already exists.
func (s *Storage) Has(ctx context.Context, filename string) (bool, error) {
	// Start a new span.
	span, _ := opentracing.StartSpanFromContext(ctx, "Storage.Filesystem::Has")
	defer span.Finish()

	span.SetTag("filename", filename)
	return exists(s.getPath(filename))
}

// Get returns the resource.
func (s *Storage) Get(ctx context.Context, filename string) ([]byte, string, error) {
	// Start a new span.
	span, _ := opentracing.StartSpanFromContext(ctx, "Storage.Filesystem::Get")
	defer span.Finish()

	span.SetTag("filename", filename)

	f, err := os.Open(s.getPath(filename))
	if err != nil {
		log.WithError(err).Error("Failed to open resource")
		raven.CaptureErrorAndWait(err, nil)
		return nil, "", err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		log.WithError(err).Error("Failed to read resource")
		raven.CaptureErrorAndWait(err, nil)
		return nil, "", err
	}

	// Since we don't store mime type with the file, let's try to figure it out.
	return data, http.DetectContentType(data), nil
}
