package capxml

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"strings"
)

// A ResponseType represents the type of action recommended.
//proteus:generate
type ResponseType int

const (
	// ResponseTypeUnknown represents an unknown response type.
	ResponseTypeUnknown ResponseType = iota

	// ResponseTypeShelter represents taking shelter in place or per the instruction.
	ResponseTypeShelter

	// ResponseTypeEvacuate represents replocating as instructed in the instruction.
	ResponseTypeEvacuate

	// ResponseTypePrepare represents making preparations per the intruction.
	ResponseTypePrepare

	// ResponseTypeExecute represents executing a pre-planned activity identified in the instruction.
	ResponseTypeExecute

	// ResponseTypeAvoid represents avoding the subject event as per the instruction.
	ResponseTypeAvoid

	// ResponseTypeMonitor represents attending to information sources as described in the instruction.
	ResponseTypeMonitor

	// ResponseTypeAssess represents evaluating the information in the mssage.
	ResponseTypeAssess

	// ResponseTypeAllClear represents that the subject event no longer poses a threat or concern and any follow on action is described in the instruction.
	ResponseTypeAllClear

	// ResponseTypeNone represents that no action recommended.
	ResponseTypeNone
)

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (responseType *ResponseType) UnmarshalText(data []byte) error {
	str := strings.ToLower(string(data))

	if str == "shelter" {
		*responseType = ResponseTypeShelter
	} else if str == "evacuate" {
		*responseType = ResponseTypeEvacuate
	} else if str == "prepare" {
		*responseType = ResponseTypePrepare
	} else if str == "execute" {
		*responseType = ResponseTypeExecute
	} else if str == "avoid" {
		*responseType = ResponseTypeAvoid
	} else if str == "monitor" {
		*responseType = ResponseTypeMonitor
	} else if str == "assess" {
		*responseType = ResponseTypeAssess
	} else if str == "allclear" || str == "all clear" {
		*responseType = ResponseTypeAllClear
	} else if str == "none" {
		*responseType = ResponseTypeNone
	} else {
		return errors.New("Unknown ResponseType value: " + str)
	}

	return nil
}

// UnmarshalXML implements the xml.Unmarshaler interface.
func (responseType *ResponseType) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var data []byte
	if err := d.DecodeElement(&data, &start); err != nil {
		return err
	}
	return responseType.UnmarshalText(data)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (responseType *ResponseType) UnmarshalJSON(b []byte) error {
	var value int

	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}

	*responseType = ResponseType(value)
	return nil
}

// MarshalXML implements the xml.Marshaler interface.
func (responseType ResponseType) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	str := responseType.String()

	// Some values differ in XML
	if str == "All Clear" {
		str = "AllClear"
	}

	return e.EncodeElement(str, start)
}

// String returns a ResponseType as a string.
func (responseType ResponseType) String() string {
	if responseType == ResponseTypeShelter {
		return "Shelter"
	} else if responseType == ResponseTypeEvacuate {
		return "Evacuate"
	} else if responseType == ResponseTypePrepare {
		return "Prepare"
	} else if responseType == ResponseTypeExecute {
		return "Execute"
	} else if responseType == ResponseTypeAvoid {
		return "Avoid"
	} else if responseType == ResponseTypeMonitor {
		return "Monitor"
	} else if responseType == ResponseTypeAssess {
		return "Assess"
	} else if responseType == ResponseTypeAllClear {
		return "All Clear"
	} else if responseType == ResponseTypeNone {
		return "None"
	} else if responseType == ResponseTypeUnknown {
		return "Unknown"
	}

	return ""
}
