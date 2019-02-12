package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/alerting/alerts/pkg/alerts"
	"github.com/alerting/alerts/pkg/cap"
	raven "github.com/getsentry/raven-go"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes"
	"github.com/olivere/elastic"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Storage defines an Elasticsearch alerts storage.
type Storage struct {
	Client *elastic.Client
	Index  string
}

// NewStorage creates the storage.
func NewStorage(url string, index string, opts ...elastic.ClientOptionFunc) (alerts.Storage, error) {
	// Generate ElasticSearch client.
	opts = append(opts, elastic.SetURL(url))
	client, err := elastic.NewClient(opts...)
	if err != nil {
		return nil, err
	}

	// Configure the ElasticSearch index.
	log.WithField("index", index).Debug("Ensuring the index exists")
	exists, err := client.IndexExists(index).Do(context.Background())
	if err != nil {
		log.WithError(err).Error("Failed to ensure index exists")
		raven.CaptureErrorAndWait(err, map[string]string{
			"index": index,
		})
		return nil, err
	}

	// Create the index, if it doesn't exist.
	if !exists {
		log.WithField("index", index).Info("Creating index")

		_, err := client.CreateIndex(index).BodyString(mapping).Do(context.Background())
		if err != nil {
			log.WithError(err).Error("Failed to create index")
			raven.CaptureErrorAndWait(err, map[string]string{
				"index": index,
			})
			return nil, err
		}
	}

	// Return storage.
	return &Storage{
		Client: client,
		Index:  index,
	}, nil
}

func generateRangeQuery(ctx context.Context, field string, ts *alerts.TimeConditions) (*elastic.RangeQuery, error) {
	span := opentracing.SpanFromContext(ctx)

	rq := elastic.NewRangeQuery(field)

	if ts.Gte != nil {
		t, err := ptypes.Timestamp(ts.Gte)
		if err != nil {
			return nil, err
		}

		span.SetTag(fmt.Sprintf("query.%s.gte", field), t)
		rq.Gte(t)
	}

	if ts.Gt != nil {
		t, err := ptypes.Timestamp(ts.Gt)
		if err != nil {
			return nil, err
		}

		span.SetTag(fmt.Sprintf("query.%s.te", field), t)
		rq.Gt(t)
	}

	if ts.Lte != nil {
		t, err := ptypes.Timestamp(ts.Lte)
		if err != nil {
			return nil, err
		}

		span.SetTag(fmt.Sprintf("query.%s.lte", field), t)
		rq.Lte(t)
	}

	if ts.Lt != nil {
		t, err := ptypes.Timestamp(ts.Lt)
		if err != nil {
			return nil, err
		}

		span.SetTag(fmt.Sprintf("query.%s.lt", field), t)
		rq.Lt(t)
	}

	return rq, nil
}

// Add adds the alert to ElasticSearch.
func (s *Storage) Add(ctx context.Context, alert *cap.Alert) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "Add")
	defer span.Finish()

	bulk := s.Client.Bulk().Index(s.Index).Type("doc")

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

	// Run bulk actions
	if bulk.NumberOfActions() > 0 {
		_, err := bulk.Do(ctx)
		if err != nil {
			log.WithError(err).Error("Failed to run bulk actions")
			raven.CaptureError(err, nil)
			ext.Error.Set(span, true)
			span.LogFields(otlog.Error(err))
			return err
		}
	}

	return nil
}

