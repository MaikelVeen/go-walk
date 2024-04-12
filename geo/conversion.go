package geo

import "math"

const (
	tileSize          = 256.0
	initialResolution = 2 * math.Pi * 6378137 / tileSize
	originShift       = 2 * math.Pi * 6378137 / 2
)

// Resolution calculates the resolution (meters/pixel) for given zoom level (measured at Equator)
func Resolution(zoom int) float64 {
	return initialResolution / math.Pow(2, float64(zoom))
}

// LatLonToMeters converts given lat/lon in WGS84 Datum to XY in Spherical Mercator EPSG:900913
func LatLonToMeters(lat, lon float64) (float64, float64) {
	x := lon * originShift / 180
	y := math.Log(math.Tan((90+lat)*math.Pi/360)) / (math.Pi / 180)
	y = y * originShift / 180
	return x, y
}

// MetersToPixels converts EPSG:900913 to pixel coordinates in given zoom level
func MetersToPixels(x, y float64, zoom int) (float64, float64) {
	res := Resolution(zoom)
	px := (x + originShift) / res
	py := (y + originShift) / res
	return px, py
}

// LatLonToPixels converts given lat/lon in WGS84 Datum to pixel coordinates in given zoom level
func LatLonToPixels(lat, lon float64, zoom int) (float64, float64) {
	x, y := LatLonToMeters(lat, lon)
	return MetersToPixels(x, y, zoom)
}

func OriginPixels(zoom int) (float64, float64) {
	return LatLonToPixels(0, 0, zoom)
}
