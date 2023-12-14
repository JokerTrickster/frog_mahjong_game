package _interface

import "github.com/labstack/echo/v4"

type ISignupAuthHandler interface {
	Signup(c echo.Context) error
}
