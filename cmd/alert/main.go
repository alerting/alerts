package main

import (
	"context"
	"io"
	"log"
	"net"
	"os"

	"github.com/alerting/alerts/cap"
	"github.com/alerting/alerts/rpc"
	"github.com/olivere/elastic"
	"github.com/urfave/cli"
	"google.golang.org/grpc"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

var client *elastic.Client
var esIndex string

//var conn stan.Conn

func logError(span opentracing.Span, err error) {
	ext.Error.Set(span, true)
	span.LogFields(otlog.Error(err))
	log.Print(err)
}

type alertServer struct{}

func (s *alertServer) HasAlert(ctx context.Context, r *rpc.AlertRequest) (*rpc.AlertExistsResponse, error) {
	exists, err := exists(ctx, client, esIndex, r.Id)

	return &rpc.AlertExistsResponse{
		Exists: exists,
	}, err
}

func (s *alertServer) GetAlert(ctx context.Context, r *rpc.AlertRequest) (*cap.Alert, error) {
	return get(ctx, client, esIndex, r.Id)
}

func (s *alertServer) FindAlerts(ctx context.Context, r *rpc.AlertsRequest) (*rpc.AlertsResponse, error) {
	return find(ctx, client, esIndex, r)
}

func (s *alertServer) AddAlert(ctx context.Context, alert *cap.Alert) (*rpc.Empty, error) {
	err := save(ctx, client, esIndex, alert)
	return &rpc.Empty{}, err
}

func (s *alertServer) AddAlerts(stream rpc.AlertService_AddAlertsServer) error {
	alerts := make([]*cap.Alert, 0)
	for {
		alert, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		alerts = append(alerts, alert)

		// If we have more than 20 alerts, save them.
		if len(alerts) > 20 {
			if err := save(stream.Context(), client, esIndex, alerts...); err != nil {
				return err
			}

			alerts = make([]*cap.Alert, 0)
		}
	}

	// Save any remaining alerts.
	if len(alerts) > 0 {
		if err := save(stream.Context(), client, esIndex, alerts...); err != nil {
			return err
		}
	}

	return stream.SendAndClose(&rpc.Empty{})
}

func newAlertServer() *alertServer {
	return &alertServer{}
}

func serve(c *cli.Context) error {
	var err error

	client, err = elastic.NewClient(
		elastic.SetURL(c.GlobalString("elastic-url")),
		elastic.SetSniff(c.GlobalBoolT("elastic-sniff")))
	if err != nil {
		return err
	}

	if err = setup(context.Background(), client, esIndex); err != nil {
		return err
	}

	// Setup NATS
	/*
		conn, err = stan.Connect(
			"test-cluster",
			"alert-api",
			stan.NatsURL("nats://127.0.0.1:4223"))
		if err != nil {
			return err
		}
	*/

	// Setup gRPC
	server := grpc.NewServer(
		grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer())),
		grpc.StreamInterceptor(otgrpc.OpenTracingStreamServerInterceptor(opentracing.GlobalTracer())),
	)
	rpc.RegisterAlertServiceServer(server, newAlertServer())

	listener, err := net.Listen("tcp", c.String("listen"))
	if err != nil {
		return err
	}

	log.Printf("Serving on %s", c.String("listen"))
	return server.Serve(listener)
}

func main() {
	// Setup tracing
	cfg, err := jaegercfg.FromEnv()
	if err != nil {
		log.Fatalf("Unable to initialize tracer: %s, using default configuration\n", err)
	}

	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		log.Printf("Unable to initialize tracer: %s\n", err)

		tracer = opentracing.NoopTracer{}
	}
	if closer != nil {
		defer closer.Close()
	}

	opentracing.SetGlobalTracer(tracer)

	// App
	app := cli.NewApp()

	app.Name = "alert"
	app.Usage = "Alert API service"
	app.Authors = []cli.Author{
		{
			Name:  "Zachary Seguin",
			Email: "zachary@zacharyseguin.ca",
		},
	}
	app.Copyright = "Copyright (c) 2018 Zachary Seguin"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "elastic-url",
			Usage:  "URL for elasticsearch",
			EnvVar: "ELASTIC_URL",
			Value:  "http://localhost:9200",
		},
		cli.StringFlag{
			Name:        "elastic-index",
			Usage:       "Index for elasticsearch",
			EnvVar:      "ELASTIC_INDEX",
			Value:       "alerts",
			Destination: &esIndex,
		},
		cli.BoolTFlag{
			Name:   "elastic-sniff",
			Usage:  "Sniff for elasticsearch endpoints",
			EnvVar: "ELASTIC_SNIFF",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:   "serve",
			Usage:  "Serves the alert service",
			Action: serve,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "listen",
					Usage:  "host:port to listen on",
					EnvVar: "LISTEN",
					Value:  ":2400",
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
