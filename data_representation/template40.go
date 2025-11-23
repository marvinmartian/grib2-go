package data_representation

import (
	"fmt"
	"math"

	"github.com/segmed/openjpeg/gojp2"
	"github.com/skysparq/grib2-go/record"
	u "github.com/skysparq/grib2-go/utility"
)

// Template40 contains the fields for Grid point data - JPEG 2000 code stream format
type Template40 struct {
	ReferenceValue         float32
	BinaryScaleFactor      int
	DecimalScaleFactor     int
	BitDepth               int
	OriginalFieldType      int
	CompressionType        int
	TargetCompressionRatio int
}

// Parse fills in the template from the provided section
func (t Template40) Parse(section record.Section5) (record.DataRepresentationDefinition, error) {
	err := checkSectionNum(section, 40)
	if err != nil {
		return t, err
	}

	data := section.DataRepresentationTemplateData
	t.ReferenceValue = u.Float32(data[0:4])
	t.BinaryScaleFactor = u.SignAndMagnitudeInt16(data[4:6])
	t.DecimalScaleFactor = u.SignAndMagnitudeInt16(data[6:8])
	t.BitDepth = int(data[8])
	t.OriginalFieldType = int(data[9])
	t.CompressionType = int(data[10])
	t.TargetCompressionRatio = int(data[11])
	return t, nil
}

// DecimalScale returns the decimal scale factor of the record. The decimal scale factor is used to shift the
// decimal point of a decoded value to the correct position.
func (t Template40) DecimalScale() int {
	return t.DecimalScaleFactor
}

// GetValues unpacks the record's data into the original values
func (t Template40) GetValues(rec record.Record) ([]float64, error) {
	// Decode JPEG 2000 data
	width, height, pixelData, err := gojp2.DecodeJ2KImage(rec.Data.Data)
	if err != nil {
		return nil, fmt.Errorf("error decoding JPEG 2000 data: %w", err)
	}

	// Check that we have the expected number of points
	expectedPoints := rec.Grid.TotalPoints
	actualPoints := width * height
	if actualPoints != expectedPoints {
		return nil, fmt.Errorf("JPEG 2000 decoded data has %d points but expected %d", actualPoints, expectedPoints)
	}

	// Convert to float64 values
	values := make([]float64, expectedPoints)
	idx := 0
	for _, row := range pixelData {
		for _, pixel := range row {
			// Apply the GRIB2 unpacking formula: Y = (R + X * 2^E) / 10^D
			scaledValue := u.Unpack(float64(t.ReferenceValue), pixel, t.BinaryScaleFactor, t.DecimalScaleFactor)
			values[idx] = scaledValue
			idx++
		}
	}

	// Apply bitmap masking for missing values
	bmpR, err := NewBitmapReader(rec)
	if err != nil {
		return nil, fmt.Errorf("error creating bitmap reader: %w", err)
	}
	for i := range values {
		if bmpR.IsMissing(i) {
			values[i] = math.NaN()
		}
	}

	return values, nil
}
