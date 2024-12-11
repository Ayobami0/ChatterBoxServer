package middleware

import (
	"fmt"
	"net/http"

	"github.com/Ayobami0/chatter_box_server/src/utils"
	"github.com/labstack/echo/v4"
)

func WSAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.QueryParam("t")

		if token == "" {
			return utils.ErrorJson(c, http.StatusUnauthorized, "token param 't' is required")
		}

		c.Request().Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		return next(c)
	}
}
