package main

import (
	"rakamin/controllers"
	"rakamin/database"
	"rakamin/helpers"
	"rakamin/middlewares"
	"rakamin/models"
	"rakamin/router"

	"github.com/gin-gonic/gin"
)

func main() {
	configApp := helpers.GetConfig()
	mysqlConfig := database.ConfigDB{
		Username: configApp.Mysql.User,
		Password: configApp.Mysql.Pass,
		Host:     configApp.Mysql.Host,
		Port:     configApp.Mysql.Port,
		Database: configApp.Mysql.Name,
	}
	mysqlDB := mysqlConfig.ConfigDB()
	authMiddleware := middlewares.NewAuthorizationMiddleware(configApp.JWT.Secret, configApp.JWT.Expired)

	userRepo := models.NewUserRepository(mysqlDB)
	userController := controllers.NewUserController(userRepo, authMiddleware)
	photoRepo := models.NewPhotoRepository(mysqlDB)
	photoController := controllers.NewPhotoController(photoRepo, authMiddleware)
	
	r := gin.Default()
	router := router.ControllerList{
		AuthMiddleware:  authMiddleware,
		UserController:  *userController,
		PhotoController: *photoController,
	}
	
	router.RouteRegister(r)

	r.Run()
}
