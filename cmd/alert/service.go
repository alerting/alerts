package main

import (
	"context"

	"github.com/alerting/alerts/cap"
	"github.com/alerting/alerts/rpc"
)

// AlertService represents the alert service.
type AlertService interface {
	Has(ctx context.Context, request *rpc.AlertRequest) (bool, error)
	Get(ctx context.Context, request *rpc.AlertRequest) (*cap.Alert, error)
	Find(ctx context.Context, criteria *rpc.AlertsRequest) (*rpc.AlertsResponse, error)
	Add(ctx context.Context, alerts ...*cap.Alert) error
	Supersede(ctx context.Context, request *rpc.AlertRequest) error
	IsSuperseded(ctx context.Context, request *rpc.AlertRequest) (bool, error)
}
