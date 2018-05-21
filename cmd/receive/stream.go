package main

import (
	"context"
	"encoding/xml"
	"log"
	"net"
	"time"

	"github.com/alerting/alerts/cap/xml"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"
	"github.com/urfave/cli"
)

func receive(ctx context.Context, alert *capxml.Alert) error {
	span := opentracing.StartSpan("stream_receive")
	defer span.Finish()

	span.SetTag("sender", alert.Sender)
	span.SetTag("sent", alert.Sent)
	span.SetTag("identifier", alert.Identifier)
	span.SetTag("id", alert.ID())

	log.Printf("Got alert: %s,%s,%s (%s)", alert.Sender, alert.Sent.FormatCAP(), alert.Identifier, alert.ID())

	if _, err := processAlert(opentracing.ContextWithSpan(ctx, span), machineryServer, alert); err != nil {
		ext.Error.Set(span, true)
		span.LogFields(otlog.Error(err))
		return err
	}

	return nil
}

func stream(c *cli.Context) error {
	if c.NArg() != 1 {
		cli.ShowCommandHelpAndExit(c, c.Command.FullName(), 1)
	}

	// Initialize the task processor
	var err error
	machineryServer, err = createMachineryServer(c)
	if err != nil {
		return err
	}

	// Connect to the streaming service
	conn, err := net.Dial("tcp", c.Args().First())
	if err != nil {
		return err
	}

	log.Printf("Conected to %s\n", c.Args().First())

	// Create the XML decoder
	decoder := xml.NewDecoder(conn)

	for {
		// Set the timeout
		conn.SetReadDeadline(time.Now().Add(time.Duration(time.Duration(c.Int("timeout")) * time.Second)))

		var alert capxml.Alert
		if err = decoder.Decode(&alert); err != nil {
			return err
		}

		if err = receive(context.Background(), &alert); err != nil {
			return err
		}
	}

	return nil
}
