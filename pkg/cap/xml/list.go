package capxml

import (
	"encoding/xml"
	"strings"
)

// A List represents a list of string values.
type List []string

// UnmarshalXML implements the xml.Unmarshaler interface.
func (m *List) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var str string
	err := d.DecodeElement(&str, &start)

	if err != nil {
		return err
	}

	// Ignore empty string
	if str == "" {
		*m = make([]string, 0)
		return nil
	}

	// TODO: Seperate by any whitespace, ignoring escaped text (")
	*m = strings.Split(str, " ")

	return nil
}
