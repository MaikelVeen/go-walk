package gpx

// Segment is a segment of a Track, containing one or many TrackPoints.
type Segment struct {
	Points []Point `xml:"trkpt"`
}
