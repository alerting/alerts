package alerts

import (
	context "golang.org/x/net/context"
	cap "zacharyseguin.ca/alerts/pkg/cap"
)

// Storage defines an interface for alert storage.
type Storage interface {
	Add(ctx context.Context, alert *cap.Alert) error
	Get(ctx context.Context, reference *cap.Reference) (*cap.Alert, error)
	Has(ctx context.Context, reference *cap.Reference) (bool, error)
	Find(ctx context.Context, criteria *FindCriteria) (*FindResult, error)

	Supersede(ctx context.Context, reference *cap.Reference) error
	IsSuperseded(ctx context.Context, reference *cap.Reference) (bool, error)
}
