package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc/credentials"

	"github.com/alerting/alerts/cap"
	"github.com/alerting/alerts/rpc"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
)

var grpcConn *grpc.ClientConn

func logError(span opentracing.Span, err error) {
	ext.Error.Set(span, true)
	span.LogFields(otlog.Error(err))
	log.Print(err)
}

func getTimeQuery(query *url.Values, prefix string) (*rpc.TimeQuery, error) {
	timeQuery := rpc.TimeQuery{}
	if val := query.Get(prefix + "_gte"); val != "" {
		t, err := time.Parse(time.RFC3339, val)
		if err != nil {
			return nil, err
		}
		timeQuery.Gte, err = ptypes.TimestampProto(t)
		if err != nil {
			return nil, err
		}
	}
	if val := query.Get(prefix + "_gt"); val != "" {
		t, err := time.Parse(time.RFC3339, val)
		if err != nil {
			return nil, err
		}
		timeQuery.Gt, err = ptypes.TimestampProto(t)
		if err != nil {
			return nil, err
		}
	}
	if val := query.Get(prefix + "_lte"); val != "" {
		t, err := time.Parse(time.RFC3339, val)
		if err != nil {
			return nil, err
		}
		timeQuery.Lte, err = ptypes.TimestampProto(t)
		if err != nil {
			return nil, err
		}
	}
	if val := query.Get(prefix + "_lt"); val != "" {
		t, err := time.Parse(time.RFC3339, val)
		if err != nil {
			return nil, err
		}
		timeQuery.Lt, err = ptypes.TimestampProto(t)
		if err != nil {
			return nil, err
		}
	}

	if timeQuery.Gte != nil || timeQuery.Gt != nil || timeQuery.Lte != nil || timeQuery.Lt != nil {
		return &timeQuery, nil
	}

	return nil, nil
}

