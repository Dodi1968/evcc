package core

import (
	"fmt"
	"math"

	"github.com/evcc-io/evcc/api"
	"github.com/evcc-io/evcc/core/types"
)

const geoLocationRadius = 50 // Maximum vehicle distance from loadpoint (m)

// isVehicleAtHome checks whether vehicle is at home (geofencing)
// false: if vehicle position is available and outside the defined radius
// true: in all other cases, even in cases of error or if position is not available
func (site *Site) isVehicleAtHome(vehicle api.Vehicle) bool {
	geoLocation := site.GetGeoLocation()

	if !geoLocation.Enabled || vehicle == nil {
		return true
	}

	if err := validateGeoLocation(geoLocation); err != nil { // validate again (e.g. for yaml config)
		site.log.ERROR.Println(err)
		return true
	}

	vp, ok := api.Cap[api.VehiclePosition](vehicle)
	if !ok {
		site.log.DEBUG.Println("vehicle does not support position tracking")
		return true
	}

	lat, lon, err := vp.Position()
	if err != nil {
		site.log.INFO.Printf("vehicle position: %v", err)
		return true
	}

	// {lat: 0, lon: 0} is a point in the atlantic ocean, not accessible by car
	// zero values indicate that the vehicle is not sending any valid position data
	if lat == 0 && lon == 0 {
		site.log.INFO.Println("vehicle position: not available")
		return true
	}

	site.log.DEBUG.Printf("vehicle position: lat %.4f, lon %.4f", lat, lon)

	atHome := isAtHome(geoLocation, lat, lon)

	site.log.DEBUG.Printf(
		"vehicle distance from loadpoint: %.1fm (radius: %vm, atHome=%v)",
		distance(geoLocation.Lat, geoLocation.Lon, lat, lon)*1e3,
		geoLocationRadius,
		atHome,
	)

	return atHome
}

// validate geolocation settings
func validateGeoLocation(geoLocation types.GeoLocation) error {
	if geoLocation.Enabled {
		if (geoLocation.Lat == 0 && geoLocation.Lon == 0) || // geolocation enabled without setting coordinates
			math.Abs(geoLocation.Lat) > 90 ||
			math.Abs(geoLocation.Lon) > 180 {
			return fmt.Errorf("invalid geolocation settings: %+v", geoLocation)
		}
	}
	return nil
}

// given a valid geolocation and vehicle coords, is it at home?
func isAtHome(geoLocation types.GeoLocation, lat, lon float64) bool {
	d := distance(geoLocation.Lat, geoLocation.Lon, lat, lon) * 1e3
	return d <= geoLocationRadius
}

// calculate the distance between two points on a globe
func distance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371 // km

	// Differences in radiant
	dLat := (lat2 - lat1) * math.Pi / 180.0
	dLon := (lon2 - lon1) * math.Pi / 180.0

	lat1 = lat1 * math.Pi / 180.0
	lat2 = lat2 * math.Pi / 180.0

	// Haversine formula
	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(lat1)*math.Cos(lat2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return c * earthRadius
}
