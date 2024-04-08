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
	  <trk>
		<name>Rotterdam Walking</name>
		<type>walking</type>
		<trkseg>
		  <trkpt lat="51.9237274490296840667724609375" lon="4.4737290032207965850830078125">
			<ele>23</ele>
			<time>2024-04-08T19:23:26.000Z</time>
			<extensions>
			  <ns3:TrackPointExtension>
				<ns3:hr>119</ns3:hr>
			  </ns3:TrackPointExtension>
			</extensions>
		  </trkpt>
		</trkseg>
	  </trk>
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
				Tracks: []Track{{
					Name: "Rotterdam Walking",
					Type: "walking",
					Segments: []TrackSegment{
						{
							Points: []TrackPoint{
								TrackPoint{
									Latitude:  51.9237274490296840667724609375,
									Longitude: 4.4737290032207965850830078125,
								},
							},
						},
					},
				}},
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
