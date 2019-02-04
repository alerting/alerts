package resources

import (
	"io"

	context "golang.org/x/net/context"
)

// Storage defines an interface for resources storage.
type Storage interface {
	Add(ctx context.Context, filename string, mimetype string, reader io.Reader) error

	Has(ctx context.Context, filename string) (bool, error)

	Get(ctx context.Context, filename string) ([]byte, string, error)
}
