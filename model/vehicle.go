package model

import (
	"math"
	"time"
)

// package const
const ()

// package errors
var (
	defaultSpeedMeasurement = "kmph"
)

// VehicleDef is a vehicle definition struct
type VehicleDef struct {
	ID        string
	Latitude  float64
	Longitude float64
	Speed     struct {
		Num         int
		Measurement string
	}
	UpdatedAt time.Time
}

// Vehicle is a vehicle model
type Vehicle struct {
	state State
}

// New return new Vehicle model
func New(s State) *Vehicle {
	return &Vehicle{s}
}

// Update updates vehicle state
func (v *Vehicle) Update(id string, lat, long float64, t time.Time) VehicleDef {
	vehicle := VehicleDef{
		ID:        id,
		Latitude:  lat,
		Longitude: long,
		UpdatedAt: t,
	}

	existed, ok := v.state.One(id)
	if ok {
		vehicle.Speed.Measurement = defaultSpeedMeasurement
		vehicle.Speed.Num = existed.speed(lat, long, t)
	}
	v.state.Update(vehicle)

	return vehicle
}

// All returns list of all vehicle ID
func (v *Vehicle) All() []string {
	return v.state.All()
}

func (v *VehicleDef) speed(lat, long float64, t time.Time) int {
	distance := v.distanceKM(lat, long)

	if distance == 0 {
		return 0
	}

	timePass := t.Second() - v.UpdatedAt.Second()

	v.Speed.Measurement = defaultSpeedMeasurement
	return int(distance/float64(timePass)*360) * 1
}

func (v VehicleDef) distanceKM(lat, long float64) float64 {
	const PI float64 = 3.141592653589793

	radlat1 := float64(PI * v.Latitude / 180)
	radlat2 := float64(PI * lat / 180)

	theta := float64(v.Longitude - long)
	radtheta := float64(PI * theta / 180)

	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)

	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / PI
	return dist * 60 * 1.1515 * 1.609344 * 1000
}
