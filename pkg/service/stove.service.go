package service

import (
	"log"
	"time"
)

const (
	turningOnStoveString  = "Turning On stove element"
	turningOffStoveString = "Turning Off stove element"
	stoveUpdateInterval        = time.Duration(1) * time.Second
	heatingRatePerSec     = 10                      // heating rate per second
	coolingRatePerSec     = 3                       // cooling rate per second
	minTemp               = 0 + coolingRatePerSec   // minimum temperature
	rangingMaxTemp        = 220 - heatingRatePerSec // ranging maximum temperature
	rangingMinTemp        = 180 + coolingRatePerSec // ranging minimum temperature
	normalOpTemp          = 200                     // normal operating temperature
)

type StoveService interface {
	On()
	Off()
	Temperature() int
	IsOn() bool
}

type stoveService struct {
	mainSwitch       MainSwitchService
	isOn             bool
	isRegulating     bool
	isRegulatingDown bool
	temperature      int
	ticker           *time.Ticker
}

func NewStoveService(mainSwitch MainSwitchService) StoveService {
	ticker := time.NewTicker(stoveUpdateInterval)
	ticker.Stop()

	stoveService := &stoveService{
		mainSwitch:       mainSwitch,
		isOn:             false,
		isRegulating:     false,
		isRegulatingDown: false,
		temperature:      0,
		ticker:           ticker,
	}

	mainSwitch.Connect(stoveService)

	return stoveService
}

// IsOn returns state of the stove
func (s *stoveService) IsOn() bool {
	return s.isOn
}

// Off switches off the stove
func (s *stoveService) Off() {
	s.turnOff()
}

// turnOff for turning off the stove
func (s *stoveService) turnOn() {
	log.Println(turningOnStoveString)
	s.temperature = 0
	s.isOn = true
	s.ticker.Reset(stoveUpdateInterval)
}

// turnOff for turning off the stove
func (s *stoveService) turnOff() {
	log.Println(turningOffStoveString)
	s.isOn = false
}

// toggleRegulating for switching on and off
// isRegulating and isRegulatingDown
func (s *stoveService) toggleRegulating() {
	if s.isOn {
		if s.temperature > normalOpTemp {
			s.isRegulating = true
		}

		if s.temperature > rangingMaxTemp {
			s.isRegulatingDown = true
		}

		if s.temperature < rangingMinTemp {
			s.isRegulatingDown = false
		}
	} else {
		s.isRegulating = false
		s.isRegulatingDown = false
	}
}

// regulatingHC - regulating heating and
// cooling above normalOpTemp
// for oscillating between rangingMaxTemp
// and rangingMinTemp
func (s *stoveService) regulatingHC() {
	if s.isRegulating {
		if s.isRegulatingDown {
			s.temperature -= coolingRatePerSec
		} else {
			s.temperature += heatingRatePerSec
		}
	}
}

// normalHC - normal heating and cooling
// responsible for heating and cooling the stove
// will continue increasing the temperature
// if isRegulating is not switched on by
// regulatingHC
func (s *stoveService) normalHC() {
	// is stove is on
	// and not yet regulating
	// increase temp
	if s.isOn {
		if !s.isRegulating {
			s.temperature += heatingRatePerSec
		}
	}

	// when stove is off
	// deduct temp until minTemp
	if !s.isOn {
		if s.temperature >= minTemp {
			s.temperature -= coolingRatePerSec
		} else {
			s.ticker.Stop()
			s.temperature = 0
		}
	}
}

// On switches on the stove
func (s *stoveService) On() {
	if !s.isOn && s.mainSwitch.IsOn() {
		s.turnOn()

		go func() {
			for {
				select {
				case <-s.ticker.C:
					s.normalHC()
					s.toggleRegulating()
					s.regulatingHC()
					log.Println("Stove is", s.IsOn(), "at", s.Temperature())
				}
			}
		}()
	}
}

// Temperature returns the current stove temperature
func (s *stoveService) Temperature() int {
	return s.temperature
}
