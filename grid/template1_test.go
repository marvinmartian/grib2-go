package grid_test

import (
	"reflect"
	"testing"

	"github.com/skysparq/grib2-go/grid"
	"github.com/skysparq/grib2-go/record"
)

func TestTemplate1Parse(t *testing.T) {
	// Create a mock section 3 with template 1 data
	// This is a minimal test data for rotated lat/lon template
	data := []byte{
		// Earth shape
		6,
		// Radius scale factor
		0,
		// Radius scale value (4 bytes)
		0, 0, 0, 0,
		// Major axis scale factor
		0,
		// Major axis scale value (4 bytes)
		0, 0, 0, 0,
		// Minor axis scale factor
		0,
		// Minor axis scale value (4 bytes)
		0, 0, 0, 0,
		// Points along parallel (4 bytes)
		0, 0, 0, 10,
		// Points along meridian (4 bytes)
		0, 0, 0, 10,
		// Basic angle (4 bytes)
		0, 0, 0, 0,
		// Subdivisions (4 bytes)
		0, 0, 0, 1,
		// First latitude (4 bytes, signed)
		0, 0, 0, 0,
		// First longitude (4 bytes, signed)
		0, 0, 0, 0,
		// Resolution and component flags
		0,
		// Last latitude (4 bytes, signed)
		0, 0, 0, 0,
		// Last longitude (4 bytes, signed)
		0, 0, 0, 0,
		// Parallel increment (4 bytes)
		0, 0, 0, 1,
		// Meridian increment (4 bytes)
		0, 0, 0, 1,
		// Scanning mode
		0,
		// South pole latitude (4 bytes, signed) - -90 degrees = -90000000
		0x85, 0x5D, 0x4A, 0x80,
		// South pole longitude (4 bytes, signed)
		0, 0, 0, 0,
		// Angle of rotation (4 bytes, signed)
		0, 0, 0, 0,
	}

	section := record.Section3{
		GridDefinitionTemplateNumber: 1,
		GridDefinitionTemplateData:   data,
	}

	template, err := grid.Template1{}.Parse(section)
	if err != nil {
		t.Fatal(err)
	}

	expected := grid.Template1{
		EarthShape:                  6,
		RadiusScaleFactor:           0,
		RadiusScaleValue:            0,
		MajorAxisScaleFactor:        0,
		MajorAxisScaleValue:         0,
		MinorAxisScaleFactor:        0,
		MinorAxisScaleValue:         0,
		PointsAlongParallel:         10,
		PointsAlongMeridian:         10,
		BasicAngle:                  0,
		Subdivisions:                1,
		FirstLatitude:               0,
		FirstLongitude:              0,
		ResolutionAndComponentFlags: 0,
		LastLatitude:                0,
		LastLongitude:               0,
		ParallelIncrement:           1,
		MeridianIncrement:           1,
		ScanningMode:                0,
		SouthPoleLatitude:           -90 * 1000000, // GRIB2 scaling
		SouthPoleLongitude:          0,
		AngleOfRotation:             0,
	}

	if typed := template.(grid.Template1); !reflect.DeepEqual(expected, typed) {
		t.Fatalf("expected\n%+v\nbut got\n%+v", expected, typed)
	}
}

func TestTemplate1Points(t *testing.T) {
	// Test with a simple 2x2 grid with south pole at -90,0 (no rotation)
	template := grid.Template1{
		PointsAlongParallel: 2,
		PointsAlongMeridian: 2,
		FirstLatitude:       0,
		FirstLongitude:      0,
		ParallelIncrement:   1000000, // 1 degree
		MeridianIncrement:   1000000, // 1 degree
		ScanningMode:        0,
		SouthPoleLatitude:   -90000000, // -90 degrees
		SouthPoleLongitude:  0,
		AngleOfRotation:     0,
	}

	points, err := template.Points()
	if err != nil {
		t.Fatal(err)
	}

	expected := 4 // 2x2 grid
	actual := len(points.Lats)
	if actual != expected {
		t.Fatalf("expected %v latitude points but got %v", expected, actual)
	}

	actual = len(points.Lngs)
	if actual != expected {
		t.Fatalf("expected %v longitude points but got %v", expected, actual)
	}

	// With south pole at -90,0, and identity rotation, (0,0) in rotated coords -> (0,0) in standard coords
	if value := points.Lats[0]; value < -0.1 || value > 0.1 {
		t.Fatalf("expected first latitude to be ~0 but got %v", value)
	}

	if value := points.Lngs[0]; value < -0.1 || value > 0.1 {
		t.Fatalf("expected first longitude to be ~0 but got %v", value)
	}
}