// Get returns the alert described by reference from ElasticSearch.
func (s *Storage) Get(ctx context.Context, reference *cap.Reference) (*cap.Alert, error) {
	// Start a new span.
	span, sctx := opentracing.StartSpanFromContext(ctx, "Storage.ElasticSearch::Get")
	defer span.Finish()

	// Save information about the request.
	if reference.Id != "" {
		span.SetTag("reference.id", reference.Id)
	} else {
		span.SetTag("reference.identifier", reference.Identifier)
		span.SetTag("reference.sender", reference.Sender)
		span.SetTag("reference.sent", reference.Sent)
	}

	id := reference.Id
	if id == "" {
		id = reference.ID()
	}

	span.LogEvent("Fetching alert")
	log.WithFields(log.Fields{
		"id":    id,
		"index": s.Index,
	}).Debug("Fetching alert from ElasticSearch")
	item, err := s.Client.Get().Index(s.Index).Type("doc").Id(id).Do(sctx)
	log.Debug("Response received")

	if err != nil {
		log.WithError(err).Error("Unable to fetch alert from ElasticSearch")
		raven.CaptureError(err, map[string]string{
			"id":    id,
			"index": s.Index,
		})
		ext.Error.Set(span, true)
		span.LogFields(otlog.Error(err))

		if elastic.IsNotFound(err) {
			return nil, status.Error(codes.NotFound, "Alert was not found")
		}
		return nil, err
	}

	// If we were sent an invalid ID (ie. of an info block), but matched a document, then return error.
	if item.Routing != "" {
		err := fmt.Errorf("Invalid alert identifier: %s", id)
		log.WithError(err).Error("Unable to fetch alert from ElasticSearch")
		raven.CaptureError(err, nil)
		ext.Error.Set(span, true)
		span.LogFields(otlog.Error(err))
		return nil, err
	}

	// Start fetching the alert.
	span.LogEvent("Processing alert")

	var alert cap.Alert

	dec := &jsonpb.Unmarshaler{
		AllowUnknownFields: true,
	}
	err = dec.Unmarshal(bytes.NewReader(*item.Source), &alert)
	if err != nil {
		log.WithError(err).Error("Unable to unmarshal alert")
		raven.CaptureError(err, nil)
		ext.Error.Set(span, true)
		span.LogFields(otlog.Error(err))
		return nil, err
	}

	// Generate the query to find infos.
	span.LogEvent("Fetching infos")
	log.Debug("Finding infos associated with the alert")
	q := elastic.NewParentIdQuery("info", id)

	search := s.Client.Search().Index(s.Index).Type("doc")
	search = search.Query(q)
	// Realistically, there shouldn't be more than ~4 (Canada NAAD)
	search = search.From(0).Size(20)
	search = search.Sort("_id", true)

	res, err := search.Do(sctx)
	log.Debug("Response received")
	if err != nil {
		log.WithError(err).Error("Failed to find infos")
		raven.CaptureError(err, nil)
		ext.Error.Set(span, true)
		span.LogFields(otlog.Error(err))
		return nil, err
	}

	span.LogEvent("Processing infos")
	if res.Hits != nil {
		if res.TotalHits() != int64(len(res.Hits.Hits)) {
			log.Warnf("Got %d of %d infos", len(res.Hits.Hits), res.TotalHits)
			raven.CaptureMessage("Did not load all infos", map[string]string{
				"id":     id,
				"total":  strconv.FormatInt(res.TotalHits(), 10),
				"loaded": strconv.Itoa(len(res.Hits.Hits)),
			})
		}
		log.Debugf("Got %d infos for alert", res.TotalHits())

		for _, hit := range res.Hits.Hits {
			var info cap.Info
			err = dec.Unmarshal(bytes.NewReader(*hit.Source), &info)
			if err != nil {
				log.WithError(err).Error("Failed to unmarshal info")
				raven.CaptureError(err, nil)
				ext.Error.Set(span, true)
				span.LogFields(otlog.Error(err))
				return nil, err
			}

			alert.Infos = append(alert.Infos, &info)
		}
	} else {
		log.WithField("id", id).Warn("No infos found for alert")
	}

	return &alert, nil
}

// Has returns whether or not an alert exists for the given reference.
func (s *Storage) Has(ctx context.Context, reference *cap.Reference) (bool, error) {
	// Start a new span.
	span, sctx := opentracing.StartSpanFromContext(ctx, "Storage.ElasticSearch::Has")
	defer span.Finish()

	// Save information about the request.
	if reference.Id != "" {
		span.SetTag("reference.id", reference.Id)
	} else {
		span.SetTag("reference.identifier", reference.Identifier)
		span.SetTag("reference.sender", reference.Sender)
		span.SetTag("reference.sent", reference.Sent)
	}

	id := reference.Id
	if id == "" {
		id = reference.ID()
	}

	// Check if the item exists.
	item := elastic.NewMultiGetItem().
		Index(s.Index).
		Type("doc").
		Id(id).
		FetchSource(elastic.NewFetchSourceContext(false))
	res, err := s.Client.MultiGet().Add(item).Do(sctx)
	if err != nil {
		log.WithError(err).Error("Failed to get from ElasticSearch")
		raven.CaptureError(err, nil)
		ext.Error.Set(span, true)
		span.LogFields(otlog.Error(err))
		return false, err
	}

	span.SetTag("result", res.Docs[0].Found)
	return res.Docs[0].Found, nil
}

