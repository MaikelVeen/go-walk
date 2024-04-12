package cmd

import (
	"math"
	"os"
	"path/filepath"

	"github.com/MaikelVeen/go-walk/geo"
	"github.com/MaikelVeen/go-walk/gpx"
	"github.com/fogleman/gg"
	"github.com/spf13/cobra"
)

// VisualiseCommand represents a command for visualizing a directory.
type VisualiseCommand struct {
	Command *cobra.Command

	// Directory is the path to the directory to be visualized.
	Directory string
	// OutputFilename is the name of the output file for the visualization.
	OutputFilename string
	// ZoomLevel is the level of zoom for the visualization.
	ZoomLevel int
}

// NewVisualiseCommand returns a new instance of the ExtractCommand.
func NewVisualiseCommand() *VisualiseCommand {
	ec := &VisualiseCommand{}
	ec.Command = &cobra.Command{
		Use:   "visualise [-d Dir | -o Output | -z Zoom]",
		Short: "Visualise coordinates parsed from gpx files",
		RunE:  ec.Run,
		Args:  cobra.NoArgs,
	}

	ec.Command.Flags().IntVarP(&ec.ZoomLevel, "zoom", "z", 16, "Determines the zoom level of the final projection")
	ec.Command.Flags().StringVarP(&ec.Directory, "dir", "d", "./data", "The path of the directory where the GXP files are located")
	// TODO: OutputFilename
	return ec
}

// Run executes the extract command and returns any error ecountered.
func (c *VisualiseCommand) Run(cmd *cobra.Command, args []string) error {
	data, err := readGPXFilesInFolder(c.Directory)
	if err != nil {
		return err
	}

	var points []gpx.LatLng
	for _, d := range data {
		points = append(points, d.Points()...)
	}

	newOriginX, newOriginY := geo.OriginPixels(c.ZoomLevel)
	minX, minY := newOriginX, newOriginY

	pixels := []struct{ X, Y float64 }{}
	for _, point := range points {
		x, y := geo.LatLonToPixels(point.Latitude, point.Longitude, c.ZoomLevel)
		x -= newOriginX
		y -= newOriginY

		if x < minX {
			minX = x
		}
		if y < minY {
			minY = y
		}

		pixels = append(pixels, struct{ X, Y float64 }{X: x, Y: y})
	}

	// Normalize points to ensure there are no negative values
	for i := range pixels {
		pixels[i].X -= minX
		pixels[i].Y -= minY
	}

	// Determine the size of the image based on transformed coordinates
	maxX, maxY := 0.0, 0.0
	for _, p := range pixels {
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}

	for i, p := range pixels {
		pixels[i].Y = maxY - p.Y // invert the Y-axis
	}

	// TODO: Extract to Draw func
	// Create a new image context
	dc := gg.NewContext(int(maxX), int(maxY))
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)
	dc.SetLineWidth(5) // TODO: Make configurable

	if len(pixels) > 0 {
		dc.MoveTo(pixels[0].X, pixels[0].Y) // move to the start point
		for i := 0; i < len(pixels)-1; i++ {
			if distance(pixels[i].X, pixels[i].Y, pixels[i+1].X, pixels[i+1].Y) > 20 {
				dc.StrokePreserve()
				dc.MoveTo(pixels[i+1].X, pixels[i+1].Y)
			} else {
				dc.LineTo(pixels[i+1].X, pixels[i+1].Y)
			}
		}
		dc.StrokePreserve()
	}

	err = dc.SavePNG("output.png")
	if err != nil {
		return err
	}

	return nil
}

func distance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2))
}

// readGPXFilesInFolder processes GPX files from a specified folder.
// It returns a slice of GPX structs and any error encountered.
func readGPXFilesInFolder(folderPath string) ([]*gpx.GPX, error) {
	files, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, err
	}

	var gpxArray []*gpx.GPX
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".gpx" {
			g, err := readGPXFile(filepath.Join(folderPath, file.Name()))
			if err != nil {
				return nil, err
			}

			gpxArray = append(gpxArray, g)
		}
	}

	return gpxArray, nil
}

// readGPXFile reads a GPX file from the specified file path and returns a parsed GPX object.
// If the file cannot be opened or if there is an error parsing the GPX file, an error is returned.
func readGPXFile(filePath string) (*gpx.GPX, error) {
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