func getAlerts(w http.ResponseWriter, r *http.Request) {
	span := opentracing.SpanFromContext(r.Context())
	ctx := opentracing.ContextWithSpan(context.Background(), span)

	c := rpc.NewAlertServiceClient(grpcConn)
	req := &rpc.AlertsRequest{}

	// Assemble the search request
	query := r.URL.Query()

	if val := query.Get("start"); val != "" {
		ival, err := strconv.ParseInt(val, 10, 32)
		if err != nil {
			logError(span, err)
			w.WriteHeader(400)
			w.Write([]byte("Invalid start value"))
			return
		}

		req.Start = int32(ival)
	}

	if val := query.Get("count"); val != "" {
		ival, err := strconv.ParseInt(val, 10, 32)
		if err != nil {
			logError(span, err)
			w.WriteHeader(400)
			w.Write([]byte("Invalid count value"))
			return
		}
		req.Count = int32(ival)
	}

	if val := query.Get("sort"); val != "" {
		req.Sort = strings.Split(val, ",")
	}

	if val := query.Get("superseded"); val != "" {
		b, err := strconv.ParseBool(val)
		if err != nil {
			logError(span, err)
			w.WriteHeader(400)
			w.Write([]byte("Invalid superseded value"))
			return
		}

		if b {
			req.Superseded = true
		} else {
			req.NotSuperseded = true
		}
	}

	if val := query.Get("status"); val != "" {
		if status, ok := cap.Alert_Status_value[strings.ToUpper(val)]; ok {
			req.Status = cap.Alert_Status(status)
		} else {
			w.WriteHeader(400)
			w.Write([]byte("Invalid status value"))
			return
		}
	}

	if val := query.Get("messageType"); val != "" {
		if msgType, ok := cap.Alert_MessageType_value[strings.ToUpper(val)]; ok {
			req.MessageType = cap.Alert_MessageType(msgType)
		} else {
			w.WriteHeader(400)
			w.Write([]byte("Invalid messageType value"))
			return
		}
	}

	if val := query.Get("scope"); val != "" {
		if scope, ok := cap.Alert_Scope_value[strings.ToUpper(val)]; ok {
			req.Scope = cap.Alert_Scope(scope)
		} else {
			w.WriteHeader(400)
			w.Write([]byte("Invalid scope value"))
			return
		}
	}

	if val := query.Get("language"); val != "" {
		req.Language = val
	}

	if val := query.Get("certainty"); val != "" {
		if certainty, ok := cap.Info_Certainty_value[strings.ToUpper(val)]; ok {
			req.Certainty = cap.Info_Certainty(certainty)
		} else {
			w.WriteHeader(400)
			w.Write([]byte("Invalid certainty value"))
			return
		}
	}

	if val := query.Get("severity"); val != "" {
		if severity, ok := cap.Info_Severity_value[strings.ToUpper(val)]; ok {
			req.Severity = cap.Info_Severity(severity)
		} else {
			w.WriteHeader(400)
			w.Write([]byte("Invalid severity value"))
			return
		}
	}

	if val := query.Get("urgency"); val != "" {
		if urgency, ok := cap.Info_Urgency_value[strings.ToUpper(val)]; ok {
			req.Urgency = cap.Info_Urgency(urgency)
		} else {
			w.WriteHeader(400)
			w.Write([]byte("Invalid urgency value"))
			return
		}
	}

	var err error
	req.Effective, err = getTimeQuery(&query, "effective")
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Invalid effective value"))
		return
	}

	req.Onset, err = getTimeQuery(&query, "onset")
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Invalid onset value"))
		return
	}

	req.Expires, err = getTimeQuery(&query, "expires")
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Invalid expires value"))
		return
	}

	if val := query.Get("headline"); val != "" {
		req.Headline = val
	}

	if val := query.Get("description"); val != "" {
		req.Description = val
	}

	if val := query.Get("instruction"); val != "" {
		req.Instruction = val
	}

	if val := query.Get("area_description"); val != "" {
		req.AreaDescription = val
	}

	if val := query.Get("point"); val != "" {
		components := strings.Split(val, ",")
		if len(components) != 2 {
			w.WriteHeader(400)
			w.Write([]byte("Invalid point value"))
			return
		}

		lat, err := strconv.ParseFloat(components[0], 64)
		if err != nil {
			logError(span, err)
			w.WriteHeader(400)
			w.Write([]byte("Invalid point.lat value"))
			return
		}

		lon, err := strconv.ParseFloat(components[1], 64)
		if err != nil {
			logError(span, err)
			w.WriteHeader(400)
			w.Write([]byte("Invalid point.lon value"))
			return
		}

		req.Point = &rpc.Coordinate{
			Lat: lat,
			Lon: lon,
		}
	}

	// Do the find
	alerts, err := c.FindAlerts(ctx, req)
	if err != nil {
		logError(span, err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	enc := jsonpb.Marshaler{
		EmitDefaults: alerts.Total == 0,
	}
	if err := enc.Marshal(w, alerts); err != nil {
		logError(span, err)
	}
}

func getAlert(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	span := opentracing.SpanFromContext(r.Context())
	ctx := opentracing.ContextWithSpan(context.Background(), span)

	c := rpc.NewAlertServiceClient(grpcConn)

	alert, err := c.GetAlert(ctx, &rpc.AlertRequest{
		Id: vars["id"],
	})
	if err != nil {
		logError(span, err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	enc := jsonpb.Marshaler{}
	if err := enc.Marshal(w, alert); err != nil {
		logError(span, err)
	}
}

func serve(c *cli.Context) error {
	var err error

	// gRPC
	var security grpc.DialOption

	if !c.GlobalBool("alerts-insecure") {
		rootCAs, err := x509.SystemCertPool()
		if err != nil {
			return err
		}

		if c.GlobalIsSet("alerts-ca-cert") {
			rootCAs = x509.NewCertPool()

			certs, err := ioutil.ReadFile(c.GlobalString("alerts-ca-cert"))
			if err != nil {
				return err
			}

			ok := rootCAs.AppendCertsFromPEM(certs)
			if !ok {
				log.Printf("No certificates imported from %s", c.GlobalString("alerts-ca-cert"))
			}
		}

		security = grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			RootCAs:            rootCAs,
			InsecureSkipVerify: false,
		}))
	} else {
		security = grpc.WithInsecure()
	}

	grpcConn, err = grpc.Dial(c.GlobalString("alerts-host"),
		grpc.WithMaxMsgSize(1024*1024*1024),
		security,
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
		grpc.WithStreamInterceptor(otgrpc.OpenTracingStreamClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		return err
	}
	defer grpcConn.Close()

	// Config the api
	mux := mux.NewRouter()
	mux.HandleFunc("/alerts", getAlerts).Methods("GET")
	mux.HandleFunc("/alerts/{id}", getAlert).Methods("GET")

	log.Printf("Listening on %s", c.String("listen"))

	handler := nethttp.Middleware(opentracing.GlobalTracer(), mux)
	if c.IsSet("cert") || c.IsSet("key") {
		return http.ListenAndServeTLS(c.String("listen"), c.String("cert"), c.String("key"), handler)
	} else {
		return http.ListenAndServe(c.String("listen"), handler)
	}
}

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

	// App
	app := cli.NewApp()

	app.Name = "alert"
	app.Usage = "API Gateway"
	app.Authors = []cli.Author{
		{
			Name:  "Zachary Seguin",
			Email: "zachary@zacharyseguin.ca",
		},
	}
	app.Copyright = "Copyright (c) 2018 Zachary Seguin"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "alerts-host",
			Usage:  "host:port for alerts service",
			EnvVar: "ALERTS_HOST",
			Value:  "localhost:2400",
		},
		cli.BoolFlag{
			Name:   "alerts-insecure",
			Usage:  "Use insecure grpc connection (no TLS)",
			EnvVar: "ALERTS_INSECURE",
		},
		cli.StringFlag{
			Name:   "alerts-ca-cert",
			Usage:  "CA Certificate for the alerts service (default uses system certs)",
			EnvVar: "ALERTS_CA_CERT",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:   "serve",
			Usage:  "Serves the API gateway service",
			Action: serve,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "listen",
					Usage:  "host:port to listen on",
					EnvVar: "LISTEN",
					Value:  ":8080",
				},
				cli.StringFlag{
					Name:   "cert",
					Usage:  "SSL certificate",
					EnvVar: "CERT",
				},
				cli.StringFlag{
					Name:   "key",
					Usage:  "SSL key",
					EnvVar: "KEY",
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
