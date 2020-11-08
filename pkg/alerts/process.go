package alerts

import (
	opentracing "github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	context "golang.org/x/net/context"
	cap "zacharyseguin.ca/alerts/pkg/cap"
)

// cleanupAlert cleans alert values.
func cleanupAlert(ctx context.Context, alert *cap.Alert) error {
	// Create a new span.
	span, _ := opentracing.StartSpanFromContext(ctx, "cleanupAlert")
	defer span.Finish()

	// Don't trust provided Id values, instead re-calculate them.
	log.Debug("Updating alert ID and references IDs")
	alert.Id = alert.ID()

	for _, ref := range alert.References {
		ref.Id = ref.ID()
	}

	return nil
}
