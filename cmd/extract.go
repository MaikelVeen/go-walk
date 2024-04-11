package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/MaikelVeen/go-walk/gpx"
	"github.com/spf13/cobra"
)

type ExtractCommand struct {
	*cobra.Command
}

// NewExtractCommand returns a new instance of the ExtractCommand.
func NewExtractCommand() *ExtractCommand {
	ec := &ExtractCommand{}
	ec.Command = &cobra.Command{
		Use:   "extract <folder>",
		Short: "Extracts coordinates from all .gpx files in a folder",
		RunE:  ec.Run,
		Args:  cobra.ExactArgs(1),
	}
	return ec
}

// Run executes the extract command and returns any error ecountered.
func (c *ExtractCommand) Run(cmd *cobra.Command, args []string) error {
	path := args[0]

	data, err := readGPXFilesInFolder(path)
	if err != nil {
		return err
	}

	fmt.Println(len(data))

	return nil
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
