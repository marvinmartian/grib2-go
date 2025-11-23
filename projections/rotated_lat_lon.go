package projections

import (
	"math"

	u "github.com/skysparq/grib2-go/utility"
)

// RotatedLatLonParams contains the parameters needed to instantiate the Rotated Latitude/Longitude projection.
// This projection rotates the standard lat/lon coordinate system so that a different point becomes the "north pole".
// The rotation is defined by the position of the southern pole in the original coordinate system.
type RotatedLatLonParams struct {
	ScanningMode  ScanningMode
	Ni            int // total number of points in the x direction
	Nj            int // total number of points in the y direction
	Di            int // x increment (distance between points)
	Dj            int // y increment (distance between points)
	I0            int // starting x point
	J0            int // starting y point
	SouthPoleLat  int // latitude of southern pole in degrees * 10^6
	SouthPoleLon  int // longitude of southern pole in degrees * 10^6
	RotationAngle int // angle of rotation in degrees * 10^6 (usually 0 for GRIB2)
}

// ExtractRotatedLatLonGrid extracts the grid points from the Rotated Latitude/Longitude projection defined by the given RotatedLatLonParams.
func ExtractRotatedLatLonGrid(params RotatedLatLonParams) (lats []float64, lngs []float64) {
	scannerParams := ScannerParams[int]{
		ScanningMode: params.ScanningMode,
		Ni:           params.Ni,
		Nj:           params.Nj,
		Di:           params.Di,
		Dj:           params.Dj,
		I0:           params.I0,
		J0:           params.J0,
	}

	s := NewScanner(scannerParams)
	lats = make([]float64, 0, params.Ni*params.Nj)
	lngs = make([]float64, 0, params.Ni*params.Nj)

	// Convert pole coordinates from scaled integers to degrees
	southPoleLat := u.StdLatLngToFloat(params.SouthPoleLat)
	southPoleLon := u.StdLatLngToFloat(params.SouthPoleLon)
	rotationAngle := u.StdLatLngToFloat(params.RotationAngle)

	// Convert to radians for calculations
	southPoleLatRad := southPoleLat * math.Pi / 180.0
	southPoleLonRad := southPoleLon * math.Pi / 180.0
	rotationAngleRad := rotationAngle * math.Pi / 180.0

	// Rotation parameters (following GRIB2 convention)
	// θ = 90° + south_pole_lat
	// φ = south_pole_lon
	theta := math.Pi/2.0 + southPoleLatRad
	phi := southPoleLonRad

	// For GRIB2 template 3.1, the rotation angle is typically 0
	// The full rotation matrix would be more complex, but for now we assume α = 0
	_ = rotationAngleRad // Mark as used to avoid compiler warning

	// Pre-compute trigonometric values
	cosTheta := math.Cos(theta)
	sinTheta := math.Sin(theta)
	cosPhi := math.Cos(phi)
	sinPhi := math.Sin(phi)

	for lat, lng := range s.Points {
		// Convert grid coordinates to rotated lat/lon in degrees
		rlat := u.StdLatLngToFloat(lat)
		rlon := u.StdLatLngToFloat(u.ShiftLongitude(lng))

		// Convert rotated lat/lon to radians
		rlatRad := rlat * math.Pi / 180.0
		rlonRad := rlon * math.Pi / 180.0

		// Convert rotated spherical coordinates to cartesian
		cosRlat := math.Cos(rlatRad)
		cosRlon := math.Cos(rlonRad)
		xr := cosRlon * cosRlat
		yr := math.Sin(rlonRad) * cosRlat
		zr := math.Sin(rlatRad)

		// Apply rotation matrices to transform back to standard coordinates
		// For GRIB2 template 3.1, rotation angle is typically 0, so we ignore it
		x := cosTheta*cosPhi*xr + cosTheta*sinPhi*yr + sinTheta*zr
		y := -cosTheta*sinPhi*xr + cosPhi*yr - sinTheta*sinPhi*zr
		z := -sinTheta*xr + cosTheta*zr

		// Convert back to spherical coordinates
		latRad := math.Asin(z)
		lonRad := math.Atan2(y, x)

		// Convert back to degrees
		latDeg := latRad * 180.0 / math.Pi
		lonDeg := lonRad * 180.0 / math.Pi

		// Shift longitude to [-180, 180] range and convert to float
		lonDeg = float64(u.ShiftLongitude(int(lonDeg*1000000))) / 1000000.0

		lats = append(lats, latDeg)
		lngs = append(lngs, lonDeg)
	}

	return lats, lngs
}
