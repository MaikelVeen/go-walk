package gpx

import "math"

// LatLng is a single point within a segment.
type LatLng struct {
	Latitude  float64 `xml:"lat,attr"`
	Longitude float64 `xml:"lon,attr"`
}

func (a LatLng) Distance(b LatLng) float64 {
	return math.Hypot(a.Longitude-b.Longitude, a.Latitude-b.Latitude)
}
