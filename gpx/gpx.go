package main

import (
	"encoding/xml"
	"io"
)

// GPX represents the root element of a GPS Exchange Format File.
type GPX struct {
	XMLName xml.Name `xml:"gpx"`
	Creator string   `xml:"creator,attr"`
	Version string   `xml:"version,attr"`
}

// UnmarshalGPX reads and parses the GPX data from the provided io.ReadCloser.
// It returns a pointer to the parsed GPX struct and any error encountered during parsing.
func UnmarshalGPX(r io.ReadCloser) (*GPX, error) {
	defer r.Close()

	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var gpx GPX
	if err := xml.Unmarshal(data, &gpx); err != nil {
		return nil, err
	}

	return &gpx, nil
}
