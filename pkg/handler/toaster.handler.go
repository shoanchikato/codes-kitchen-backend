package handler

import (
	ss "codes-kitchen/pkg/service"
	"encoding/json"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type ToasterHandler interface {
	Switch() fiber.Handler
	IsOn() fiber.Handler
}

type toasterHandler struct {
	service ss.ToasterService
	ticker  *time.Ticker
}

func NewToasterHandler(service ss.ToasterService) ToasterHandler {
	ticker := time.NewTicker(refreshInterval)

	return &toasterHandler{
		service: service,
		ticker:  ticker,
	}
}

// Switch switches toaster on or off
func (t *toasterHandler) Switch() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var turnSwitch TurnSwitch
		c.BodyParser(&turnSwitch)

		if turnSwitch.TurnOn {
			t.service.Toast()

			log.Println("Turning on toaster")
			return c.JSON(TurnSwitch{t.service.IsOn()})
		} else {
			t.service.Cancel()

			log.Println("Popped toast and switching off toaster")
			return c.JSON(TurnSwitch{t.service.IsOn()})
		}
	}
}

// IsOn open a websocket broadcasting if the toaster is 
// on or off
func (t *toasterHandler) IsOn() func(*fiber.Ctx) error {
	return websocket.New(func(c *websocket.Conn) {

	outer:
		for {
			err := t.readMessage(c)
			if err != nil {
				break outer
			}

			for {
				select {
				case <-t.ticker.C:
					err := t.writeMessage(c)
					if err != nil {
						break outer
					}
				}
			}
		}
	})
}

func (t *toasterHandler) writeMessage(c *websocket.Conn) error {
	data, err := json.Marshal(map[string]any{
		"isOn": t.service.IsOn(),
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

func (t *toasterHandler) readMessage(c *websocket.Conn) error {
	mt, msg, err := c.ReadMessage()
	if err != nil {
		log.Println("Error reading message from client", err)
		return err
	}
	log.Printf("Received message from client: %s %d", msg, mt)

	return nil
}
