package route

import (
	"myapp/controller"

	"github.com/gin-gonic/gin"
)

func RegisterRouter(r *gin.Engine) {
	v1 := r.Group("api/react/v1")
	// user
	v1.GET("users", controller.GetAllUsers)

	v1.POST("users/signup", controller.Signup)
	v1.POST("users/signin", controller.Signin)
	v1.POST("users/signout", controller.Signout)

	v1.POST("users/addToFavorites", controller.AddToFavorites)
	v1.GET("users/getFavorites", controller.GetFavorites)
}
