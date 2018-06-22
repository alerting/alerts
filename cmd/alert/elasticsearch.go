package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/golang/protobuf/ptypes"

	"github.com/alerting/alerts/cap"
	"github.com/alerting/alerts/rpc"
	"github.com/golang/protobuf/jsonpb"
	"github.com/olivere/elastic"
	opentracing "github.com/opentracing/opentracing-go"
)

var mapping = `{
    "settings": {
      "analysis": {
        "analyzer": {
          "folding": {
            "tokenizer": "standard",
            "filter": ["lowercase", "asciifolding"]
          }
        },
        "normalizer": {
          "keyword_normalizer": {
            "type": "custom",
            "filter": ["lowercase", "asciifolding"]
          }
        }
      }
    },
    "mappings": {
      "doc": {
        "dynamic": false,
        "properties": {
          "identifier": { "type": "keyword", "normalizer": "keyword_normalizer" },
          "sender": { "type": "keyword", "normalizer": "keyword_normalizer" },
          "sent": { "type": "date" },
          "status": { "type": "keyword", "normalizer": "keyword_normalizer" },
          "message_type": { "type": "keyword", "normalizer": "keyword_normalizer" },
          "source": { "type": "keyword", "normalizer": "keyword_normalizer" },
          "scope": { "type": "keyword", "normalizer": "keyword_normalizer" },
          "restriction": { "type": "text" },
          "addresses": { "type": "keyword", "normalizer": "keyword_normalizer" },
          "codes": { "type": "keyword", "normalizer": "keyword_normalizer" },
          "note": { "type": "text", "analyzer": "folding" },
          "references": {
            "type": "nested",
            "dynamic": false,
            "properties": {
              "sender": { "type": "keyword", "normalizer": "keyword_normalizer" },
              "sent": { "type": "date" },
				  "indentifier": { "type": "keyword", "normalizer": "keyword_normalizer" },
				  "id": { "type": "keyword", "normalizer": "keyword_normalizer" }
            }
          },
          "incidents": { "type": "keyword", "normalizer": "keyword_normalizer" },
          "superseded": { "type": "boolean" },

          "language": { "type": "keyword", "normalizer": "keyword_normalizer" },
          "categories": { "type": "keyword", "normalizer": "keyword_normalizer" },
          "event": { "type": "keyword", "normalizer": "keyword_normalizer" },
          "response_types": { "type": "keyword", "normalizer": "keyword_normalizer" },
          "urgency": { "type": "keyword", "normalizer": "keyword_normalizer" },
          "severity": { "type": "keyword", "normalizer": "keyword_normalizer" },
          "certainty": { "type": "keyword", "normalizer": "keyword_normalizer" },
          "audience": { "type": "keyword", "normalizer": "keyword_normalizer" },
          "event_codes": { "type": "object" },
          "effective": { "type": "date" },
          "onset": { "type": "date" },
          "expires": { "type": "date" },
          "sender_name": { "type": "keyword", "normalizer": "keyword_normalizer" },
          "headline": { "type": "text", "analyzer": "folding" },
          "description": { "type": "text", "analyzer": "folding" },
          "instruction": { "type": "text", "analyzer": "folding" },
          "web": { "type": "keyword", "normalizer": "keyword_normalizer" },
          "contact": { "type": "keyword", "normalizer": "keyword_normalizer" },
          "parameters": { "type": "object" },
          "resources": {
            "type": "nested",
            "dynamic": false,
            "properties": {
              "description": { "type": "text", "analyzer": "folding" },
              "mime_type": { "type": "keyword", "normalizer": "keyword_normalizer" },
              "size": { "type": "integer" },
              "uri": { "type": "keyword", "normalizer": "keyword_normalizer" },
              "derefUri": { "type": "binary" },
              "digest": { "type": "keyword", "normalizer": "keyword_normalizer" }
            }
          },
          "areas": {
            "type": "nested",
            "dynamic": false,
            "properties": {
              "description": { "type": "text", "analyzer": "folding" },
              "polygons": { "type": "geo_shape", "ignore_malformed": true },
              "circles": { "type": "geo_shape", "ignore_malformed": true },
              "geocodes": { "type": "object" },
              "altitude": { "type": "float" },
              "ceiling": { "type": "float" }
            }
          },

          "_object": {
            "type": "join",
            "relations": {
              "alert": "info"
            }
          }
        }
      }
    }
  }`

// AlertElasticsearchServiceConfig represents the configuration for the alert service.
type AlertElasticsearchServiceConfig struct {
	// Tracing
	tracer *opentracing.Tracer

	// Elasticsearch
	Client *elastic.Client
	Index  string
}

