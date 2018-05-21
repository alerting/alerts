package main

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/backends"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/alerting/alerts/cap/xml"
	"github.com/urfave/cli"
)

func createMachineryServer(c *cli.Context) (*machinery.Server, error) {
	var conf config.Config

	conf.Broker = c.GlobalString("machinery-broker")
	conf.DefaultQueue = c.GlobalString("machinery-queue")
	conf.ResultBackend = c.GlobalString("machinery-result-backend")
	conf.ResultsExpireIn = c.GlobalInt("machinery-result-expiry")

	return machinery.NewServer(&conf)
}

func processAlert(ctx context.Context, server *machinery.Server, alert *capxml.Alert) (*backends.AsyncResult, error) {
	// Convert the alert to JSON
	var b bytes.Buffer

	encoder := json.NewEncoder(&b)
	if err := encoder.Encode(&alert); err != nil {
		return nil, err
	}

	// TODO: Pick retry count / handle failures
	task := tasks.Signature{
		UUID: alert.ID(),
		Name: "process_alert",
		Args: []tasks.Arg{
			{
				Type:  "string",
				Value: b.String(),
			},
		},
		RetryCount: 20,
	}

	return server.SendTaskWithContext(ctx, &task)
}

func fetchMissingAlert(ctx context.Context, server *machinery.Server, reference *capxml.Reference) (*backends.AsyncResult, error) {
	// Convert the reference to JSON
	var b bytes.Buffer

	encoder := json.NewEncoder(&b)
	if err := encoder.Encode(&reference); err != nil {
		return nil, err
	}

	// TODO: Pick retry count / handle failures
	task := tasks.Signature{
		UUID: reference.ID(),
		Name: "fetch_reference",
		Args: []tasks.Arg{
			{
				Type:  "string",
				Value: b.String(),
			},
		},
		RetryCount: 20,
	}

	return server.SendTaskWithContext(ctx, &task)
}
