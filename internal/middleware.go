package internal

import (
	"github.com/labstack/echo"
	"net/http"
)

func (s *Server) TokenValidator(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("token")
		if token != s.ValidToken {
			return echo.NewHTTPError(http.StatusForbidden, "token is invalid")
		}
		return next(c)
	}
}
