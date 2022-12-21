package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/mohnaofal/rest-go-jwt/models"
)

// jwtCustomClaims are custom claims extending default ones.
// See https://github.com/golang-jwt/jwt for more examples
type JwtCustomClaims struct {
	jwt.StandardClaims
	UserID int `json:"user_id"`
}

func GenerateJwt(user *models.User) string {
	// Set custom claims
	claims := JwtCustomClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
		int(user.ID),
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return ``
	}

	return t
}

func VerifyJwt() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenString := strings.TrimPrefix(c.Request().Header.Get(echo.HeaderAuthorization), "Bearer ")

			if len(tokenString) == 0 {
				return echo.NewHTTPError(http.StatusUnauthorized, "authorization is invalid")
			}

			// Parse takes the token string and a function for looking up the key. The latter is especially
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// Don't forget to validate the alg is what you expect:
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}

				// Secret is a []byte containing your secret, e.g. []byte("my_secret_key")
				return []byte(os.Getenv("SECRET")), nil
			})
			if err != nil {
				switch ve := err.(type) {
				case *jwt.ValidationError:
					switch ve.Errors {
					case jwt.ValidationErrorExpired, jwt.ValidationErrorNotValidYet:
						return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized expired")
					case jwt.ValidationErrorMalformed:
						return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized malformed")
					}
				}
				return echo.NewHTTPError(http.StatusUnauthorized, "token parsing error")
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				c.Set("user_id", int(claims["user_id"].(float64)))
			} else {
				return echo.NewHTTPError(http.StatusUnauthorized, "authorization is invalid")
			}

			return next(c)
		}
	}
}
