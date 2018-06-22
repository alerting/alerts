package main

import (
	"bytes"
	"context"
	"strings"
	"time"

	"github.com/golang/protobuf/jsonpb"

	"github.com/alerting/alerts/cap"
	"github.com/olivere/elastic"
)

type InfoFinder struct {
	client *elastic.Client
	index  string

	parentId     string
	superseded   *bool
	parentFields map[string]string
	termFields   map[string]string
	textFields   map[string]string
	effective    map[string]time.Time
	expires      map[string]time.Time
	onset        map[string]time.Time
	area         string
	point        *elastic.GeoPoint

	start int
	count int

	sort []string
}

type InfoHit struct {
	Id      string    `json:"id"`
	AlertId string    `json:"alert_id"`
	Info    *cap.Info `json:"info"`
}

type InfoResults struct {
	TotalHits int64      `json:"total_hits"`
	Hits      []*InfoHit `json:"hits"`
}

func NewInfoFinder(client *elastic.Client, index string) *InfoFinder {
	return &InfoFinder{
		client:       client,
		index:        index,
		superseded:   nil,
		parentFields: make(map[string]string),
		termFields:   make(map[string]string),
		textFields:   make(map[string]string),
		effective:    make(map[string]time.Time),
		expires:      make(map[string]time.Time),
		onset:        make(map[string]time.Time),
		start:        -1,
		count:        -1,
		sort:         make([]string, 0),
	}
}

/** FILTERS **/
func (f *InfoFinder) AlertId(id string) *InfoFinder {
	f.parentId = id
	return f
}

func (f *InfoFinder) Superseded(superseded bool) *InfoFinder {
	f.superseded = &superseded
	return f
}

/*
func (f *InfoFinder) Status(status cap.Status) *InfoFinder {
	f.parentFields["status"] = status.String()
	return f
}

func (f *InfoFinder) MessageType(messageType cap.MessageType) *InfoFinder {
	f.parentFields["message_type"] = messageType.String()
	return f
}

func (f *InfoFinder) Scope(scope cap.Scope) *InfoFinder {
	f.parentFields["scope"] = scope.String()
	return f
}
*/

func (f *InfoFinder) Language(language string) *InfoFinder {
	f.termFields["language"] = language
	return f
}

/*
func (f *InfoFinder) Certainty(certainty cap.Certainty) *InfoFinder {
	f.termFields["certainty"] = certainty.String()
	return f
}

func (f *InfoFinder) Severity(severity cap.Severity) *InfoFinder {
	f.termFields["severity"] = severity.String()
	return f
}

func (f *InfoFinder) Urgency(urgency cap.Urgency) *InfoFinder {
	f.termFields["urgency"] = urgency.String()
	return f
}
*/

func (f *InfoFinder) Headline(headline string) *InfoFinder {
	f.textFields["headline"] = headline
	return f
}

func (f *InfoFinder) Description(description string) *InfoFinder {
	f.textFields["description"] = description
	return f
}

func (f *InfoFinder) Instruction(instruction string) *InfoFinder {
	f.textFields["instruction"] = instruction
	return f
}

func (f *InfoFinder) EffectiveGte(t time.Time) *InfoFinder {
	f.effective["gte"] = t
	return f
}

func (f *InfoFinder) EffectiveGt(t time.Time) *InfoFinder {
	f.effective["gt"] = t
	return f
}

func (f *InfoFinder) EffectiveLte(t time.Time) *InfoFinder {
	f.effective["lte"] = t
	return f
}

func (f *InfoFinder) EffectiveLt(t time.Time) *InfoFinder {
	f.effective["lt"] = t
	return f
}

func (f *InfoFinder) ExpiresGte(t time.Time) *InfoFinder {
	f.expires["gte"] = t
	return f
}

func (f *InfoFinder) ExpiresGt(t time.Time) *InfoFinder {
	f.expires["gt"] = t
	return f
}

func (f *InfoFinder) ExpiresLte(t time.Time) *InfoFinder {
	f.expires["lte"] = t
	return f
}

func (f *InfoFinder) ExpiresLt(t time.Time) *InfoFinder {
	f.expires["lt"] = t
	return f
}

func (f *InfoFinder) OnsetGte(t time.Time) *InfoFinder {
	f.onset["gte"] = t
	return f
}

func (f *InfoFinder) OnsetGt(t time.Time) *InfoFinder {
	f.onset["gt"] = t
	return f
}

func (f *InfoFinder) OnsetLte(t time.Time) *InfoFinder {
	f.onset["lte"] = t
	return f
}

func (f *InfoFinder) OnsetLt(t time.Time) *InfoFinder {
	f.onset["lt"] = t
	return f
}

func (f *InfoFinder) Area(area string) *InfoFinder {
	f.area = area
	return f
}

func (f *InfoFinder) Point(lat, lon float64) *InfoFinder {
	f.point = elastic.GeoPointFromLatLon(lat, lon)
	return f
}

/** PAGINATION **/
func (f *InfoFinder) Start(start int) *InfoFinder {
	f.start = start
	return f
}

func (f *InfoFinder) Count(count int) *InfoFinder {
	f.count = count
	return f
}

/** SORTING **/
func (f *InfoFinder) Sort(fields ...string) *InfoFinder {
	f.sort = append(f.sort, fields...)
	return f
}