// Find returns alerts matching the search criteria. NOTE: Results are per Info block.
func (s *Storage) Find(ctx context.Context, criteria *alerts.FindCriteria) (*alerts.FindResult, error) {
	// Start a new span.
	span, sctx := opentracing.StartSpanFromContext(ctx, "Storage.ElasticSearch::Find")
	defer span.Finish()

	search := s.Client.Search(s.Index).Type("doc")

	// Pagination
	if criteria.Start > 0 {
		log.WithField("value", criteria.Start).Debug("start")
		span.SetTag("start", criteria.Start)
		search = search.From(int(criteria.Start))
	}

	if criteria.Count > 0 {
		log.WithField("value", criteria.Count).Debug("count")
		span.SetTag("count", criteria.Count)
		search = search.Size(int(criteria.Count))
	}

	// Sort
	if len(criteria.Sort) == 0 {
		// Sort by effective, descending (if no sort provided)
		criteria.Sort = []string{"-effective"}
	}
	log.WithField("value", criteria.Sort).Debug("sort")
	span.SetTag("sort", criteria.Sort)

	for _, field := range criteria.Sort {
		asc := true
		if strings.HasPrefix(field, "-") {
			field = field[1:]
			asc = false
		}
		search = search.Sort(field, asc)
	}

	// Generate the query
	query := elastic.NewBoolQuery()

	// Alert filters
	alertQuery := elastic.NewBoolQuery()

	// We don't want referenced alerts.
	alertQuery = alertQuery.Must(elastic.NewExistsQuery("status"))

	if criteria.Superseded {
		log.WithField("value", true).Debug("superseded")
		span.SetTag("superseded", true)
		alertQuery = alertQuery.Must(elastic.NewTermQuery("superseded", true))
	}

	if criteria.NotSuperseded {
		log.WithField("value", false).Debug("superseded")
		span.SetTag("superseded", false)
		alertQuery = alertQuery.MustNot(elastic.NewTermQuery("superseded", true))
	}

	if criteria.Status != cap.Alert_STATUS_UNKNOWN {
		log.WithField("value", criteria.Status).Debug("status")
		span.SetTag("status", criteria.Status)
		alertQuery = alertQuery.Must(elastic.NewTermQuery("status", criteria.Status.String()))
	}

	if criteria.MessageType != cap.Alert_MESSAGE_TYPE_UNKNOWN {
		log.WithField("value", criteria.MessageType).Debug("messageType")
		span.SetTag("messageType", criteria.MessageType)
		alertQuery = alertQuery.Must(elastic.NewTermQuery("messageType", criteria.MessageType.String()))
	}

	if criteria.Scope != cap.Alert_SCOPE_UNKNOWN {
		log.WithField("value", criteria.Scope).Debug("scope")
		span.SetTag("scope", criteria.Scope)
		alertQuery = alertQuery.Must(elastic.NewTermQuery("scope", criteria.Scope.String()))
	}

	query = query.Must(elastic.NewHasParentQuery("alert", alertQuery).InnerHit(elastic.NewInnerHit().FetchSource(true)))

	// Info filters
	if criteria.Language != "" {
		log.WithField("value", criteria.Language).Debug("language")
		span.SetTag("language", criteria.Language)

		if strings.ContainsAny(criteria.Language, "*?") {
			query = query.Must(elastic.NewWildcardQuery("language", criteria.Language))
		} else {
			query = query.Must(elastic.NewTermQuery("language", criteria.Language))
		}
	}

	if criteria.Certainty != cap.Info_CERTAINTY_UNKNOWN {
		log.WithField("value", criteria.Certainty).Debug("certainty")
		span.SetTag("certainty", criteria.Certainty)
		query = query.Must(elastic.NewTermQuery("certainty", criteria.Certainty.String()))
	}

	if criteria.Severity != cap.Info_SEVERITY_UNKNOWN {
		log.WithField("value", criteria.Severity).Debug("severity")
		span.SetTag("severity", criteria.Severity)
		query = query.Must(elastic.NewTermQuery("severity", criteria.Severity.String()))
	}

	if criteria.Urgency != cap.Info_URGENCY_UNKNOWN {
		log.WithField("value", criteria.Urgency).Debug("urgency")
		span.SetTag("urgency", criteria.Urgency)
		query = query.Must(elastic.NewTermQuery("urgency", criteria.Urgency.String()))
	}

	if criteria.Effective != nil {
		log.WithField("value", criteria.Effective).Debug("effective")
		rangeQuery, err := generateRangeQuery(sctx, "effective", criteria.Effective)
		if err != nil {
			log.WithError(err).Error("Failed to generate Effective range query")
			raven.CaptureError(err, nil)
			ext.Error.Set(span, true)
			span.LogFields(otlog.Error(err))
			return nil, err
		}

		query = query.Must(rangeQuery)
	}

	if criteria.Expires != nil {
		log.WithField("value", criteria.Expires).Debug("expires")
		rangeQuery, err := generateRangeQuery(sctx, "expires", criteria.Expires)
		if err != nil {
			log.WithError(err).Error("Failed to generate Expires range query")
			raven.CaptureError(err, nil)
			ext.Error.Set(span, true)
			span.LogFields(otlog.Error(err))
			return nil, err
		}

		query = query.Must(rangeQuery)
	}

	if criteria.Onset != nil {
		log.WithField("value", criteria.Expires).Debug("onset")
		rangeQuery, err := generateRangeQuery(sctx, "onset", criteria.Onset)
		if err != nil {
			log.WithError(err).Error("Failed to generate Onset range query")
			raven.CaptureError(err, nil)
			ext.Error.Set(span, true)
			span.LogFields(otlog.Error(err))
			return nil, err
		}

		query = query.Must(rangeQuery)
	}

	if criteria.Headline != "" {
		log.WithField("value", criteria.Headline).Debug("headline")
		span.SetTag("headline", criteria.Headline)
		query = query.Must(elastic.NewQueryStringQuery(criteria.Headline).Field("headline"))
	}

	if criteria.Description != "" {
		log.WithField("value", criteria.Description).Debug("description")
		span.SetTag("description", criteria.Description)
		query = query.Must(elastic.NewQueryStringQuery(criteria.Description).Field("description"))
	}

	if criteria.Instruction != "" {
		log.WithField("value", criteria.Instruction).Debug("instruction")
		span.SetTag("instruction", criteria.Description)
		query = query.Must(elastic.NewQueryStringQuery(criteria.Instruction).Field("instruction"))
	}

	// Info.Area
	areaQuery := elastic.NewBoolQuery()

	if criteria.AreaDescription != "" {
		log.WithField("value", criteria.AreaDescription).Debug("areas.description")
		span.SetTag("areas.description", criteria.AreaDescription)
		areaQuery = areaQuery.Must(elastic.NewQueryStringQuery(criteria.AreaDescription).Field("areas.description"))
	}

	if criteria.Point != nil {
		log.WithField("value", criteria.Point).Debug("point")

		areaQuery = areaQuery.Must(elastic.NewBoolQuery().
			Should(NewGeoShapeQuery("areas.polygons").SetPoint(criteria.Point.Lat, criteria.Point.Lon)).
			Should(NewGeoShapeQuery("areas.circles").SetPoint(criteria.Point.Lat, criteria.Point.Lon)))
	}

	query = query.Must(elastic.NewNestedQuery("areas", areaQuery).InnerHit(elastic.NewInnerHit().FetchSource(false)))

	// Do the search
	search = search.Query(query)
	search = search.FetchSourceContext(elastic.NewFetchSourceContext(true).Exclude("resources.derefUri"))

	results, err := search.Do(sctx)
	if err != nil {
		log.WithError(err).Error("Failed to execute search")
		raven.CaptureError(err, nil)
		ext.Error.Set(span, true)
		span.LogFields(otlog.Error(err))
		return nil, err
	}

	response := alerts.FindResult{
		Total: results.TotalHits(),
		Hits:  make([]*alerts.Hit, 0),
	}

	unmarshaller := &jsonpb.Unmarshaler{
		AllowUnknownFields: true,
	}

	for _, hit := range results.Hits.Hits {
		var alert cap.Alert
		var info cap.Info

		if err := unmarshaller.Unmarshal(bytes.NewReader(*hit.InnerHits["alert"].Hits.Hits[0].Source), &alert); err != nil {
			log.WithError(err).Error("Failed to unmarshal alert")
			raven.CaptureError(err, nil)
			ext.Error.Set(span, true)
			span.LogFields(otlog.Error(err))
			return nil, err
		}

		if err := unmarshaller.Unmarshal(bytes.NewReader(*hit.Source), &info); err != nil {
			log.WithError(err).Error("Failed to unmarshal info")
			raven.CaptureError(err, nil)
			ext.Error.Set(span, true)
			span.LogFields(otlog.Error(err))
			return nil, err
		}

		response.Hits = append(response.Hits, &alerts.Hit{
			Id:    hit.Id,
			Alert: &alert,
			Info:  &info,
		})
	}

	return &response, nil
}

