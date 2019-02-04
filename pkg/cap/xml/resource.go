package capxml

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/xml"
)

// A Resource represents a reference to additional information related to the Info segment.
// These are references to digital assets such as images or audio.
type Resource struct {
	XMLName xml.Name `xml:"resource" json:"-"`

	Description string `xml:"resourceDesc" json:"description"`
	MimeType    string `xml:"mimeType" json:"mime_type"`
	Size        int    `xml:"size" json:"size"`
	URI         string `xml:"uri" json:"uri"`
	Digest      string `xml:"digest" json:"digest"`
	DerefURI    string `xml:"derefUri" json:"deref_uri"`
}

// Checksum returns the checksum value for the resource, if available.
// The Resource.Digest value is returned if provided. If no digest is available,
// it is calculated from the DerefURI value (if available). Otherwise,
// an empty string is returned.
func (res *Resource) Checksum() string {
	// If we have the digest, return it.
	if res.Digest != "" {
		return res.Digest
	}

	// If we have the contents, sha1sum that
	if res.DerefURI != "" {
		hash := sha1.New()
		hash.Write([]byte(res.DerefURI))
		return hex.EncodeToString(hash.Sum(nil))
	}

	return ""
}
