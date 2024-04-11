package gpx

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPointDistance(t *testing.T) {
	pointA := Point{
		Latitude:  51.924373269285844,
		Longitude: 4.469002690910358,
	}

	pointB := Point{
		Latitude:  51.92308253810835,
		Longitude: 4.469793977934499,
	}

	expectedDistance := math.Hypot(pointA.Longitude-pointB.Longitude, pointA.Latitude-pointB.Latitude)
	actualDistance := pointA.Distance(pointB)

	assert.Equal(t, expectedDistance, actualDistance)
}
