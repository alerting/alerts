package capxml

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"strings"
)

// A Severity represents the severity of an event.
//proteus:generate
type Severity int

const (
	// SeverityUnknown represents an event with unknown certainty.
	SeverityUnknown Severity = iota

	// SeverityExtreme represents an event with extraordinary threat to life or property.
	SeverityExtreme

	// SeveritySevere represents an event with significant threat to life or property.
	SeveritySevere

	// SeverityModerate represents an event with possible threat to life or property.
	SeverityModerate

	// SeverityMinor represents an event with minimal to no known threat to life or property.
	SeverityMinor
)

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (severity *Severity) UnmarshalText(data []byte) error {
	str := strings.ToLower(string(data))

	if str == "extreme" {
		*severity = SeverityExtreme
	} else if str == "severe" {
		*severity = SeveritySevere
	} else if str == "moderate" {
		*severity = SeverityModerate
	} else if str == "minor" {
		*severity = SeverityMinor
	} else if str == "unknown" {
		*severity = SeverityUnknown
	} else {
		return errors.New("Unknown Severity value: " + str)
	}

	return nil
}

// UnmarshalXML implements the xml.Unmarshaler interface.
func (severity *Severity) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var data []byte
	if err := d.DecodeElement(&data, &start); err != nil {
		return err
	}
	return severity.UnmarshalText(data)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (severity *Severity) UnmarshalJSON(b []byte) error {
	var value int

	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}

	*severity = Severity(value)
	return nil
}

// MarshalXML implements the xml.Marshaler interface.
func (severity Severity) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(severity.String(), start)
}

// String returns a Severity as a string.
func (severity Severity) String() string {
	if severity == SeverityExtreme {
		return "Extreme"
	} else if severity == SeveritySevere {
		return "Severe"
	} else if severity == SeverityModerate {
		return "Moderate"
	} else if severity == SeverityMinor {
		return "Minor"
	} else if severity == SeverityUnknown {
		return "Unknown"
	}

	return ""
}
