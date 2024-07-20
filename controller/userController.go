package controller

import (
	"myapp/auth"
	"myapp/config"
	"myapp/ent"
	"myapp/ent/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Signup(c *gin.Context) {
	var user_input ent.User
	err := c.BindJSON(&user_input)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	u, err := config.Client.User.Create().
		SetFirstName(user_input.FirstName).
		SetLastName(user_input.LastName).
		SetEmail(user_input.Email).
		SetPassword(user_input.Password).
		SetAge(user_input.Age).Save(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	token, err := auth.GenerateToken(c, u.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "user": u, "token": token})
}

func Signin(c *gin.Context) {
	var data struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := c.BindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, data)
		return
	}

	u, err := config.Client.User.Query().Where(user.EmailEQ(data.Email)).Only(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	if data.Password != u.Password {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "wrong password"})
		return
	}

	u.Token, err = auth.GenerateToken(c, u.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "user": u, "token": u.Token})
}

func Signout(c *gin.Context) {
	var requestBody struct {
		Token string `json:"token"`
	}
	err := c.BindJSON(&requestBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := auth.ValidateToken(c, requestBody.Token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	err = auth.DeleteToken(c, u.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "signed out successfully"})
}

func GetAllUsers(c *gin.Context) {
	users, err := config.Client.User.Query().All(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "users": users})
}
