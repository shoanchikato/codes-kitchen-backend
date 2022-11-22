package service

const (
	onString  = "On"
	offString = "Off"
)

type Appliance interface {
	On()
	Off()
}
