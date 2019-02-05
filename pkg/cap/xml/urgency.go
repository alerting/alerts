package capxml

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"strings"
)

// A Urgency represents the urgency of an event.
//proteus:generate
type Urgency int

const (
	// UrgencyUnknown represents unknown urgency.
	UrgencyUnknown Urgency = iota

	// UrgencyImmediate represents that responsive action should be taken immediately.
	UrgencyImmediate

	// UrgencyExpected represents that responsive action should be taken soon (within the next hour).
	UrgencyExpected

	// UrgencyFuture represents that responsive action should be taken in the near future.
	UrgencyFuture

	// UrgencyPast represents that responsive action is no longer required.
	UrgencyPast
)

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (urgency *Urgency) UnmarshalText(data []byte) error {
	str := strings.ToLower(string(data))

	if str == "immediate" {
		*urgency = UrgencyImmediate
	} else if str == "expected" {
		*urgency = UrgencyExpected
	} else if str == "future" {
		*urgency = UrgencyFuture
	} else if str == "past" {
		*urgency = UrgencyPast
	} else if str == "unknown" {
		*urgency = UrgencyUnknown
	} else {
		return errors.New("Unknown Urgency value: " + str)
	}

	return nil
}

// UnmarshalXML implements the xml.Unmarshaler interface.
func (urgency *Urgency) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var data []byte
	if err := d.DecodeElement(&data, &start); err != nil {
		return err
	}
	return urgency.UnmarshalText(data)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (urgency *Urgency) UnmarshalJSON(b []byte) error {
	var value int

	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}

	*urgency = Urgency(value)
	return nil
}

// MarshalXML implements the xml.Marshaler interface.
func (urgency Urgency) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(urgency.String(), start)
}

// String returns a Urgency as a string.
func (urgency Urgency) String() string {
	if urgency == UrgencyImmediate {
		return "Immediate"
	} else if urgency == UrgencyExpected {
		return "Expected"
	} else if urgency == UrgencyFuture {
		return "Future"
	} else if urgency == UrgencyPast {
		return "Past"
	} else if urgency == UrgencyUnknown {
		return "Unknown"
	}

	return ""
}
