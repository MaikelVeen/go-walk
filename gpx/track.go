package gpx

// Track is a track in a GPX file, containing one or many segments.
type Track struct {
	Name     string    `xml:"name"`
	Type     string    `xml:"type"`
	Segments []Segment `xml:"trkseg"`
}
