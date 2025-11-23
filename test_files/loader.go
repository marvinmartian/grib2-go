package test_files

import (
	"fmt"
	"os"
	"path/filepath"
)

type TestFile string

const (
	SingleRecordProdDef8                    TestFile = `.test_files/single-grib2-record-prod-def-8.grb2`
	SingleRecordProdDef0                    TestFile = `.test_files/single-grib2-record-prod-def-0.grb2`
	SingleRecordDataDef0Bitmap              TestFile = `.test_files/single-grib2-record-data-def-0-bitmap.grb2`
	FullGfsFile                             TestFile = `.test_files/full-gfs-file.grb2`
	SingleRecordDataDef40                   TestFile = `.test_files/single-grib2-record-data-def-40.grb2`
	SingleRecordDataDef41                   TestFile = `.test_files/single-grib2-record-data-def-41.grb2`
	SingleRecordGridDef30                   TestFile = `.test_files/single-grib2-record-grid-def-30.grb2`
	SingleRecordGridDef40                   TestFile = `.test_files/single-grib2-record-grid-def-40.grb2`
	MrmsCompositeRefl                       TestFile = `.test_files/MRMS_MergedReflectivityQCComposite_00.50_20251005-081241.grb2`
	MrmsLghtngProb                          TestFile = `.test_files/MRMS_LightningProbabilityNext30minGrid_scale_1_20251005-113039.grb2`
	MrmsAzShear                             TestFile = `.test_files/MRMS_MergedAzShear_0-2kmAGL_00.50_20251005-112817.grb2`
	SingleRecordDataDef3Bitmap              TestFile = `.test_files/single-grib2-data-def-3-bitmap.grb2`
	SingleRecordDataDef3PrimaryMissingValue TestFile = `.test_files/single-grib2-data-def-3-primary-missing-value-mgnmnt.grb2`
	Template31TestFile                      TestFile = `.test_files/20251123T00Z_MSC_HRDPS_ABSV_ISBL_0700_RLatLon0.0225_PT000H.grib2`
)

func Load(file TestFile) (int, *os.File, error) {
	wd, err := os.Getwd()
	if err != nil {
		return 0, nil, fmt.Errorf("error loading test file: %v", err)
	}
	for {
		if wd == `/` {
			return 0, nil, fmt.Errorf("did not find project root directory, project directory must be called 'grib2-go'")
		}
		if filepath.Base(wd) == `grib2-go` {
			break
		}
		wd = filepath.Dir(wd)
	}

	path := filepath.Join(wd, string(file))

	stat, err := os.Stat(path)
	if err != nil {
		return 0, nil, fmt.Errorf("error loading test file: %v", err)
	}
	size := int(stat.Size())

	f, err := os.Open(path)
	if err != nil {
		return size, f, fmt.Errorf("error loading test file: %v", err)
	}
	return size, f, nil
}
