package router

import (
	"rakamin/controllers"
	"rakamin/middlewares"

	"github.com/gin-gonic/gin"
)

type ControllerList struct {
	AuthMiddleware  *middlewares.AuthorizationMiddleware
	UserController  controllers.UserController
	PhotoController controllers.PhotoController
}

func (cl *ControllerList) RouteRegister(g *gin.Engine) {
	g.Static("/public/images", "./public/images")
	apiV1 := g.Group("api/v1")

	user := apiV1.Group("/users")
	user.POST("/register", cl.UserController.Register)
	user.POST("/login", cl.UserController.Login)
	user.GET("/", cl.UserController.GetUserById, cl.AuthMiddleware.Authorization())
	user.PUT("/", cl.UserController.UpdateUserById, cl.AuthMiddleware.Authorization())
	user.DELETE("/", cl.UserController.DeleteUserById, cl.AuthMiddleware.Authorization())

	photo := apiV1.Group("/photos", cl.AuthMiddleware.Authorization())
	photo.GET("/", cl.PhotoController.GetPhotos)
	photo.POST("/", cl.PhotoController.Upload)
	photo.PUT("/:photoId", cl.PhotoController.UpdatePhotoById)
	photo.DELETE("/:photoId", cl.PhotoController.DeletePhotoById)
}
