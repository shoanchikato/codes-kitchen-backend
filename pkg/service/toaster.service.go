package service

import (
	"log"
	"time"
)

const (
	turningOnToasterString  = "Turning on toaster"
	turningOffToasterString = "Popped toast and switching off toaster"
	toastDuration         = time.Duration(45) * time.Second
)

type ToasterService interface {
	Toast()
	Cancel()
	IsOn() bool
	On()
	Off()
}

type toasterService struct {
	mainSwitch MainSwitchService
	isOn   bool
	ticker *time.Ticker
}

func NewToasterService(mainSwitch MainSwitchService) ToasterService {
	ticker := time.NewTicker(toastDuration)
	ticker.Stop()

	toasterService := &toasterService{
		mainSwitch: mainSwitch,
		isOn:   false,
		ticker: ticker,
	}

	mainSwitch.Connect(toasterService)

	return toasterService
}

// Cancel pops the bread and switches
// of the toaster
func (t *toasterService) Cancel() {
	t.turnOff()
}

// IsOn returns the state of the toaster
func (t *toasterService) IsOn() bool {
	return t.isOn
}

func (t *toasterService) turnOn() {
	t.isOn = true
	t.ticker.Reset(toastDuration)
	log.Println(turningOnToasterString)
}

func (t *toasterService) turnOff() {
	t.isOn = false
	log.Println(turningOffToasterString)
}

// Toast switches on the toaster
func (t *toasterService) Toast() {

	go func() {
		if !t.isOn && t.mainSwitch.IsOn() {
			t.turnOn()

			for {
				select {
				case <-t.ticker.C:
					t.turnOff()
					return
				}
			}
		}
	}()
}

func(t *toasterService) On() {
	t.Toast()
}

func (t *toasterService) Off() {
	t.Cancel()
}
