package models

import (
	"errors"
	"rakamin/helpers"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        int        `gorm:"primaryKey"`
	Username  string     `gorm:"not null"`
	Email     string     `gorm:"not null;unique"`
	Password  string     `gorm:"not null"`
	Photo     []Photo    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt *time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt *time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

type UserDBConnectionRepository struct {
	Conn *gorm.DB
}

type UserRepository interface {
	GetByEmail(email string) (user User, err error)
	Register(user User) (err error)
	GetById(id int) (user User, err error)
	UpdateById(id int, user User) (err error)
	DeleteById(id int) (err error)
}

func NewUserRepository(conn *gorm.DB) UserRepository {
	return &UserDBConnectionRepository{
		Conn: conn,
	}
}

func (repository *UserDBConnectionRepository) Register(user User) (err error) {
	user.Password, err = helpers.Hash(user.Password)
	if err != nil {
		return
	}

	_, err = repository.GetByEmail(user.Email)
	if err == nil {
		err = errors.New("duplicate email")
		return
	}

	err = repository.Conn.Create(&user).Error

	return
}

func (repository *UserDBConnectionRepository) GetByEmail(email string) (user User, err error) {
	err = repository.Conn.Where("email = ?", email).First(&user).Error

	return
}

func (repository *UserDBConnectionRepository) GetById(id int) (user User, err error) {
	err = repository.Conn.Where("id = ?", id).First(&user).Error

	return
}

func (repository *UserDBConnectionRepository) UpdateById(id int, user User) (err error) {
	if user.Password != "" {
		user.Password, err = helpers.Hash(user.Password)
		if err != nil {
			return
		}
	}
	err = repository.Conn.Where("id = ?", id).Updates(&user).Error

	return
}

func (repository *UserDBConnectionRepository) DeleteById(id int) (err error) {
	err = repository.Conn.Where("id = ?", id).Delete(&User{}).Error

	return
}
