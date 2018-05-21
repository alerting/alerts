package main

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/alerting/alerts/cap/xml"
	"github.com/urfave/cli"
)

type alertSystem interface {
	Fetch(ctx context.Context, reference *capxml.Reference) (*capxml.Alert, error)
}

func createSystem(c *cli.Context) (alertSystem, error) {
	switch c.GlobalString("system") {
	case "naads":
		fetchBaseURLs := make([]*url.URL, 0)

		for _, baseURL := range c.StringSlice("fetch-url") {
			fetchBaseURL, err := url.Parse(baseURL)
			if err != nil {
				return nil, err
			}

			fetchBaseURLs = append(fetchBaseURLs, fetchBaseURL)
		}

		return &naads{
			BaseURLs: fetchBaseURLs,
		}, nil
	}
	return nil, errors.New("Unknown system")
}

type naads struct {
	BaseURLs []*url.URL
}

func clean(str string) string {
	str = strings.Replace(str, "-", "_", -1)
	str = strings.Replace(str, "+", "p", -1)
	str = strings.Replace(str, ":", "_", -1)

	return str
}

func (s *naads) Fetch(ctx context.Context, reference *capxml.Reference) (*capxml.Alert, error) {
	// Generate the URL
	u, err := url.Parse(fmt.Sprintf("%s/%sI%s.xml",
		reference.Sent.Format("2006-01-02"),
		clean(reference.Sent.FormatCAP()),
		clean(reference.Identifier)))

	if err != nil {
		return nil, err
	}

	for _, baseURL := range s.BaseURLs {
		u = baseURL.ResolveReference(u)

		log.Println("Fetching", u.String())
		res, err := http.Get(u.String())
		if err != nil {
			log.Println(err)
			continue
		}
		defer res.Body.Close()

		log.Printf("Got status code %d", res.StatusCode)
		if res.StatusCode != 200 {
			continue
		}

		d := xml.NewDecoder(res.Body)

		var alert capxml.Alert
		if err = d.Decode(&alert); err != nil {
			log.Println(err)
			continue
		}

		return &alert, nil
	}

	return nil, errors.New("Unable to fetch alert")
}
