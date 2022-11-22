package handler

import (
	ss "codes-kitchen/pkg/service"
	"encoding/json"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type StoveHandler interface {
	Switch() fiber.Handler
	Temperature() fiber.Handler
}

type stoveHandler struct {
	service ss.StoveService
	ticker  *time.Ticker
}

func NewStoveHandler(service ss.StoveService) StoveHandler {
	ticker := time.NewTicker(refreshInterval)
	ticker.Stop()

	return &stoveHandler{
		service: service,
		ticker:  ticker,
	}
}

// Switch switches stove on or off
func (s *stoveHandler) Switch() fiber.Handler {
	return func(c *fiber.Ctx) error {
		turnSwitch := new(TurnSwitch)
		if err := c.BodyParser(turnSwitch); err != nil {
			return c.Status(400).JSON(err.Error())
		}

		if turnSwitch.TurnOn {
			s.service.On()
			s.ticker.Reset(refreshInterval)

			log.Println("switched on stove")
			return c.JSON(TurnSwitch{s.service.IsOn()})
		} else {
			s.service.Off()
			
			log.Println("switched off stove")
			return c.JSON(TurnSwitch{s.service.IsOn()})
		}
	}
}

// Temperature open a websocket broadcasting stove temperature
// values
func (s *stoveHandler) Temperature() func(*fiber.Ctx) error {
	return websocket.New(func(c *websocket.Conn) {

	outer:
		for {
			err := s.readMessage(c)
			if err != nil {
				break outer
			}

			for {
				select {
				case <-s.ticker.C:
					err := s.writeMessage(c)
					if err != nil {
						break outer
					}
				}
			}
		}
	})
}

func (s *stoveHandler) writeMessage(c *websocket.Conn) error {
	data, err := json.Marshal(map[string]any{
		"isOn":        s.service.IsOn(),
		"temperature": s.service.Temperature(),
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

func (s *stoveHandler) readMessage(c *websocket.Conn) error {
	mt, msg, err := c.ReadMessage()
	if err != nil {
		log.Println("Error reading message from client", err)
		return err
	}
	log.Printf("Received message from client: %s %d", msg, mt)

	return nil
}
