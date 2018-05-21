package main

import (
	"context"
	"encoding/xml"
	"log"
	"os"

	"github.com/alerting/alerts/cap"
	"github.com/alerting/alerts/rpc"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"google.golang.org/grpc"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

var grpcConn *grpc.ClientConn

func logError(span opentracing.Span, err error) {
	ext.Error.Set(span, true)
	span.LogFields(otlog.Error(err))
	log.Print(err)
}

func loadAlert(ctx context.Context, client rpc.AlertService_AddAlertsClient, filename string) error {
	span := opentracing.SpanFromContext(ctx)

	span.LogEventWithPayload("filename", filename)

	// Read the file from the filesystem
	f, err := os.Open(filename)
	if err != nil {
		logError(span, err)
		return err
	}

	var alert cap.Alert
	xmlDecoder := xml.NewDecoder(f)
	if err := xmlDecoder.Decode(&alert); err != nil {
		logError(span, err)
		return err
	}

	return client.Send(&alert)
}

func main() {
	// Setup tracing
	cfg, err := jaegercfg.FromEnv()
	if err != nil {
		log.Fatalf("Unable to initialize tracer: %s\n", err)
	}

	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		log.Fatalf("Unable to initialize tracer: %s\n", err)
	}
	defer closer.Close()

	opentracing.SetGlobalTracer(tracer)

	// Setup the gRPC client
	grpcConn, err = grpc.Dial("localhost:2400",
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(tracer)),
		grpc.WithStreamInterceptor(otgrpc.OpenTracingStreamClientInterceptor(tracer)),
	)

	// Loop
	span := opentracing.StartSpan("import")
	defer span.Finish()

	// Create client
	ctx := opentracing.ContextWithSpan(context.Background(), span)
	client, err := rpc.NewAlertServiceClient(grpcConn).AddAlerts(ctx)
	if err != nil {
		logError(span, err)
		log.Fatal(err)
	}

	for _, filename := range os.Args[1:] {
		loadAlert(ctx, client, filename)
	}

	_, err = client.CloseAndRecv()
	if err != nil {
		log.Fatal(err)
	}
}
