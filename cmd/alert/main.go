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

func logError(span opentracing.Span, err error) {
	ext.Error.Set(span, true)
	span.LogFields(otlog.Error(err))
	log.Print(err)
}

type alertServer struct {
	Service AlertService
}

func (s *alertServer) HasAlert(ctx context.Context, r *rpc.AlertRequest) (*rpc.AlertExistsResponse, error) {
	exists, err := s.Service.Has(ctx, r)

	return &rpc.AlertExistsResponse{
		Exists: exists,
	}, err
}

func (s *alertServer) GetAlert(ctx context.Context, r *rpc.AlertRequest) (*cap.Alert, error) {
	return s.Service.Get(ctx, r)
}

func (s *alertServer) FindAlerts(ctx context.Context, r *rpc.AlertsRequest) (*rpc.AlertsResponse, error) {
	return s.Service.Find(ctx, r)
}

func (s *alertServer) processAlert(ctx context.Context, alert *cap.Alert) error {
	// Do some cleanup to ensure the IDs are correct
	alert.Id = alert.ID()

	for _, ref := range alert.References {
		ref.Id = ref.ID()
	}

	// Check if the alert has been superseded
	isSuperseded, err := s.Service.IsSuperseded(ctx, &rpc.AlertRequest{
		Id: alert.Id,
	})
	if err != nil {
		log.Println("error: IsSupserseded")
		//return err
	}
	alert.Superseded = isSuperseded

	// Supersede referenced alerts
	log.Println(alert.MessageType)
	if alert.MessageType == cap.Alert_UPDATE || alert.MessageType == cap.Alert_CANCEL {
		for _, ref := range alert.References {
			err := s.Service.Supersede(ctx, &rpc.AlertRequest{
				Id: ref.Id,
			})
			if err != nil {
				// TODO: If 404 ignore, else error
				log.Println(err)
			}
		}
	}

	return nil
}

func (s *alertServer) AddAlert(ctx context.Context, alert *cap.Alert) (*rpc.Empty, error) {
	s.processAlert(ctx, alert)
	err := s.Service.Add(ctx, alert)
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

		s.processAlert(stream.Context(), alert)

		alerts = append(alerts, alert)

		// If we have more than 20 alerts, save them.
		if len(alerts) > 20 {
			if err := s.Service.Add(stream.Context(), alerts...); err != nil {
				return err
			}

			alerts = make([]*cap.Alert, 0)
		}
	}

	// Save any remaining alerts.
	if len(alerts) > 0 {
		if err := s.Service.Add(stream.Context(), alerts...); err != nil {
			return err
		}
	}

	return stream.SendAndClose(&rpc.Empty{})
}

func newAlertServer(service AlertService) *alertServer {
	return &alertServer{
		Service: service,
	}
}

func serve(c *cli.Context) error {
	var err error

	client, err := elastic.NewClient(
		elastic.SetURL(c.GlobalString("elastic-url")),
		elastic.SetSniff(c.GlobalBoolT("elastic-sniff")),
		elastic.SetHealthcheck(c.GlobalBoolT("elastic-healthcheck")))
	if err != nil {
		return err
	}

	service, err := NewElasticsearchAlertService(context.Background(), &AlertElasticsearchServiceConfig{
		Client: client,
		Index:  c.GlobalString("elastic-index"),
	})
	if err != nil {
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
	rpc.RegisterAlertServiceServer(server, newAlertServer(service))

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
			Name:   "elastic-index",
			Usage:  "Index for elasticsearch",
			EnvVar: "ELASTIC_INDEX",
			Value:  "alerts",
		},
		cli.BoolTFlag{
			Name:   "elastic-sniff",
			Usage:  "Sniff for elasticsearch endpoints",
			EnvVar: "ELASTIC_SNIFF",
		},
		cli.BoolTFlag{
			Name:   "elastic-healthcheck",
			Usage:  "Healthcheck for elasticsearch endpoints",
			EnvVar: "ELASTIC_HEALTHCHECK",
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
