package gpx

import "math"

// Point is a single point within a segment.
type Point struct {
	Latitude  float64 `xml:"lat,attr"`
	Longitude float64 `xml:"lon,attr"`
}

func (a Point) Distance(b Point) float64 {
	return math.Hypot(a.Longitude-b.Longitude, a.Latitude-b.Latitude)
}