/** FIND **/
func (f *InfoFinder) Find() (*InfoResults, error) {
	search := f.client.Search().Index(f.index).Type("doc")
	search = f.query(search)
	search = f.pagination(search)
	search = f.sorting(search)

	res, err := search.Do(context.Background())
	if err != nil {
		return nil, err
	}

	// Process results
	results := InfoResults{
		TotalHits: res.Hits.TotalHits,
		Hits:      make([]*InfoHit, 0),
	}

	for _, hit := range res.Hits.Hits {
		var info cap.Info
		if err = (&jsonpb.Unmarshaler{
			AllowUnknownFields: true,
		}).Unmarshal(bytes.NewReader(*hit.Source), &info); err != nil {
			return nil, err
		}

		infoHit := InfoHit{
			Id:      hit.Id,
			AlertId: hit.Routing,
			Info:    &info,
		}

		results.Hits = append(results.Hits, &infoHit)
	}

	return &results, nil
}

func (f *InfoFinder) query(service *elastic.SearchService) *elastic.SearchService {
	q := elastic.NewBoolQuery()

	// Parent filter
	if f.parentId != "" || len(f.parentFields) > 0 || f.superseded != nil {
		if f.parentId != "" {
			q = q.Must(elastic.NewParentIdQuery("info", f.parentId))
		}

		if len(f.parentFields) > 0 || f.superseded != nil {
			pq := elastic.NewBoolQuery()

			if f.superseded != nil {
				if *f.superseded {
					pq = pq.Must(elastic.NewTermQuery("superseded", true))
				} else {
					pq = pq.MustNot(elastic.NewTermQuery("superseded", true))
				}
				for k, v := range f.parentFields {
					pq = pq.Must(elastic.NewTermQuery(k, v))
				}
			}

			q = q.Must(elastic.NewHasParentQuery("alert", pq))
		}

		if f.superseded != nil {

		}
	} else {
		q = q.Must(elastic.NewHasParentQuery("alert", elastic.NewMatchAllQuery()))
	}

	// Filter on termFields
	if len(f.termFields) > 0 {
		for k, v := range f.termFields {
			q = q.Must(elastic.NewTermQuery(k, v))
		}
	}

	// Filter on textFields
	if len(f.textFields) > 0 {
		for k, v := range f.textFields {
			q = q.Must(elastic.NewQueryStringQuery(v).Field(k))
		}
	}

	// Filter on times
	if len(f.effective) > 0 {
		rq := elastic.NewRangeQuery("effective")

		if val, ok := f.effective["gte"]; ok {
			rq.Gte(val)
		}

		if val, ok := f.effective["gt"]; ok {
			rq.Gt(val)
		}

		if val, ok := f.effective["lte"]; ok {
			rq.Lte(val)
		}

		if val, ok := f.effective["lt"]; ok {
			rq.Lt(val)
		}

		q = q.Must(rq)
	}

	if len(f.expires) > 0 {
		rq := elastic.NewRangeQuery("expires")

		if val, ok := f.expires["gte"]; ok {
			rq.Gte(val)
		}

		if val, ok := f.expires["gt"]; ok {
			rq.Gt(val)
		}

		if val, ok := f.expires["lte"]; ok {
			rq.Lte(val)
		}

		if val, ok := f.expires["lt"]; ok {
			rq.Lt(val)
		}

		q = q.Must(rq)
	}

	if len(f.onset) > 0 {
		rq := elastic.NewRangeQuery("onset")

		if val, ok := f.onset["gte"]; ok {
			rq.Gte(val)
		}

		if val, ok := f.onset["gt"]; ok {
			rq.Gt(val)
		}

		if val, ok := f.onset["lte"]; ok {
			rq.Lte(val)
		}

		if val, ok := f.onset["lt"]; ok {
			rq.Lt(val)
		}

		q = q.Must(rq)
	}

	// Filter on area
	if f.area != "" || f.point != nil {
		aq := elastic.NewBoolQuery()

		if f.area != "" {
			aq = aq.Must(elastic.NewQueryStringQuery(f.area).Field("areas.description"))
		}

		if f.point != nil {
			pq := elastic.NewBoolQuery()
			pq = pq.Should(NewGeoShapeQuery("areas.polygons").SetPoint(f.point.Lat, f.point.Lon))
			pq = pq.Should(NewGeoShapeQuery("areas.circles").SetPoint(f.point.Lat, f.point.Lon))

			aq = aq.Must(pq)
		}

		nq := elastic.NewNestedQuery("areas", aq)
		nq.InnerHit(elastic.NewInnerHit().FetchSourceContext(elastic.NewFetchSourceContext(false)))

		q = q.Must(nq)
	}

	service = service.Query(q)
	return service
}

func (f *InfoFinder) pagination(service *elastic.SearchService) *elastic.SearchService {
	if f.start >= 0 {
		service = service.From(f.start)
	}

	if f.count >= 0 {
		service = service.Size(f.count)
	}

	return service
}

func (f *InfoFinder) sorting(service *elastic.SearchService) *elastic.SearchService {
	if len(f.sort) == 0 {
		service = service.Sort("_score", false)
		return service
	}

	// Prefix of "-" means to sort descending.
	for _, field := range f.sort {
		asc := true
		if strings.HasPrefix(field, "-") {
			field = field[1:]
			asc = false
		}

		service = service.Sort(field, asc)
	}

	return service
}
