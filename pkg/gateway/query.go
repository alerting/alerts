package gateway

import (
	"context"
	"errors"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/alerting/alerts/pkg/alerts"
	"github.com/alerting/alerts/pkg/cap"
	raven "github.com/getsentry/raven-go"
	"github.com/golang/protobuf/ptypes"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"
	log "github.com/sirupsen/logrus"
)

func getTimeQuery(ctx context.Context, query url.Values, prefix string) (*alerts.TimeConditions, error) {
	timeQuery := alerts.TimeConditions{}
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

func getAlertsQuery(ctx context.Context, query url.Values) (*alerts.FindCriteria, error) {
	span := opentracing.SpanFromContext(ctx)

	req := &alerts.FindCriteria{}

	if val := query.Get("start"); val != "" {
		ival, err := strconv.ParseInt(val, 10, 32)
		if err != nil {
			log.WithError(err).Error("Failed to get start")
			raven.CaptureError(err, nil)
			ext.Error.Set(span, true)
			span.LogFields(otlog.Error(err))
			return nil, err
		}

		req.Start = int32(ival)
	}

	if val := query.Get("count"); val != "" {
		ival, err := strconv.ParseInt(val, 10, 32)
		if err != nil {
			log.WithError(err).Error("Failed to get count")
			raven.CaptureError(err, nil)
			ext.Error.Set(span, true)
			span.LogFields(otlog.Error(err))
			return nil, err
		}
		req.Count = int32(ival)
	}

	if val := query.Get("sort"); val != "" {
		req.Sort = strings.Split(val, ",")
	}

	if val := query.Get("fields"); val != "" {
		req.Fields = strings.Split(val, ",")
	}

	if val := query.Get("superseded"); val != "" {
		b, err := strconv.ParseBool(val)
		if err != nil {
			log.WithError(err).Error("Failed to get superseded")
			raven.CaptureError(err, nil)
			ext.Error.Set(span, true)
			span.LogFields(otlog.Error(err))
			return nil, err
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
			return nil, errors.New("Invalid status value")
		}
	}

	if val := query.Get("messageType"); val != "" {
		if msgType, ok := cap.Alert_MessageType_value[strings.ToUpper(val)]; ok {
			req.MessageType = cap.Alert_MessageType(msgType)
		} else {
			return nil, errors.New("Invalid messageType value")
		}
	}

	if val := query.Get("scope"); val != "" {
		if scope, ok := cap.Alert_Scope_value[strings.ToUpper(val)]; ok {
			req.Scope = cap.Alert_Scope(scope)
		} else {
			return nil, errors.New("Invalid scope value")
		}
	}

	if val := query.Get("language"); val != "" {
		req.Language = val
	}

	if val := query.Get("certainty"); val != "" {
		if certainty, ok := cap.Info_Certainty_value[strings.ToUpper(val)]; ok {
			req.Certainty = cap.Info_Certainty(certainty)
		} else {
			return nil, errors.New("Invalid certainty value")
		}
	}

	if val := query.Get("severity"); val != "" {
		if severity, ok := cap.Info_Severity_value[strings.ToUpper(val)]; ok {
			req.Severity = cap.Info_Severity(severity)
		} else {
			return nil, errors.New("Invalid severity value")
		}
	}

	if val := query.Get("urgency"); val != "" {
		if urgency, ok := cap.Info_Urgency_value[strings.ToUpper(val)]; ok {
			req.Urgency = cap.Info_Urgency(urgency)
		} else {
			return nil, errors.New("Invalid urgency value")
		}
	}

	var err error
	req.Effective, err = getTimeQuery(ctx, query, "effective")
	if err != nil {
		return nil, err
	}

	req.Onset, err = getTimeQuery(ctx, query, "onset")
	if err != nil {
		return nil, err
	}

	req.Expires, err = getTimeQuery(ctx, query, "expires")
	if err != nil {
		return nil, err
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
			return nil, err
		}

		lat, err := strconv.ParseFloat(components[0], 64)
		if err != nil {
			log.WithError(err).Error("Failed to get point.lat")
			raven.CaptureError(err, nil)
			ext.Error.Set(span, true)
			span.LogFields(otlog.Error(err))
			return nil, err
		}

		lon, err := strconv.ParseFloat(components[1], 64)
		if err != nil {
			log.WithError(err).Error("Failed to get point.lon")
			raven.CaptureError(err, nil)
			ext.Error.Set(span, true)
			span.LogFields(otlog.Error(err))
			return nil, err
		}

		req.Point = &alerts.Coordinate{
			Lat: lat,
			Lon: lon,
		}
	}

	return req, nil
}
