module github.com/marvinmartian/grib2-go

go 1.24

require (
	github.com/segmed/openjpeg/gojp2 v0.0.0-20230309051650-206b95b2b9a4
	github.com/skysparq/grib2-go v0.4.8
)

require (
	github.com/google/uuid v1.3.0 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20220927061507-ef77025ab5aa // indirect
	golang.org/x/sys v0.2.0 // indirect
	modernc.org/libc v1.21.5 // indirect
	modernc.org/mathutil v1.5.0 // indirect
	modernc.org/memory v1.4.0 // indirect
)

replace github.com/skysparq/grib2-go => ./
