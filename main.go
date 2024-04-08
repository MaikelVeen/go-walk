package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/MaikelVeen/go-walk/gpx"
	"github.com/lmittmann/tint"
)

func main() {
	logger := slog.New(
		tint.NewHandler(os.Stdout, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.TimeOnly,
		}),
	)

	logger.Info("Hello World!")

	walkData, err := processGPXFiles()
	if err != nil {
		slog.Error(err.Error())
	}

	// TODO: Create function to get all points of a GPX.
	for _, dat := range walkData {
		for _, track := range dat.Tracks {
			for _, segment := range track.Segments {
				for _, point := range segment.Points {
					logger.Info("coords", "lat", point.Latitude, "lng", point.Longitude)
				}
			}
		}
	}
}

// TODO: Make the folder an argument.
func processGPXFiles() ([]*gpx.GPX, error) {
	files, err := os.ReadDir("data")
	if err != nil {
		return nil, err
	}

	var gpxArray []*gpx.GPX
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".gpx" {
			gpxFile, err := os.Open(fmt.Sprintf("./data/%s", file.Name()))
			if err != nil {
				return nil, err
			}

			g, err := gpx.UnmarshalGPX(gpxFile)
			if err != nil {
				return nil, err
			}

			gpxArray = append(gpxArray, g)
		}
	}

	return gpxArray, nil
}
