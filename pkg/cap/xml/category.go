package capxml

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"strings"
)

// A Category represents an alert category.
//proteus:generate
type Category int

const (
	// CategoryUnknown represents an unknown event category.
	CategoryUnknown Category = iota

	// CategoryGeophysical represents a geophysical event.
	CategoryGeophysical

	// CategoryMeteorological represents a meteorological event.
	CategoryMeteorological

	// CategorySafety represents a general emergency and public safety event.
	CategorySafety

	// CategorySecurity represents a law enforcement, military, homeland and local/private security event.
	CategorySecurity

	// CategoryRescue represents a rescue and recovert event.
	CategoryRescue

	// CategoryFire represents a fire suppression and rescue event.
	CategoryFire

	// CategoryHealth represents a medical and public health event.
	CategoryHealth

	// CategoryEnvironment represents a pollutiohn and other environmental event.
	CategoryEnvironment

	// CategoryTransport represents a public and private transportation event.
	CategoryTransport

	// CategoryInfrastructure represents a utility, telecommunication or other non-transport infrastructure event.
	CategoryInfrastructure

	// CategoryCBRNE represents a Chemical, Biological, Radiological, Nuclear or High-Yield explosive threat or attack event.
	CategoryCBRNE

	// CategoryOther represnts other events.
	CategoryOther
)

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (category *Category) UnmarshalText(data []byte) error {
	value := strings.ToLower(string(data))

	if value == "geo" || value == "geophysical" {
		*category = CategoryGeophysical
	} else if value == "met" || value == "meteorological" {
		*category = CategoryMeteorological
	} else if value == "safety" {
		*category = CategorySafety
	} else if value == "security" {
		*category = CategorySecurity
	} else if value == "rescue" {
		*category = CategoryRescue
	} else if value == "fire" {
		*category = CategoryFire
	} else if value == "health" {
		*category = CategoryHealth
	} else if value == "env" || value == "environment" {
		*category = CategoryEnvironment
	} else if value == "transport" {
		*category = CategoryTransport
	} else if value == "infra" || value == "infrastructure" {
		*category = CategoryInfrastructure
	} else if value == "cbrne" {
		*category = CategoryCBRNE
	} else if value == "other" {
		*category = CategoryOther
	} else {
		return errors.New("Unknown Category value: " + value)
	}

	return nil
}

// UnmarshalXML implements the xml.Unmarshaler interface.
func (category *Category) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var data []byte
	if err := d.DecodeElement(&data, &start); err != nil {
		return err
	}
	return category.UnmarshalText(data)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (category *Category) UnmarshalJSON(b []byte) error {
	var value int

	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}

	*category = Category(value)
	return nil
}

// MarshalXML implements the xml.Marshaler interface.
func (category Category) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	str := category.String()

	// Some of the values are shortened in the xml version
	if str == "Geophysical" {
		str = "Geo"
	} else if str == "Meteorological" {
		str = "Met"
	} else if str == "Environment" {
		str = "Env"
	} else if str == "Infrastructure" {
		str = "Infra"
	}

	return e.EncodeElement(str, start)
}

// String returns a Category as a string.
func (category Category) String() string {
	if category == CategoryGeophysical {
		return "Geophysical"
	} else if category == CategoryMeteorological {
		return "Meteorological"
	} else if category == CategorySafety {
		return "Safety"
	} else if category == CategorySecurity {
		return "Security"
	} else if category == CategoryRescue {
		return "Rescue"
	} else if category == CategoryFire {
		return "Fire"
	} else if category == CategoryHealth {
		return "Health"
	} else if category == CategoryEnvironment {
		return "Environment"
	} else if category == CategoryTransport {
		return "Transport"
	} else if category == CategoryInfrastructure {
		return "Infrastructure"
	} else if category == CategoryCBRNE {
		return "CBRNE"
	} else if category == CategoryOther {
		return "Other"
	} else if category == CategoryUnknown {
		return "Unknown"
	}

	return ""
}
