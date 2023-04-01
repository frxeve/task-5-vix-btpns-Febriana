package controllers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"rakamin/app"
	"rakamin/helpers"
	"rakamin/middlewares"
	"rakamin/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PhotoController struct {
	photoRepo      models.PhotoRepository
	AuthMiddleware *middlewares.AuthorizationMiddleware
}

func NewPhotoController(photoRepo models.PhotoRepository, authMiddleware *middlewares.AuthorizationMiddleware) *PhotoController {
	return &PhotoController{
		photoRepo:      photoRepo,
		AuthMiddleware: authMiddleware,
	}
}

func (controller *PhotoController) Upload(g *gin.Context) {
	var (
		request app.PhotoRequest
		err     error
		id      int
	)

	id, err = controller.AuthMiddleware.GetUserId(g)
	if err != nil {
		return
	}

	err = g.ShouldBind(&request)
	if err != nil {
		response := helpers.NewErrorResponse(err)
		g.AbortWithStatusJSON(http.StatusBadRequest, response)

		return
	}

	src, err := request.Photo.Open()
	if err != nil {
		response := helpers.NewErrorResponse(errors.New("can't open file"))
		g.AbortWithStatusJSON(http.StatusInternalServerError, response)

		return
	}

	defer src.Close()

	fileBytes, err := ioutil.ReadAll(src)
	if err != nil {
		response := helpers.NewErrorResponse(errors.New("invalid file"))
		g.AbortWithStatusJSON(http.StatusInternalServerError, response)

		return
	}

	filetype := http.DetectContentType(fileBytes)
	if filetype != "image/jpeg" && filetype != "image/jpg" &&
		filetype != "image/gif" && filetype != "image/png" {
		response := helpers.NewErrorResponse(errors.New("invalid file type"))
		g.AbortWithStatusJSON(http.StatusInternalServerError, response)

		return
	}

	fileFormat, err := mime.ExtensionsByType(filetype)
	if err != nil {
		response := helpers.NewErrorResponse(errors.New("can't read File Type"))
		g.AbortWithStatusJSON(http.StatusInternalServerError, response)

		return
	}

	path := fmt.Sprintf("%s%s%s", "public/images/", helpers.GetUUID(), fileFormat[len(fileFormat)-1])

	err = controller.photoRepo.Insert(models.Photo{
		Title:    request.Title,
		Caption:  request.Caption,
		PhotoURL: path,
		UserID:   id,
	})
	if err != nil {
		response := helpers.NewErrorResponse(err)
		g.JSON(http.StatusInternalServerError, response)

		return
	}

	g.SaveUploadedFile(request.Photo, path)

	response := helpers.NewSuccessInsertResponse(nil)
	g.JSON(http.StatusCreated, response)
}

func (controller *PhotoController) GetPhotos(g *gin.Context) {
	var (
		err   error
		id    int
		photo app.Photos
		res   app.GetAllPhotoByIdResponse
	)

	id, err = controller.AuthMiddleware.GetUserId(g)
	if err != nil {
		return
	}

	data, err := controller.photoRepo.GetAllByUserId(id)
	if err != nil {
		response := helpers.NewErrorResponse(err)
		g.JSON(http.StatusInternalServerError, response)

		return
	}

	for _, value := range data {
		photo.ID = value.ID
		photo.Title = value.Title
		photo.Caption = value.Caption
		photo.PhotoURL = value.PhotoURL
		photo.UserID = value.UserID
		res.Photos = append(res.Photos, photo)
	}

	response := helpers.NewSuccessResponse(res)
	g.JSON(http.StatusOK, response)
}

func (controller *PhotoController) UpdatePhotoById(g *gin.Context) {
	var (
		err error
		id  int
		req app.UpdatePhotoByIdRequest
	)

	id, err = controller.AuthMiddleware.GetUserId(g)
	if err != nil {
		return
	}

	err = g.ShouldBind(&req)
	if err != nil {
		response := helpers.NewErrorResponse(err)
		g.AbortWithStatusJSON(http.StatusBadRequest, response)

		return
	}

	pid := g.Param("photoId")

	photoId, _ := strconv.Atoi(pid)
	if err != nil {
		response := helpers.NewErrorResponse(errors.New("cant parse photo id"))
		g.AbortWithStatusJSON(http.StatusBadRequest, response)

		return
	}
	src, err := req.Photo.Open()
	if err != nil {
		response := helpers.NewErrorResponse(errors.New("can't open file"))
		g.AbortWithStatusJSON(http.StatusInternalServerError, response)

		return
	}

	defer src.Close()

	fileBytes, err := ioutil.ReadAll(src)
	if err != nil {
		response := helpers.NewErrorResponse(errors.New("invalid file"))
		g.AbortWithStatusJSON(http.StatusInternalServerError, response)

		return
	}

	filetype := http.DetectContentType(fileBytes)
	if filetype != "image/jpeg" && filetype != "image/jpg" &&
		filetype != "image/gif" && filetype != "image/png" {
		response := helpers.NewErrorResponse(errors.New("invalid file type"))
		g.AbortWithStatusJSON(http.StatusInternalServerError, response)

		return
	}

	fileFormat, err := mime.ExtensionsByType(filetype)
	if err != nil {
		response := helpers.NewErrorResponse(errors.New("can't read File Type"))
		g.AbortWithStatusJSON(http.StatusInternalServerError, response)

		return
	}

	path := fmt.Sprintf("%s%s%s", "public/images/", helpers.GetUUID(), fileFormat[len(fileFormat)-1])

	err = controller.photoRepo.UpdatePhotoById(models.Photo{
		ID:       photoId,
		Title:    req.Title,
		Caption:  req.Caption,
		PhotoURL: path,
		UserID:   id,
	})
	if err != nil {
		response := helpers.NewErrorResponse(err)
		g.JSON(http.StatusInternalServerError, response)

		return
	}

	g.SaveUploadedFile(req.Photo, path)
	
	response := helpers.NewSuccessResponse(nil)
	g.JSON(http.StatusOK, response)
}

func (controller *PhotoController) DeletePhotoById(g *gin.Context) {
	var (
		err error
		id  int
		req app.DeletePhotoByIdRequest
	)

	id, err = controller.AuthMiddleware.GetUserId(g)
	if err != nil {
		return
	}

	err = g.ShouldBindUri(&req)
	if err != nil {
		response := helpers.NewErrorResponse(err)
		g.AbortWithStatusJSON(http.StatusBadRequest, response)

		return
	}
	err = controller.photoRepo.DeletePhotoById(id, req.ID)
	if err != nil {
		response := helpers.NewErrorResponse(err)
		g.JSON(http.StatusInternalServerError, response)

		return
	}

	response := helpers.NewSuccessResponse(nil)
	g.JSON(http.StatusOK, response)
}
