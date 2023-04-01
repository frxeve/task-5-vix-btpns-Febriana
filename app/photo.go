package app

import "mime/multipart"

type PhotoRequest struct {
	Photo   *multipart.FileHeader `form:"file" binding:"required"`
	Title   string                `form:"title" binding:"required"`
	Caption string                `form:"caption" binding:"required"`
}

type GetAllPhotoByIdResponse struct {
	Photos []Photos `json:"photos"`
}

type Photos struct {
	ID       int
	Title    string
	Caption  string
	PhotoURL string
	UserID   int
}

type UpdatePhotoByIdRequest struct {
	Photo   *multipart.FileHeader `form:"file" binding:"required"`
	Title   string                `form:"title" binding:"required"`
	Caption string                `form:"caption" binding:"required"`
}

type DeletePhotoByIdRequest struct {
	ID int `uri:"photoId" binding:"required"`
}
