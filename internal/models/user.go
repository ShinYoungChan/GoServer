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

func LoginCheck(username, password string) (*User, error) {
	var user User
	if err := DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}

	// 2. 비밀번호 비교 (암호화된 값 vs 입력된 값)
	// bcrypt.CompareHashAndPassword(암호화된비번, 입력받은비번)
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		// 비밀번호가 틀린 경우
		return nil, err
	}

	// 3. 로그인 성공 시 유저 객체 반환
	return &user, nil
}