// Supersede marks the reference as superseded.
func (s *Storage) Supersede(ctx context.Context, reference *cap.Reference) error {
	// Start a new span.
	span, sctx := opentracing.StartSpanFromContext(ctx, "Storage.ElasticSearch::Supersede")
	defer span.Finish()

	// Save information about the request.
	if reference.Id != "" {
		span.SetTag("reference.id", reference.Id)
	} else {
		span.SetTag("reference.identifier", reference.Identifier)
		span.SetTag("reference.sender", reference.Sender)
		span.SetTag("reference.sent", reference.Sent)
	}

	id := reference.Id
	if id == "" {
		id = reference.ID()
	}

	log.WithField("id", id).Debug("Superseding alert")

	update := map[string]interface{}{
		"superseded": true,
	}
	_, err := s.Client.Update().Index(s.Index).Type("doc").Id(id).Doc(update).Do(sctx)
	if err != nil {
		log.WithError(err).Error("ElasticSearch update failed")
		raven.CaptureError(err, nil)
		ext.Error.Set(span, true)
		span.LogFields(otlog.Error(err))
		return err
	}
	return nil
}

// IsSuperseded returns whether or not the referenced alert has been superseded.
func (s *Storage) IsSuperseded(ctx context.Context, reference *cap.Reference) (bool, error) {
	// Start a new span.
	span, sctx := opentracing.StartSpanFromContext(ctx, "Storage.ElasticSearch::IsSuperseded")
	defer span.Finish()

	// Save information about the request.
	if reference.Id != "" {
		span.SetTag("reference.id", reference.Id)
	} else {
		span.SetTag("reference.identifier", reference.Identifier)
		span.SetTag("reference.sender", reference.Sender)
		span.SetTag("reference.sent", reference.Sent)
	}

	id := reference.Id
	if id == "" {
		id = reference.ID()
	}

	log.WithField("id", id).Debug("Checking if alert has been superseded")
	search := s.Client.Search().Index(s.Index).Type("doc")
	search = search.Size(1)
	search = search.Query(elastic.NewNestedQuery("references", elastic.NewTermQuery("references.id", id)))
	search = search.FetchSource(false)

	res, err := search.Do(sctx)
	log.Debug("Got response")

	if err != nil {
		log.WithError(err).Error("ElasticSearch query failed")
		raven.CaptureError(err, nil)
		ext.Error.Set(span, true)
		span.LogFields(otlog.Error(err))
		return false, err
	}

	span.SetTag("result", res.TotalHits() > 0)
	return res.TotalHits() > 0, nil
}
