package cap

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"

	capxml "github.com/alerting/alerts/pkg/cap/xml"
	"github.com/golang/protobuf/jsonpb"
)

func (alert *Alert) ID() string {
	hash := sha1.New()
	hash.Write([]byte(fmt.Sprintf("%s,%d,%s", alert.Sender, alert.Sent.GetSeconds(), alert.Identifier)))
	return hex.EncodeToString(hash.Sum(nil))
}

func (alert *Alert) Reference() *Reference {
	return &Reference{
		Identifier: alert.Identifier,
		Sender:     alert.Sender,
		Sent:       alert.Sent,
	}
}

func (r *Reference) ID() string {
	hash := sha1.New()
	hash.Write([]byte(fmt.Sprintf("%s,%d,%s", r.Sender, r.Sent.GetSeconds(), r.Identifier)))
	return hex.EncodeToString(hash.Sum(nil))
}

func (alert *Alert) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	raw := capxml.Alert{}

	// Decode into capxml.Alert
	if err := d.DecodeElement(&raw, &start); err != nil {
		return err
	}

	// Now jsonify
	b, _ := json.Marshal(raw)
	log.Println(string(b))

	jd := jsonpb.Unmarshaler{
		AllowUnknownFields: true,
	}
	return jd.Unmarshal(bytes.NewBuffer(b), alert)
}

// Checksum returns the checksum value for the resource, if available.
// The Resource.Digest value is returned if provided. If no digest is available,
// it is calculated from the DerefUri value (if available). Otherwise,
// an empty string is returned.
func (res *Resource) Checksum() string {
	// If we have the digest, return it.
	if res.Digest != "" {
		return res.Digest
	}

	// If we have the contents, sha1sum that
	if len(res.DerefUri) > 0 {
		hash := sha1.New()
		hash.Write(res.DerefUri)
		return hex.EncodeToString(hash.Sum(nil))
	}

	return ""
}
