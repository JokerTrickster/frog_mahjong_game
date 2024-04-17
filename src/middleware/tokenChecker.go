package middleware

import (
	"context"
	"main/utils"
	"main/utils/db/mysql"
	"time"

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
		// db에서 유효한 토큰인지 체크
		if CheckDBAccessToken(uID, accessToken) != nil {
			return utils.ErrorMsg(ctx, utils.ErrBadParameter, utils.Trace(), "invalid access token", utils.ErrFromClient)
		}

		// set token data to Context
		c.Set("uID", uID)
		c.Set("email", email)

		return next(c)

	}
}

func CheckDBAccessToken(uID uint, accessToken string) error {
	token := mysql.Tokens{
		UserID:      uID,
		AccessToken: accessToken,
	}
	err := mysql.GormMysqlDB.Model(&token).Where("user_id = ? AND access_token = ?", uID, accessToken).First(&token).Error
	if err != nil {
		return utils.ErrorMsg(context.TODO(), utils.ErrBadToken, utils.Trace(), "invalid access token", utils.ErrFromClient)
	}
	now := time.Now()
	if token.RefreshExpiredAt < now.Unix() {
		return utils.ErrorMsg(context.TODO(), utils.ErrBadToken, utils.Trace(), "refresh token expired", utils.ErrFromClient)
	}

	return nil
}
