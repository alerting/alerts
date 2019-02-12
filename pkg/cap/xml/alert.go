package capxml

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/xml"
	"fmt"
)

// An Alert represents an alert issued through a Common Alerting Protocol (CAP) compliant system.
//proteus:generate
type Alert struct {
	XMLName xml.Name `xml:"alert" json:"-"`

	Identifier  string      `xml:"identifier" json:"identifier"`
	Sender      string      `xml:"sender" json:"sender"`
	Sent        Time        `xml:"sent" json:"sent"`
	Status      Status      `xml:"status" json:"status"`
	MessageType MessageType `xml:"msgType" json:"message_type"`
	Source      string      `xml:"source" json:"source"`
	Scope       Scope       `xml:"scope" json:"scope"`
	Restriction string      `xml:"restriction" json:"restriction"`
	Addresses   List        `xml:"addresses" json:"addresses"`
	Codes       []string    `xml:"code" json:"codes"`
	Note        string      `xml:"note" json:"note"`
	References  References  `xml:"references" json:"references"`
	Incidents   List        `xml:"incidents" json:"incidents"`
	Infos       []*Info     `xml:"info" json:"infos"`
	Superseded  bool        `xml:"-" json:"superseded"`
}

// ID returns the alert's identification string.
// The identification string is a sha1sum of sender,unix(sent),identifier.
func (alert *Alert) ID() string {
	hash := sha1.New()
	hash.Write([]byte(fmt.Sprintf("%s,%d,%s", alert.Sender, alert.Sent.Unix(), alert.Identifier)))
	return hex.EncodeToString(hash.Sum(nil))
}

// Supersede marks an alert as superseded.
func (alert *Alert) Supersede() {
	alert.Superseded = true
}
