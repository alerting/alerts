package capxml

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"strings"
)

// A MessageType represents the nature of the alert message.
//proteus:generate
type MessageType int

const (
	// MessageTypeUnknown represents an unknown message type.
	MessageTypeUnknown MessageType = iota

	// MessageTypeAlert represents initial information requiring attention by targeted recipients.
	MessageTypeAlert

	// MessageTypeUpdate represents updates and supercedes the earlier message(s) identified in Alert.References.
	MessageTypeUpdate

	// MessageTypeCancel represents a cancellation of the earlier message(s) identified in Alert.References.
	MessageTypeCancel

	// MessageTypeAck represents an acknowledgment receipt and acceptance of the message(s) identified in Alert.References.
	MessageTypeAck

	// MessageTypeError represents a rejection of the message(s) identified in Alert.References.
	MessageTypeError
)

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (messageType *MessageType) UnmarshalText(data []byte) error {
	str := strings.ToLower(string(data))

	if str == "alert" {
		*messageType = MessageTypeAlert
	} else if str == "update" {
		*messageType = MessageTypeUpdate
	} else if str == "cancel" {
		*messageType = MessageTypeCancel
	} else if str == "ack" {
		*messageType = MessageTypeAck
	} else if str == "error" {
		*messageType = MessageTypeError
	} else {
		return errors.New("Unknown MessageType value")
	}

	return nil
}

// UnmarshalXML implements the xml.Unmarshaler interface.
func (messageType *MessageType) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var data []byte
	if err := d.DecodeElement(&data, &start); err != nil {
		return err
	}
	return messageType.UnmarshalText(data)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (messageType *MessageType) UnmarshalJSON(b []byte) error {
	var value int

	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}

	*messageType = MessageType(value)
	return nil
}

// MarshalXML implements the xml.Marshaler interface.
func (messageType MessageType) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(messageType.String(), start)
}

// String returns a MessageType as a string.
func (messageType MessageType) String() string {
	if messageType == MessageTypeAlert {
		return "Alert"
	} else if messageType == MessageTypeUpdate {
		return "Update"
	} else if messageType == MessageTypeCancel {
		return "Cancel"
	} else if messageType == MessageTypeAck {
		return "Ack"
	} else if messageType == MessageTypeError {
		return "Error"
	} else if messageType == MessageTypeUnknown {
		return "Unknown"
	}

	return ""
}
