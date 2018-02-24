package handlers

import (
	"encoding/base64"
	"strings"

	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type (
	Auth struct {
		jwtSecret string
		user      string
		password  string
	}
	User struct {
		jwt.StandardClaims
		Name string
	}
)

func NewAuth(jwtSecret, user, password string) *Auth {
	return &Auth{jwtSecret, user, password}
}

func (a *Auth) Authorized(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := a.User(c)
		if err != nil {
			return err
		}
		return next(c)
	}
}

func (a *Auth) Login(c echo.Context) error {
	user := c.FormValue("user")
	password := c.FormValue("password")
	if err := a.checkUser(c, user, password); err != nil {
		return err
	}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), &User{
		Name: user,
	})
	tokenstring, err := token.SignedString([]byte(a.jwtSecret))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, echo.Map{
		"token": tokenstring,
	})
}

func (a *Auth) checkUser(c echo.Context, user, password string) error {
	if user == a.user && password == a.password {
		return nil
	}
	return echo.NewHTTPError(http.StatusForbidden, "user or password incorrect.")
}

func (a *Auth) User(c echo.Context) (*User, error) {
	authorization := c.Request().Header.Get("Authorization")
	if authorization == "" {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "not logged in")
	}
	typeAndAuth := strings.Split(authorization, " ")
	if len(typeAndAuth) != 2 {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Invalid Authorization header.")
	}
	t := typeAndAuth[0]
	au := typeAndAuth[1]
	switch t {
	case "Bearer":
		user := User{}
		token, err := jwt.ParseWithClaims(au, &user, func(token *jwt.Token) (interface{}, error) {
			return []byte(a.jwtSecret), nil
		})
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusBadRequest, "Can not parse token")
		}
		if !token.Valid {
			return nil, echo.NewHTTPError(http.StatusBadRequest, "Invalid token")
		}
		return &user, nil
	case "Basic":
		userAndPassword, err := base64.URLEncoding.DecodeString(au)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusBadRequest, "Invalid Basic Authorization.")
		}
		userAndPasswordString := string(userAndPassword)
		userAndPasswordArray := strings.Split(userAndPasswordString, ":")
		if len(userAndPasswordArray) != 2 {
			return nil, echo.NewHTTPError(http.StatusBadRequest, "Invalid Basic Authorization.")
		}
		if err := a.checkUser(c, userAndPasswordArray[0], userAndPasswordArray[1]); err != nil {
			return nil, err
		}
		return &User{
			Name: userAndPasswordArray[0],
		}, nil
	}
	return nil, echo.NewHTTPError(http.StatusBadRequest, "Unsupported Authorization header.")
}

func (a *Auth) Info(c echo.Context) error {
	u, err := a.User(c)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, u)
}
