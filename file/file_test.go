package file_test

import (
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/skysparq/grib2-go/file"
	"github.com/skysparq/grib2-go/templates"
	"github.com/skysparq/grib2-go/test_files"
)

func TestLoadGribFile(t *testing.T) {
	size, r, err := test_files.Load(test_files.FullGfsFile)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = r.Close() }()

	var peak uint64
	done := make(chan struct{})
	go func() {
		var m runtime.MemStats
		for {
			select {
			case <-done:
				return
			default:
				runtime.ReadMemStats(&m)
				if m.Alloc > peak {
					peak = m.Alloc
				}
				time.Sleep(1 * time.Millisecond)
			}
		}
	}()

	grib := file.NewGribFile(r, templates.Version33())
	gridDefs := make(map[int]int)
	prodDefs := make(map[int]int)
	dataRepDefs := make(map[int]int)
	totalLength := 0
	i := 1 // message number starts at 1
	for indexed, recErr := range grib.Records {
		if recErr != nil {
			t.Fatal(recErr)
		}
		if indexed.MessageNumber != indexed.MessageNumber {
			t.Fatalf(`expected message number %d but got %d`, i, indexed.MessageNumber)
		}
		rec := indexed.Record
		gridDefs[rec.Grid.GridDefinitionTemplateNumber]++
		prodDefs[rec.Product.ProductDefinitionTemplateNumber]++
		dataRepDefs[rec.DataRepresentation.DataRepresentationTemplateNumber]++
		totalLength += rec.Indicator.GribLength
		i++
	}

	if totalLength != size {
		t.Fatalf(`expected file size %d but got %d`, size, totalLength)
	}
	close(done)
	t.Logf("Peak memory usage (Alloc): %d bytes (%.2f MB)", peak, float64(peak)/1024/1024)
	t.Logf(`grid definition templates: %v, product definition templates: %v, data representation templates: %v`, gridDefs, prodDefs, dataRepDefs)
}

func TestExtractRecordFromFile(t *testing.T) {
	t.Skip(`only run when you need to isolate a record from an existing grib file for testing purposes`)
	_, r, err := test_files.Load(`.test_files/hrrr.t00z.wrfnatf01.grib2`)
	if err != nil {
		t.Fatal(err)
	}
	recBytes, err := file.ExtractRecordBytes(r, 1082)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(`../.test_files/hrrr-1082.grb2`, recBytes, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLoadTemplate31GribFile(t *testing.T) {
	size, r, err := test_files.Load(test_files.Template31TestFile)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = r.Close() }()

	grib := file.NewGribFile(r, templates.Version33())
	foundTemplate31 := false
	totalLength := 0
	for indexed, recErr := range grib.Records {
		if recErr != nil {
			t.Fatal(recErr)
		}
		rec := indexed.Record
		if rec.Grid.GridDefinitionTemplateNumber == 1 {
			foundTemplate31 = true
			gridDef, err := rec.Grid.Definition()
			if err != nil {
				t.Fatalf("error getting grid definition: %v", err)
			}
			points, err := gridDef.Points()
			if err != nil {
				t.Fatalf("error getting grid points: %v", err)
			}
			if len(points.Lats) == 0 || len(points.Lngs) == 0 {
				t.Fatalf("expected non-empty lat/lon arrays for template 3.1, got %d lats, %d lngs", len(points.Lats), len(points.Lngs))
			}
		}
		totalLength += rec.Indicator.GribLength
	}
	if !foundTemplate31 {
		t.Fatalf("did not find any records with grid definition template 3.1 in %s", test_files.Template31TestFile)
	}
	if totalLength != size {
		t.Fatalf("expected file size %d but got %d", size, totalLength)
	}
}
