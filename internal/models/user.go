package models

import (
	"golang.org/x/crypto/bcrypt"
)

// 회원가입: 새 사용자 생성
func CreateUser(username, password string) error {
	// 1. 비밀번호 암호화 (해싱)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 2. DB 저장
	user := User{Username: username, Password: string(hashedPassword)}
	return DB.Create(&user).Error
}

// ID 찾기
func GetUserByUsername(username string) (*User, error) {
	var user User
	if err := DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
