package capxml

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"strings"
)

// A Scope represents the intended distribution of the alert message.
//proteus:generate
type Scope int

const (
	// ScopeUnknown represents an unknown scope.
	ScopeUnknown Scope = iota

	// ScopePublic represents an alert for general dissemination to unrestricted audiences.
	ScopePublic

	// ScopeRestricted represents an alert for dissemination only to users with a known operational requirement.
	ScopeRestricted

	// ScopePrivate represents an alert for dissemination only to specified addresses.
	ScopePrivate
)

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (scope *Scope) UnmarshalText(data []byte) error {
	str := strings.ToLower(string(data))

	if str == "public" {
		*scope = ScopePublic
	} else if str == "restricted" {
		*scope = ScopeRestricted
	} else if str == "private" {
		*scope = ScopePrivate
	} else {
		return errors.New("Unknown Scope value")
	}

	return nil
}

// UnmarshalXML implements the xml.Unmarshaler interface.
func (scope *Scope) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var data []byte
	if err := d.DecodeElement(&data, &start); err != nil {
		return err
	}
	return scope.UnmarshalText(data)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (scope *Scope) UnmarshalJSON(b []byte) error {
	var value int

	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}

	*scope = Scope(value)
	return nil
}

// MarshalXML implements the xml.Marshaler interface.
func (scope Scope) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(scope.String(), start)
}

// String returns a Scope as a string.
func (scope Scope) String() string {
	if scope == ScopePublic {
		return "Public"
	} else if scope == ScopeRestricted {
		return "Restricted"
	} else if scope == ScopePrivate {
		return "Private"
	} else if scope == ScopeUnknown {
		return "Unknown"
	}

	return ""
}
