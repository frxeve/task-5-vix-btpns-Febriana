package models

import "gorm.io/gorm"

type Photo struct {
	ID       int    `gorm:"primary_key;auto_increment"`
	Title    string `gorm:"not null"`
	Caption  string `gorm:"not null"`
	PhotoURL string `gorm:"not null"`
	UserID   int    `gorm:"not null"`
	User     *User
}

type PhotoDBConnectionRepository struct {
	Conn *gorm.DB
}

type PhotoRepository interface {
	Insert(photo Photo) (err error)
	GetAllByUserId(id int) (photos []Photo, err error)
	UpdatePhotoById(photo Photo) (err error)
	DeletePhotoById(userId, photoId int) (err error)
}

func NewPhotoRepository(conn *gorm.DB) PhotoRepository {
	return &PhotoDBConnectionRepository{
		Conn: conn,
	}
}

func (repository *PhotoDBConnectionRepository) Insert(photo Photo) (err error) {
	err = repository.Conn.Create(&photo).Error

	return
}

func (repository *PhotoDBConnectionRepository) GetAllByUserId(id int) (photos []Photo, err error) {
	err = repository.Conn.Where("user_id = ?", id).Find(&photos).Error

	return
}

func (repository *PhotoDBConnectionRepository) UpdatePhotoById(photo Photo) (err error) {
	err = repository.Conn.Where("id = ? AND user_id = ?", photo.ID, photo.UserID).Updates(&photo).Error

	return
}

func (repository *PhotoDBConnectionRepository) DeletePhotoById(userId, photoId int) (err error) {
	err = repository.Conn.Where("id = ? AND user_id = ?", photoId, userId).Delete(&Photo{}).Error

	return
}
