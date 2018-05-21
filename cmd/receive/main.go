package main

import (
	"log"
	"os"

	machinery "github.com/RichardKnop/machinery/v1"
	"github.com/alerting/alerts/rpc"
	opentracing "github.com/opentracing/opentracing-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
)

var grpcConn *grpc.ClientConn
var alertServiceClient rpc.AlertServiceClient
var machineryServer *machinery.Server

func main() {
	// Setup tracing
	cfg, err := jaegercfg.FromEnv()
	if err != nil {
		log.Fatalf("Unable to initialize tracer: %s\n", err)
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

	app := cli.NewApp()

	// Application information
	app.Name = "receive"
	app.Usage = "Receive alerts from a streaming TCP system"
	app.Authors = []cli.Author{
		{
			Name:  "Zachary Seguin",
			Email: "zachary@zacharyseguin.ca",
		},
	}
	app.Copyright = "Copyright (c) 2018 Zachary Seguin"

	// Command flags
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "alerts-host",
			Usage:  "host:port for alerts service",
			EnvVar: "ALERTS_HOST",
			Value:  "localhost:2400",
		},
		cli.StringFlag{
			Name:   "system, s",
			Usage:  "NAAD system (e.g., naads)",
			EnvVar: "ALERTS_SYSTEM",
		},
		cli.StringFlag{
			Name:   "machinery-broker",
			Usage:  "Broker URL",
			EnvVar: "MACHINERY_BROKER",
		},
		cli.StringFlag{
			Name:   "machinery-queue",
			Usage:  "Default queue",
			EnvVar: "MACHINERY_QUEUE",
		},
		cli.StringFlag{
			Name:   "machinery-result-backend",
			Usage:  "Result backend",
			EnvVar: "MACHINERY_RESULT_BACKEND",
		},
		cli.IntFlag{
			Name:   "machinery-result-expiry",
			Usage:  "Result expirty time (in Seconds)",
			EnvVar: "MACHINERY_RESULT_EXPIRY",
			Value:  3600,
		},
	}

	// Commands
	app.Commands = []cli.Command{
		{
			Name:      "stream",
			Usage:     "Follow the TCP stream and process incoming alerts",
			ArgsUsage: "host:port",
			Action:    stream,
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:   "timeout",
					Usage:  "Timeout for incoming data on the stream",
					EnvVar: "ALERTS_TIMEOUT",
					Value:  80,
				},
			},
		},
		{
			Name:   "process",
			Usage:  "Processes alerts received over the stream and fetches any missing alerts",
			Action: process,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "tag",
					Usage:  "Worker tag (should be unique per worker)",
					EnvVar: "TAG",
				},
				cli.IntFlag{
					Name:   "concurrency",
					Usage:  "Concurrency",
					EnvVar: "CONCURRENCY",
				},
				cli.StringSliceFlag{
					Name:   "fetch-url",
					Usage:  "Base URL(s) for fetching missed alerts",
					EnvVar: "ALERTS_FETCH_URL",
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
