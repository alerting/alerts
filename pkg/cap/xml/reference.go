package capxml

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"
)

// A Reference represents a the identifying information of another alert.
type Reference struct {
	Sender     string `json:"sender"`
	Identifier string `json:"identifier"`
	Sent       Time   `json:"sent"`
}

// References is an array of Reference.
type References []*Reference

// ID returns the ID of the referenced alert.
func (reference *Reference) ID() string {
	hash := sha1.New()
	hash.Write([]byte(fmt.Sprintf("%s,%d,%s", reference.Sender, reference.Sent.Unix(), reference.Identifier)))
	return hex.EncodeToString(hash.Sum(nil))
}

// MarshalJSON implements the json.Marshaler interface.
func (reference Reference) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender     string `json:"sender"`
		Identifier string `json:"identifier"`
		Sent       Time   `json:"sent"`
		ID         string `json:"id"`
	}{
		Sender:     reference.Sender,
		Identifier: reference.Identifier,
		Sent:       reference.Sent,
		ID:         reference.ID(),
	})
}

// UnmarshalXML implements the xml.Unmarshaler interface.
func (m *References) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var str string
	if err := d.DecodeElement(&str, &start); err != nil {
		return err
	}

	// Ignore empty strings
	if str == "" {
		return nil
	}

	// Create the list if it doesn't already exist
	references := strings.Split(str, " ")

	for _, reference := range references {
		components := strings.Split(reference, ",")

		var sent Time
		if err := sent.UnmarshalText([]byte(components[2])); err != nil {
			return err
		}

		*m = append(*m, &Reference{
			Sender:     components[0],
			Identifier: components[1],
			Sent:       sent,
		})
	}

	return nil
}

// MarshalXML implements the xml.Marshaler interface.
func (m *References) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	values := make([]string, len(*m))

	for i, ref := range *m {
		values[i] = fmt.Sprintf("%s,%s,%s", ref.Sender, ref.Identifier, ref.Sent.FormatCAP())
	}

	return e.EncodeElement(strings.Join(values, " "), start)
}
