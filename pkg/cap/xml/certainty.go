package capxml

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"strings"
)

// A Certainty represents the certainty of an event.
//proteus:generate
type Certainty int

const (
	// CertaintyUnknown represents an event with unknown certainty.
	CertaintyUnknown Certainty = iota

	// CertaintyObserved represents an event determined to have occurred or to be ongoing.
	CertaintyObserved

	// CertaintyLikely represents an event that is likely (p > ~50%).
	CertaintyLikely

	// CertaintyPossible represents an event that is possible but not likely (p <= ~50%)
	CertaintyPossible

	// CertaintyUnlikely represents an event that is unlikely to occure (p ~ 0)
	CertaintyUnlikely
)

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (certainty *Certainty) UnmarshalText(data []byte) error {
	str := strings.ToLower(string(data))

	if str == "observed" {
		*certainty = CertaintyObserved
	} else if str == "likely" || str == "verylikely" || str == "very likely" {
		*certainty = CertaintyLikely
	} else if str == "possible" {
		*certainty = CertaintyObserved
	} else if str == "unlikely" {
		*certainty = CertaintyUnlikely
	} else if str == "unknown" {
		*certainty = CertaintyUnknown
	} else {
		return errors.New("Unknown Certainty value: " + str)
	}

	return nil
}

// UnmarshalXML implements the xml.Unmarshaler interface.
func (certainty *Certainty) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var data []byte
	if err := d.DecodeElement(&data, &start); err != nil {
		return err
	}
	return certainty.UnmarshalText(data)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (certainty *Certainty) UnmarshalJSON(b []byte) error {
	var value int

	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}

	*certainty = Certainty(value)
	return nil
}

// MarshalXML implements the xml.Marshaler interface.
func (certainty Certainty) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(certainty.String(), start)
}

// String returns a Certainty as a string.
func (certainty Certainty) String() string {
	if certainty == CertaintyObserved {
		return "Observed"
	} else if certainty == CertaintyLikely {
		return "Likely"
	} else if certainty == CertaintyPossible {
		return "Possible"
	} else if certainty == CertaintyUnlikely {
		return "Unlikley"
	} else if certainty == CertaintyUnknown {
		return "Unknown"
	}

	return ""
}
