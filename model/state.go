//go:generate mockgen -destination mock/mock_user.go github.com/fleetState/model VehicleDef

package model

// State is an interface of vehicle state store methods
type State interface {
	All() []string
	One(id string) (VehicleDef, bool)
	Update(v VehicleDef)
}