type AlertElasticsearchService struct {
	config *AlertElasticsearchServiceConfig
}

func NewElasticsearchAlertService(ctx context.Context, config *AlertElasticsearchServiceConfig) (AlertService, error) {
	if config.Client == nil {
		return nil, errors.New("No Elasticsearch client provided")
	}

	if config.Index == "" {
		return nil, errors.New("No Elasticsearch index provided")
	}

	service := AlertElasticsearchService{
		config: config,
	}

	service.setup(ctx)
	return &service, nil
}

func (service *AlertElasticsearchService) setup(ctx context.Context) error {
	exists, err := service.config.Client.IndexExists(service.config.Index).Do(ctx)
	if err != nil {
		return err
	}

	if !exists {
		_, err := service.config.Client.CreateIndex(service.config.Index).BodyString(mapping).Do(ctx)
		return err
	}

	return nil
}

func (service *AlertElasticsearchService) Add(ctx context.Context, alerts ...*cap.Alert) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "Add")
	defer span.Finish()

	bulk := service.config.Client.Bulk().Index(service.config.Index).Type("doc")

	for _, alert := range alerts {
		span.LogEventWithPayload("alert", fmt.Sprintf("%s %s,%s,%s", alert.Id, alert.Identifier, alert.Sender, alert.Sent))

		// Conver the alert to a map
		// Convert to map[string]interface{}
		var alertMap map[string]interface{}
		b, _ := (&jsonpb.Marshaler{}).MarshalToString(alert)
		json.Unmarshal([]byte(b), &alertMap)

		// Remove the infos item, as we need to add those seperately
		delete(alertMap, "infos")

		// Setup the child-parent relationship
		alertMap["_object"] = map[string]string{
			"name": "alert",
		}

		// Index the alert
		bulk.Add(
			elastic.NewBulkIndexRequest().
				Id(alert.Id).
				Doc(alertMap))

		// Index the infos
		for indx, info := range alert.Infos {
			// Some cleanup
			if info.Effective == nil {
				info.Effective = alert.Sent
			}

			var infoMap map[string]interface{}
			b, _ := (&jsonpb.Marshaler{}).MarshalToString(info)
			json.Unmarshal([]byte(b), &infoMap)

			// Setup Parent
			infoMap["_object"] = map[string]string{
				"name":   "info",
				"parent": alert.Id,
			}

			bulk.Add(
				elastic.NewBulkIndexRequest().
					Id(fmt.Sprintf("%s:%d", alert.Id, indx)).
					Routing(alert.Id).
					Doc(infoMap))
		}

		if bulk.NumberOfActions() > 100 {
			_, err := bulk.Do(ctx)
			if err != nil {
				logError(span, err)
				return err
			}
		}
	}

	// Cleanup any remaining bulk actions
	if bulk.NumberOfActions() > 0 {
		_, err := bulk.Do(ctx)
		if err != nil {
			logError(span, err)
			return err
		}
	}

	return nil
}

func (service *AlertElasticsearchService) Has(ctx context.Context, request *rpc.AlertRequest) (bool, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "Has")
	defer span.Finish()

	span.SetTag("alert.id", request.Id)

	item := elastic.NewMultiGetItem().Index(service.config.Index).Type("doc").Id(request.Id)
	res, err := service.config.Client.MultiGet().Add(item).Do(ctx)
	if err != nil {
		logError(span, err)
		return false, err
	}

	span.SetTag("alert.exists", res.Docs[0].Found)
	return res.Docs[0].Found, nil
}

func (service *AlertElasticsearchService) Get(ctx context.Context, request *rpc.AlertRequest) (*cap.Alert, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "Get")
	defer span.Finish()

	span.SetTag("alert.id", request.Id)

	span.LogEvent("Get alert")
	item, err := service.config.Client.Get().
		Index(service.config.Index).
		Type("doc").
		Id(request.Id).
		Do(ctx)

	if err != nil {
		logError(span, err)
		return nil, err
	}

	// If we got the id of an info block, return not found
	if item.Routing != "" {
		err := errors.New("Unable to find alert with id: " + request.Id)
		logError(span, err)
		return nil, err
	}

	// Fetch the alert itself
	var alert cap.Alert
	err = (&jsonpb.Unmarshaler{
		AllowUnknownFields: true,
	}).Unmarshal(bytes.NewReader(*item.Source), &alert)
	if err != nil {
		logError(span, err)
		return nil, err
	}

	span.LogEvent("Get infos")

	// Fetch the children (ie. infos)
	finder := NewInfoFinder(service.config.Client, service.config.Index)
	finder = finder.AlertId(request.Id)
	finder = finder.Sort("_id")

	infos, err := finder.Find()
	if err != nil {
		logError(span, err)
		return nil, err
	}

	for _, hit := range infos.Hits {
		alert.Infos = append(alert.Infos, hit.Info)
	}

	return &alert, nil
}

