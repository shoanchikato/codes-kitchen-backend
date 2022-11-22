package service

import "log"

const (
	mainSwitchOnString  = "Main switch is on"
	mainSwitchOffString = "Main switch is off"
)

type MainSwitchService interface {
	On()
	Off()
	IsOn() bool
	Connect(Appliance)
}

type mainSwitchService struct {
	isOn     bool
	isOnChan chan bool
	appliances []Appliance
}

func NewMainSwitchService() MainSwitchService {
	return &mainSwitchService{
		isOn: false,
	}
}

// On turns on the main switch
func (m *mainSwitchService) On() {
	m.isOn = true
	log.Println(mainSwitchOnString)
}

// Off turns off the main switch
func (m *mainSwitchService) Off() {
	m.isOn = false
	log.Println(mainSwitchOffString)

	for _, applicance := range(m.appliances) {
		applicance.Off()
	}
}

// IsOn returns the is on value
func (m *mainSwitchService) IsOn() bool {
	return m.isOn
}

// SubscribeIsOn adds an appliance to the list of 
// appliance that are using the main switch
func (m *mainSwitchService) Connect(appliance Appliance) {
	log.Println("added an appliance", appliance)
	m.appliances = append(m.appliances, appliance)
}
