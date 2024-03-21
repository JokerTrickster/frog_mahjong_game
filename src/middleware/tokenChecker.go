package middleware

import (
	"main/utils"

	"github.com/labstack/echo/v4"
)

// CheckJWT : check user's jwt token from "token" header value
func TokenChecker(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		// get jwt Token
		accessToken := c.Request().Header.Get("tkn")
		if accessToken == "" {
			return utils.ErrorMsg(ctx, utils.ErrBadParameter, utils.Trace(), "no access token in header", utils.ErrFromClient)
		}

		// verify & get Data
		uID, email, err := utils.ValidateAndParseAccessToken(accessToken)
		if err != nil {
			return utils.ErrorMsg(ctx, utils.ErrBadParameter, utils.Trace(), "invalid access token", utils.ErrFromClient)
		}

		// set token data to Context
		c.Set("uID", uID)
		c.Set("email", email)

		return next(c)

	}
}
