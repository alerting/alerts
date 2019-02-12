package capxml

import (
	"encoding/xml"
	"strings"
	"time"
)

const (
	// TimeFormat is the time format used for dates in the CAP protocol.
	TimeFormat = "2006-01-02T15:04:05-07:00"
)

// We re-implement the Time structure
// simply to overwrite how the timezone is output
// since the cap standard requires -00:00 instead of +00:00.

// A Time represents an instant in time.
type Time struct {
	time.Time
}

// FormatCAP returns the time, formatted to Common Alert Protocol standards.
func (time Time) FormatCAP() string {
	str := time.Format(TimeFormat)
	return strings.Replace(str, "+00:00", "-00:00", 1)
}

// MarshalXML implements the xml.Marshaler interface.
func (time Time) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(time.FormatCAP(), start)
}
