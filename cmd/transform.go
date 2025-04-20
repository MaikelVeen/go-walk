package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	geojson "github.com/paulmach/go.geojson"
	"github.com/spf13/cobra"
)

// TransformCommand represents a command for transforming GPX data to GeoJSON.
type TransformCommand struct {
	Command *cobra.Command

	// Directory is the path to the directory containing GPX files.
	Directory string
	// OutputDir is the path to the directory where GeoJSON files will be saved.
	OutputDir string
}

// NewTransformCommand returns a new instance of the TransformCommand.
func NewTransformCommand() *TransformCommand {
	tc := &TransformCommand{}
	tc.Command = &cobra.Command{
		Use:          "transform [-d Dir | -O OutputDir]",
		Short:        "Transform each GPX file in a directory into a separate GeoJSON file",
		RunE:         tc.Run,
		Args:         cobra.NoArgs,
		SilenceUsage: true,
	}

	tc.Command.Flags().StringVarP(&tc.Directory, "dir", "d", ".", "The path of the directory where the GPX files are located")
	tc.Command.Flags().StringVarP(&tc.OutputDir, "outdir", "O", ".", "The directory to save the output GeoJSON files")
	return tc
}

// Run executes the transform command.
func (c *TransformCommand) Run(cmd *cobra.Command, args []string) error {
	absInputDir, err := filepath.Abs(c.Directory)
	if err != nil {
		return fmt.Errorf("invalid input directory %s: %w", c.Directory, err)
	}

	absOutputDir, err := filepath.Abs(c.OutputDir)
	if err != nil {
		return fmt.Errorf("invalid output directory %s: %w", c.OutputDir, err)
	}

	if err := os.MkdirAll(absOutputDir, 0755); err != nil {
		return fmt.Errorf("cannot create output directory %s: %w", absOutputDir, err)
	}

	files, err := os.ReadDir(absInputDir)
	if err != nil {
		return fmt.Errorf("cannot read input directory %s: %w", absInputDir, err)
	}

	processedFiles := 0
	totalPoints := 0
	hadErrors := false
	var successfulFiles []string

	fmt.Printf("Transforming GPX files from %s to GeoJSON in %s...\n", absInputDir, absOutputDir)

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(strings.ToLower(file.Name()), ".gpx") {
			continue
		}

		inputFilePath := filepath.Join(absInputDir, file.Name())
		baseName := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
		outputFilename := baseName + ".geojson"

		points, err := transformSingleGPX(inputFilePath, absOutputDir, outputFilename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error processing %s: %v\n", file.Name(), err)
			hadErrors = true
			continue // Process next file
		}

		if points > 0 {
			processedFiles++
			totalPoints += points
			successfulFiles = append(successfulFiles, outputFilename)
		}
	}

	if len(successfulFiles) > 0 {
		manifestPath := filepath.Join(absOutputDir, "manifest.json")
		manifestData, err := json.MarshalIndent(successfulFiles, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to marshal manifest file: %v\n", err)
			hadErrors = true
		} else {
			err = os.WriteFile(manifestPath, manifestData, 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Failed to write manifest file %s: %v\n", manifestPath, err)
				hadErrors = true
			} else {
				fmt.Printf("Generated manifest file: %s\n", manifestPath)
			}
		}
	}

	fmt.Println("--------------------")
	fmt.Printf("Summary: Processed %d GPX file(s), total %d points.\n", processedFiles, totalPoints)
	if hadErrors {
		fmt.Println("Warning: Some files could not be processed successfully (see errors above).")
		return errors.New("one or more files failed transformation")
	}

	if processedFiles == 0 && !hadErrors {
		fmt.Println("No valid GPX files with track points found to process.")
	}

	fmt.Println("Transformation complete.")
	return nil
}

// transformSingleGPX reads a single GPX file, transforms it to GeoJSON,
// and writes it to the output directory.
// It returns the number of points processed and any error encountered.
func transformSingleGPX(inputPath, outputDir, outputFilename string) (int, error) {
	gpxData, err := readGPXFile(inputPath)
	if err != nil {
		return 0, fmt.Errorf("failed to read/parse GPX: %w", err)
	}

	points := gpxData.Points()
	if len(points) == 0 {
		return 0, nil
	}

	coordinates := make([][]float64, len(points))
	for i, p := range points {
		coordinates[i] = []float64{p.Longitude, p.Latitude}
	}

	geometry := geojson.NewLineStringGeometry(coordinates)
	feature := geojson.NewFeature(geometry)
	featureCollection := geojson.NewFeatureCollection()
	featureCollection.AddFeature(feature)

	jsonData, err := json.MarshalIndent(featureCollection, "", "  ")
	if err != nil {
		return 0, fmt.Errorf("failed to marshal GeoJSON: %w", err)
	}

	outputFilePath := filepath.Join(outputDir, outputFilename)

	err = os.WriteFile(outputFilePath, jsonData, 0644)
	if err != nil {
		return 0, fmt.Errorf("failed to write output file %s: %w", outputFilePath, err)
	}

	return len(points), nil
}
