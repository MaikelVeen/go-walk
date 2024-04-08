package main

import (
	"encoding/xml"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalGPX(t *testing.T) {
	t.Parallel()

	validGPXData := `<?xml version="1.0" encoding="UTF-8"?>
<gpx creator="TestCreator" version="1.1">
</gpx>`
	invalidXMLData := `<?xml version="1.0" encoding="UTF-8"?>
<gpx creator="TestCreator" version="1.1">`
	emptyGPXData := ``

	type args struct {
		r io.ReadCloser
	}
	tests := []struct {
		name    string
		args    args
		want    *GPX
		wantErr bool
	}{
		{
			name: "Valid GPX data",
			args: args{
				r: io.NopCloser(strings.NewReader(validGPXData)),
			},
			want: &GPX{
				XMLName: xml.Name{Local: "gpx"},
				Creator: "TestCreator",
				Version: "1.1",
			},
			wantErr: false,
		},
		{
			name: "Invalid XML data",
			args: args{
				r: io.NopCloser(strings.NewReader(invalidXMLData)),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Empty GPX data",
			args: args{
				r: io.NopCloser(strings.NewReader(emptyGPXData)),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnmarshalGPX(tt.args.r)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.EqualValues(t, tt.want, got)
		})
	}
}
