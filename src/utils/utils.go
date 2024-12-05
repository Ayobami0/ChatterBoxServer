package utils

import (
	"github.com/Ayobami0/chatter_box_server/src/model"
	"github.com/labstack/echo/v4"
)


func ErrorJson(c echo.Context, code int, message string) error {
	return c.JSON(code, model.NewHttpError(code, message))
}
