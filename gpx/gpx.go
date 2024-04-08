package gpx

import (
	"encoding/xml"
	"io"
)

// GPX is the root element of a GPS Exchange Format File.
type GPX struct {
	XMLName xml.Name `xml:"gpx"`
	Creator string   `xml:"creator,attr"`
	Version string   `xml:"version,attr"`
	Tracks  []Track  `xml:"trk"`
}

// Track is a track in a GPX file, containing one or many segments.
type Track struct {
	Name     string         `xml:"name"`
	Type     string         `xml:"type"`
	Segments []TrackSegment `xml:"trkseg"`
}

// TrackSegement is a segement of a Track, containing one or many TrackPoints.
type TrackSegment struct {
	Points []TrackPoint `xml:"trkpt"`
}

// TrackPoint is a single point within a segment.
type TrackPoint struct {
	Latitude  float64 `xml:"lat,attr"`
	Longitude float64 `xml:"lon,attr"`
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
