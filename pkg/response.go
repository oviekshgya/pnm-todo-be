package pkg

import "github.com/gofiber/fiber/v2"

type InitialResponse struct {
	Ctx *fiber.Ctx
}

type DataResponseSuccess struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Error   bool        `json:"error"`
}

type DataResponseError struct {
	Message string `json:"message"`
	Error   bool   `json:"error"`
}

func (c *InitialResponse) Respose(httpcode int, message string, error bool, data interface{}) error {
	if !error {
		return c.Ctx.Status(httpcode).JSON(&DataResponseSuccess{
			Data:    data,
			Message: message,
			Error:   error,
		})
	} else {
		return c.Ctx.Status(httpcode).JSON(&DataResponseError{
			Message: message,
			Error:   error,
		})
	}
}
