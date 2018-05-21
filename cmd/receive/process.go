package main

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/opentracing/opentracing-go"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"

	"github.com/alerting/alerts/cap"
	"github.com/alerting/alerts/rpc"
	"github.com/golang/protobuf/jsonpb"
	"google.golang.org/grpc"

	"github.com/alerting/alerts/cap/xml"
	"github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"
	"github.com/urfave/cli"
)

var system alertSystem

func doProcessAlert(ctx context.Context, alertJSON string) error {
	refSpan := opentracing.SpanFromContext(ctx)
	var span opentracing.Span
	if refSpan != nil {
		span, ctx = opentracing.StartSpanFromContext(ctx, "process_alert", opentracing.FollowsFrom(refSpan.Context()))
	} else {
		span, ctx = opentracing.StartSpanFromContext(ctx, "process_alert")
	}
	defer span.Finish()

	// Decode the alert into an alert object
	var xmlAlert capxml.Alert
	var alert cap.Alert

	// Decode teh capxml.Alert version
	if err := json.Unmarshal([]byte(alertJSON), &xmlAlert); err != nil {
		ext.Error.Set(span, true)
		span.LogFields(otlog.Error(err))
		log.Println(err)
		return err
	}

	// Decode the cap.Alert version
	jd := jsonpb.Unmarshaler{
		AllowUnknownFields: true,
	}
	if err := jd.Unmarshal(strings.NewReader(alertJSON), &alert); err != nil {
		ext.Error.Set(span, true)
		span.LogFields(otlog.Error(err))
		log.Println(err)
		return err
	}

	// Start processing the alert
	log.Printf("Processing %s,%s,%s (%s)", xmlAlert.Sender, xmlAlert.Sent.FormatCAP(), xmlAlert.Identifier, xmlAlert.ID())

	// Check that the referenced alerts exist
	log.Println("Confirming referenced alerts are present")

	for _, reference := range xmlAlert.References {
		res, err := alertServiceClient.HasAlert(ctx, &rpc.AlertRequest{
			Id: reference.ID(),
		})

		if err != nil {
			ext.Error.Set(span, true)
			span.LogFields(otlog.Error(err))
			log.Println(err)
			return err
		}

		if !res.Exists {
			log.Printf("Alert %s,%s,%s (%s) missing", reference.Sender, reference.Sent.FormatCAP(), reference.Identifier, reference.ID())

			if _, err := fetchMissingAlert(ctx, machineryServer, reference); err != nil {
				ext.Error.Set(span, true)
				span.LogFields(otlog.Error(err))
				log.Println(err)
				return err
			}
		}
	}

	// Save the alert if we have an ALERT, UPDATE or CANCEL
	if (alert.Status == cap.Alert_ACTUAL || alert.Status == cap.Alert_EXCERCISE || alert.Status == cap.Alert_TEST) && (alert.MessageType == cap.Alert_ALERT || alert.MessageType == cap.Alert_UPDATE || alert.MessageType == cap.Alert_CANCEL) {
		log.Println("Adding alert")

		// Save the alert
		if _, err := alertServiceClient.AddAlert(ctx, &alert); err != nil {
			log.Println(err)
			return err
		}

		// Supersede the referenced alerts
	}

	log.Printf("Done processing %s,%s,%s (%s)", xmlAlert.Sender, xmlAlert.Sent.FormatCAP(), xmlAlert.Identifier, xmlAlert.ID())
	return nil
}

func doFetchReference(ctx context.Context, referenceJSON string) error {
	refSpan := opentracing.SpanFromContext(ctx)
	var span opentracing.Span
	if refSpan != nil {
		span, ctx = opentracing.StartSpanFromContext(ctx, "fetch_reference", opentracing.FollowsFrom(refSpan.Context()))
	} else {
		span, ctx = opentracing.StartSpanFromContext(ctx, "fetch_reference")
	}
	defer span.Finish()

	var reference *capxml.Reference
	if err := json.Unmarshal([]byte(referenceJSON), &reference); err != nil {
		return err
	}

	log.Printf("Fetching %s,%s,%s (%s)", reference.Sender, reference.Sent.FormatCAP(), reference.Identifier, reference.ID())

	alert, err := system.Fetch(ctx, reference)
	if err != nil {
		return err
	}

	_, err = processAlert(ctx, machineryServer, alert)
	return err
}

func process(c *cli.Context) error {
	// Initialize gRPC
	var err error
	grpcConn, err = grpc.Dial(c.GlobalString("alerts-host"),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
		grpc.WithStreamInterceptor(otgrpc.OpenTracingStreamClientInterceptor(opentracing.GlobalTracer())),
	)

	// Create client
	alertServiceClient = rpc.NewAlertServiceClient(grpcConn)

	// Create the system
	system, err = createSystem(c)
	if err != nil {
		return err
	}

	// Initialize the task processor
	machineryServer, err = createMachineryServer(c)
	if err != nil {
		return err
	}

	tasks := map[string]interface{}{
		"process_alert":   doProcessAlert,
		"fetch_reference": doFetchReference,
	}

	if err = machineryServer.RegisterTasks(tasks); err != nil {
		return err
	}

	// Start the worker
	worker := machineryServer.NewWorker(c.String("tag"), c.Int("concurrency"))
	return worker.Launch()
}
