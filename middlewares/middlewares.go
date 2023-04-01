package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"rakamin/helpers"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type JwtCustomClaims struct {
	ID int `json:"id"`
	jwt.StandardClaims
}

type AuthorizationMiddleware struct {
	jwtSecret       string
	ExpiresDuration int
}

func NewAuthorizationMiddleware(jwtSecret string, expired int) *AuthorizationMiddleware {
	return &AuthorizationMiddleware{
		jwtSecret:       jwtSecret,
		ExpiresDuration: expired,
	}
}

func (a *AuthorizationMiddleware) Authorization() gin.HandlerFunc {
	return func(g *gin.Context) {
		authHeader := g.GetHeader("Authorization")
		if authHeader == "" {
			response := helpers.NewErrorResponse(errors.New("token tidak ditemukan"))
			g.JSON(http.StatusUnauthorized, response)

			return
		}
		_, err := a.ValidateToken(authHeader)
		if err != nil {
			response := helpers.NewErrorResponse(errors.New("token tidak valid"))
			g.JSON(http.StatusUnauthorized, response)

			return
		}
	}
}

func (a *AuthorizationMiddleware) GenerateToken(userID int) string {
	claims := &JwtCustomClaims{
		userID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(int64(a.ExpiresDuration))).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(a.jwtSecret))
	if err != nil {
		fmt.Println("error signed token :", err)
	}
	return t
}

func (a *AuthorizationMiddleware) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t_ *jwt.Token) (interface{}, error) {
		if _, ok := t_.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpexted signing method %v", t_.Header["alg"])
		}
		return []byte(a.jwtSecret), nil
	})
}

func (a *AuthorizationMiddleware) GetUserId(g *gin.Context) (id int, err error) {
	token := g.Request.Header.Get("Authorization")
	t, err := a.ValidateToken(token)
	if err != nil {
		return
	}

	claims, valid := t.Claims.(jwt.MapClaims)
	if !valid {
		return
	}

	id, _ = strconv.Atoi(fmt.Sprintf("%v", claims["id"]))

	return
}
