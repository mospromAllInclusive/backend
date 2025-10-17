package auth

import (
	"backend/src/domains/entities"
	"backend/src/services"
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const tokenExpirationTime = 7 * 24 * time.Hour

type service struct {
	jwtKey       []byte
	usersService services.IUsersService
}

func NewService(usersService services.IUsersService) services.IAuthService {
	jwtKeyString := os.Getenv("JWT_KEY")
	return &service{
		usersService: usersService,
		jwtKey:       []byte(jwtKeyString),
	}
}

func (s *service) JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		claims := &claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return s.jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Next()
	}
}

func (s *service) Login(ctx context.Context, email string, password string) (*entities.User, string, error) {
	user, err := s.usersService.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, "", err
	}

	if user.Password != password {
		return nil, "", ErrorWrongPassword{}
	}

	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *service) generateToken(userID int64) (string, error) {
	expirationTime := time.Now().Add(tokenExpirationTime)
	claims := &claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtKey)
}

type claims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}