func generateRangeQuery(ctx context.Context, field string, ts *rpc.TimeQuery) (*elastic.RangeQuery, error) {
	span := opentracing.SpanFromContext(ctx)

	rq := elastic.NewRangeQuery(field)

	if ts.Gte != nil {
		t, err := ptypes.Timestamp(ts.Gte)
		if err != nil {
			logError(span, err)
			return nil, err
		}

		span.SetTag(fmt.Sprintf("query.%s.gte", field), t)
		rq.Gte(t)
	}

	if ts.Gt != nil {
		t, err := ptypes.Timestamp(ts.Gt)
		if err != nil {
			logError(span, err)
			return nil, err
		}

		span.SetTag(fmt.Sprintf("query.%s.te", field), t)
		rq.Gt(t)
	}

	if ts.Lte != nil {
		t, err := ptypes.Timestamp(ts.Lte)
		if err != nil {
			logError(span, err)
			return nil, err
		}

		span.SetTag(fmt.Sprintf("query.%s.lte", field), t)
		rq.Lte(t)
	}

	if ts.Lt != nil {
		t, err := ptypes.Timestamp(ts.Lt)
		if err != nil {
			logError(span, err)
			return nil, err
		}

		span.SetTag(fmt.Sprintf("query.%s.lt", field), t)
		rq.Lt(t)
	}

	return rq, nil
}

