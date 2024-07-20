package controller

import (
	"myapp/auth"
	"myapp/config"
	"myapp/ent/movie"
	"myapp/ent/user"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AddToFavorites(c *gin.Context) {
	var input struct {
		MovieName string    `json:"movieName"`
		ImgURL    string    `json:"imgUrl"`
		UserID    uuid.UUID `json:"userID"`
		MovieID   string    `json:"movieID"`
	}
	err := c.BindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
		return
	}

	userID, err := auth.ValidateTokenAndGetUserID(c, token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	if userID != input.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "User ID mismatch"})
		return
	}

	u, err := config.Client.User.Get(c, input.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	movieIDInt, err := strconv.Atoi(input.MovieID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	m, err := config.Client.Movie.Query().Where(movie.ID(movieIDInt)).Only(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}

	_, err = u.Update().AddFavoriteMovies(m).Save(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add movie to favorites"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Movie added to favorites successfully"})
}

type FavoriteMovieResponse struct {
	ID        string `json:"_id"`
	MovieName string `json:"movieName"`
	ImgURL    string `json:"imgUrl"`
	UserID    string `json:"userID"`
	MovieID   string `json:"movieID"`
}

func GetFavorites(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
		return
	}

	// Validate token and get user ID
	userID, err := auth.ValidateTokenAndGetUserID(c, token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Fetch user's favorite movies
	u, err := config.Client.User.Query().
		Where(user.ID(userID)).
		WithFavoriteMovies().
		Only(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	/* // Prepare response
	   var favorites []FavoriteMovieResponse
	   for _, movie := range u.Edges.FavoriteMovies {
	       favorites = append(favorites, FavoriteMovieResponse{
	           ID:        movie.ID.String(),
	           MovieName: movie.MovieName,
	           ImgURL:    movie.ImgUrl,
	           UserID:    userID,
	           MovieID:   movie.Edges.ID.String(),
	       })
	   } */

	c.JSON(http.StatusOK, gin.H{"message": "success", "Favorites": u})
}
