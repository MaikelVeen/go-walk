package main

import (
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

	walkData, err := ReadGPXFilesInFolder("./data")
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

// ReadGPXFilesInFolder processes GPX files from a specified folder.
// It returns a slice of GPX structs and any error encountered.
func ReadGPXFilesInFolder(folderPath string) ([]*gpx.GPX, error) {
	files, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, err
	}

	var gpxArray []*gpx.GPX
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".gpx" {
			g, err := ReadGPXFile(filepath.Join(folderPath, file.Name()))
			if err != nil {
				return nil, err
			}

			gpxArray = append(gpxArray, g)
		}
	}

	return gpxArray, nil
}

// ReadGPXFile reads a GPX file from the specified file path and returns a parsed GPX object.
// If the file cannot be opened or if there is an error parsing the GPX file, an error is returned.
func ReadGPXFile(filePath string) (*gpx.GPX, error) {
	gpxFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	g, err := gpx.UnmarshalGPX(gpxFile)
	if err != nil {
		return nil, err
	}

	return g, nil
}
