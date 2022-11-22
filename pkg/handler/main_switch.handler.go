package handler

import (
	ss "codes-kitchen/pkg/service"
	"encoding/json"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type MainSwitchHandler interface {
	Switch() fiber.Handler
	IsOn() fiber.Handler
}

type mainSwitchHandler struct {
	service ss.MainSwitchService
	ticker  *time.Ticker
}

func NewMainSwitchHandler(service ss.MainSwitchService) MainSwitchHandler {
	ticker := time.NewTicker(refreshInterval)

	return &mainSwitchHandler{
		service: service,
		ticker:  ticker,
	}
}

// Switch switches main switch on or off
func (m *mainSwitchHandler) Switch() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var turnSwitch TurnSwitch
		c.BodyParser(&turnSwitch)

		if turnSwitch.TurnOn {
			m.service.On()

			log.Println("switched on mainSwitch")
			return c.JSON(TurnSwitch{m.service.IsOn()})
		} else {
			m.service.Off()

			log.Println("switched off mainSwitch")
			return c.JSON(TurnSwitch{m.service.IsOn()})
		}
	}
}

// IsOn open a websocket broading is the main switch is 
// on or off
func (m *mainSwitchHandler) IsOn() func(*fiber.Ctx) error {
	return websocket.New(func(c *websocket.Conn) {

	outer:
		for {
			err := m.readMessage(c)
			if err != nil {
				break outer
			}

			for {
				select {
				case <-m.ticker.C:
					err := m.writeMessage(c)
					if err != nil {
						break outer
					}
				}
			}
		}
	})
}

func (m *mainSwitchHandler) writeMessage(c *websocket.Conn) error {
	data, err := json.Marshal(map[string]any{
		"isOn": m.service.IsOn(),
	})
	if err != nil {
		log.Println("Error marshalling data", err)
		return err
	}

	err = c.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		log.Println("Error writing message:", err)
		return err
	}

	return nil
}

func (m *mainSwitchHandler) readMessage(c *websocket.Conn) error {
	mt, msg, err := c.ReadMessage()
	if err != nil {
		log.Println("Error reading message from client", err)
		return err
	}
	log.Printf("Received message from client: %s %d", msg, mt)

	return nil
}
