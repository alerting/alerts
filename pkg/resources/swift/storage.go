package swift

import (
	"context"
	"io"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/alerting/alerts/pkg/resources"
	"github.com/ncw/swift"
	opentracing "github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
)

type Storage struct {
	connection *swift.Connection
	container  string
}

// NewStorage creates the storage.
func NewStorage(username, domain, apiKey, authURL, region, tenant, tenantDomain, container string) (resources.Storage, error) {
	conn := &swift.Connection{
		UserName:     username,
		Domain:       domain,
		ApiKey:       apiKey,
		AuthUrl:      authURL,
		Region:       region,
		Tenant:       tenant,
		TenantDomain: tenantDomain,
		AuthVersion:  3,
	}

	log.Println(conn)

	// Authenticate
	err := conn.Authenticate()
	if err != nil {
		return nil, err
	}

	// Attempt to load the container
	_, _, err = conn.Container(container)
	if err == swift.ContainerNotFound {
		log.WithField("container", container).Println("Creating container")
		err = conn.ContainerCreate(container, nil)
	}

	if err != nil {
		return nil, err
	}

	return &Storage{
		connection: conn,
		container:  container,
	}, nil
}

func (s *Storage) Add(ctx context.Context, filename string, mimetype string, reader io.Reader) error {
	// Start a new span.
	span, _ := opentracing.StartSpanFromContext(ctx, "Storage.Swift::Add")
	defer span.Finish()

	// Create the object
	_, err := s.connection.ObjectPut(s.container, filename, reader, false, "", mimetype, nil)
	return err
}

func (s *Storage) Has(ctx context.Context, filename string) (bool, error) {
	// Start a new span.
	span, _ := opentracing.StartSpanFromContext(ctx, "Storage.Swift::Has")
	defer span.Finish()

	_, _, err := s.connection.Object(s.container, filename)
	if err != nil {
		if err == swift.ObjectNotFound {
			return false, nil
		}

		log.WithError(err).Error("Failed to get object")
		return false, err
	}
	return true, nil
}

func (s *Storage) Get(ctx context.Context, filename string) ([]byte, string, error) {
	// Start a new span.
	span, _ := opentracing.StartSpanFromContext(ctx, "Storage.Swift::Get")
	defer span.Finish()

	obj, _, err := s.connection.Object(s.container, filename)
	if err != nil {
		log.WithError(err).Error("Failed to get object")

		if err == swift.ObjectNotFound {
			return nil, "", status.Error(codes.NotFound, err.Error())
		}

		return nil, "", err
	}

	b, err := s.connection.ObjectGetBytes(s.container, filename)
	if err != nil {
		log.WithError(err).Error("Failed to load object")
	}

	return b, obj.ContentType, nil
}
