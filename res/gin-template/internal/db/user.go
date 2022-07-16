package db

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username       string `gorm:"type:varchar(31);uniqueIndex"`
	HashedPassword string `gorm:"type:varchar(127);not null"`
	FullName       string `json:"full_name"`
	Email          string `json:"email"`
}

func (u *User) TableName() string {
	return "users"
}

func CreateUser(username, hashedPassword, fullname, email string) (*User, error) {
	user := &User{
		Username:       username,
		HashedPassword: hashedPassword,
		FullName:       fullname,
		Email:          email,
	}

	if err := db.Write.Create(user).Error; err != nil {
		return nil, errors.Wrapf(err, "failed to create user")
	}

	return user, nil
}

func GetUser(username string) (*User, error) {
	user := &User{}

	if err := db.Read.Where("username = ?", username).First(user).Error; err != nil {
		return nil, errors.Wrapf(err, "failed to get user %v", username)
	}

	return user, nil
}
