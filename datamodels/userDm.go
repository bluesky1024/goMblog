package datamodels

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	Id           int64
	Uid          int64
	NickName     string
	Password     string
	Telephone    string
	Email        string
	ProfileImage string
	FollowsCount int64
	FriendsCount int64
	CreateTime   time.Time `xorm:"created"`
	UpdateTime   time.Time `xorm:"updated"`
}

func (u User) IsValid() bool {
	return u.Id > 0
}

func GeneratePassword(userPassword string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash[:]), nil
}

func ValidatePassword(userPassword string, hashed string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(userPassword)); err != nil {
		return false, err
	}
	return true, nil
}
