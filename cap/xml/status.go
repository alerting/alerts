package capxml

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"strings"
)

// A Status represents the appropriate handling of the alert message.
//proteus:generate
type Status int

const (
	// StatusUnknown represents an unknown status.
	StatusUnknown Status = iota

	// StatusActual represents an alert actionable by all targeted recipients.
	StatusActual

	// StatusExcercise represents an alert actionable only by desginated exercise participants.
	StatusExcercise

	// StatusSystem represents a message that supports alert network internal functions.
	StatusSystem

	// StatusTest represents technical testing only. All recipients should disregard.
	StatusTest

	// StatusDraft represents a perliminary template or draft. The alert is not actionable in its current form.
	StatusDraft
)

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (status *Status) UnmarshalText(data []byte) error {
	str := strings.ToLower(string(data))

	if str == "actual" {
		*status = StatusActual
	} else if str == "excercise" {
		*status = StatusExcercise
	} else if str == "system" {
		*status = StatusSystem
	} else if str == "test" {
		*status = StatusTest
	} else if str == "draft" {
		*status = StatusDraft
	} else {
		return errors.New("Unknown Status value: " + str)
	}

	return nil
}

// UnmarshalXML implements the xml.Marshaler interface.
func (status *Status) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var data []byte
	if err := d.DecodeElement(&data, &start); err != nil {
		return err
	}
	return status.UnmarshalText(data)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (status *Status) UnmarshalJSON(b []byte) error {
	var value int

	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}

	*status = Status(value)
	return nil
}

// MarshalXML implements the xml.Marshaler interface.
func (status Status) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(status.String(), start)
}

// String returns a Status as a string.
func (status Status) String() string {
	if status == StatusActual {
		return "Actual"
	} else if status == StatusExcercise {
		return "Excercise"
	} else if status == StatusSystem {
		return "System"
	} else if status == StatusTest {
		return "Test"
	} else if status == StatusDraft {
		return "Draft"
	} else if status == StatusUnknown {
		return "Unknown"
	}

	return ""
}
