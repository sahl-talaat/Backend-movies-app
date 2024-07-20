package auth

import (
	"errors"
	"myapp/config"
	"myapp/ent"
	"myapp/ent/user"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GenerateToken(c *gin.Context, userID uuid.UUID) (string, error) {
	token := uuid.New().String()

	u, err := config.Client.User.UpdateOneID(userID).SetToken(token).Save(c)
	if err != nil {
		return "", nil
	}
	return u.Token, nil
}

func ValidateToken(c *gin.Context, token string) (*ent.User, error) {
	u, err := config.Client.User.Query().Where(user.Token(token)).Only(c)
	if err != nil {
		return nil, errors.New("invalid token")
	}
	return u, nil
}

func RefreshToken(c *gin.Context, userID uuid.UUID) (string, error) {
	return GenerateToken(c, userID)
}

func DeleteToken(c *gin.Context, userID uuid.UUID) error {
	_, err := config.Client.User.UpdateOneID(userID).SetToken("").Save(c)
	return err
}

func ValidateTokenAndGetUserID(c *gin.Context, token string) (uuid.UUID, error) {
	userID, err := config.Client.User.Query().Where(user.TokenEQ(token)).OnlyID(c)
	if err != nil {
		return userID, nil
	}
	return uuid.Nil, errors.New("token validation failed")
}
