package handlers

import "github.com/labstack/echo"

func parseRequest(c echo.Context, i interface{}) error {
	if err := c.Bind(i); err != nil {
		return err
	}
	if err := c.Validate(i); err != nil {
		return err
	}

	return nil
}