func (service *AlertElasticsearchService) Find(ctx context.Context, r *rpc.AlertsRequest) (*rpc.AlertsResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "find")
	defer span.Finish()

	rctx := opentracing.ContextWithSpan(ctx, span)

	search := service.config.Client.Search().Index(service.config.Index).Type("doc")

	// Pagination
	if r.Start > 0 {
		span.SetTag("start", r.Start)
		search = search.From(int(r.Start))
	}

	if r.Count > 0 {
		span.SetTag("count", r.Count)
		search = search.Size(int(r.Count))
	}

	// Sort
	if len(r.Sort) == 0 {
		// By default, sort by effective descending
		r.Sort = []string{"-effective"}
	}

	span.SetTag("sort", r.Sort)

	for _, field := range r.Sort {
		asc := true
		if strings.HasPrefix(field, "-") {
			field = field[1:]
			asc = false
		}
		search = search.Sort(field, asc)
	}

	// Generate the search query
	q := elastic.NewBoolQuery()

	// Alert filters
	pq := elastic.NewBoolQuery()
	pq = pq.Must(elastic.NewExistsQuery("status")) // Don't match on referenced alerts

	if r.Superseded {
		span.SetTag("query.superseded", true)
		pq = pq.Must(elastic.NewTermQuery("superseded", true))
	}

	if r.NotSuperseded {
		span.SetTag("query.superseded", false)
		pq = pq.MustNot(elastic.NewTermQuery("superseded", true))
	}

	if r.Status != cap.Alert_STATUS_UNKNOWN {
		span.SetTag("query.status", r.Status)
		pq = pq.Must(elastic.NewTermQuery("status", r.Status.String()))
	}

	if r.MessageType != cap.Alert_MESSAGE_TYPE_UNKNOWN {
		span.SetTag("query.message_type", r.MessageType)
		pq = pq.Must(elastic.NewTermQuery("message_type", r.MessageType.String()))
	}

	if r.Scope != cap.Alert_SCOPE_UNKNOWN {
		span.SetTag("query.scope", r.Scope)
		pq = pq.Must(elastic.NewTermQuery("scope", r.Scope.String()))
	}

	q = q.Must(elastic.NewHasParentQuery("alert", pq).
		InnerHit(elastic.NewInnerHit().FetchSource(true)))

	// Info filters
	if r.Language != "" {
		span.SetTag("query.language", r.Language)
		q = q.Must(elastic.NewTermQuery("language", r.Language))
	}

	if r.Certainty != cap.Info_CERTAINTY_UNKNOWN {
		span.SetTag("query.certainty", r.Certainty)
		q = q.Must(elastic.NewTermQuery("certainty", r.Certainty.String()))
	}

	if r.Severity != cap.Info_SEVERITY_UNKNOWN {
		span.SetTag("query.severity", r.Severity)
		q = q.Must(elastic.NewTermQuery("severity", r.Severity.String()))
	}

	if r.Urgency != cap.Info_URGENCY_UNKNOWN {
		span.SetTag("query.urgency", r.Urgency)
		q = q.Must(elastic.NewTermQuery("urgency", r.Urgency.String()))
	}

	if r.Headline != "" {
		span.SetTag("query.headline", r.Headline)
		q = q.Must(elastic.NewQueryStringQuery(r.Headline).Field("headline"))
	}

	if r.Effective != nil {
		rq, err := generateRangeQuery(rctx, "effective", r.Effective)
		if err != nil {
			logError(span, err)
			return nil, err
		}

		q = q.Must(rq)
	}

	if r.Onset != nil {
		rq, err := generateRangeQuery(rctx, "onset", r.Onset)
		if err != nil {
			logError(span, err)
			return nil, err
		}

		q = q.Must(rq)
	}

	if r.Expires != nil {
		rq, err := generateRangeQuery(rctx, "expires", r.Expires)
		if err != nil {
			logError(span, err)
			return nil, err
		}

		q = q.Must(rq)
	}

	if r.Description != "" {
		span.SetTag("query.description", r.Description)
		q = q.Must(elastic.NewQueryStringQuery(r.Description).Field("description"))
	}

	if r.Instruction != "" {
		span.SetTag("query.instruction", r.Instruction)
		q = q.Must(elastic.NewQueryStringQuery(r.Instruction).Field("instruction"))
	}

	// info.area
	aq := elastic.NewBoolQuery()

	if r.AreaDescription != "" {
		span.SetTag("query.area.description", r.AreaDescription)
		aq = aq.Must(elastic.NewQueryStringQuery(r.AreaDescription).Field("areas.description"))
	}

	if r.Point != nil {
		span.SetTag("query.area.point", r.Point)

		pq := elastic.NewBoolQuery()
		pq = pq.Should(NewGeoShapeQuery("areas.polygons").SetPoint(r.Point.Lat, r.Point.Lon))
		pq = pq.Should(NewGeoShapeQuery("areas.circles").SetPoint(r.Point.Lat, r.Point.Lon))

		aq = aq.Must(pq)
	}

	q = q.Must(elastic.
		NewNestedQuery("areas", aq).
		InnerHit(elastic.
			NewInnerHit().
			FetchSourceContext(elastic.
				NewFetchSourceContext(false))))

	// Do the search
	search = search.Query(q)
	search = search.FetchSourceContext(elastic.NewFetchSourceContext(true).Exclude("resources.derefUri"))

	results, err := search.Do(ctx)
	if err != nil {
		return nil, err
	}

	response := rpc.AlertsResponse{
		Total: results.TotalHits(),
		Hits:  make([]*rpc.AlertHit, 0),
	}

	um := &jsonpb.Unmarshaler{
		AllowUnknownFields: true,
	}

	// Get alerts
	for _, hit := range results.Hits.Hits {
		var info cap.Info
		var alert cap.Alert

		if err := um.Unmarshal(bytes.NewReader(*hit.Source), &info); err != nil {
			return nil, err
		}

		if err := um.Unmarshal(bytes.NewReader(*hit.InnerHits["alert"].Hits.Hits[0].Source), &alert); err != nil {
			return nil, err
		}

		response.Hits = append(response.Hits, &rpc.AlertHit{
			Id:    hit.Id,
			Info:  &info,
			Alert: &alert,
		})
	}

	return &response, nil
}

func (service *AlertElasticsearchService) Supersede(ctx context.Context, request *rpc.AlertRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "Has")
	defer span.Finish()

	span.SetTag("alert.id", request.Id)

	update := map[string]interface{}{
		"superseded": true,
	}
	_, err := service.config.Client.Update().Index(service.config.Index).Type("doc").Id(request.Id).Doc(update).Do(ctx)
	if err != nil {
		logError(span, err)
		return err
	}
	return nil
}

func (service *AlertElasticsearchService) IsSuperseded(ctx context.Context, request *rpc.AlertRequest) (bool, error) {
	search := service.config.Client.Search().Index(service.config.Index).Type("doc")
	search = search.Size(1)
	search = search.Query(elastic.NewNestedQuery("references", elastic.NewTermQuery("references.id", request.Id)))
	search = search.FetchSource(false)

	results, err := search.Do(ctx)
	if err != nil {
		return false, err
	}

	return results.TotalHits() > 0, nil
}
