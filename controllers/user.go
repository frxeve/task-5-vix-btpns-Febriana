package controllers

import (
	"errors"
	"net/http"
	"rakamin/app"
	"rakamin/helpers"
	"rakamin/middlewares"
	"rakamin/models"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userRepo       models.UserRepository
	AuthMiddleware *middlewares.AuthorizationMiddleware
}

func NewUserController(userRepo models.UserRepository, authMiddleware *middlewares.AuthorizationMiddleware) *UserController {
	return &UserController{
		userRepo:       userRepo,
		AuthMiddleware: authMiddleware,
	}
}

func (controller *UserController) Register(g *gin.Context) {
	var (
		request app.RegisterRequest
		err     error
	)

	err = g.ShouldBind(&request)
	if err != nil {
		response := helpers.NewErrorResponse(err)
		g.JSON(http.StatusBadRequest, response)

		return
	}

	err = controller.userRepo.Register(models.User{
		Username: request.Username,
		Email:    request.Email,
		Password: request.Password,
	})

	if err != nil {
		response := helpers.NewErrorResponse(err)
		g.JSON(http.StatusInternalServerError, response)

		return
	}

	response := helpers.NewSuccessInsertResponse(nil)
	g.JSON(http.StatusCreated, response)
}

func (controller *UserController) Login(g *gin.Context) {
	var (
		request  app.LoginRequest
		response app.LoginResponse
		err      error
	)

	err = g.ShouldBind(&request)
	if err != nil {
		response := helpers.NewErrorResponse(err)
		g.JSON(http.StatusBadRequest, response)

		return
	}

	data, err := controller.userRepo.GetByEmail(request.Email)
	if err != nil {
		response := helpers.NewErrorResponse(err)
		g.JSON(http.StatusInternalServerError, response)

		return
	}

	if !helpers.ValidateHash(request.Password, data.Password) {
		if err != nil {
			response := helpers.NewErrorResponse(errors.New("wrong email and password"))
			g.JSON(http.StatusUnauthorized, response)

			return
		}
	}

	token := controller.AuthMiddleware.GenerateToken(data.ID)

	response.Token = token
	res := helpers.NewSuccessResponse(response)
	g.JSON(http.StatusOK, res)
}

func (controller *UserController) GetUserById(g *gin.Context) {
	var (
		id  int
		err error
		res app.GetUserByIdResponse
	)

	id, err = controller.AuthMiddleware.GetUserId(g)
	if err != nil {
		return
	}

	data, err := controller.userRepo.GetById(id)
	if err != nil {
		response := helpers.NewErrorResponse(err)
		g.JSON(http.StatusInternalServerError, response)

		return
	}

	res.Username = data.Username
	res.Email = data.Email
	res.CreatedAt = *data.CreatedAt

	response := helpers.NewSuccessResponse(res)
	g.JSON(http.StatusOK, response)
}

func (controller *UserController) UpdateUserById(g *gin.Context) {
	var (
		id  int
		err error
		req app.UpdateUserByIdRequest
	)

	id, err = controller.AuthMiddleware.GetUserId(g)
	if err != nil {
		return
	}

	err = g.ShouldBind(&req)
	if err != nil {
		response := helpers.NewErrorResponse(err)
		g.JSON(http.StatusBadRequest, response)

		return
	}

	err = controller.userRepo.UpdateById(id, models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		response := helpers.NewErrorResponse(err)
		g.JSON(http.StatusInternalServerError, response)

		return
	}

	response := helpers.NewSuccessResponse(nil)
	g.JSON(http.StatusOK, response)
}

func (controller *UserController) DeleteUserById(g *gin.Context) {
	var (
		id  int
		err error
	)

	id, err = controller.AuthMiddleware.GetUserId(g)
	if err != nil {
		return
	}

	err = controller.userRepo.DeleteById(id)
	if err != nil {
		response := helpers.NewErrorResponse(err)
		g.JSON(http.StatusInternalServerError, response)

		return
	}

	response := helpers.NewSuccessResponse(nil)
	g.JSON(http.StatusOK, response)
}
