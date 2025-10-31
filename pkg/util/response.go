package util

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/quantsmithapp/datastation-backend/pkg/errors"
)

type HttpResponse struct {
	Success   bool        `json:"success"`
	ErrorCode string      `json:"error_code,omitempty"`
	Result    interface{} `json:"result,omitempty"`
}

func Response(c *fiber.Ctx, statusCode int, errorCode string, result interface{}) error {
	success := statusCode == http.StatusOK || statusCode == http.StatusCreated || statusCode == http.StatusAccepted
	payload := HttpResponse{
		Success:   success,
		ErrorCode: errorCode,
		Result:    result,
	}

	return c.Status(statusCode).JSON(payload)
}

func ResponseOK(c *fiber.Ctx, result interface{}) error {
	return Response(c, http.StatusOK, "", result)
}

func ResponseCreated(c *fiber.Ctx, result interface{}) error {
	return Response(c, http.StatusCreated, "", result)
}

func ResponseError(c *fiber.Ctx, err error) error {
	switch e := err.(type) {
	case errors.AppError:
		return Response(c, e.StatusCode, e.ErrorCode, nil)
	default:
		return Response(c, http.StatusInternalServerError, "APP500", nil)
	}
}
